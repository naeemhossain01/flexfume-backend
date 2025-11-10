package services

import (
	"errors"

	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"gorm.io/gorm"
)

var (
	ErrDiscountNotFound          = errors.New("discount not found")
	ErrDiscountProductRequired   = errors.New("product ID is required for discount")
	ErrDiscountPercentageInvalid = errors.New("discount percentage must be between 0 and 100")
	ErrDiscountPriceRequired     = errors.New("discount price is required and must be greater than 0")
	ErrDiscountAlreadyExists     = errors.New("discount already exists for this product")
)

// DiscountService handles discount-related business logic
type DiscountService struct {
	productService ProductServiceInterface
}

// NewDiscountService creates a new discount service
func NewDiscountService(productService ProductServiceInterface) *DiscountService {
	return &DiscountService{
		productService: productService,
	}
}

// AddDiscounts creates multiple discounts
func (s *DiscountService) AddDiscounts(discounts []models.Discount) ([]models.Discount, error) {
	if len(discounts) == 0 {
		return nil, errors.New("no discounts provided")
	}

	// Validate each discount
	for i := range discounts {
		if discounts[i].ProductID == "" {
			return nil, ErrDiscountProductRequired
		}

		if discounts[i].DiscountPrice <= 0 {
			return nil, ErrDiscountPriceRequired
		}

		if discounts[i].Percentage < 0 || discounts[i].Percentage > 100 {
			return nil, ErrDiscountPercentageInvalid
		}

		// Verify product exists
		_, err := s.productService.GetProductByID(discounts[i].ProductID)
		if err != nil {
			return nil, err
		}

		// Check if discount already exists for this product
		var existingDiscount models.Discount
		if err := database.GetDB().Where("product_id = ?", discounts[i].ProductID).First(&existingDiscount).Error; err == nil {
			return nil, ErrDiscountAlreadyExists
		}
	}

	// Create all discounts
	if err := database.GetDB().Create(&discounts).Error; err != nil {
		return nil, err
	}

	// Load product relationships
	for i := range discounts {
		if err := database.GetDB().Preload("Product").First(&discounts[i], "id = ?", discounts[i].ID).Error; err != nil {
			return nil, err
		}
	}

	return discounts, nil
}

// UpdateDiscounts updates multiple discounts by product ID
func (s *DiscountService) UpdateDiscounts(discounts []models.Discount) ([]models.Discount, error) {
	if len(discounts) == 0 {
		return nil, errors.New("no discounts provided")
	}

	var updatedDiscounts []models.Discount

	// Update each discount
	for _, item := range discounts {
		if item.ProductID == "" {
			return nil, ErrDiscountProductRequired
		}

		if item.DiscountPrice <= 0 {
			return nil, ErrDiscountPriceRequired
		}

		if item.Percentage < 0 || item.Percentage > 100 {
			return nil, ErrDiscountPercentageInvalid
		}

		// Find discount by product ID
		var discount models.Discount
		if err := database.GetDB().Where("product_id = ?", item.ProductID).First(&discount).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, ErrDiscountNotFound
			}
			return nil, err
		}

		// Update discount fields
		discount.DiscountPrice = item.DiscountPrice
		discount.Percentage = item.Percentage
		if err := database.GetDB().Save(&discount).Error; err != nil {
			return nil, err
		}

		// Load product relationship
		if err := database.GetDB().Preload("Product").First(&discount, "id = ?", discount.ID).Error; err != nil {
			return nil, err
		}

		updatedDiscounts = append(updatedDiscounts, discount)
	}

	return updatedDiscounts, nil
}

// GetAllDiscounts retrieves all discounts
func (s *DiscountService) GetAllDiscounts() ([]models.Discount, error) {
	var discounts []models.Discount
	if err := database.GetDB().Preload("Product").Find(&discounts).Error; err != nil {
		return nil, err
	}
	return discounts, nil
}

// GetDiscountByProductID retrieves a discount by product ID
func (s *DiscountService) GetDiscountByProductID(productID string) (*models.Discount, error) {
	if productID == "" {
		return nil, errors.New("product ID is required")
	}

	var discount models.Discount
	if err := database.GetDB().Preload("Product").Where("product_id = ?", productID).First(&discount).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrDiscountNotFound
		}
		return nil, err
	}

	return &discount, nil
}

// DeleteDiscount deletes a discount by ID
func (s *DiscountService) DeleteDiscount(discountID string) error {
	if discountID == "" {
		return errors.New("discount ID is required")
	}

	// Check if discount exists
	var discount models.Discount
	if err := database.GetDB().First(&discount, "id = ?", discountID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrDiscountNotFound
		}
		return err
	}

	// Delete discount
	if err := database.GetDB().Delete(&discount).Error; err != nil {
		return err
	}

	return nil
}

// Ensure DiscountService implements DiscountServiceInterface
var _ DiscountServiceInterface = (*DiscountService)(nil)
