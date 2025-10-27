package services

import (
	"errors"
	"strings"

	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategoryNameRequired  = errors.New("category name is required")
	ErrCategoryAlreadyExists = errors.New("category with this name already exists")
)

// CategoryService handles category-related business logic
type CategoryService struct{}

// NewCategoryService creates a new category service
func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(category *models.Category) (*models.Category, error) {
	if category == nil {
		return nil, errors.New("category cannot be nil")
	}

	// Validate required fields
	category.Name = strings.TrimSpace(category.Name)
	if category.Name == "" {
		return nil, ErrCategoryNameRequired
	}

	// Check if category with same name already exists (excluding soft-deleted)
	var existingCategory models.Category
	err := database.GetDB().Where("name = ? AND deleted_at IS NULL", category.Name).First(&existingCategory).Error
	if err == nil {
		return nil, ErrCategoryAlreadyExists
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Create category
	if err := database.GetDB().Create(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory updates an existing category
func (s *CategoryService) UpdateCategory(categoryID string, updatedCategory *models.Category) (*models.Category, error) {
	if categoryID == "" {
		return nil, errors.New("category ID is required")
	}

	if updatedCategory == nil {
		return nil, errors.New("category data cannot be nil")
	}

	// Get existing category
	category, err := s.GetCategoryByID(categoryID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if updatedCategory.Name != "" {
		category.Name = strings.TrimSpace(updatedCategory.Name)
		
		// Check if new name conflicts with existing category (excluding soft-deleted)
		var existingCategory models.Category
		err := database.GetDB().Where("name = ? AND id != ? AND deleted_at IS NULL", category.Name, categoryID).First(&existingCategory).Error
		if err == nil {
			return nil, ErrCategoryAlreadyExists
		} else if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	if updatedCategory.Description != "" {
		category.Description = strings.TrimSpace(updatedCategory.Description)
	}

	// Save changes
	if err := database.GetDB().Save(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategoryByID retrieves a category by ID
func (s *CategoryService) GetCategoryByID(categoryID string) (*models.Category, error) {
	if categoryID == "" {
		return nil, errors.New("category ID is required")
	}

	var category models.Category
	if err := database.GetDB().First(&category, "id = ?", categoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	return &category, nil
}

// GetAllCategories retrieves all categories
func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := database.GetDB().Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// DeleteCategory deletes a category by ID
func (s *CategoryService) DeleteCategory(categoryID string) error {
	if categoryID == "" {
		return errors.New("category ID is required")
	}

	// Check if category exists
	_, err := s.GetCategoryByID(categoryID)
	if err != nil {
		return err
	}

	// Delete category
	if err := database.GetDB().Delete(&models.Category{}, "id = ?", categoryID).Error; err != nil {
		return err
	}

	return nil
}

// Ensure CategoryService implements CategoryServiceInterface
var _ CategoryServiceInterface = (*CategoryService)(nil)
