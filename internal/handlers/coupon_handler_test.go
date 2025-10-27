package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func setupCouponTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	err = db.AutoMigrate(&models.Coupon{}, &models.CouponUsage{}, &models.User{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	database.DB = db
	return db
}

func teardownCouponTestDB() {
	database.DB = nil
}

func createTestCouponInDB(t *testing.T, db *gorm.DB, code string, couponType string, amount float64) *models.Coupon {
	coupon := &models.Coupon{
		Code:             code,
		CouponType:       couponType,
		Amount:           amount,
		MinOrderAmount:   100.00,
		MaxAmountApplied: 50.00,
		ExpirationTime:   time.Now().Add(24 * time.Hour),
		UsageLimit:       5,
		Active:           true,
	}
	if err := db.Create(coupon).Error; err != nil {
		t.Fatalf("Failed to create test coupon: %v", err)
	}
	return coupon
}

func TestCouponHandler_CreateCoupon(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupCouponTestDB(t)
	defer teardownCouponTestDB()

	couponService := services.NewCouponService()
	handler := NewCouponHandler(couponService)

	t.Run("Success - Create valid percentage coupon", func(t *testing.T) {
		reqBody := CreateCouponRequest{
			Code:          "SAVE20",
			DiscountType:  "PERCENTAGE",
			DiscountValue: 20.0,
			MinOrderValue: 100.0,
			MaxDiscount:   50.0,
			ValidUntil:    time.Now().Add(30 * 24 * time.Hour),
			UsageLimit:    10,
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/coupon", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateCoupon(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))
		assert.NotNil(t, response["data"])

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "SAVE20", data["code"])
		assert.Equal(t, "PERCENTAGE", data["couponType"])
		assert.Equal(t, 20.0, data["amount"])
	})

	t.Run("Success - Create valid fixed coupon", func(t *testing.T) {
		reqBody := CreateCouponRequest{
			Code:          "FLAT50",
			DiscountType:  "FIXED",
			DiscountValue: 50.0,
			MinOrderValue: 200.0,
			MaxDiscount:   50.0,
			ValidUntil:    time.Now().Add(30 * 24 * time.Hour),
			UsageLimit:    5,
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/coupon", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateCoupon(c)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "FLAT50", data["code"])
		assert.Equal(t, "FIXED", data["couponType"])
	})

	t.Run("Error - Duplicate coupon code", func(t *testing.T) {
		// Create first coupon
		createTestCouponInDB(t, db, "DUPLICATE", "PERCENTAGE", 10.0)

		// Try to create duplicate
		reqBody := CreateCouponRequest{
			Code:          "DUPLICATE",
			DiscountType:  "PERCENTAGE",
			DiscountValue: 15.0,
			MinOrderValue: 100.0,
			MaxDiscount:   50.0,
			ValidUntil:    time.Now().Add(30 * 24 * time.Hour),
			UsageLimit:    5,
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/coupon", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateCoupon(c)

		assert.Equal(t, http.StatusConflict, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.False(t, response["success"].(bool))
	})

	t.Run("Error - Invalid discount type", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"code":          "INVALID",
			"discountType":  "INVALID_TYPE",
			"discountValue": 20.0,
			"minOrderValue": 100.0,
			"maxDiscount":   50.0,
			"validUntil":    time.Now().Add(30 * 24 * time.Hour),
			"usageLimit":    5,
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/coupon", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateCoupon(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCouponHandler_UpdateCoupon(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupCouponTestDB(t)
	defer teardownCouponTestDB()

	couponService := services.NewCouponService()
	handler := NewCouponHandler(couponService)

	t.Run("Success - Update coupon", func(t *testing.T) {
		coupon := createTestCouponInDB(t, db, "UPDATE_TEST", "PERCENTAGE", 10.0)

		active := false
		reqBody := UpdateCouponRequest{
			DiscountValue: 25.0,
			Active:        &active,
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: coupon.ID}}
		c.Request = httptest.NewRequest("PUT", "/api/v1/coupon/"+coupon.ID, bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateCoupon(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, 25.0, data["amount"])
		assert.False(t, data["active"].(bool))
	})

	t.Run("Error - Coupon not found", func(t *testing.T) {
		reqBody := UpdateCouponRequest{
			DiscountValue: 25.0,
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "non-existent-id"}}
		c.Request = httptest.NewRequest("PUT", "/api/v1/coupon/non-existent-id", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateCoupon(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCouponHandler_DeleteCoupon(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupCouponTestDB(t)
	defer teardownCouponTestDB()

	couponService := services.NewCouponService()
	handler := NewCouponHandler(couponService)

	t.Run("Success - Delete coupon without usage", func(t *testing.T) {
		coupon := createTestCouponInDB(t, db, "DELETE_TEST", "PERCENTAGE", 10.0)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: coupon.ID}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/coupon/"+coupon.ID, nil)

		handler.DeleteCoupon(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))
		assert.Equal(t, "Coupon deleted successfully", response["data"])
	})

	t.Run("Error - Coupon not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "non-existent-id"}}
		c.Request = httptest.NewRequest("DELETE", "/api/v1/coupon/non-existent-id", nil)

		handler.DeleteCoupon(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCouponHandler_GetCouponByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupCouponTestDB(t)
	defer teardownCouponTestDB()

	couponService := services.NewCouponService()
	handler := NewCouponHandler(couponService)

	t.Run("Success - Get coupon by ID", func(t *testing.T) {
		coupon := createTestCouponInDB(t, db, "GET_TEST", "PERCENTAGE", 15.0)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: coupon.ID}}
		c.Request = httptest.NewRequest("GET", "/api/v1/coupon/"+coupon.ID, nil)

		handler.GetCouponByID(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "GET_TEST", data["code"])
		assert.Equal(t, 15.0, data["amount"])
	})

	t.Run("Error - Coupon not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "non-existent-id"}}
		c.Request = httptest.NewRequest("GET", "/api/v1/coupon/non-existent-id", nil)

		handler.GetCouponByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCouponHandler_GetAllCoupons(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupCouponTestDB(t)
	defer teardownCouponTestDB()

	couponService := services.NewCouponService()
	handler := NewCouponHandler(couponService)

	t.Run("Success - Get all coupons", func(t *testing.T) {
		createTestCouponInDB(t, db, "COUPON1", "PERCENTAGE", 10.0)
		createTestCouponInDB(t, db, "COUPON2", "FIXED", 50.0)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/coupon/all", nil)

		handler.GetAllCoupons(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.True(t, response["success"].(bool))

		data := response["data"].([]interface{})
		assert.GreaterOrEqual(t, len(data), 2)
	})
}
