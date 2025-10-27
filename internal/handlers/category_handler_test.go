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

func setupCategoryTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	err = db.AutoMigrate(&models.Category{}, &models.Product{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	database.DB = db
	return db
}

func teardownCategoryTestDB() {
	database.DB = nil
}

func createTestCategoryInDB(t *testing.T, db *gorm.DB, name, description string) *models.Category {
	category := &models.Category{
		Name:        name,
		Description: description,
	}
	if err := db.Create(category).Error; err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}
	return category
}

func TestCategoryHandler_CreateCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupCategoryTestDB(t)
	defer teardownCategoryTestDB()

	categoryService := services.NewCategoryService()
	handler := NewCategoryHandler(categoryService)

	t.Run("Success - Create valid category", func(t *testing.T) {
		reqBody := CreateCategoryRequest{
			Name:        "Electronics",
			Description: "Electronic devices",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/category/add", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateCategory(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))
		assert.NotNil(t, response["data"])

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "Electronics", data["name"])
		assert.Equal(t, "Electronic devices", data["description"])
	})

	t.Run("Error - Missing required name", func(t *testing.T) {
		reqBody := CreateCategoryRequest{
			Description: "Test description",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/category/add", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateCategory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response["success"].(bool))
		assert.NotNil(t, response["error"])
	})

	t.Run("Error - Duplicate category name", func(t *testing.T) {
		// Create existing category
		createTestCategoryInDB(t, db, "Books", "Book collection")

		reqBody := CreateCategoryRequest{
			Name:        "Books",
			Description: "Another books category",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/category/add", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateCategory(c)

		assert.Equal(t, http.StatusConflict, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response["success"].(bool))
	})

	t.Run("Error - Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/category/add", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateCategory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCategoryHandler_UpdateCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupCategoryTestDB(t)
	defer teardownCategoryTestDB()

	categoryService := services.NewCategoryService()
	handler := NewCategoryHandler(categoryService)

	t.Run("Success - Update category", func(t *testing.T) {
		existing := createTestCategoryInDB(t, db, "Old Name", "Old description")

		reqBody := UpdateCategoryRequest{
			Name:        "New Name",
			Description: "New description",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PUT", "/api/v1/category/update/"+existing.ID, bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: existing.ID}}

		handler.UpdateCategory(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "New Name", data["name"])
		assert.Equal(t, "New description", data["description"])
	})

	t.Run("Error - Category not found", func(t *testing.T) {
		reqBody := UpdateCategoryRequest{
			Name: "Test",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PUT", "/api/v1/category/update/non-existent", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "non-existent"}}

		handler.UpdateCategory(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Error - Duplicate name", func(t *testing.T) {
		cat1 := createTestCategoryInDB(t, db, "Category1", "Description 1")
		createTestCategoryInDB(t, db, "Category2", "Description 2")

		reqBody := UpdateCategoryRequest{
			Name: "Category2", // Try to use existing name
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PUT", "/api/v1/category/update/"+cat1.ID, bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: cat1.ID}}

		handler.UpdateCategory(c)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestCategoryHandler_GetAllCategories(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupCategoryTestDB(t)
	defer teardownCategoryTestDB()

	categoryService := services.NewCategoryService()
	handler := NewCategoryHandler(categoryService)

	t.Run("Success - Get all categories", func(t *testing.T) {
		createTestCategoryInDB(t, db, "Category1", "Description 1")
		createTestCategoryInDB(t, db, "Category2", "Description 2")
		createTestCategoryInDB(t, db, "Category3", "Description 3")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/category/all", nil)

		handler.GetAllCategories(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		data := response["data"].([]interface{})
		assert.Len(t, data, 3)
	})

	t.Run("Success - Empty list", func(t *testing.T) {
		db.Exec("DELETE FROM categories")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/category/all", nil)

		handler.GetAllCategories(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		data := response["data"].([]interface{})
		assert.Len(t, data, 0)
	})
}

func TestCategoryHandler_GetCategoryByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupCategoryTestDB(t)
	defer teardownCategoryTestDB()

	categoryService := services.NewCategoryService()
	handler := NewCategoryHandler(categoryService)

	t.Run("Success - Get category by ID", func(t *testing.T) {
		category := createTestCategoryInDB(t, db, "Test Category", "Test description")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/category/get/"+category.ID, nil)
		c.Params = gin.Params{{Key: "id", Value: category.ID}}

		handler.GetCategoryByID(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "Test Category", data["name"])
		assert.Equal(t, "Test description", data["description"])
	})

	t.Run("Error - Category not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/category/get/non-existent", nil)
		c.Params = gin.Params{{Key: "id", Value: "non-existent"}}

		handler.GetCategoryByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response["success"].(bool))
	})
}

func TestCategoryHandler_DeleteCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupCategoryTestDB(t)
	defer teardownCategoryTestDB()

	categoryService := services.NewCategoryService()
	handler := NewCategoryHandler(categoryService)

	t.Run("Success - Delete category", func(t *testing.T) {
		category := createTestCategoryInDB(t, db, "ToDelete", "Will be deleted")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("DELETE", "/api/v1/category/delete/"+category.ID, nil)
		c.Params = gin.Params{{Key: "id", Value: category.ID}}

		handler.DeleteCategory(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))
		assert.Equal(t, "Deleted", response["data"])

		// Verify deletion
		var count int64
		db.Model(&models.Category{}).Where("id = ?", category.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("Error - Category not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("DELETE", "/api/v1/category/delete/non-existent", nil)
		c.Params = gin.Params{{Key: "id", Value: "non-existent"}}

		handler.DeleteCategory(c)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response["success"].(bool))
	})
}

