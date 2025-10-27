package models

import (
	"time"
)

// CouponUsage represents a coupon usage record in the system
// Note: This table does not support soft deletes (no deleted_at column)
type CouponUsage struct {
	ID               string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CouponID         string    `gorm:"type:uuid;not null" json:"couponId"`
	UserID           string    `gorm:"type:uuid;not null" json:"userId"`
	UsageCount       int       `gorm:"default:0" json:"usageCount"`
	DiscountedAmount float64   `gorm:"type:decimal(10,2)" json:"discountedAmount"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	// DeletedAt field removed - table doesn't have deleted_at column
	CreatedBy        string    `gorm:"-" json:"-"` // Excluded from DB operations
	UpdatedBy        string    `gorm:"-" json:"-"` // Excluded from DB operations

	// Relationships
	Coupon *Coupon `gorm:"foreignKey:CouponID" json:"coupon,omitempty"`
	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for the CouponUsage model
func (CouponUsage) TableName() string {
	return "coupon_usages"
}

// CouponUsageResponse represents the coupon usage data returned in API responses
type CouponUsageResponse struct {
	ID               string    `json:"id"`
	CouponID         string    `json:"couponId"`
	UserID           string    `json:"userId"`
	UsageCount       int       `json:"usageCount"`
	DiscountedAmount float64   `json:"discountedAmount"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// ToResponse converts a CouponUsage model to CouponUsageResponse
func (cu *CouponUsage) ToResponse() CouponUsageResponse {
	return CouponUsageResponse{
		ID:               cu.ID,
		CouponID:         cu.CouponID,
		UserID:           cu.UserID,
		UsageCount:       cu.UsageCount,
		DiscountedAmount: cu.DiscountedAmount,
		CreatedAt:        cu.CreatedAt,
		UpdatedAt:        cu.UpdatedAt,
	}
}
