package services

import (
	"errors"
	"math"
	"time"

	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"gorm.io/gorm"
)

var (
	ErrCouponUsageNotFound    = errors.New("coupon usage not found")
	ErrCouponUsageLimitExceeded = errors.New("coupon usage limit exceeded")
	ErrCouponInactive         = errors.New("coupon is not active")
	ErrCouponExpired          = errors.New("coupon has expired")
	ErrMinOrderAmountNotMet   = errors.New("minimum order amount not met")
	ErrInvalidCartItems       = errors.New("invalid cart items")
)

// CouponUsageService handles coupon usage-related business logic
type CouponUsageService struct {
	couponService  CouponServiceInterface
	discountService DiscountServiceInterface
}

// NewCouponUsageService creates a new coupon usage service
func NewCouponUsageService(couponService CouponServiceInterface, discountService DiscountServiceInterface) *CouponUsageService {
	return &CouponUsageService{
		couponService:  couponService,
		discountService: discountService,
	}
}

// ApplyCoupon applies a coupon to cart items and returns the final discounted amount
func (s *CouponUsageService) ApplyCoupon(cartIDs []string, couponCode string, userID string) (float64, error) {
	if len(cartIDs) == 0 {
		return 0, ErrInvalidCartItems
	}

	if userID == "" {
		return 0, errors.New("user ID is required")
	}

	// Get coupon by code
	coupon, err := s.couponService.GetCouponByCode(couponCode)
	if err != nil {
		return 0, err
	}

	// Validate coupon
	if !coupon.Active {
		return 0, ErrCouponInactive
	}

	if time.Now().After(coupon.ExpirationTime) {
		return 0, ErrCouponExpired
	}

	// Check if user has already used this coupon and if usage limit is exceeded
	var couponUsage models.CouponUsage
	err = database.GetDB().Where("coupon_id = ? AND user_id = ?", coupon.ID, userID).First(&couponUsage).Error
	
	if err == nil {
		// Coupon usage exists, check limit
		if couponUsage.UsageCount >= coupon.UsageLimit {
			return 0, ErrCouponUsageLimitExceeded
		}
	} else if err != gorm.ErrRecordNotFound {
		return 0, err
	}

	// Get cart items with product and discount information
	cartItems, err := s.getCartItemsData(cartIDs)
	if err != nil {
		return 0, err
	}

	// Calculate total amount after product discounts
	productTotalAmount := s.calculateCartTotal(cartItems)

	// Check minimum order amount
	if productTotalAmount < coupon.MinOrderAmount {
		return 0, ErrMinOrderAmountNotMet
	}

	// Calculate coupon discount amount
	couponDiscountAmount := s.calculateCouponDiscount(coupon, productTotalAmount)

	// Apply max discount limit if set
	if coupon.MaxAmountApplied > 0 && couponDiscountAmount > coupon.MaxAmountApplied {
		couponDiscountAmount = coupon.MaxAmountApplied
	}

	// Calculate final amount
	finalAmount := productTotalAmount - couponDiscountAmount

	// Round to 2 decimal places
	finalAmount = math.Round(finalAmount*100) / 100
	couponDiscountAmount = math.Round(couponDiscountAmount*100) / 100

	// Save or update coupon usage
	if err == gorm.ErrRecordNotFound {
		// Create new coupon usage
		couponUsage = models.CouponUsage{
			CouponID:         coupon.ID,
			UserID:           userID,
			UsageCount:       1,
			DiscountedAmount: couponDiscountAmount,
		}
		if err := database.GetDB().Create(&couponUsage).Error; err != nil {
			return 0, err
		}
	} else {
		// Update existing coupon usage
		couponUsage.UsageCount++
		couponUsage.DiscountedAmount = couponDiscountAmount
		if err := database.GetDB().Save(&couponUsage).Error; err != nil {
			return 0, err
		}
	}

	return finalAmount, nil
}

