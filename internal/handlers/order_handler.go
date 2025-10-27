package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
)

// OrderHandler handles order-related requests
type OrderHandler struct {
	orderService *services.OrderService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// PlaceOrder creates a new order
func (h *OrderHandler) PlaceOrder(c *gin.Context) {
	// Get user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Error:   true,
			Message: "User not authenticated",
		})
		return
	}

	var req models.OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Set user ID from JWT context
	req.UserID = userID.(string)

	order, err := h.orderService.PlaceOrder(req)
	if err != nil {
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: "User not found",
			})
			return
		}
		if err == services.ErrInsufficientStock || strings.Contains(err.Error(), "insufficient stock") {
			c.JSON(http.StatusBadRequest, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: "Failed to place order: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: order.ToOrderInfo(),
	})
}

// UpdateOrderStatus updates the status of an order (Admin only)
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("orderId")
	status := c.Query("status")

	if status == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "Status parameter is required",
		})
		return
	}

	// Convert status to uppercase to match enum
	status = strings.ToUpper(status)

	err := h.orderService.UpdateOrderStatus(orderID, status)
	if err != nil {
		switch err {
		case services.ErrOrderNotFound:
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: "Order not found",
			})
		case services.ErrInvalidOrderStatus:
			c.JSON(http.StatusBadRequest, APIResponse{
				Error:   true,
				Message: "Invalid order status. Valid statuses: PENDING, CONFIRMED, PROCESSING, SHIPPED, DELIVERED, CANCELLED, RETURNED",
			})
		default:
			c.JSON(http.StatusInternalServerError, APIResponse{
				Error:   true,
				Message: "Failed to update order status: " + err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: "Order status changes successfully",
	})
}

// FilterOrders filters orders by various criteria (Admin only)
func (h *OrderHandler) FilterOrders(c *gin.Context) {
	// Parse query parameters
	status := c.Query("status")
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	pageStr := c.DefaultQuery("page", "0")
	sizeStr := c.DefaultQuery("size", "1000")

	// Parse page and size
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "Invalid page parameter",
		})
		return
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "Invalid size parameter",
		})
		return
	}

	// Parse dates if provided
	var startDate, endDate *time.Time
	if startDateStr != "" {
		parsedStartDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{
				Error:   true,
				Message: "Invalid startDate format. Use ISO 8601 format (e.g., 2024-01-01T00:00:00Z)",
			})
			return
		}
		startDate = &parsedStartDate
	}

	if endDateStr != "" {
		parsedEndDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{
				Error:   true,
				Message: "Invalid endDate format. Use ISO 8601 format (e.g., 2024-12-31T23:59:59Z)",
			})
			return
		}
		endDate = &parsedEndDate
	}

	// Convert status to uppercase if provided
	if status != "" {
		status = strings.ToUpper(status)
	}

	// Filter orders
	response, err := h.orderService.FilterOrders(status, startDate, endDate, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: "Failed to filter orders: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: response,
	})
}

// GetOrderHistory retrieves order history for the authenticated user
func (h *OrderHandler) GetOrderHistory(c *gin.Context) {
	// Get user ID from JWT context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Error:   true,
			Message: "User not authenticated",
		})
		return
	}

	userID := userIDInterface.(string)

	orders, err := h.orderService.GetOrderHistory(userID)
	if err != nil {
		if err == services.ErrUserNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: "Failed to retrieve order history: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: orders,
	})
}
