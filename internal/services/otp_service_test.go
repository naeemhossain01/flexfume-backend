package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOTPService_GenerateOTP(t *testing.T) {
	service := NewOTPService()

	t.Run("Success - Generate OTP", func(t *testing.T) {
		otp := service.GenerateOTP()

		assert.NotEmpty(t, otp)
		assert.Len(t, otp, 6) // Assuming 6-digit OTP
	})
}

func TestOTPService_ValidateOTP(t *testing.T) {
	service := NewOTPService()

	t.Run("Success - Validate correct OTP", func(t *testing.T) {
		phoneNumber := "+1234567890"
		otp := service.GenerateOTP()

		// Store OTP
		service.StoreOTP(phoneNumber, otp)

		// Validate
		isValid := service.ValidateOTP(phoneNumber, otp)

		assert.True(t, isValid)
	})

	t.Run("Error - Validate incorrect OTP", func(t *testing.T) {
		phoneNumber := "+1234567890"
		otp := service.GenerateOTP()

		// Store OTP
		service.StoreOTP(phoneNumber, otp)

		// Validate with wrong OTP
		isValid := service.ValidateOTP(phoneNumber, "000000")

		assert.False(t, isValid)
	})

	t.Run("Error - Validate non-existent phone number", func(t *testing.T) {
		isValid := service.ValidateOTP("+9999999999", "123456")

		assert.False(t, isValid)
	})
}
