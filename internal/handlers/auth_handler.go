package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/auth"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	jwtManager  *auth.JWTManager
	otpService  services.OTPServiceInterface
	userService services.UserServiceInterface
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(jwtManager *auth.JWTManager, otpService services.OTPServiceInterface, userService services.UserServiceInterface) *AuthHandler {
	return &AuthHandler{
		jwtManager:  jwtManager,
		otpService:  otpService,
		userService: userService,
	}
}

// VerifyOTPRequest represents the request body for verifying OTP
type VerifyOTPRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	OTP         string `json:"otp" binding:"required"`
}

// RegisterRequest represents the request body for completing registration
type RegisterRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email"`
	Password    string `json:"password" binding:"required,min=8"`
	Role        string `json:"role"` // Optional, defaults to USER in service
}

// LoginRequest represents the request body for login
type LoginRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

// Login handles user login (matches Spring Boot UserServiceImpl.loginUser)
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	// Get user by phone number
	user, err := h.userService.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	// Verify password (matches Spring Boot's passwordEncoder.matches)
	if !auth.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "Username or password is incorrect",
		})
		return
	}

	// Generate JWT token (matches Spring Boot's jwtUtils.generateToken)
	token, err := h.jwtManager.GenerateToken(user.ID, user.PhoneNumber, user.Role)
	if err != nil {
		// JWT generation failure is an internal server error
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: "Can't generate token",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: token,
	})
}
