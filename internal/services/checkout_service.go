package services

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/auth"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"gorm.io/gorm"
)

// CheckoutService handles checkout-related business logic
type CheckoutService struct {
	otpService     OTPServiceInterface
	userService    UserServiceInterface
	jwtManager     *auth.JWTManager
}

// NewCheckoutService creates a new checkout service
func NewCheckoutService(otpService OTPServiceInterface, userService UserServiceInterface, jwtManager *auth.JWTManager) *CheckoutService {
	return &CheckoutService{
		otpService:  otpService,
		userService: userService,
		jwtManager:  jwtManager,
	}
}

// VerifyOTPAndHandleUser verifies OTP and creates/updates user for checkout
func (s *CheckoutService) VerifyOTPAndHandleUser(request *models.CheckoutOTPVerifyRequest) (*models.CheckoutOTPResponse, error) {
	// First verify the OTP
	if err := s.otpService.VerifyOTP(request.PhoneNumber, request.OTP); err != nil {
		return nil, err
	}

	// Check if user exists with this phone number
	existingUser, err := s.userService.GetUserByPhoneNumber(request.PhoneNumber)
	
	response := &models.CheckoutOTPResponse{}
	
	if err == nil && existingUser != nil {
		// User exists - update their information
		updateData := map[string]interface{}{
			"name": request.Name,
		}
		if request.Email != "" {
			updateData["email"] = request.Email
		}
		
		updatedUser, err := s.userService.UpdateUser(existingUser.ID, updateData)
		if err != nil {
			return nil, err
		}

		// Update or create address
		if err := s.updateOrCreateAddress(updatedUser.ID, request); err != nil {
			return nil, err
		}

		// Populate response
		response.NewUser = false
		response.UserID = updatedUser.ID
		response.UserName = updatedUser.Name
		response.UserEmail = updatedUser.Email
		response.UserPhoneNumber = updatedUser.PhoneNumber

		// Generate JWT token
		token, err := s.jwtManager.GenerateToken(updatedUser.ID, updatedUser.PhoneNumber, updatedUser.Role)
		if err != nil {
			return nil, err
		}
		response.Token = token

	} else if err == ErrUserNotFound {
		// Create new user
		newUser := &models.User{
			Name:        request.Name,
			Email:       request.Email,
			PhoneNumber: request.PhoneNumber,
			Role:        "USER",
			CreatedBy:   "SYSTEM",
		}

		// Generate a random password for the user
		tempPassword, err := generateRandomPassword(12)
		if err != nil {
			return nil, err
		}
		
		hashedPassword, err := auth.HashPassword(tempPassword)
		if err != nil {
			return nil, err
		}
		newUser.Password = hashedPassword

		// Create user in database
		if err := database.GetDB().Create(newUser).Error; err != nil {
			return nil, err
		}

		// Create address
		if err := s.updateOrCreateAddress(newUser.ID, request); err != nil {
			return nil, err
		}

		// Populate response
		response.NewUser = true
		response.UserID = newUser.ID
		response.UserName = newUser.Name
		response.UserEmail = newUser.Email
		response.UserPhoneNumber = newUser.PhoneNumber

		// Generate JWT token
		token, err := s.jwtManager.GenerateToken(newUser.ID, newUser.PhoneNumber, newUser.Role)
		if err != nil {
			return nil, err
		}
		response.Token = token
	} else {
		// Some other error occurred
		return nil, err
	}

	return response, nil
}

// updateOrCreateAddress updates or creates an address for the user
func (s *CheckoutService) updateOrCreateAddress(userID string, request *models.CheckoutOTPVerifyRequest) error {
	var address models.Address
	
	// Try to find existing address
	err := database.GetDB().Where("user_id = ?", userID).First(&address).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create new address
		address = models.Address{
			UserID:      userID,
			FullName:    request.Name,
			PhoneNumber: request.PhoneNumber,
			Address:     request.Address,
			IsDefault:   true, // First address is default
		}
		return database.GetDB().Create(&address).Error
	} else if err != nil {
		return err
	}

	// Update existing address
	address.FullName = request.Name
	address.PhoneNumber = request.PhoneNumber
	address.Address = request.Address
	
	return database.GetDB().Save(&address).Error
}

// generateRandomPassword generates a random password
func generateRandomPassword(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// Ensure CheckoutService implements CheckoutServiceInterface
var _ CheckoutServiceInterface = (*CheckoutService)(nil)
