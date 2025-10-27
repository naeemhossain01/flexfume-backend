package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/auth"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

// Mock OTP Service for handler testing
type MockCheckoutOTPService struct {
	shouldFailSend   bool
	shouldFailVerify bool
}

func (m *MockCheckoutOTPService) GenerateOTP() (string, error) {
	return "1234", nil
}

func (m *MockCheckoutOTPService) SendOTP(phoneNumber, otpType string) error {
	if m.shouldFailSend {
		return assert.AnError
	}
	return nil
}

func (m *MockCheckoutOTPService) VerifyOTP(phoneNumber, otp string) error {
	if m.shouldFailVerify {
		return services.ErrInvalidOTP
	}
	return nil
}

func (m *MockCheckoutOTPService) MarkPhoneAsVerified(phoneNumber, otp string) error {
	return nil
}

func (m *MockCheckoutOTPService) IsPhoneVerified(phoneNumber string) (bool, error) {
	return true, nil
}

func (m *MockCheckoutOTPService) CheckPhoneNotVerified(phoneNumber string) error {
	return nil
}

func (m *MockCheckoutOTPService) RequirePhoneVerified(phoneNumber string) error {
	return nil
}

func (m *MockCheckoutOTPService) VerifyResetPasswordOTP(phoneNumber, otp string) error {
	return nil
}

