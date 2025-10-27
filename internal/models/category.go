package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a product category in the system
type Category struct {
	ID          string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy   string         `json:"-"`
	UpdatedBy   string         `json:"-"`
	
	// Relationships
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

// TableName specifies the table name for the Category model
func (Category) TableName() string {
	return "categories"
}

// CategoryResponse represents the category data returned in API responses
type CategoryResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ToResponse converts a Category model to CategoryResponse
func (c *Category) ToResponse() CategoryResponse {
	return CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
