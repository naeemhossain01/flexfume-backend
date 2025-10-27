package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"gorm.io/gorm"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrInvalidOrderStatus = errors.New("invalid order status")
	ErrInsufficientStock  = errors.New("insufficient stock")
)

// OrderService handles order-related business logic
type OrderService struct {
	userService         *UserService
	productService      *ProductService
	discountService     *DiscountService
	deliveryCostService *DeliveryCostService
	couponService       *CouponService
}

// NewOrderService creates a new order service
func NewOrderService(userService *UserService, productService *ProductService, discountService *DiscountService, deliveryCostService *DeliveryCostService, couponService *CouponService) *OrderService {
	return &OrderService{
		userService:         userService,
		productService:      productService,
		discountService:     discountService,
		deliveryCostService: deliveryCostService,
		couponService:       couponService,
	}
}

// PlaceOrder creates a new order
func (s *OrderService) PlaceOrder(req models.OrderRequest) (*models.Order, error) {
	// Verify user exists
	user, err := s.userService.GetUserByID(req.UserID)
	if err != nil {
		return nil, err
	}

	// Get delivery cost
	deliveryCost, err := s.deliveryCostService.GetDeliveryCostByID(req.DeliveryCostID)
	if err != nil {
		if err == ErrDeliveryCostNotFound {
			return nil, errors.New("delivery cost not found")
		}
		return nil, err
	}

	// Create order items and calculate prices
	var orderItems []models.OrderItem
	var subtotal float64 = 0

	for _, item := range req.Items {
		// Get product details
		product, err := s.productService.GetProductByID(item.ProductID)
		if err != nil {
			if err == ErrProductNotFound {
				return nil, fmt.Errorf("product %s not found", item.ProductID)
			}
			return nil, err
		}

		// Check stock availability
		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("product %s: %w", product.ProductName, ErrInsufficientStock)
		}

		// Calculate price with discount
		itemPrice := product.Price
		discount, err := s.discountService.GetDiscountByProductID(item.ProductID)
		if err == nil && discount != nil {
			// Apply discount percentage
			discountAmount := itemPrice * float64(discount.Percentage) / 100.0
			itemPrice = itemPrice - discountAmount
		}

		// Calculate total price for this item (price * quantity)
		totalItemPrice := itemPrice * float64(item.Quantity)
		subtotal += totalItemPrice

		// Create order item with calculated price per unit
		orderItem := models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     itemPrice, // Store the discounted price per unit
			CreatedBy: user.PhoneNumber,
			UpdatedBy: user.PhoneNumber,
		}
		orderItems = append(orderItems, orderItem)
	}

	// Calculate total amount with delivery cost
	totalAmount := subtotal + deliveryCost.Cost

	// Apply coupon if provided
	var couponID *string
	var discountAmount float64 = 0

	if req.CouponCode != "" {
		coupon, err := s.couponService.GetCouponByCode(req.CouponCode)
		if err != nil {
			if err == ErrCouponNotFound {
				return nil, errors.New("coupon not found")
			}
			return nil, err
		}

		// Validate coupon
		if err := s.couponService.ValidateCoupon(coupon); err != nil {
			return nil, fmt.Errorf("coupon validation failed: %w", err)
		}

		// Check minimum order amount
		if subtotal < coupon.MinOrderAmount {
			return nil, fmt.Errorf("order amount %.2f is less than minimum required %.2f for this coupon", subtotal, coupon.MinOrderAmount)
		}

		// Calculate discount based on coupon type
		if coupon.CouponType == "PERCENTAGE" {
			discountAmount = subtotal * coupon.Amount / 100.0
		} else if coupon.CouponType == "FIXED" {
			discountAmount = coupon.Amount
		}

		// Apply max discount limit if set
		if coupon.MaxAmountApplied > 0 && discountAmount > coupon.MaxAmountApplied {
			discountAmount = coupon.MaxAmountApplied
		}

		// Ensure discount doesn't exceed subtotal
		if discountAmount > subtotal {
			discountAmount = subtotal
		}

		// Apply discount to total
		totalAmount = totalAmount - discountAmount
		couponID = &coupon.ID
	}

	// Validate shipping address
	if req.ShippingAddress == nil {
		return nil, errors.New("shipping address is required")
	}
	if req.ShippingAddress.Address == "" {
		return nil, errors.New("address is required")
	}
	if req.ShippingAddress.Area == "" {
		return nil, errors.New("area is required")
	}

	// Set payment method, default to CASH_ON_DELIVERY if not provided
	paymentMethod := req.PaymentMethod
	if paymentMethod == "" {
		paymentMethod = "CASH_ON_DELIVERY"
	}

	// Create order
	order := models.Order{
		UserID:         req.UserID,
		Address:        req.ShippingAddress.Address,
		Area:           req.ShippingAddress.Area,
		TotalAmount:    totalAmount,
		PaymentMethod:  paymentMethod,
		PaymentStatus:  models.PaymentStatusPending,
		DeliveryStatus: models.DeliveryStatusPending,
		DeliveryCost:   deliveryCost.Cost,
		DiscountAmount: discountAmount,
		CouponID:       couponID,
		CreatedBy:      user.PhoneNumber,
		UpdatedBy:      user.PhoneNumber,
	}

	// Start transaction
	tx := database.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Save order
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Save order items
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
		if err := tx.Create(&orderItems[i]).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Update product stock
		if err := tx.Model(&models.Product{}).
			Where("id = ?", orderItems[i].ProductID).
			Update("stock", gorm.Expr("stock - ?", orderItems[i].Quantity)).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Load relationships for response
	if err := database.GetDB().
		Preload("User").
		Preload("OrderItems.Product.Category").
		First(&order, "id = ?", order.ID).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

