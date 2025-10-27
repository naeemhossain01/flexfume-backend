package models

import (
	"time"

	"gorm.io/gorm"
)

// PaymentStatus represents the payment status of an order
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusPaid      PaymentStatus = "PAID"
	PaymentStatusFailed    PaymentStatus = "FAILED"
	PaymentStatusRefunded  PaymentStatus = "REFUNDED"
)

// DeliveryStatus represents the delivery status of an order
type DeliveryStatus string

const (
	DeliveryStatusPending    DeliveryStatus = "PENDING"
	DeliveryStatusConfirmed  DeliveryStatus = "CONFIRMED"
	DeliveryStatusProcessing DeliveryStatus = "PROCESSING"
	DeliveryStatusShipped    DeliveryStatus = "SHIPPED"
	DeliveryStatusDelivered  DeliveryStatus = "DELIVERED"
	DeliveryStatusCancelled  DeliveryStatus = "CANCELLED"
	DeliveryStatusReturned   DeliveryStatus = "RETURNED"
)

// Order represents an order in the system
type Order struct {
	ID             string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID         string         `gorm:"type:uuid;not null" json:"userId"`
	Address        string         `gorm:"type:text;not null" json:"address"`
	Area           string         `gorm:"type:varchar(255);not null" json:"area"`
	TotalAmount    float64        `gorm:"type:decimal(10,2);not null" json:"totalAmount"`
	PaymentMethod  string         `gorm:"type:varchar(50);not null" json:"paymentMethod"`
	PaymentStatus  PaymentStatus  `gorm:"type:varchar(50);not null;default:'PENDING'" json:"paymentStatus"`
	DeliveryStatus DeliveryStatus `gorm:"type:varchar(50);not null;default:'PENDING'" json:"deliveryStatus"`
	CouponID       *string        `gorm:"type:uuid" json:"couponId,omitempty"`
	DiscountAmount float64        `gorm:"type:decimal(10,2);default:0" json:"discountAmount"`
	DeliveryCost   float64        `gorm:"type:decimal(10,2);default:0" json:"deliveryCost"`
	Notes          string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy      string         `gorm:"-" json:"-"` // Excluded from DB operations
	UpdatedBy      string         `gorm:"-" json:"-"` // Excluded from DB operations

	// Relationships
	User       *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"orderItems,omitempty"`
}

// TableName specifies the table name for the Order model
func (Order) TableName() string {
	return "orders"
}

// OrderItem represents an item in an order
// Note: This table does not support soft deletes (no deleted_at column)
type OrderItem struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	Price     float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	OrderID   string    `gorm:"type:uuid;not null" json:"orderId"`
	ProductID string    `gorm:"type:uuid;not null" json:"productId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	// DeletedAt field removed - table doesn't have deleted_at column
	CreatedBy string    `gorm:"-" json:"-"` // Excluded from DB operations
	UpdatedBy string    `gorm:"-" json:"-"` // Excluded from DB operations

	// Relationships
	Order   *Order   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// TableName specifies the table name for the OrderItem model
func (OrderItem) TableName() string {
	return "order_items"
}

// OrderItemRequest represents the request for creating an order item
type OrderItemRequest struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
	// Price is calculated server-side based on product price and discount
}

// ShippingAddressRequest represents simplified shipping address in order request
type ShippingAddressRequest struct {
	Address string `json:"address" binding:"required"`
	Area    string `json:"area" binding:"required"`
}

// OrderRequest represents the request for creating an order
type OrderRequest struct {
	UserID          string                  `json:"userId,omitempty"` // Set by handler from JWT context, not required in payload
	Items           []OrderItemRequest      `json:"items" binding:"required,min=1"`
	ShippingAddress *ShippingAddressRequest `json:"shippingAddress" binding:"required"`
	PaymentMethod   string                  `json:"paymentMethod,omitempty"`
	DeliveryCostID  string                  `json:"deliveryCostId" binding:"required"` // ID to fetch delivery cost
	CouponCode      string                  `json:"couponCode,omitempty"`              // Optional coupon code
	// TotalAmount, DeliveryCost, and item prices are calculated server-side
}

// OrderItemInfo represents order item information in responses
type OrderItemInfo struct {
	ID          string           `json:"id"`
	Quantity    int              `json:"quantity"`
	Price       float64          `json:"price"`
	ProductInfo *ProductResponse `json:"productInfo,omitempty"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
}

// OrderInfo represents order information in responses
type OrderInfo struct {
	ID             string          `json:"id"`
	Address        string          `json:"address"`
	Area           string          `json:"area"`
	TotalAmount    float64         `json:"totalAmount"`
	PaymentMethod  string          `json:"paymentMethod"`
	PaymentStatus  string          `json:"paymentStatus"`
	DeliveryStatus string          `json:"deliveryStatus"`
	DiscountAmount float64         `json:"discountAmount"`
	DeliveryCost   float64         `json:"deliveryCost"`
	UserInfo       *UserResponse   `json:"userInfo,omitempty"`
	OrderItemList  []OrderItemInfo `json:"orderItemInfoList,omitempty"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

// OrderResponse represents paginated order response
type OrderResponse struct {
	OrderInfoList []OrderInfo `json:"orderInfoList"`
	TotalPage     int         `json:"totalPage"`
	TotalElements int64       `json:"totalElements"`
}

// ToOrderItemInfo converts OrderItem to OrderItemInfo
func (oi *OrderItem) ToOrderItemInfo() OrderItemInfo {
	info := OrderItemInfo{
		ID:        oi.ID,
		Quantity:  oi.Quantity,
		Price:     oi.Price,
		CreatedAt: oi.CreatedAt,
		UpdatedAt: oi.UpdatedAt,
	}

	if oi.Product != nil {
		productResp := oi.Product.ToResponse()
		info.ProductInfo = &productResp
	}

	return info
}

// ToOrderInfo converts Order to OrderInfo
func (o *Order) ToOrderInfo() OrderInfo {
	info := OrderInfo{
		ID:             o.ID,
		Address:        o.Address,
		Area:           o.Area,
		TotalAmount:    o.TotalAmount,
		PaymentMethod:  o.PaymentMethod,
		PaymentStatus:  string(o.PaymentStatus),
		DeliveryStatus: string(o.DeliveryStatus),
		DiscountAmount: o.DiscountAmount,
		DeliveryCost:   o.DeliveryCost,
		CreatedAt:      o.CreatedAt,
		UpdatedAt:      o.UpdatedAt,
	}

	if o.User != nil {
		userResp := o.User.ToResponse()
		info.UserInfo = &userResp
	}

	if len(o.OrderItems) > 0 {
		info.OrderItemList = make([]OrderItemInfo, len(o.OrderItems))
		for i, item := range o.OrderItems {
			info.OrderItemList[i] = item.ToOrderItemInfo()
		}
	}

	return info
}
