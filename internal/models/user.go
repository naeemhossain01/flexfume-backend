package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID          string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Email       string         `gorm:"uniqueIndex" json:"email"`
	PhoneNumber string         `gorm:"column:phone_number;uniqueIndex;not null" json:"phoneNumber"`
	Password    string         `gorm:"not null" json:"-"` // "-" prevents password from being serialized to JSON
	Role        string         `gorm:"type:varchar(50);default:'USER'" json:"role"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
	CreatedBy   string         `gorm:"column:created_by" json:"-"`
	UpdatedBy   string         `gorm:"column:updated_by" json:"-"`
	
	// Relationships (matches Spring Boot User entity)
	Address *Address `gorm:"foreignKey:UserID" json:"address,omitempty"`
	Orders  []Order  `gorm:"foreignKey:UserID" json:"orders,omitempty"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// UserResponse represents the user data returned in API responses
// Matches Spring Boot's UserInfo structure
type UserResponse struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Email       string        `json:"email"`
	PhoneNumber string        `json:"phoneNumber"`
	Role        string        `json:"role,omitempty"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt,omitempty"`
	AddressInfo *AddressInfo  `json:"addressInfo,omitempty"`
	OrderList   []OrderInfo   `json:"orderInfoList,omitempty"`
}

// ToResponse converts a User model to UserResponse
func (u *User) ToResponse() UserResponse {
	resp := UserResponse{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Role:        u.Role,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
	
	// Include address if loaded
	if u.Address != nil {
		addressInfo := u.Address.ToAddressInfo()
		resp.AddressInfo = &addressInfo
	}
	
	// Include orders if loaded
	if len(u.Orders) > 0 {
		resp.OrderList = make([]OrderInfo, len(u.Orders))
		for i, order := range u.Orders {
			resp.OrderList[i] = order.ToOrderInfo()
		}
	}
	
	return resp
}
