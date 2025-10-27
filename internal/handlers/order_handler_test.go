package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func setupOrderTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	err = db.AutoMigrate(&models.Order{}, &models.OrderItem{}, &models.User{}, &models.Product{}, &models.Address{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	database.DB = db
	return db
}

func teardownOrderTestDB() {
	database.DB = nil
}

func TestOrderHandler_GetOrderByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupOrderTestDB(t)
	defer teardownOrderTestDB()

	orderService := services.NewOrderService()
	handler := NewOrderHandler(orderService)

	t.Run("Success - Get order by ID", func(t *testing.T) {
		// Create test user
		user := &models.User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Name:     "Test User",
		}
		db.Create(user)

		// Create test order
		order := &models.Order{
			UserID:         user.ID,
			TotalAmount:    100.00,
			PaymentMethod:  "COD",
			PaymentStatus:  "PENDING",
			DeliveryStatus: "PENDING",
		}
		db.Create(order)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: order.ID}}
		c.Request = httptest.NewRequest("GET", "/api/v1/order/"+order.ID, nil)

		handler.GetOrderByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Error - Order not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "non-existent-id"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/order/non-existent-id", nil)

		handler.GetOrderByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestOrderHandler_GetUserOrders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupOrderTestDB(t)
	defer teardownOrderTestDB()

	orderService := services.NewOrderService()
	handler := NewOrderHandler(orderService)

	t.Run("Success - Get user orders", func(t *testing.T) {
		// Create test user
		user := &models.User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Name:     "Test User",
		}
		db.Create(user)

		// Create test orders
		db.Create(&models.Order{UserID: user.ID, TotalAmount: 100.00, PaymentMethod: "COD", PaymentStatus: "PENDING", DeliveryStatus: "PENDING"})
		db.Create(&models.Order{UserID: user.ID, TotalAmount: 200.00, PaymentMethod: "COD", PaymentStatus: "PENDING", DeliveryStatus: "PENDING"})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", user.ID)
		c.Request = httptest.NewRequest("GET", "/api/v1/order/user", nil)

		handler.GetUserOrders(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Error - Missing user ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/order/user", nil)

		handler.GetUserOrders(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
