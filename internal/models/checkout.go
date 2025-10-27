package models

// CheckoutOTPRequest represents the request to send OTP for checkout
type CheckoutOTPRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
}

// CheckoutOTPVerifyRequest represents the request to verify OTP and complete checkout
type CheckoutOTPVerifyRequest struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	OTP         string `json:"otp" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email"`
	Address     string `json:"address" binding:"required"`
}

// CheckoutOTPResponse represents the response after successful checkout OTP verification
type CheckoutOTPResponse struct {
	Token            string `json:"token"`
	NewUser          bool   `json:"newUser"`
	UserID           string `json:"userId"`
	UserName         string `json:"userName"`
	UserEmail        string `json:"userEmail"`
	UserPhoneNumber  string `json:"userPhoneNumber"`
}
