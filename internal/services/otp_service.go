package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"
)

const (
	// OTP types
	OTPTypeRegistration = "USER_REGISTRATION"
	OTPTypePasswordReset = "PASSWORD_RESET"
	OTPTypeCheckout = "CHECKOUT"

	// Redis key prefixes (must match Spring Boot for compatibility)
	RedisOTPPrefix = "OTP_"
	RedisOTPResetPasswordPrefix = "RESET_"
	RedisPhoneVerifiedPrefix = "VERIFIED_"

	// OTP expiration times
	OTPExpirationRegistration = 5 * time.Minute  // 300 seconds
	OTPExpirationPasswordReset = 10 * time.Minute // 600 seconds
	PhoneVerificationExpiration = 30 * time.Minute // 1800 seconds
)

var (
	ErrInvalidOTP = errors.New("OTP is not valid")
	ErrOTPExpired = errors.New("OTP has expired")
	ErrPhoneNotVerified = errors.New("phone number not verified")
	ErrPhoneAlreadyVerified = errors.New("Phone number already verified")
	ErrOTPAlreadySent = errors.New("Otp already send. Please try again after 5 minutes.")
)

// OTPServiceInterface defines the interface for OTP operations
type OTPServiceInterface interface {
	GenerateOTP() (string, error)
	SendOTP(phoneNumber, otpType string) error
	VerifyOTP(phoneNumber, otp string) error
	MarkPhoneAsVerified(phoneNumber, otp string) error
	IsPhoneVerified(phoneNumber string) (bool, error)
	CheckPhoneNotVerified(phoneNumber string) error
	RequirePhoneVerified(phoneNumber string) error
	VerifyResetPasswordOTP(phoneNumber, otp string) error
}

// OTPService handles OTP operations
type OTPService struct {
	redisService RedisServiceInterface
	smsService   SMSServiceInterface
}

// OTPVerification represents OTP verification data stored in Redis
// This matches the Spring Boot OtpRequest structure for compatibility
type OTPVerification struct {
	PhoneNumber string `json:"phoneNumber"`
	OTP         string `json:"otp"`
	Verified    bool   `json:"verified"`
}

// NewOTPService creates a new OTP service
func NewOTPService(redisService RedisServiceInterface, smsService SMSServiceInterface) *OTPService {
	return &OTPService{
		redisService: redisService,
		smsService:   smsService,
	}
}

