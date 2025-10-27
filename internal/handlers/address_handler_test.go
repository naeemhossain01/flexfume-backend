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
)

func setupAddressTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	err = db.AutoMigrate(&models.User{}, &models.Address{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	database.DB = db
	return db
}

func teardownAddressTestDB() {
	database.DB = nil
}

func createTestUserForAddress(t *testing.T, db *gorm.DB, name, email, phoneNumber string) *models.User {
	user := &models.User{
		Name:        name,
		Email:       email,
		PhoneNumber: phoneNumber,
		Password:    "hashedpassword",
		Role:        "USER",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

func createTestAddressInDB(t *testing.T, db *gorm.DB, userID, buildingName, road, area, city string) *models.Address {
	address := &models.Address{
		UserID:       userID,
		BuildingName: buildingName,
		Road:         road,
		Area:         area,
		City:         city,
	}
	if err := db.Create(address).Error; err != nil {
		t.Fatalf("Failed to create test address: %v", err)
	}
	return address
}

func TestAddressHandler_AddAddress(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAddressTestDB(t)
	defer teardownAddressTestDB()

	userService := services.NewUserService(nil)
	addressService := services.NewAddressService(userService)
	handler := NewAddressHandler(addressService)

	t.Run("Success - Add valid address", func(t *testing.T) {
		user := createTestUserForAddress(t, db, "John Doe", "john@example.com", "1234567890")

		reqBody := models.AddressRequest{
			UserID:       user.ID,
			BuildingName: "Building A",
			Road:         "Main Street",
			Area:         "Downtown",
			City:         "New York",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/address", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddAddress(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))
		assert.NotNil(t, response["data"])

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "Building A", data["buildingName"])
		assert.Equal(t, "Main Street", data["road"])
		assert.Equal(t, "Downtown", data["area"])
		assert.Equal(t, "New York", data["city"])
	})

	t.Run("Error - Missing required userId", func(t *testing.T) {
		reqBody := models.AddressRequest{
			BuildingName: "Building A",
			Road:         "Main Street",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/address", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddAddress(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response["success"].(bool))
		assert.NotNil(t, response["error"])
	})

	t.Run("Error - User not found", func(t *testing.T) {
		reqBody := models.AddressRequest{
			UserID:       "non-existent-user-id",
			BuildingName: "Building A",
			Road:         "Main Street",
			Area:         "Downtown",
			City:         "New York",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/address", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddAddress(c)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response["success"].(bool))
	})

	t.Run("Error - Address already exists for user", func(t *testing.T) {
		user := createTestUserForAddress(t, db, "Jane Doe", "jane@example.com", "0987654321")
		createTestAddressInDB(t, db, user.ID, "Old Building", "Old Street", "Old Area", "Old City")

		reqBody := models.AddressRequest{
			UserID:       user.ID,
			BuildingName: "New Building",
			Road:         "New Street",
			Area:         "New Area",
			City:         "New City",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/address", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddAddress(c)

		assert.Equal(t, http.StatusConflict, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response["success"].(bool))
	})

	t.Run("Error - Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/address", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.AddAddress(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAddressHandler_UpdateAddress(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAddressTestDB(t)
	defer teardownAddressTestDB()

	userService := services.NewUserService(nil)
	addressService := services.NewAddressService(userService)
	handler := NewAddressHandler(addressService)

	t.Run("Success - Update address", func(t *testing.T) {
		user := createTestUserForAddress(t, db, "John Smith", "john.smith@example.com", "1111111111")
		address := createTestAddressInDB(t, db, user.ID, "Old Building", "Old Street", "Old Area", "Old City")

		reqBody := models.AddressRequest{
			UserID:       user.ID,
			BuildingName: "New Building",
			Road:         "New Street",
			Area:         "New Area",
			City:         "New City",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PUT", "/api/v1/address/"+address.ID, bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: address.ID}}

		handler.UpdateAddress(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "New Building", data["buildingName"])
		assert.Equal(t, "New Street", data["road"])
		assert.Equal(t, "New Area", data["area"])
		assert.Equal(t, "New City", data["city"])
	})

	t.Run("Success - Partial update", func(t *testing.T) {
		user := createTestUserForAddress(t, db, "Jane Smith", "jane.smith@example.com", "2222222222")
		address := createTestAddressInDB(t, db, user.ID, "Building X", "Street X", "Area X", "City X")

		reqBody := models.AddressRequest{
			UserID: user.ID,
			City:   "Updated City",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PUT", "/api/v1/address/"+address.ID, bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: address.ID}}

		handler.UpdateAddress(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "Building X", data["buildingName"]) // Unchanged
		assert.Equal(t, "Updated City", data["city"])       // Changed
	})

	t.Run("Error - Address not found", func(t *testing.T) {
		user := createTestUserForAddress(t, db, "Test User", "test@example.com", "3333333333")

		reqBody := models.AddressRequest{
			UserID: user.ID,
			City:   "Test City",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PUT", "/api/v1/address/non-existent", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "non-existent"}}

		handler.UpdateAddress(c)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response["success"].(bool))
	})

	t.Run("Error - Wrong user trying to update", func(t *testing.T) {
		user1 := createTestUserForAddress(t, db, "User One", "user1@example.com", "4444444444")
		user2 := createTestUserForAddress(t, db, "User Two", "user2@example.com", "5555555555")
		address := createTestAddressInDB(t, db, user1.ID, "Building", "Street", "Area", "City")

		reqBody := models.AddressRequest{
			UserID: user2.ID, // Different user
			City:   "Hacked City",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PUT", "/api/v1/address/"+address.ID, bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: address.ID}}

		handler.UpdateAddress(c)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response["success"].(bool))
	})

	t.Run("Error - Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("PUT", "/api/v1/address/some-id", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "some-id"}}

		handler.UpdateAddress(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAddressHandler_GetAddressByUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupAddressTestDB(t)
	defer teardownAddressTestDB()

	userService := services.NewUserService(nil)
	addressService := services.NewAddressService(userService)
	handler := NewAddressHandler(addressService)

	t.Run("Success - Get address by user ID", func(t *testing.T) {
		user := createTestUserForAddress(t, db, "Alice Johnson", "alice@example.com", "6666666666")
		address := createTestAddressInDB(t, db, user.ID, "Building B", "Second Street", "Uptown", "Boston")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/address/"+user.ID, nil)
		c.Params = gin.Params{{Key: "userId", Value: user.ID}}

		handler.GetAddressByUser(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, address.ID, data["id"])
		assert.Equal(t, "Building B", data["buildingName"])
		assert.Equal(t, "Second Street", data["road"])
		assert.Equal(t, "Uptown", data["area"])
		assert.Equal(t, "Boston", data["city"])
	})

	t.Run("Error - Address not found for user", func(t *testing.T) {
		user := createTestUserForAddress(t, db, "Bob Wilson", "bob@example.com", "7777777777")
		// No address created for this user

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/address/"+user.ID, nil)
		c.Params = gin.Params{{Key: "userId", Value: user.ID}}

		handler.GetAddressByUser(c)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response["success"].(bool))
	})

	t.Run("Error - Invalid user ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/address/non-existent-user", nil)
		c.Params = gin.Params{{Key: "userId", Value: "non-existent-user"}}

		handler.GetAddressByUser(c)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response["success"].(bool))
	})
}
