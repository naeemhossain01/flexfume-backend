package models

import (
	"time"

	"gorm.io/gorm"
)

// DeliveryCost represents delivery cost configuration in the system
type DeliveryCost struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Location  string         `gorm:"column:location" json:"location"`
	Service   string         `gorm:"column:service" json:"service"`
	Cost      float64        `gorm:"column:cost" json:"cost"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy string         `json:"-"`
	UpdatedBy string         `json:"-"`
}

// TableName specifies the table name for the DeliveryCost model
func (DeliveryCost) TableName() string {
	return "delivery_costs"
}

// DeliveryCostRequest represents the request body for creating/updating delivery cost
type DeliveryCostRequest struct {
	Location string   `json:"location" binding:"required"`
	Service  string   `json:"service"`
	Cost     *float64 `json:"cost" binding:"required,gte=0"`
}

// DeliveryCostInfo represents delivery cost information in responses
type DeliveryCostInfo struct {
	ID        string    `json:"id"`
	Location  string    `json:"location"`
	Service   string    `json:"service"`
	Cost      float64   `json:"cost"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToDeliveryCostInfo converts DeliveryCost to DeliveryCostInfo
func (d *DeliveryCost) ToDeliveryCostInfo() DeliveryCostInfo {
	return DeliveryCostInfo{
		ID:        d.ID,
		Location:  d.Location,
		Service:   d.Service,
		Cost:      d.Cost,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
