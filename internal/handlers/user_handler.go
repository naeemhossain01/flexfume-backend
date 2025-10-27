package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userService *services.UserService
	otpService  *services.OTPService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *services.UserService, otpService *services.OTPService) *UserHandler {
	return &UserHandler{
		userService: userService,
		otpService:  otpService,
	}
}

// UpdateUserRequest represents the request body for updating user information
type UpdateUserRequest struct {
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

// ChangePasswordRequest represents the request body for changing password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

// ResetPasswordRequest represents the request body for resetting password
type ResetPasswordRequest struct {
	PhoneNumber     string `json:"phoneNumber" binding:"required"`
	OTP             string `json:"otp" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

// GetUserByID returns user details by user ID
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.userService.GetUserByID(userID)
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
			Message: "Failed to retrieve user",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: user.ToResponse(),
	})
}

// GetProfile returns the authenticated user's profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Error:   true,
			Message: "User not authenticated",
		})
		return
	}

	user, err := h.userService.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Error:   true,
			Message: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: user.ToResponse(),
	})
}

// UpdateUser updates user information
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	// Get authenticated user ID from context
	authenticatedUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Error:   true,
			Message: "User not authenticated",
		})
		return
	}

	// Get user role from context
	userRole, _ := c.Get("user_role")

	// Check authorization: user can only update their own profile unless they're an admin
	if authenticatedUserID.(string) != userID && userRole != "ADMIN" {
		c.JSON(http.StatusForbidden, APIResponse{
			Error:   true,
			Message: "You are not authorized to update this user's information",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Convert request to map for selective updates
	updateData := make(map[string]interface{})
	if req.Name != "" {
		updateData["name"] = req.Name
	}
	if req.Email != "" {
		updateData["email"] = req.Email
	}
	if req.PhoneNumber != "" {
		updateData["phoneNumber"] = req.PhoneNumber
	}

	user, err := h.userService.UpdateUser(userID, updateData)
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
			Message: "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: user.ToResponse(),
	})
}

// GetAllUsers returns all users (Admin only)
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: "Failed to retrieve users",
		})
		return
	}

	// Convert to response format
	userResponses := make([]interface{}, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: userResponses,
	})
}

// ChangePassword changes the authenticated user's password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Error:   true,
			Message: "User not authenticated",
		})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	err := h.userService.ChangePassword(userID.(string), req.CurrentPassword, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		switch err {
		case services.ErrPasswordMismatch:
			c.JSON(http.StatusBadRequest, APIResponse{
				Error:   true,
				Message: "New password and confirm password do not match",
			})
		case services.ErrCurrentPasswordInvalid:
			c.JSON(http.StatusBadRequest, APIResponse{
				Error:   true,
				Message: "Current password is incorrect",
			})
		case services.ErrSamePassword:
			c.JSON(http.StatusBadRequest, APIResponse{
				Error:   true,
				Message: "New password must be different from current password",
			})
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: "User not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, APIResponse{
				Error:   true,
				Message: "Failed to change password",
			})
		}
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: "Password changed successfully",
	})
}

// ResetPasswordRequest handles password reset OTP request
func (h *UserHandler) ResetPasswordRequest(c *gin.Context) {
	phoneNumber := c.Query("phoneNumber")
	otpType := c.Query("type")

	if phoneNumber == "" || otpType == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "phoneNumber and type are required",
		})
		return
	}

	// Fix URL encoding issue: '+' in URL query params is decoded as space by Gin
	// Trim spaces and add '+' prefix if missing
	phoneNumber = strings.TrimSpace(phoneNumber)
	if !strings.HasPrefix(phoneNumber, "+") && len(phoneNumber) > 0 {
		phoneNumber = "+" + phoneNumber
	}

	// Verify user exists
	_, err := h.userService.GetUserByPhoneNumber(phoneNumber)
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
			Message: "Failed to process request",
		})
		return
	}

	// Send OTP
	if err := h.otpService.SendOTP(phoneNumber, services.OTPTypePasswordReset); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: "Failed to send OTP",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: "Password OTP sent",
	})
}

// ResetPassword resets user password using OTP
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Verify OTP
	if err := h.otpService.VerifyResetPasswordOTP(req.PhoneNumber, req.OTP); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "Invalid or expired OTP",
		})
		return
	}

	// Reset password
	err := h.userService.ResetPassword(req.PhoneNumber, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		switch err {
		case services.ErrPasswordMismatch:
			c.JSON(http.StatusBadRequest, APIResponse{
				Error:   true,
				Message: "New password and confirm password do not match",
			})
		case services.ErrUserNotFound:
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: "User not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, APIResponse{
				Error:   true,
				Message: "Failed to reset password",
			})
		}
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: "Password reset successful",
	})
}
