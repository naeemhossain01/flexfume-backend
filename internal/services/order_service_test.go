package services

import (
	"testing"

	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

// setupOrderTestDB creates an in-memory SQLite database for testing
func setupOrderTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.Address{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Coupon{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Set the global DB for testing
	database.DB = db

	return db
}

// teardownOrderTestDB closes the database connection
func teardownOrderTestDB() {
	database.DB = nil
}

func TestOrderService_CreateOrder(t *testing.T) {
	db := setupOrderTestDB(t)
	defer teardownOrderTestDB()

	service := NewOrderService()

	t.Run("Success - Create order", func(t *testing.T) {
		// Create test user
		user := &models.User{
			Email:    "test@example.com",
			Password: "hashedpassword",
			Name:     "Test User",
		}
		db.Create(user)

		// Create test address
		address := &models.Address{
			UserID:      user.ID,
			FullName:    "Test User",
			PhoneNumber: "+1234567890",
			District:    "Dhaka",
			Address:     "123 Test St",
		}
		db.Create(address)

		// Create test product
		product := &models.Product{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       100.00,
			Stock:       10,
		}
		db.Create(product)

		// Create order
		order := &models.Order{
			UserID:         user.ID,
			Address:        "123 Test St",
			Area:           "Dhaka",
			TotalAmount:    100.00,
			PaymentMethod:  "COD",
			PaymentStatus:  "PENDING",
			DeliveryStatus: "PENDING",
		}

		result, err := service.CreateOrder(order)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.ID)
	})
}

func TestOrderService_GetOrderByID(t *testing.T) {
	db := setupOrderTestDB(t)
	defer teardownOrderTestDB()

	service := NewOrderService()

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
			Address:        "456 Test Ave",
			Area:           "Chittagong",
			TotalAmount:    100.00,
			PaymentMethod:  "COD",
			PaymentStatus:  "PENDING",
			DeliveryStatus: "PENDING",
		}
		db.Create(order)

		result, err := service.GetOrderByID(order.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, order.ID, result.ID)
	})

	t.Run("Error - Order not found", func(t *testing.T) {
		result, err := service.GetOrderByID("non-existent-id")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
