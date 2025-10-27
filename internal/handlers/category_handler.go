package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
)

// CategoryHandler handles category-related requests
type CategoryHandler struct {
	categoryService services.CategoryServiceInterface
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryService services.CategoryServiceInterface) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategoryRequest represents the request body for creating a category
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateCategoryRequest represents the request body for updating a category
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateCategory creates a new category (Admin only)
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	createdCategory, err := h.categoryService.CreateCategory(category)
	if err != nil {
		if err == services.ErrCategoryAlreadyExists {
			c.JSON(http.StatusConflict, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: createdCategory.ToResponse(),
	})
}

// UpdateCategory updates an existing category (Admin only)
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	categoryID := c.Param("id")

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	updatedCategory, err := h.categoryService.UpdateCategory(categoryID, category)
	if err != nil {
		if err == services.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		if err == services.ErrCategoryAlreadyExists {
			c.JSON(http.StatusConflict, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: updatedCategory.ToResponse(),
	})
}

// GetAllCategories retrieves all categories
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	// Convert to response format
	responses := make([]models.CategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = category.ToResponse()
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: responses,
	})
}

// GetCategoryByID retrieves a category by ID
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	categoryID := c.Param("id")

	category, err := h.categoryService.GetCategoryByID(categoryID)
	if err != nil {
		if err == services.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: category.ToResponse(),
	})
}

// DeleteCategory deletes a category (Admin only)
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	categoryID := c.Param("id")

	err := h.categoryService.DeleteCategory(categoryID)
	if err != nil {
		if err == services.ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: "Deleted",
	})
}
