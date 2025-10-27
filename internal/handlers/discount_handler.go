package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
)

// DiscountHandler handles discount-related requests
type DiscountHandler struct {
	discountService services.DiscountServiceInterface
}

// NewDiscountHandler creates a new discount handler
func NewDiscountHandler(discountService services.DiscountServiceInterface) *DiscountHandler {
	return &DiscountHandler{
		discountService: discountService,
	}
}

// AddDiscountRequest represents a single discount in the request
type AddDiscountRequest struct {
	ProductID          string `json:"productId" binding:"required"`
	DiscountPercentage int    `json:"discountPercentage" binding:"required,min=0,max=100"`
}

// AddDiscounts adds discounts to products (Admin only)
func (h *DiscountHandler) AddDiscounts(c *gin.Context) {
	var req []AddDiscountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	if len(req) == 0 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "at least one discount is required",
		})
		return
	}

	// Convert to discount models
	discounts := make([]models.Discount, len(req))
	for i, item := range req {
		discounts[i] = models.Discount{
			ProductID:  item.ProductID,
			Percentage: item.DiscountPercentage,
		}
	}

	createdDiscounts, err := h.discountService.AddDiscounts(discounts)
	if err != nil {
		if err == services.ErrProductNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		if err == services.ErrDiscountAlreadyExists || 
		   err == services.ErrDiscountPercentageInvalid ||
		   err == services.ErrDiscountProductRequired {
			c.JSON(http.StatusBadRequest, APIResponse{
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

	// Convert to response format
	responses := make([]models.DiscountResponse, len(createdDiscounts))
	for i, discount := range createdDiscounts {
		responses[i] = discount.ToResponse()
	}

	c.JSON(http.StatusCreated, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: responses,
	})
}

// UpdateDiscounts updates product discounts (Admin only)
func (h *DiscountHandler) UpdateDiscounts(c *gin.Context) {
	var req map[string]int
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	if len(req) == 0 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "at least one discount is required",
		})
		return
	}

	updatedDiscounts, err := h.discountService.UpdateDiscounts(req)
	if err != nil {
		if err == services.ErrDiscountNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		if err == services.ErrDiscountPercentageInvalid {
			c.JSON(http.StatusBadRequest, APIResponse{
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

	// Convert to response format
	responses := make([]models.DiscountResponse, len(updatedDiscounts))
	for i, discount := range updatedDiscounts {
		responses[i] = discount.ToResponse()
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: responses,
	})
}

// GetAllDiscounts retrieves all product discounts
func (h *DiscountHandler) GetAllDiscounts(c *gin.Context) {
	discounts, err := h.discountService.GetAllDiscounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	// Convert to response format
	responses := make([]models.DiscountResponse, len(discounts))
	for i, discount := range discounts {
		responses[i] = discount.ToResponse()
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: responses,
	})
}

// DeleteDiscount deletes a product discount (Admin only)
func (h *DiscountHandler) DeleteDiscount(c *gin.Context) {
	discountID := c.Param("id")

	err := h.discountService.DeleteDiscount(discountID)
	if err != nil {
		if err == services.ErrDiscountNotFound {
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
		Response: "DELETED",
	})
}
