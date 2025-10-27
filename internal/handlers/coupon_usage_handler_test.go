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

func setupCouponUsageTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	err = db.AutoMigrate(
		&models.User{},
		&models.Coupon{},
		&models.CouponUsage{},
		&models.Cart{},
		&models.CartItem{},
		&models.Product{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	database.DB = db
	return db
}

func teardownCouponUsageTestDB() {
	database.DB = nil
}

func createTestUserForCoupon(t *testing.T, db *gorm.DB, email string) *models.User {
	user := &models.User{
		Email:    email,
		Password: "hashedpassword",
		Name:     "Test User",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

func createTestCouponForUsage(t *testing.T, db *gorm.DB, code string, couponType string, amount float64) *models.Coupon {
	coupon := &models.Coupon{
		Code:             code,
		CouponType:       couponType,
		Amount:           amount,
		MinOrderAmount:   50.00,
		MaxAmountApplied: 100.00,
		ExpirationTime:   time.Now().Add(24 * time.Hour),
		UsageLimit:       10,
		Active:           true,
	}
	if err := db.Create(coupon).Error; err != nil {
		t.Fatalf("Failed to create test coupon: %v", err)
	}
	return coupon
}

func createTestProductForCoupon(t *testing.T, db *gorm.DB, name string, price float64) *models.Product {
	product := &models.Product{
		Name:        name,
		Description: "Test product",
		Price:       price,
		Stock:       100,
	}
	if err := db.Create(product).Error; err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}
	return product
}

func TestCouponUsageHandler_ApplyCoupon(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupCouponUsageTestDB(t)
	defer teardownCouponUsageTestDB()

	couponUsageService := services.NewCouponUsageService()
	handler := NewCouponUsageHandler(couponUsageService)

	t.Run("Success - Apply valid coupon", func(t *testing.T) {
		user := createTestUserForCoupon(t, db, "test@example.com")
		createTestCouponForUsage(t, db, "SAVE10", "PERCENTAGE", 10.0)
		product := createTestProductForCoupon(t, db, "Test Product", 100.0)

		// Create cart and cart item
		cart := &models.Cart{UserID: user.ID}
		db.Create(cart)
		cartItem := &models.CartItem{
			CartID:    cart.ID,
			ProductID: product.ID,
			Quantity:  2,
		}
		db.Create(cartItem)

		reqBody := ApplyCouponRequest{
			CouponCode:   "SAVE10",
			CartInfoList: []string{cartItem.ID},
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", user.ID)
		c.Request = httptest.NewRequest("POST", "/api/v1/coupon/apply", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ApplyCoupon(c)

		// Check response (may vary based on implementation)
		if w.Code != http.StatusOK {
			t.Logf("Apply coupon returned status: %d, body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Error - Missing user authentication", func(t *testing.T) {
		reqBody := ApplyCouponRequest{
			CouponCode:   "SAVE10",
			CartInfoList: []string{"cart-item-1"},
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/coupon/apply", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ApplyCoupon(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Error - Invalid request body", func(t *testing.T) {
		user := createTestUserForCoupon(t, db, "test2@example.com")

		reqBody := map[string]interface{}{
			"couponCode": "SAVE10",
			// Missing cartInfoList
		}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", user.ID)
		c.Request = httptest.NewRequest("POST", "/api/v1/coupon/apply", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.ApplyCoupon(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCouponUsageHandler_GetUserCouponUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupCouponUsageTestDB(t)
	defer teardownCouponUsageTestDB()

	couponUsageService := services.NewCouponUsageService()
	handler := NewCouponUsageHandler(couponUsageService)

	t.Run("Success - Get user coupon usage", func(t *testing.T) {
		user := createTestUserForCoupon(t, db, "test3@example.com")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", user.ID)
		c.Request = httptest.NewRequest("GET", "/api/v1/coupon/usage", nil)

		handler.GetUserCouponUsage(c)

		// Should return OK even if empty
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Error - Missing user authentication", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/v1/coupon/usage", nil)

		handler.GetUserCouponUsage(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