// GenerateOTP generates a random 4-digit OTP
func (s *OTPService) GenerateOTP() (string, error) {
	max := big.NewInt(10000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%04d", n.Int64()), nil
}

// SendOTP sends an OTP to the phone number
func (s *OTPService) SendOTP(phoneNumber, otpType string) error {
	// Use phone number as-is to match Spring Boot behavior
	// Check if OTP was already sent and still valid (rate limiting)
	if err := s.checkOTPNotRecentlySent(phoneNumber, otpType); err != nil {
		return err
	}

	// Generate OTP
	otp, err := s.GenerateOTP()
	if err != nil {
		return err
	}

	// Generate SMS message based on type
	message := s.generateSMSMessage(otp, otpType)

	// Send SMS
	if err := s.smsService.SendSMS(phoneNumber, message); err != nil {
		return err
	}

	// Store OTP in Redis
	return s.storeOTP(phoneNumber, otp, otpType)
}

// VerifyOTP verifies the OTP for a phone number
func (s *OTPService) VerifyOTP(phoneNumber, otp string) error {
	key := RedisOTPPrefix + phoneNumber

	// Get stored OTP from Redis
	storedOTP, err := s.redisService.GetString(key)
	if err != nil {
		return ErrInvalidOTP
	}

	// Compare OTPs
	if storedOTP != otp {
		return ErrInvalidOTP
	}

	return nil
}

// MarkPhoneAsVerified marks a phone number as verified
// This matches Spring Boot's markPhoneNumberAsValid behavior which stores the OtpRequest object
func (s *OTPService) MarkPhoneAsVerified(phoneNumber, otp string) error {
	key := RedisPhoneVerifiedPrefix + phoneNumber
	verification := OTPVerification{
		PhoneNumber: phoneNumber,
		OTP:         otp,
		Verified:    true,
	}

	return s.redisService.Set(key, verification, PhoneVerificationExpiration)
}

// IsPhoneVerified checks if a phone number is verified
func (s *OTPService) IsPhoneVerified(phoneNumber string) (bool, error) {
	key := RedisPhoneVerifiedPrefix + phoneNumber

	var verification OTPVerification
	err := s.redisService.Get(key, &verification)
	if err != nil {
		return false, nil // Not verified or expired
	}

	return verification.Verified, nil
}

// CheckPhoneNotVerified ensures phone is not already verified (for registration)
func (s *OTPService) CheckPhoneNotVerified(phoneNumber string) error {
	verified, err := s.IsPhoneVerified(phoneNumber)
	if err != nil {
		return err
	}

	if verified {
		return ErrPhoneAlreadyVerified
	}

	return nil
}

// RequirePhoneVerified ensures phone is verified (for registration completion)
func (s *OTPService) RequirePhoneVerified(phoneNumber string) error {
	verified, err := s.IsPhoneVerified(phoneNumber)
	if err != nil {
		return err
	}

	if !verified {
		return ErrPhoneNotVerified
	}

	return nil
}

// VerifyResetPasswordOTP verifies OTP for password reset
func (s *OTPService) VerifyResetPasswordOTP(phoneNumber, otp string) error {
	key := RedisOTPResetPasswordPrefix + phoneNumber

	storedOTP, err := s.redisService.GetString(key)
	if err != nil {
		return ErrInvalidOTP
	}

	if storedOTP != otp {
		return ErrInvalidOTP
	}

	return nil
}

// storeOTP stores OTP in Redis with appropriate expiration
func (s *OTPService) storeOTP(phoneNumber, otp, otpType string) error {
	var key string
	var expiration time.Duration

	switch otpType {
	case OTPTypePasswordReset:
		key = RedisOTPResetPasswordPrefix + phoneNumber
		expiration = OTPExpirationPasswordReset
	default:
		key = RedisOTPPrefix + phoneNumber
		expiration = OTPExpirationRegistration
	}

	return s.redisService.Set(key, otp, expiration)
}

// checkOTPNotRecentlySent checks if an OTP was recently sent (rate limiting)
// This prevents SMS spam and matches Spring Boot behavior
func (s *OTPService) checkOTPNotRecentlySent(phoneNumber, otpType string) error {
	var key string
	
	switch otpType {
	case OTPTypePasswordReset:
		key = RedisOTPResetPasswordPrefix + phoneNumber
	default:
		key = RedisOTPPrefix + phoneNumber
	}
	
	// Check if OTP exists in Redis
	existingOTP, err := s.redisService.GetString(key)
	if err == nil && existingOTP != "" {
		// OTP still exists and hasn't expired - rate limit
		return ErrOTPAlreadySent
	}
	
	// No existing OTP or it expired - OK to send new one
	return nil
}

// generateSMSMessage generates SMS message based on OTP type
func (s *OTPService) generateSMSMessage(otp, otpType string) string {
	switch otpType {
	case OTPTypeRegistration:
		return fmt.Sprintf("Your FlexFume registration OTP is: %s. Valid for 5 minutes.", otp)
	case OTPTypePasswordReset:
		return fmt.Sprintf("Your FlexFume password reset OTP is: %s. Valid for 10 minutes.", otp)
	case OTPTypeCheckout:
		return fmt.Sprintf("Your FlexFume checkout OTP is: %s. Valid for 5 minutes.", otp)
	default:
		return fmt.Sprintf("Your FlexFume OTP is: %s", otp)
	}
}

// normalizePhoneNumber normalizes phone number for Redis key consistency
// Removes '+' prefix and trims whitespace to match Spring Boot behavior
func normalizePhoneNumber(phoneNumber string) string {
	// Trim whitespace
	normalized := strings.TrimSpace(phoneNumber)
	// Remove '+' prefix if present
	normalized = strings.TrimPrefix(normalized, "+")
	// Trim again in case there were spaces after '+'
	normalized = strings.TrimSpace(normalized)
	return normalized
}
