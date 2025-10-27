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

// setupProductTestDB creates an in-memory SQLite database for testing
func setupProductTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.Product{}, &models.Category{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Set the global DB for testing
	database.DB = db

	return db
}

// teardownProductTestDB closes the database connection
func teardownProductTestDB() {
	database.DB = nil
}

func TestProductService_CreateProduct(t *testing.T) {
	db := setupProductTestDB(t)
	defer teardownProductTestDB()

	service := NewProductService()

	t.Run("Success - Create product", func(t *testing.T) {
		product := &models.Product{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       100.00,
			Stock:       10,
		}

		result, err := service.CreateProduct(product)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, "Test Product", result.Name)
	})

	t.Run("Error - Nil product", func(t *testing.T) {
		result, err := service.CreateProduct(nil)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestProductService_GetProductByID(t *testing.T) {
	db := setupProductTestDB(t)
	defer teardownProductTestDB()

	service := NewProductService()

	t.Run("Success - Get product by ID", func(t *testing.T) {
		// Create test product
		product := &models.Product{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       100.00,
			Stock:       10,
		}
		db.Create(product)

		result, err := service.GetProductByID(product.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, product.ID, result.ID)
	})

	t.Run("Error - Product not found", func(t *testing.T) {
		result, err := service.GetProductByID("non-existent-id")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestProductService_GetAllProducts(t *testing.T) {
	db := setupProductTestDB(t)
	defer teardownProductTestDB()

	service := NewProductService()

	t.Run("Success - Get all products", func(t *testing.T) {
		// Create test products
		db.Create(&models.Product{Name: "Product 1", Price: 100.00, Stock: 10})
		db.Create(&models.Product{Name: "Product 2", Price: 200.00, Stock: 20})

		result, err := service.GetAllProducts()

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
	})
}
