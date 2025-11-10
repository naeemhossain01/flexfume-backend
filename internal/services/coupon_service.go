package services

import (
	"errors"
	"strings"
	"time"

	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"gorm.io/gorm"
)

var (
	ErrCouponNotFound       = errors.New("coupon not found")
	ErrCouponCodeRequired   = errors.New("coupon code is required")
	ErrCouponTypeRequired   = errors.New("coupon type is required")
	ErrCouponAmountRequired = errors.New("coupon amount is required")
	ErrInvalidCouponType    = errors.New("invalid coupon type, must be PERCENTAGE or FIXED")
	ErrCouponCodeExists     = errors.New("coupon code already exists")
	ErrCouponHasUsage       = errors.New("cannot delete coupon that has been used")
)

// CouponService handles coupon-related business logic
type CouponService struct{}

// NewCouponService creates a new coupon service
func NewCouponService() *CouponService {
	return &CouponService{}
}

// CreateCoupon creates a new coupon
func (s *CouponService) CreateCoupon(coupon *models.Coupon) (*models.Coupon, error) {
	if coupon == nil {
		return nil, errors.New("coupon cannot be nil")
	}

	// Validate required fields
	coupon.Code = strings.TrimSpace(coupon.Code)
	if coupon.Code == "" {
		return nil, ErrCouponCodeRequired
	}

	coupon.CouponType = strings.ToUpper(strings.TrimSpace(coupon.CouponType))
	if coupon.CouponType == "" {
		return nil, ErrCouponTypeRequired
	}

	if coupon.CouponType != "PERCENTAGE" && coupon.CouponType != "FIXED" {
		return nil, ErrInvalidCouponType
	}

	if coupon.Amount <= 0 {
		return nil, ErrCouponAmountRequired
	}

	// Check if coupon code already exists
	var existingCoupon models.Coupon
	if err := database.GetDB().Where("code = ?", coupon.Code).First(&existingCoupon).Error; err == nil {
		return nil, ErrCouponCodeExists
	}

	// Create coupon
	if err := database.GetDB().Create(coupon).Error; err != nil {
		return nil, err
	}

	return coupon, nil
}

// UpdateCoupon updates an existing coupon
func (s *CouponService) UpdateCoupon(couponID string, updatedCoupon *models.Coupon) (*models.Coupon, error) {
	if couponID == "" {
		return nil, errors.New("coupon ID is required")
	}

	if updatedCoupon == nil {
		return nil, errors.New("coupon data cannot be nil")
	}

	// Get existing coupon
	coupon, err := s.GetCouponByID(couponID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if updatedCoupon.Code != "" {
		updatedCoupon.Code = strings.TrimSpace(updatedCoupon.Code)
		// Check if new code already exists (excluding current coupon)
		var existingCoupon models.Coupon
		if err := database.GetDB().Where("code = ? AND id != ?", updatedCoupon.Code, couponID).First(&existingCoupon).Error; err == nil {
			return nil, ErrCouponCodeExists
		}
		coupon.Code = updatedCoupon.Code
	}

	if updatedCoupon.CouponType != "" {
		updatedCoupon.CouponType = strings.ToUpper(strings.TrimSpace(updatedCoupon.CouponType))
		if updatedCoupon.CouponType != "PERCENTAGE" && updatedCoupon.CouponType != "FIXED" {
			return nil, ErrInvalidCouponType
		}
		coupon.CouponType = updatedCoupon.CouponType
	}

	if updatedCoupon.Amount > 0 {
		coupon.Amount = updatedCoupon.Amount
	}

	if updatedCoupon.MinOrderAmount >= 0 {
		coupon.MinOrderAmount = updatedCoupon.MinOrderAmount
	}

	if updatedCoupon.MaxAmountApplied >= 0 {
		coupon.MaxAmountApplied = updatedCoupon.MaxAmountApplied
	}

	if !updatedCoupon.ExpirationTime.IsZero() {
		coupon.ExpirationTime = updatedCoupon.ExpirationTime
	}

	if updatedCoupon.UsageLimit > 0 {
		coupon.UsageLimit = updatedCoupon.UsageLimit
	}

	// Update active status
	coupon.Active = updatedCoupon.Active

	// Save changes
	if err := database.GetDB().Save(coupon).Error; err != nil {
		return nil, err
	}

	return coupon, nil
}

// GetCouponByID retrieves a coupon by ID
func (s *CouponService) GetCouponByID(couponID string) (*models.Coupon, error) {
	if couponID == "" {
		return nil, errors.New("coupon ID is required")
	}

	var coupon models.Coupon
	if err := database.GetDB().First(&coupon, "id = ?", couponID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrCouponNotFound
		}
		return nil, err
	}

	return &coupon, nil
}

