package models

import (
	"time"

	"gorm.io/gorm"
)

// Address represents a user's address in the system
type Address struct {
	ID          string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      string         `gorm:"column:user_id;type:uuid;not null" json:"userId"`
	FullName    string         `gorm:"column:full_name;not null" json:"fullName"`
	PhoneNumber string         `gorm:"column:phone_number;not null" json:"phoneNumber"`
	Address     string         `gorm:"column:address;not null" json:"address"`
	IsDefault   bool           `gorm:"column:is_default;default:false" json:"isDefault"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for the Address model
func (Address) TableName() string {
	return "addresses"
}

// AddressRequest represents the request body for creating/updating an address
type AddressRequest struct {
	FullName    string `json:"fullName" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Address     string `json:"address" binding:"required"`
	IsDefault   bool   `json:"isDefault"`
}

// AddressInfo represents address information in responses
type AddressInfo struct {
	ID          string        `json:"id"`
	FullName    string        `json:"fullName"`
	PhoneNumber string        `json:"phoneNumber"`
	Address     string        `json:"address"`
	IsDefault   bool          `json:"isDefault"`
	UserInfo    *UserResponse `json:"userInfo,omitempty"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

// ToAddressInfo converts Address to AddressInfo
func (a *Address) ToAddressInfo() AddressInfo {
	info := AddressInfo{
		ID:          a.ID,
		FullName:    a.FullName,
		PhoneNumber: a.PhoneNumber,
		Address:     a.Address,
		IsDefault:   a.IsDefault,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}

	if a.User != nil {
		userResp := a.User.ToResponse()
		info.UserInfo = &userResp
	}

	return info
}
