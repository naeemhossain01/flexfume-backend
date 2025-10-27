package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
)

// CheckoutHandler handles checkout-related requests
type CheckoutHandler struct {
	otpService      services.OTPServiceInterface
	checkoutService services.CheckoutServiceInterface
}

// NewCheckoutHandler creates a new checkout handler
func NewCheckoutHandler(otpService services.OTPServiceInterface, checkoutService services.CheckoutServiceInterface) *CheckoutHandler {
	return &CheckoutHandler{
		otpService:      otpService,
		checkoutService: checkoutService,
	}
}

// SendCheckoutOTP sends OTP for checkout verification
// @Summary Send checkout OTP
// @Description Send OTP for checkout verification
// @Tags checkout
// @Accept json
// @Produce json
// @Param request body models.CheckoutOTPRequest true "Phone number"
// @Success 200 {object} map[string]interface{} "success response"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/checkout/send-otp [post]
func (h *CheckoutHandler) SendCheckoutOTP(c *gin.Context) {
	var req models.CheckoutOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	// Send OTP for checkout verification
	if err := h.otpService.SendOTP(req.PhoneNumber, services.OTPTypeCheckout); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: "OTP sent successfully for checkout verification",
	})
}

// VerifyCheckoutOTP verifies OTP and completes checkout
// @Summary Verify checkout OTP
// @Description Verify OTP and complete checkout
// @Tags checkout
// @Accept json
// @Produce json
// @Param request body models.CheckoutOTPVerifyRequest true "Checkout verification details"
// @Success 200 {object} map[string]interface{} "success response with checkout data"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/checkout/verify-otp [post]
func (h *CheckoutHandler) VerifyCheckoutOTP(c *gin.Context) {
	var req models.CheckoutOTPVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	// Verify OTP and handle user account creation/update
	response, err := h.checkoutService.VerifyOTPAndHandleUser(&req)
	if err != nil {
		if err == services.ErrInvalidOTP {
			c.JSON(http.StatusBadRequest, APIResponse{
				Error:   true,
				Message: "Invalid or expired OTP",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: "Checkout verification failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: response,
	})
}
