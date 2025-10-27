package services

import (
	"errors"
	"strings"

	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/auth"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrPasswordMismatch       = errors.New("new password and confirm password do not match")
	ErrCurrentPasswordInvalid = errors.New("current password is incorrect")
	ErrSamePassword           = errors.New("new password must be different from current password")
	ErrPhoneNumberRequired    = errors.New("phone number is required")
	ErrUserAlreadyExists      = errors.New("Phone number already registered")
)

// UserService handles user-related business logic
type UserService struct {
	otpService *OTPService
}

// NewUserService creates a new user service
func NewUserService(otpService *OTPService) *UserService {
	return &UserService{
		otpService: otpService,
	}
}

// RegisterUser creates a new user account (matches Spring Boot UserServiceImpl.registerUser)
func (s *UserService) RegisterUser(user *models.User) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.GetUserByPhoneNumber(user.PhoneNumber)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}
	
	// Only return error if it's not "user not found"
	if err != nil && err != ErrUserNotFound {
		return nil, err
	}
	
	// Set role logic: if ADMIN role is explicitly set, keep it; otherwise default to USER
	// (matches Spring Boot UserServiceImpl.registerUser lines 49-55)
	if user.Role != "" && strings.EqualFold(user.Role, "ADMIN") {
		user.Role = "ADMIN"
	} else {
		user.Role = "USER"
	}
	
	// Hash password (matches Spring Boot line 58: user.setPassword(passwordEncoder.encode(user.getPassword())))
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword
	
	// Set audit fields (matches Spring Boot line 59: user.setCreatedBy(user.getPhoneNumber()))
	user.CreatedBy = user.PhoneNumber
	
	// Save user to database
	if err := database.GetDB().Create(user).Error; err != nil {
		return nil, err
	}
	
	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := database.GetDB().First(&user, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByPhoneNumber retrieves a user by phone number
func (s *UserService) GetUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	if phoneNumber == "" {
		return nil, ErrPhoneNumberRequired
	}

	var user models.User
	db := database.GetDB()
	
	// Debug: Log the query
	db = db.Debug()
	
	if err := db.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := database.GetDB().Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(userID string, updateData map[string]interface{}) (*models.User, error) {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Update only provided fields
	if name, ok := updateData["name"].(string); ok && name != "" {
		user.Name = strings.TrimSpace(name)
	}
	if email, ok := updateData["email"].(string); ok && email != "" {
		user.Email = strings.ToLower(strings.TrimSpace(email))
	}
	if phoneNumber, ok := updateData["phoneNumber"].(string); ok && phoneNumber != "" {
		user.PhoneNumber = strings.TrimSpace(phoneNumber)
	}

	if err := database.GetDB().Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword changes the user's password
func (s *UserService) ChangePassword(userID, oldPassword, newPassword, confirmPassword string) error {
	// Validate passwords match
	if newPassword != confirmPassword {
		return ErrPasswordMismatch
	}

	// Get user
	user, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Verify old password
	if !auth.CheckPassword(oldPassword, user.Password) {
		return ErrCurrentPasswordInvalid
	}

	// Check if new password is different
	if auth.CheckPassword(newPassword, user.Password) {
		return ErrSamePassword
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	user.Password = hashedPassword
	user.UpdatedBy = user.PhoneNumber

	if err := database.GetDB().Save(user).Error; err != nil {
		return err
	}

	return nil
}

// ResetPassword resets the user's password after OTP verification
func (s *UserService) ResetPassword(phoneNumber, newPassword, confirmPassword string) error {
	// Validate passwords match
	if newPassword != confirmPassword {
		return ErrPasswordMismatch
	}

	// Get user by phone number
	user, err := s.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	user.Password = hashedPassword
	user.UpdatedBy = user.PhoneNumber

	if err := database.GetDB().Save(user).Error; err != nil {
		return err
	}

	return nil
}

// ValidateAlreadyHaveAccount checks if a user already exists with the given phone number
// Returns ErrUserAlreadyExists if user exists, nil otherwise
func (s *UserService) ValidateAlreadyHaveAccount(phoneNumber string) error {
	if phoneNumber == "" {
		return ErrPhoneNumberRequired
	}

	var user models.User
	err := database.GetDB().Where("phone_number = ?", phoneNumber).First(&user).Error
	
	if err == nil {
		// User found - they already have an account
		return ErrUserAlreadyExists
	}
	
	if err == gorm.ErrRecordNotFound {
		// User not found - this is good, they can register
		return nil
	}
	
	// Some other database error occurred
	return err
}

// Ensure UserService implements UserServiceInterface
var _ UserServiceInterface = (*UserService)(nil)