func setupCheckoutHandlerTestDB(t *testing.T) *gorm.DB {
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

func teardownCheckoutHandlerTestDB() {
	database.DB = nil
}

func createTestUserForCheckoutHandler(t *testing.T, db *gorm.DB, phoneNumber string) *models.User {
	hashedPassword, _ := auth.HashPassword("password123")
	user := &models.User{
		Name:        "Test User",
		Email:       "test@example.com",
		PhoneNumber: phoneNumber,
		Password:    hashedPassword,
		Role:        "USER",
		CreatedBy:   phoneNumber,
	}

	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return user
}

func TestCheckoutHandler_SendCheckoutOTP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupCheckoutHandlerTestDB(t)
	defer teardownCheckoutHandlerTestDB()

	mockOTPService := &MockCheckoutOTPService{}
	userService := services.NewUserService(nil)
	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	checkoutService := services.NewCheckoutService(mockOTPService, userService, jwtManager)
	handler := NewCheckoutHandler(mockOTPService, checkoutService)

	t.Run("Success - Send OTP", func(t *testing.T) {
		reqBody := models.CheckoutOTPRequest{
			PhoneNumber: "+8801712345678",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/checkout/send-otp", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.SendCheckoutOTP(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))
		assert.Equal(t, "OTP sent successfully for checkout verification", response["data"])
	})

	t.Run("Error - Missing phone number", func(t *testing.T) {
		reqBody := map[string]interface{}{}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/checkout/send-otp", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.SendCheckoutOTP(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response["success"].(bool))
		assert.NotNil(t, response["error"])
	})

	t.Run("Error - Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/checkout/send-otp", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.SendCheckoutOTP(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error - OTP service failure", func(t *testing.T) {
		mockOTPService.shouldFailSend = true
		defer func() { mockOTPService.shouldFailSend = false }()

		reqBody := models.CheckoutOTPRequest{
			PhoneNumber: "+8801712345678",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/checkout/send-otp", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.SendCheckoutOTP(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response["success"].(bool))
	})

	t.Run("Success - Different phone number formats", func(t *testing.T) {
		phoneNumbers := []string{
			"+8801712345678",
			"01712345678",
			"8801712345678",
		}

		for _, phone := range phoneNumbers {
			reqBody := models.CheckoutOTPRequest{
				PhoneNumber: phone,
			}
			body, _ := json.Marshal(reqBody)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/api/v1/checkout/send-otp", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.SendCheckoutOTP(c)

			assert.Equal(t, http.StatusOK, w.Code)
		}
	})
}

func TestCheckoutHandler_VerifyCheckoutOTP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupCheckoutHandlerTestDB(t)
	defer teardownCheckoutHandlerTestDB()

	mockOTPService := &MockCheckoutOTPService{}
	userService := services.NewUserService(nil)
	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	checkoutService := services.NewCheckoutService(mockOTPService, userService, jwtManager)
	handler := NewCheckoutHandler(mockOTPService, checkoutService)

	t.Run("Success - Verify OTP for new user", func(t *testing.T) {
		reqBody := models.CheckoutOTPVerifyRequest{
			PhoneNumber: "+8801712345678",
			OTP:         "1234",
			Name:        "John Doe",
			Email:       "john@example.com",
			Address:     "Building 123, Road 5, Gulshan, Dhaka",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/checkout/verify-otp", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.VerifyCheckoutOTP(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.NotEmpty(t, data["token"])
		assert.True(t, data["newUser"].(bool))
		assert.NotEmpty(t, data["userId"])
		assert.Equal(t, "John Doe", data["userName"])
		assert.Equal(t, "john@example.com", data["userEmail"])
		assert.Equal(t, "+8801712345678", data["userPhoneNumber"])
	})

	t.Run("Success - Verify OTP for existing user", func(t *testing.T) {
		existingUser := createTestUserForCheckoutHandler(t, db, "+8801798765432")

		reqBody := models.CheckoutOTPVerifyRequest{
			PhoneNumber: "+8801798765432",
			OTP:         "1234",
			Name:        "Updated Name",
			Email:       "updated@example.com",
			Address:     "New Building, New Road, New Area, New City",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/checkout/verify-otp", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.VerifyCheckoutOTP(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.False(t, data["newUser"].(bool))
		assert.Equal(t, existingUser.ID, data["userId"])
		assert.Equal(t, "Updated Name", data["userName"])
		assert.Equal(t, "updated@example.com", data["userEmail"])
	})

	t.Run("Success - Verify OTP with minimal data", func(t *testing.T) {
		reqBody := models.CheckoutOTPVerifyRequest{
			PhoneNumber: "+8801611111111",
			OTP:         "1234",
			Name:        "Minimal User",
			Address:     "Test Address",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/checkout/verify-otp", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.VerifyCheckoutOTP(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "Minimal User", data["userName"])
	})

	t.Run("Error - Missing required phone number", func(t *testing.T) {
		reqBody := models.CheckoutOTPVerifyRequest{
			OTP:  "1234",
			Name: "Test User",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/checkout/verify-otp", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.VerifyCheckoutOTP(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response["success"].(bool))
	})

	t.Run("Error - Missing required OTP", func(t *testing.T) {
		reqBody := models.CheckoutOTPVerifyRequest{
			PhoneNumber: "+8801712345678",
			Name:        "Test User",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/checkout/verify-otp", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.VerifyCheckoutOTP(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error - Missing required name", func(t *testing.T) {
		reqBody := models.CheckoutOTPVerifyRequest{
			PhoneNumber: "+8801712345678",
			OTP:         "1234",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/checkout/verify-otp", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.VerifyCheckoutOTP(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error - Invalid OTP", func(t *testing.T) {
		mockOTPService.shouldFailVerify = true
		defer func() { mockOTPService.shouldFailVerify = false }()

		reqBody := models.CheckoutOTPVerifyRequest{
			PhoneNumber: "+8801712345678",
			OTP:         "9999",
			Name:        "Test User",
			Email:       "test@example.com",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/checkout/verify-otp", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.VerifyCheckoutOTP(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response["success"].(bool))
		assert.Contains(t, response["error"], "Invalid or expired OTP")
	})

	t.Run("Error - Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/checkout/verify-otp", bytes.NewBuffer([]byte("invalid json")))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.VerifyCheckoutOTP(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCheckoutHandler_TokenValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupCheckoutHandlerTestDB(t)
	defer teardownCheckoutHandlerTestDB()

	mockOTPService := &MockCheckoutOTPService{}
	userService := services.NewUserService(nil)
	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	checkoutService := services.NewCheckoutService(mockOTPService, userService, jwtManager)
	handler := NewCheckoutHandler(mockOTPService, checkoutService)

	t.Run("Success - Returned token is valid", func(t *testing.T) {
		reqBody := models.CheckoutOTPVerifyRequest{
			PhoneNumber: "+8801712345678",
			OTP:         "1234",
			Name:        "Token Test",
			Email:       "token@example.com",
			Address:     "Test Address",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/checkout/verify-otp", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.VerifyCheckoutOTP(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].(map[string]interface{})
		token := data["token"].(string)

		// Validate the token
		claims, err := jwtManager.ValidateToken(token)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, "+8801712345678", claims.PhoneNumber)
		assert.Equal(t, "USER", claims.Role)
	})
}

func TestCheckoutHandler_AddressCreation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupCheckoutHandlerTestDB(t)
	defer teardownCheckoutHandlerTestDB()

	mockOTPService := &MockCheckoutOTPService{}
	userService := services.NewUserService(nil)
	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	checkoutService := services.NewCheckoutService(mockOTPService, userService, jwtManager)
	handler := NewCheckoutHandler(mockOTPService, checkoutService)

	t.Run("Success - Address is created with user", func(t *testing.T) {
		reqBody := models.CheckoutOTPVerifyRequest{
			PhoneNumber: "+8801712345678",
			OTP:         "1234",
			Name:        "Address Test",
			Email:       "address@example.com",
			Address:     "Test Building, Test Road, Test Area, Test City",
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/checkout/verify-otp", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.VerifyCheckoutOTP(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].(map[string]interface{})
		userID := data["userId"].(string)

		// Verify address was created
		var address models.Address
		err := db.Where("user_id = ?", userID).First(&address).Error
		assert.NoError(t, err)
		assert.Equal(t, "Address Test", address.FullName)
		assert.Equal(t, "Test Building, Test Road, Test Area, Test City", address.Address)
		assert.Equal(t, "+8801712345678", address.PhoneNumber)
	})
}

