package models

import (
	"time"
)

// CouponUsage represents a coupon usage record in the system
// Note: This table does not support soft deletes (no deleted_at column)
type CouponUsage struct {
	ID       string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CouponID string    `gorm:"type:uuid;not null" json:"couponId"`
	UserID   string    `gorm:"type:uuid;not null" json:"userId"`
	OrderID  *string   `gorm:"type:uuid" json:"orderId,omitempty"`
	UsedAt   time.Time `gorm:"autoCreateTime" json:"usedAt"`

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
	ID       string    `json:"id"`
	CouponID string    `json:"couponId"`
	UserID   string    `json:"userId"`
	OrderID  *string   `json:"orderId,omitempty"`
	UsedAt   time.Time `json:"usedAt"`
}

// ToResponse converts a CouponUsage model to CouponUsageResponse
func (cu *CouponUsage) ToResponse() CouponUsageResponse {
	return CouponUsageResponse{
		ID:       cu.ID,
		CouponID: cu.CouponID,
		UserID:   cu.UserID,
		OrderID:  cu.OrderID,
		UsedAt:   cu.UsedAt,
	}
}
