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

func setupDeliveryCostTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	err = db.AutoMigrate(&models.DeliveryCost{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	database.DB = db
	return db
}

func teardownDeliveryCostTestDB() {
	database.DB = nil
}

func TestDeliveryCostHandler_CreateDeliveryCost(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupDeliveryCostTestDB(t)
	defer teardownDeliveryCostTestDB()

	deliveryCostService := services.NewDeliveryCostService()
	handler := NewDeliveryCostHandler(deliveryCostService)

	t.Run("Success - Create delivery cost", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"district": "Dhaka",
			"cost":     50.00,
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/delivery-cost", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateDeliveryCost(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Error - Invalid request body", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"district": "",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/delivery-cost", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateDeliveryCost(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeliveryCostHandler_GetDeliveryCostByDistrict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupDeliveryCostTestDB(t)
	defer teardownDeliveryCostTestDB()

	deliveryCostService := services.NewDeliveryCostService()
	handler := NewDeliveryCostHandler(deliveryCostService)

	t.Run("Success - Get delivery cost by district", func(t *testing.T) {
		// Create test delivery cost
		db.Create(&models.DeliveryCost{District: "Chittagong", Cost: 60.00})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "district", Value: "Chittagong"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/delivery-cost/Chittagong", nil)

		handler.GetDeliveryCostByDistrict(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Error - District not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "district", Value: "NonExistent"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/delivery-cost/NonExistent", nil)

		handler.GetDeliveryCostByDistrict(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDeliveryCostHandler_GetAllDeliveryCosts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupDeliveryCostTestDB(t)
	defer teardownDeliveryCostTestDB()

	deliveryCostService := services.NewDeliveryCostService()
	handler := NewDeliveryCostHandler(deliveryCostService)

	t.Run("Success - Get all delivery costs", func(t *testing.T) {
		// Create test delivery costs
		db.Create(&models.DeliveryCost{District: "Dhaka", Cost: 50.00})
		db.Create(&models.DeliveryCost{District: "Chittagong", Cost: 60.00})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/delivery-cost/all", nil)

		handler.GetAllDeliveryCosts(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
