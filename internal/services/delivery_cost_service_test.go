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

// setupDeliveryCostTestDB creates an in-memory SQLite database for testing
func setupDeliveryCostTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.DeliveryCost{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Set the global DB for testing
	database.DB = db

	return db
}

// teardownDeliveryCostTestDB closes the database connection
func teardownDeliveryCostTestDB() {
	database.DB = nil
}

func TestDeliveryCostService_CreateDeliveryCost(t *testing.T) {
	db := setupDeliveryCostTestDB(t)
	defer teardownDeliveryCostTestDB()

	service := NewDeliveryCostService()

	t.Run("Success - Create delivery cost", func(t *testing.T) {
		deliveryCost := &models.DeliveryCost{
			District: "Dhaka",
			Cost:     50.00,
		}

		result, err := service.CreateDeliveryCost(deliveryCost)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, "Dhaka", result.District)
		assert.Equal(t, 50.00, result.Cost)
	})

	t.Run("Error - Nil delivery cost", func(t *testing.T) {
		result, err := service.CreateDeliveryCost(nil)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestDeliveryCostService_GetDeliveryCostByDistrict(t *testing.T) {
	db := setupDeliveryCostTestDB(t)
	defer teardownDeliveryCostTestDB()

	service := NewDeliveryCostService()

	t.Run("Success - Get delivery cost by district", func(t *testing.T) {
		// Create test delivery cost
		deliveryCost := &models.DeliveryCost{
			District: "Chittagong",
			Cost:     60.00,
		}
		db.Create(deliveryCost)

		result, err := service.GetDeliveryCostByDistrict("Chittagong")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Chittagong", result.District)
		assert.Equal(t, 60.00, result.Cost)
	})

	t.Run("Error - District not found", func(t *testing.T) {
		result, err := service.GetDeliveryCostByDistrict("NonExistent")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestDeliveryCostService_GetAllDeliveryCosts(t *testing.T) {
	db := setupDeliveryCostTestDB(t)
	defer teardownDeliveryCostTestDB()

	service := NewDeliveryCostService()

	t.Run("Success - Get all delivery costs", func(t *testing.T) {
		// Create test delivery costs
		db.Create(&models.DeliveryCost{District: "Dhaka", Cost: 50.00})
		db.Create(&models.DeliveryCost{District: "Chittagong", Cost: 60.00})

		result, err := service.GetAllDeliveryCosts()

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
	})
}