// UpdateOrderStatus updates the delivery status of an order
func (s *OrderService) UpdateOrderStatus(orderID string, status string) error {
	// Validate delivery status
	deliveryStatus := models.DeliveryStatus(status)
	validStatuses := []models.DeliveryStatus{
		models.DeliveryStatusPending,
		models.DeliveryStatusConfirmed,
		models.DeliveryStatusProcessing,
		models.DeliveryStatusShipped,
		models.DeliveryStatusDelivered,
		models.DeliveryStatusCancelled,
		models.DeliveryStatusReturned,
	}

	isValid := false
	for _, validStatus := range validStatuses {
		if deliveryStatus == validStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		return ErrInvalidOrderStatus
	}

	// Find order
	var order models.Order
	if err := database.GetDB().First(&order, "id = ?", orderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrOrderNotFound
		}
		return err
	}

	// Update delivery status
	order.DeliveryStatus = deliveryStatus
	if err := database.GetDB().Save(&order).Error; err != nil {
		return err
	}

	return nil
}

// FilterOrders filters orders based on criteria with pagination
func (s *OrderService) FilterOrders(status string, startDate, endDate *time.Time, page, size int) (*models.OrderResponse, error) {
	query := database.GetDB().Model(&models.Order{})

	// Apply filters
	if status != "" {
		query = query.Where("delivery_status = ?", status)
	}

	if startDate != nil {
		query = query.Where("created_at >= ?", startDate)
	}

	if endDate != nil {
		query = query.Where("created_at <= ?", endDate)
	}

	// Count total records
	var totalElements int64
	if err := query.Count(&totalElements).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	offset := page * size
	query = query.Order("created_at DESC").Limit(size).Offset(offset)

	// Fetch orders with relationships
	var orders []models.Order
	if err := query.
		Preload("User").
		Preload("OrderItems.Product.Category").
		Find(&orders).Error; err != nil {
		return nil, err
	}

	// Convert to response format
	orderInfoList := make([]models.OrderInfo, len(orders))
	for i, order := range orders {
		orderInfoList[i] = order.ToOrderInfo()
	}

	// Calculate total pages
	totalPages := int(totalElements) / size
	if int(totalElements)%size != 0 {
		totalPages++
	}

	response := &models.OrderResponse{
		OrderInfoList: orderInfoList,
		TotalPage:     totalPages,
		TotalElements: totalElements,
	}

	return response, nil
}

// GetOrderHistory retrieves order history for a specific user
func (s *OrderService) GetOrderHistory(userID string) ([]models.OrderInfo, error) {
	// Verify user exists
	if _, err := s.userService.GetUserByID(userID); err != nil {
		return nil, err
	}

	// Fetch orders
	var orders []models.Order
	if err := database.GetDB().
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Preload("User").
		Preload("OrderItems.Product.Category").
		Find(&orders).Error; err != nil {
		return nil, err
	}

	// Convert to response format
	orderInfoList := make([]models.OrderInfo, len(orders))
	for i, order := range orders {
		orderInfoList[i] = order.ToOrderInfo()
	}

	return orderInfoList, nil
}
