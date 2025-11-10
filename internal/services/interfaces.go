package services

import (
	"time"

	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
)

// RedisServiceInterface defines the interface for Redis operations
type RedisServiceInterface interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string, dest interface{}) error
	GetString(key string) (string, error)
	Delete(key string) error
	Exists(key string) (bool, error)
}

// SMSServiceInterface defines the interface for SMS operations
type SMSServiceInterface interface {
	SendSMS(phoneNumber, message string) error
}

// UserServiceInterface defines the interface for user operations
type UserServiceInterface interface {
	RegisterUser(user *models.User) (*models.User, error)
	GetUserByID(userID string) (*models.User, error)
	GetUserByPhoneNumber(phoneNumber string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(userID string, updateData map[string]interface{}) (*models.User, error)
	ChangePassword(userID, oldPassword, newPassword, confirmPassword string) error
	ResetPassword(phoneNumber, newPassword, confirmPassword string) error
	ValidateAlreadyHaveAccount(phoneNumber string) error
}

// CategoryServiceInterface defines the interface for category operations
type CategoryServiceInterface interface {
	CreateCategory(category *models.Category) (*models.Category, error)
	UpdateCategory(categoryID string, category *models.Category) (*models.Category, error)
	GetCategoryByID(categoryID string) (*models.Category, error)
	GetAllCategories() ([]models.Category, error)
	DeleteCategory(categoryID string) error
}

// ProductServiceInterface defines the interface for product operations
type ProductServiceInterface interface {
	CreateProduct(product *models.Product) (*models.Product, error)
	UpdateProduct(productID string, product *models.Product) (*models.Product, error)
	GetProductByID(productID string) (*models.Product, error)
	GetAllProducts() ([]models.Product, error)
	DeleteProduct(productID string) error
	GetProductsByCategory(categoryID string) ([]models.Product, error)
	SearchProducts(searchValue string) ([]models.Product, error)
	UploadProductImage(productID string, imageURL string) (*models.Product, error)
}

// CouponServiceInterface defines the interface for coupon operations
type CouponServiceInterface interface {
	CreateCoupon(coupon *models.Coupon) (*models.Coupon, error)
	UpdateCoupon(couponID string, coupon *models.Coupon) (*models.Coupon, error)
	GetCouponByID(couponID string) (*models.Coupon, error)
	GetCouponByCode(code string) (*models.Coupon, error)
	GetAllCoupons() ([]models.Coupon, error)
	DeleteCoupon(couponID string) error
	EnrichCouponWithStatistics(coupon *models.Coupon) models.CouponResponse
}

// DiscountServiceInterface defines the interface for discount operations
type DiscountServiceInterface interface {
	AddDiscounts(discounts []models.Discount) ([]models.Discount, error)
	UpdateDiscounts(discounts []models.Discount) ([]models.Discount, error)
	GetAllDiscounts() ([]models.Discount, error)
	GetDiscountByProductID(productID string) (*models.Discount, error)
	DeleteDiscount(discountID string) error
}

// AddressServiceInterface defines the interface for address operations
type AddressServiceInterface interface {
	AddAddress(address *models.Address) (*models.Address, error)
	UpdateAddress(addressID string, address *models.Address, userID string) (*models.Address, error)
	GetAddressByUserID(userID string) (*models.Address, error)
}

// DeliveryCostServiceInterface defines the interface for delivery cost operations
type DeliveryCostServiceInterface interface {
	AddCost(deliveryCost *models.DeliveryCost) (*models.DeliveryCost, error)
	UpdateCost(id string, deliveryCost *models.DeliveryCost) (*models.DeliveryCost, error)
	GetDeliveryCostByID(id string) (*models.DeliveryCost, error)
	GetAllDeliveryCosts() ([]models.DeliveryCost, error)
	GetDeliveryCostByLocation(location string) ([]models.DeliveryCost, error)
	DeleteDeliveryCost(id string) error
}

// CheckoutServiceInterface defines the interface for checkout operations
type CheckoutServiceInterface interface {
	VerifyOTPAndHandleUser(request *models.CheckoutOTPVerifyRequest) (*models.CheckoutOTPResponse, error)
}
