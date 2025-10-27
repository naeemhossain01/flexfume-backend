package services

import (
	"errors"

	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"gorm.io/gorm"
)

var (
	// ErrAddressNotFound is returned when an address is not found
	ErrAddressNotFound = errors.New("address not found")
	// ErrAddressAlreadyExists is returned when an address already exists for a user
	ErrAddressAlreadyExists = errors.New("address already exists for this user")
)

// AddressService handles address-related business logic
type AddressService struct {
	userService UserServiceInterface
}

// NewAddressService creates a new address service
func NewAddressService(userService UserServiceInterface) *AddressService {
	return &AddressService{
		userService: userService,
	}
}

// AddAddress adds a new address for a user
func (s *AddressService) AddAddress(address *models.Address) (*models.Address, error) {
	db := database.DB

	// Verify user exists
	user, err := s.userService.GetUserByID(address.UserID)
	if err != nil {
		return nil, err
	}

	// Check if address already exists for this user
	var existingAddress models.Address
	err = db.Where("user_id = ?", user.ID).First(&existingAddress).Error
	if err == nil {
		return nil, ErrAddressAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create the address
	if err := db.Create(address).Error; err != nil {
		return nil, err
	}

	// Load the user relationship
	if err := db.Preload("User").First(address, "id = ?", address.ID).Error; err != nil {
		return nil, err
	}

	return address, nil
}

// UpdateAddress updates an existing address
func (s *AddressService) UpdateAddress(addressID string, address *models.Address, userID string) (*models.Address, error) {
	db := database.DB

	// Find the existing address by ID and UserID
	var existingAddress models.Address
	if err := db.Where("id = ? AND user_id = ?", addressID, userID).First(&existingAddress).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAddressNotFound
		}
		return nil, err
	}

	// Update only non-empty fields
	updates := make(map[string]interface{})
	if address.FullName != "" {
		updates["full_name"] = address.FullName
	}
	if address.PhoneNumber != "" {
		updates["phone_number"] = address.PhoneNumber
	}
	if address.Address != "" {
		updates["address"] = address.Address
	}
	// Always update IsDefault as it's a boolean
	updates["is_default"] = address.IsDefault

	// Perform the update
	if len(updates) > 0 {
		if err := db.Model(&existingAddress).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	// Reload the address with user relationship
	if err := db.Preload("User").First(&existingAddress, "id = ?", existingAddress.ID).Error; err != nil {
		return nil, err
	}

	return &existingAddress, nil
}

// GetAddressByUserID retrieves an address by user ID
func (s *AddressService) GetAddressByUserID(userID string) (*models.Address, error) {
	db := database.DB

	var address models.Address
	if err := db.Preload("User").Where("user_id = ?", userID).First(&address).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAddressNotFound
		}
		return nil, err
	}

	return &address, nil
}
