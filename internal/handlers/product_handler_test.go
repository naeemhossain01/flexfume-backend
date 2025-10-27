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

func setupProductHandlerTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	err = db.AutoMigrate(&models.Product{}, &models.Category{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	database.DB = db
	return db
}

func teardownProductHandlerTestDB() {
	database.DB = nil
}

func TestProductHandler_CreateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupProductHandlerTestDB(t)
	defer teardownProductHandlerTestDB()

	productService := services.NewProductService()
	handler := NewProductHandler(productService)

	t.Run("Success - Create product", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":        "Test Product",
			"description": "Test Description",
			"price":       100.00,
			"stock":       10,
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/product", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateProduct(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Error - Invalid request body", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": "",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/product", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateProduct(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProductHandler_GetProductByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupProductHandlerTestDB(t)
	defer teardownProductHandlerTestDB()

	productService := services.NewProductService()
	handler := NewProductHandler(productService)

	t.Run("Success - Get product by ID", func(t *testing.T) {
		// Create test product
		product := &models.Product{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       100.00,
			Stock:       10,
		}
		db.Create(product)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: product.ID}}
		c.Request = httptest.NewRequest("GET", "/api/v1/product/"+product.ID, nil)

		handler.GetProductByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Error - Product not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "non-existent-id"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/product/non-existent-id", nil)

		handler.GetProductByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestProductHandler_GetAllProducts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupProductHandlerTestDB(t)
	defer teardownProductHandlerTestDB()

	productService := services.NewProductService()
	handler := NewProductHandler(productService)

	t.Run("Success - Get all products", func(t *testing.T) {
		// Create test products
		db.Create(&models.Product{Name: "Product 1", Price: 100.00, Stock: 10})
		db.Create(&models.Product{Name: "Product 2", Price: 200.00, Stock: 20})

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/product/all", nil)

		handler.GetAllProducts(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
