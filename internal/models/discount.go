package models

import (
	"time"

	"gorm.io/gorm"
)

// Discount represents a product discount in the system
type Discount struct {
	ID         string         `gorm:"type:uuid;primaryKey;" json:"id"`
	ProductID  string         `gorm:"type:uuid;not null;uniqueIndex" json:"productId"`
	Percentage int            `gorm:"not null" json:"percentage"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy  string         `json:"-"`
	UpdatedBy  string         `json:"-"`

	// Relationships
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// TableName specifies the table name for the Discount model
func (Discount) TableName() string {
	return "discounts"
}

// DiscountResponse represents the discount data returned in API responses
type DiscountResponse struct {
	ID          string           `json:"id"`
	ProductID   string           `json:"productId,omitempty"`
	Percentage  int              `json:"percentage"`
	ProductInfo *ProductResponse `json:"productInfo,omitempty"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
}

// ToResponse converts a Discount model to DiscountResponse
func (d *Discount) ToResponse() DiscountResponse {
	response := DiscountResponse{
		ID:         d.ID,
		ProductID:  d.ProductID,
		Percentage: d.Percentage,
		CreatedAt:  d.CreatedAt,
		UpdatedAt:  d.UpdatedAt,
	}

	if d.Product != nil {
		productResp := d.Product.ToResponse()
		response.ProductInfo = &productResp
	}

	return response
}
