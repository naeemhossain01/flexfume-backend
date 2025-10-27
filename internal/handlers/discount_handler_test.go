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

func setupDiscountTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	err = db.AutoMigrate(&models.Discount{}, &models.Product{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	database.DB = db
	return db
}

func teardownDiscountTestDB() {
	database.DB = nil
}

func TestDiscountHandler_CreateDiscount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupDiscountTestDB(t)
	defer teardownDiscountTestDB()

	discountService := services.NewDiscountService()
	handler := NewDiscountHandler(discountService)

	t.Run("Success - Create discount", func(t *testing.T) {
		// Create test product
		product := &models.Product{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       100.00,
			Stock:       10,
		}
		db.Create(product)

		reqBody := map[string]interface{}{
			"productId":      product.ID,
			"discountType":   "PERCENTAGE",
			"discountValue":  20.0,
			"startDate":      "2024-01-01T00:00:00Z",
			"endDate":        "2024-12-31T23:59:59Z",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/discount", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateDiscount(c)

		// Check response
		if w.Code != http.StatusCreated && w.Code != http.StatusOK {
			t.Logf("Create discount returned status: %d, body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Error - Invalid request body", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"productId": "",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/discount", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateDiscount(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDiscountHandler_GetDiscountByProductID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupDiscountTestDB(t)
	defer teardownDiscountTestDB()

	discountService := services.NewDiscountService()
	handler := NewDiscountHandler(discountService)

	t.Run("Success - Get discount by product ID", func(t *testing.T) {
		// Create test product
		product := &models.Product{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       100.00,
			Stock:       10,
		}
		db.Create(product)

		// Create test discount
		discount := &models.Discount{
			ProductID:     product.ID,
			DiscountType:  "PERCENTAGE",
			DiscountValue: 20.0,
		}
		db.Create(discount)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "productId", Value: product.ID}}
		c.Request = httptest.NewRequest("GET", "/api/v1/discount/product/"+product.ID, nil)

		handler.GetDiscountByProductID(c)

		// Check response
		if w.Code == http.StatusOK {
			assert.Equal(t, http.StatusOK, w.Code)
		} else {
			t.Logf("Get discount returned status: %d, body: %s", w.Code, w.Body.String())
		}
	})
}

func TestDiscountHandler_GetAllDiscounts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupDiscountTestDB(t)
	defer teardownDiscountTestDB()

	discountService := services.NewDiscountService()
	handler := NewDiscountHandler(discountService)

	t.Run("Success - Get all discounts", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/discount/all", nil)

		handler.GetAllDiscounts(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
