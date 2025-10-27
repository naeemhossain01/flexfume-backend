package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
)

// CouponHandler handles coupon-related requests
type CouponHandler struct {
	couponService services.CouponServiceInterface
}

// NewCouponHandler creates a new coupon handler
func NewCouponHandler(couponService services.CouponServiceInterface) *CouponHandler {
	return &CouponHandler{
		couponService: couponService,
	}
}

// CreateCouponRequest represents the request body for creating a coupon
type CreateCouponRequest struct {
	Code             string    `json:"code" binding:"required"`
	DiscountType     string    `json:"discountType" binding:"required,oneof=PERCENTAGE FIXED"`
	DiscountValue    float64   `json:"discountValue" binding:"required,gt=0"`
	MinOrderValue    float64   `json:"minOrderValue"`
	MaxDiscount      float64   `json:"maxDiscount"`
	ValidFrom        time.Time `json:"validFrom"`
	ValidUntil       time.Time `json:"validUntil" binding:"required"`
	UsageLimit       int       `json:"usageLimit" binding:"required,min=1"`
}

// UpdateCouponRequest represents the request body for updating a coupon
type UpdateCouponRequest struct {
	Code             string    `json:"code"`
	DiscountType     string    `json:"discountType" binding:"omitempty,oneof=PERCENTAGE FIXED"`
	DiscountValue    float64   `json:"discountValue"`
	MinOrderValue    float64   `json:"minOrderValue"`
	MaxDiscount      float64   `json:"maxDiscount"`
	ValidFrom        time.Time `json:"validFrom"`
	ValidUntil       time.Time `json:"validUntil"`
	UsageLimit       int       `json:"usageLimit"`
	Active           *bool     `json:"active"`
}

// CreateCoupon creates a new coupon (Admin only)
func (h *CouponHandler) CreateCoupon(c *gin.Context) {
	var req CreateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	coupon := &models.Coupon{
		Code:             req.Code,
		CouponType:       req.DiscountType,
		Amount:           req.DiscountValue,
		MinOrderAmount:   req.MinOrderValue,
		MaxAmountApplied: req.MaxDiscount,
		ExpirationTime:   req.ValidUntil,
		UsageLimit:       req.UsageLimit,
		Active:           true,
	}

	createdCoupon, err := h.couponService.CreateCoupon(coupon)
	if err != nil {
		if err == services.ErrCouponCodeExists {
			c.JSON(http.StatusConflict, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		if err == services.ErrInvalidCouponType || err == services.ErrCouponCodeRequired || err == services.ErrCouponAmountRequired {
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

	c.JSON(http.StatusCreated, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: createdCoupon.ToResponse(),
	})
}

// UpdateCoupon updates an existing coupon (Admin only)
func (h *CouponHandler) UpdateCoupon(c *gin.Context) {
	couponID := c.Param("id")

	var req UpdateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	coupon := &models.Coupon{
		Code:             req.Code,
		CouponType:       req.DiscountType,
		Amount:           req.DiscountValue,
		MinOrderAmount:   req.MinOrderValue,
		MaxAmountApplied: req.MaxDiscount,
		ExpirationTime:   req.ValidUntil,
		UsageLimit:       req.UsageLimit,
	}

	if req.Active != nil {
		coupon.Active = *req.Active
	}

	updatedCoupon, err := h.couponService.UpdateCoupon(couponID, coupon)
	if err != nil {
		if err == services.ErrCouponNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		if err == services.ErrCouponCodeExists || err == services.ErrInvalidCouponType {
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

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: updatedCoupon.ToResponse(),
	})
}

// DeleteCoupon deletes a coupon (Admin only)
func (h *CouponHandler) DeleteCoupon(c *gin.Context) {
	couponID := c.Param("id")

	err := h.couponService.DeleteCoupon(couponID)
	if err != nil {
		if err == services.ErrCouponNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		if err == services.ErrCouponHasUsage {
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

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: "Coupon deleted successfully",
	})
}

// GetCouponByID retrieves a coupon by ID
func (h *CouponHandler) GetCouponByID(c *gin.Context) {
	couponID := c.Param("id")

	coupon, err := h.couponService.GetCouponByID(couponID)
	if err != nil {
		if err == services.ErrCouponNotFound {
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

	// Enrich with usage statistics
	response := h.couponService.EnrichCouponWithStatistics(coupon)

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: response,
	})
}

// GetAllCoupons retrieves all coupons
func (h *CouponHandler) GetAllCoupons(c *gin.Context) {
	coupons, err := h.couponService.GetAllCoupons()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	// Enrich each coupon with usage statistics
	responses := make([]models.CouponResponse, len(coupons))
	for i, coupon := range coupons {
		responses[i] = h.couponService.EnrichCouponWithStatistics(&coupon)
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: responses,
	})
}