// RemoveCouponUsage removes a coupon usage record
func (s *CouponUsageService) RemoveCouponUsage(couponCode string, userID string) error {
	if couponCode == "" {
		return errors.New("coupon code is required")
	}

	if userID == "" {
		return errors.New("user ID is required")
	}

	// Get coupon by code
	coupon, err := s.couponService.GetCouponByCode(couponCode)
	if err != nil {
		return err
	}

	// Find coupon usage
	var couponUsage models.CouponUsage
	if err := database.GetDB().Where("coupon_id = ? AND user_id = ?", coupon.ID, userID).First(&couponUsage).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrCouponUsageNotFound
		}
		return err
	}

	// Delete coupon usage
	if err := database.GetDB().Delete(&couponUsage).Error; err != nil {
		return err
	}

	return nil
}

// GetCouponUsage retrieves a coupon usage record
func (s *CouponUsageService) GetCouponUsage(couponCode string, userID string) (*models.CouponUsage, error) {
	if couponCode == "" {
		return nil, errors.New("coupon code is required")
	}

	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get coupon by code
	coupon, err := s.couponService.GetCouponByCode(couponCode)
	if err != nil {
		return nil, err
	}

	// Find coupon usage
	var couponUsage models.CouponUsage
	if err := database.GetDB().Where("coupon_id = ? AND user_id = ?", coupon.ID, userID).First(&couponUsage).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrCouponUsageNotFound
		}
		return nil, err
	}

	return &couponUsage, nil
}

// getCartItemsData retrieves cart items with product and discount information
func (s *CouponUsageService) getCartItemsData(cartIDs []string) ([]models.CartItemData, error) {
	type QueryResult struct {
		CartID          string
		Quantity        int
		ProductPrice    float64
		DiscountPercent int
	}

	var results []QueryResult
	
	// Query to get cart items with product prices and discounts
	err := database.GetDB().Raw(`
		SELECT 
			c.id as cart_id,
			c.quantity,
			p.price as product_price,
			COALESCE(d.percentage, 0) as discount_percent
		FROM carts c
		INNER JOIN products p ON c.product_id = p.id
		LEFT JOIN discounts d ON p.id = d.product_id
		WHERE c.id IN ?
	`, cartIDs).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, ErrInvalidCartItems
	}

	// Convert to CartItemData
	cartItems := make([]models.CartItemData, len(results))
	for i, result := range results {
		cartItems[i] = models.CartItemData{
			CartID:          result.CartID,
			Quantity:        result.Quantity,
			ProductPrice:    result.ProductPrice,
			DiscountPercent: result.DiscountPercent,
		}
	}

	return cartItems, nil
}

// calculateCartTotal calculates the total cart amount after product discounts
func (s *CouponUsageService) calculateCartTotal(cartItems []models.CartItemData) float64 {
	total := 0.0

	for _, item := range cartItems {
		// Calculate product price after discount
		productPrice := item.ProductPrice
		if item.DiscountPercent > 0 {
			discountAmount := productPrice * float64(item.DiscountPercent) / 100.0
			productPrice = productPrice - discountAmount
		}

		// Multiply by quantity
		itemTotal := productPrice * float64(item.Quantity)
		total += itemTotal
	}

	return math.Round(total*100) / 100
}

// calculateCouponDiscount calculates the discount amount based on coupon type
func (s *CouponUsageService) calculateCouponDiscount(coupon *models.Coupon, amount float64) float64 {
	var discountAmount float64

	switch coupon.CouponType {
	case "PERCENTAGE":
		discountAmount = amount * coupon.Amount / 100.0
	case "FIXED":
		discountAmount = coupon.Amount
	default:
		discountAmount = 0
	}

	return math.Round(discountAmount*100) / 100
}

// Ensure CouponUsageService implements CouponUsageServiceInterface
var _ CouponUsageServiceInterface = (*CouponUsageService)(nil)
