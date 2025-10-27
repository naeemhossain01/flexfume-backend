package models

import (
	"time"

	"gorm.io/gorm"
)

// Coupon represents a coupon in the system
type Coupon struct {
	ID              string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Code            string         `gorm:"uniqueIndex;not null" json:"code"`
	CouponType      string         `gorm:"not null" json:"couponType"` // PERCENTAGE or FIXED
	Amount          float64        `gorm:"type:decimal(10,2);not null" json:"amount"`
	MinOrderAmount  float64        `gorm:"type:decimal(10,2)" json:"minOrderAmount"`
	MaxAmountApplied float64       `gorm:"type:decimal(10,2)" json:"maxAmountApplied"`
	ExpirationTime  time.Time      `json:"expirationTime"`
	UsageLimit      int            `gorm:"default:1" json:"usageLimit"`
	Active          bool           `gorm:"default:true" json:"active"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy       string         `gorm:"-" json:"-"` // Excluded from DB operations
	UpdatedBy       string         `gorm:"-" json:"-"` // Excluded from DB operations
}

// TableName specifies the table name for the Coupon model
func (Coupon) TableName() string {
	return "coupons"
}

// CouponResponse represents the coupon data returned in API responses
type CouponResponse struct {
	ID               string    `json:"id"`
	Code             string    `json:"code"`
	CouponType       string    `json:"couponType"`
	Amount           float64   `json:"amount"`
	MinOrderAmount   float64   `json:"minOrderAmount"`
	MaxAmountApplied float64   `json:"maxAmountApplied"`
	ExpirationTime   time.Time `json:"expirationTime"`
	UsageLimit       int       `json:"usageLimit"`
	Active           bool      `json:"active"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	// Usage statistics (calculated fields)
	TotalUsageCount   int     `json:"totalUsageCount,omitempty"`
	TotalSavingsAmount float64 `json:"totalSavingsAmount,omitempty"`
	UniqueUsersCount  int     `json:"uniqueUsersCount,omitempty"`
}

// ToResponse converts a Coupon model to CouponResponse
func (c *Coupon) ToResponse() CouponResponse {
	return CouponResponse{
		ID:               c.ID,
		Code:             c.Code,
		CouponType:       c.CouponType,
		Amount:           c.Amount,
		MinOrderAmount:   c.MinOrderAmount,
		MaxAmountApplied: c.MaxAmountApplied,
		ExpirationTime:   c.ExpirationTime,
		UsageLimit:       c.UsageLimit,
		Active:           c.Active,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}
