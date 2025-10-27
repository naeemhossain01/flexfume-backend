package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
)

// CouponUsageHandler handles coupon usage-related requests
type CouponUsageHandler struct {
	couponUsageService services.CouponUsageServiceInterface
}

// NewCouponUsageHandler creates a new coupon usage handler
func NewCouponUsageHandler(couponUsageService services.CouponUsageServiceInterface) *CouponUsageHandler {
	return &CouponUsageHandler{
		couponUsageService: couponUsageService,
	}
}

// ApplyCouponRequest represents the request body for applying a coupon
type ApplyCouponRequest struct {
	CouponCode   string   `json:"couponCode" binding:"required"`
	CartInfoList []string `json:"cartInfoList" binding:"required,min=1"`
}

// ApplyCoupon applies a coupon to cart items (Authenticated users)
func (h *CouponUsageHandler) ApplyCoupon(c *gin.Context) {
	var req ApplyCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Error:   true,
			Message: "user not authenticated",
		})
		return
	}

	finalAmount, err := h.couponUsageService.ApplyCoupon(req.CartInfoList, req.CouponCode, userID.(string))
	if err != nil {
		if err == services.ErrCouponNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		if err == services.ErrCouponUsageLimitExceeded || 
		   err == services.ErrCouponInactive || 
		   err == services.ErrCouponExpired || 
		   err == services.ErrMinOrderAmountNotMet ||
		   err == services.ErrInvalidCartItems {
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
		Response: finalAmount,
	})
}

// RemoveCoupon removes an applied coupon
func (h *CouponUsageHandler) RemoveCoupon(c *gin.Context) {
	couponCode := c.Query("code")
	if couponCode == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "coupon code is required",
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Error:   true,
			Message: "user not authenticated",
		})
		return
	}

	err := h.couponUsageService.RemoveCouponUsage(couponCode, userID.(string))
	if err != nil {
		if err == services.ErrCouponUsageNotFound || err == services.ErrCouponNotFound {
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
		Response: "Coupon removed",
	})
}