// GetCouponByCode retrieves a coupon by code
func (s *CouponService) GetCouponByCode(code string) (*models.Coupon, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, errors.New("coupon code is required")
	}

	var coupon models.Coupon
	if err := database.GetDB().Where("code = ?", code).First(&coupon).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrCouponNotFound
		}
		return nil, err
	}

	return &coupon, nil
}

// GetAllCoupons retrieves all coupons
func (s *CouponService) GetAllCoupons() ([]models.Coupon, error) {
	var coupons []models.Coupon
	if err := database.GetDB().Find(&coupons).Error; err != nil {
		return nil, err
	}
	return coupons, nil
}

// DeleteCoupon deletes a coupon by ID
func (s *CouponService) DeleteCoupon(couponID string) error {
	if couponID == "" {
		return errors.New("coupon ID is required")
	}

	// Check if coupon exists
	coupon, err := s.GetCouponByID(couponID)
	if err != nil {
		return err
	}

	// Check if coupon has been used
	var usageCount int64
	if err := database.GetDB().Model(&models.CouponUsage{}).Where("coupon_id = ?", couponID).Count(&usageCount).Error; err != nil {
		return err
	}

	if usageCount > 0 {
		return ErrCouponHasUsage
	}

	// Delete coupon
	if err := database.GetDB().Delete(coupon).Error; err != nil {
		return err
	}

	return nil
}

// GetCouponUsageStatistics retrieves usage statistics for a coupon
func (s *CouponService) GetCouponUsageStatistics(couponID string) (totalUsage int, totalSavings float64, uniqueUsers int, err error) {
	type UsageStats struct {
		TotalUsage   int
		TotalSavings float64
		UniqueUsers  int
	}

	var stats UsageStats
	err = database.GetDB().Model(&models.CouponUsage{}).
		Select("COUNT(*) as total_usage, 0 as total_savings, COUNT(DISTINCT user_id) as unique_users").
		Where("coupon_id = ?", couponID).
		Scan(&stats).Error

	if err != nil {
		return 0, 0, 0, err
	}

	return stats.TotalUsage, stats.TotalSavings, stats.UniqueUsers, nil
}

// EnrichCouponWithStatistics adds usage statistics to a coupon response
func (s *CouponService) EnrichCouponWithStatistics(coupon *models.Coupon) models.CouponResponse {
	response := coupon.ToResponse()
	
	totalUsage, totalSavings, uniqueUsers, err := s.GetCouponUsageStatistics(coupon.ID)
	if err == nil {
		response.TotalUsageCount = totalUsage
		response.TotalSavingsAmount = totalSavings
		response.UniqueUsersCount = uniqueUsers
	}
	
	return response
}

// ValidateCoupon checks if a coupon is valid for use
func (s *CouponService) ValidateCoupon(coupon *models.Coupon) error {
	if !coupon.Active {
		return errors.New("coupon is not active")
	}

	if time.Now().After(coupon.ExpirationTime) {
		return errors.New("coupon has expired")
	}

	return nil
}

// Ensure CouponService implements CouponServiceInterface
var _ CouponServiceInterface = (*CouponService)(nil)
