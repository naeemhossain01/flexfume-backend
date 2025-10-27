package services

import (
	"errors"
	"strings"

	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound         = errors.New("product not found")
	ErrProductNameRequired     = errors.New("product name is required")
	ErrProductPriceRequired    = errors.New("product price is required")
	ErrProductCategoryRequired = errors.New("product category is required")
	ErrNoProductsFoundByCategory = errors.New("no products found for this category")
	ErrNoProductsFoundBySearch = errors.New("no products found matching search criteria")
)

// ProductService handles product-related business logic
type ProductService struct {
	categoryService CategoryServiceInterface
}

// NewProductService creates a new product service
func NewProductService(categoryService CategoryServiceInterface) *ProductService {
	return &ProductService{
		categoryService: categoryService,
	}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(product *models.Product) (*models.Product, error) {
	if product == nil {
		return nil, errors.New("product cannot be nil")
	}

	// Validate required fields
	product.ProductName = strings.TrimSpace(product.ProductName)
	if product.ProductName == "" {
		return nil, ErrProductNameRequired
	}

	if product.Price <= 0 {
		return nil, ErrProductPriceRequired
	}

	if product.CategoryID == "" {
		return nil, ErrProductCategoryRequired
	}

	// Verify category exists
	_, err := s.categoryService.GetCategoryByID(product.CategoryID)
	if err != nil {
		return nil, err
	}

	// Create product
	if err := database.GetDB().Create(product).Error; err != nil {
		return nil, err
	}

	// Load category and discount relationships
	if err := database.GetDB().Preload("Category").Preload("Discount").First(product, "id = ?", product.ID).Error; err != nil {
		return nil, err
	}

	return product, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(productID string, updatedProduct *models.Product) (*models.Product, error) {
	if productID == "" {
		return nil, errors.New("product ID is required")
	}

	if updatedProduct == nil {
		return nil, errors.New("product data cannot be nil")
	}

	// Get existing product
	product, err := s.GetProductByID(productID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if updatedProduct.ProductName != "" {
		product.ProductName = strings.TrimSpace(updatedProduct.ProductName)
	}

	if updatedProduct.ProductCode != "" {
		product.ProductCode = strings.TrimSpace(updatedProduct.ProductCode)
	}

	if updatedProduct.Description != "" {
		product.Description = strings.TrimSpace(updatedProduct.Description)
	}

	if updatedProduct.ImageURL != "" {
		product.ImageURL = updatedProduct.ImageURL
	}

	if updatedProduct.Price > 0 {
		product.Price = updatedProduct.Price
	}

	if updatedProduct.Stock >= 0 {
		product.Stock = updatedProduct.Stock
	}

	if updatedProduct.CategoryID != "" {
		// Verify new category exists
		_, err := s.categoryService.GetCategoryByID(updatedProduct.CategoryID)
		if err != nil {
			return nil, err
		}
		product.CategoryID = updatedProduct.CategoryID
	}

	// Update KeyFeatures if provided
	if updatedProduct.KeyFeatures != nil {
		product.KeyFeatures = updatedProduct.KeyFeatures
	}

	// Update WhyChooseBenefits if provided
	if updatedProduct.WhyChooseBenefits != nil {
		product.WhyChooseBenefits = updatedProduct.WhyChooseBenefits
	}

	// Save changes
	if err := database.GetDB().Save(product).Error; err != nil {
		return nil, err
	}

	// Reload with category and discount
	if err := database.GetDB().Preload("Category").Preload("Discount").First(product, "id = ?", product.ID).Error; err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductByID retrieves a product by ID
func (s *ProductService) GetProductByID(productID string) (*models.Product, error) {
	if productID == "" {
		return nil, errors.New("product ID is required")
	}

	var product models.Product
	if err := database.GetDB().Preload("Category").Preload("Discount").First(&product, "id = ?", productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	return &product, nil
}

// GetAllProducts retrieves all products
func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	if err := database.GetDB().Preload("Category").Preload("Discount").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// DeleteProduct deletes a product by ID
func (s *ProductService) DeleteProduct(productID string) error {
	if productID == "" {
		return errors.New("product ID is required")
	}

	// Check if product exists
	_, err := s.GetProductByID(productID)
	if err != nil {
		return err
	}

	// Delete product
	if err := database.GetDB().Delete(&models.Product{}, "id = ?", productID).Error; err != nil {
		return err
	}

	return nil
}

// GetProductsByCategory retrieves all products in a specific category
func (s *ProductService) GetProductsByCategory(categoryID string) ([]models.Product, error) {
	if categoryID == "" {
		return nil, errors.New("category ID is required")
	}

	// Verify category exists
	_, err := s.categoryService.GetCategoryByID(categoryID)
	if err != nil {
		return nil, err
	}

	var products []models.Product
	if err := database.GetDB().Preload("Category").Preload("Discount").Where("category_id = ?", categoryID).Find(&products).Error; err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, ErrNoProductsFoundByCategory
	}

	return products, nil
}

// SearchProducts searches for products by name or description
func (s *ProductService) SearchProducts(searchValue string) ([]models.Product, error) {
	searchValue = strings.TrimSpace(searchValue)
	if searchValue == "" {
		return nil, errors.New("search value is required")
	}

	var products []models.Product
	searchPattern := "%" + searchValue + "%"
	
	if err := database.GetDB().Preload("Category").Preload("Discount").
		Where("product_name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern).
		Find(&products).Error; err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, ErrNoProductsFoundBySearch
	}

	return products, nil
}

// UploadProductImage updates the product's image URL
func (s *ProductService) UploadProductImage(productID string, imageURL string) (*models.Product, error) {
	if productID == "" {
		return nil, errors.New("product ID is required")
	}

	if imageURL == "" {
		return nil, errors.New("image URL is required")
	}

	// Get existing product
	product, err := s.GetProductByID(productID)
	if err != nil {
		return nil, err
	}

	// Update image URL
	product.ImageURL = imageURL

	// Save changes
	if err := database.GetDB().Save(product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

// Ensure ProductService implements ProductServiceInterface
var _ ProductServiceInterface = (*ProductService)(nil)
