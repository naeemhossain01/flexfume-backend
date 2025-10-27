package models

import (
	"time"

	"gorm.io/gorm"
)

// Cart represents a shopping cart item in the system
type Cart struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    string         `gorm:"type:uuid;not null" json:"userId"`
	ProductID string         `gorm:"type:uuid;not null" json:"productId"`
	Quantity  int            `gorm:"not null;default:1" json:"quantity"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy string         `gorm:"-" json:"-"` // Excluded from DB operations
	UpdatedBy string         `gorm:"-" json:"-"` // Excluded from DB operations

	// Relationships
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// TableName specifies the table name for the Cart model
func (Cart) TableName() string {
	return "carts"
}

// CartResponse represents the cart data returned in API responses
type CartResponse struct {
	ID          string           `json:"id"`
	UserID      string           `json:"userId"`
	ProductID   string           `json:"productId"`
	Quantity    int              `json:"quantity"`
	ProductInfo *ProductResponse `json:"productInfo,omitempty"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
}

// ToResponse converts a Cart model to CartResponse
func (c *Cart) ToResponse() CartResponse {
	response := CartResponse{
		ID:        c.ID,
		UserID:    c.UserID,
		ProductID: c.ProductID,
		Quantity:  c.Quantity,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	if c.Product != nil {
		productResp := c.Product.ToResponse()
		response.ProductInfo = &productResp
	}

	return response
}

// CartItemData represents cart item with product and discount information
type CartItemData struct {
	CartID           string
	Quantity         int
	ProductPrice     float64
	DiscountPercent  int
}
