package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/auth"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/config"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/handlers"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/middleware"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Redis service
	redisService, err := services.NewRedisService(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v (OTP features will be disabled)", err)
		// Continue without Redis for development
	} else {
		log.Println("Redis connection established successfully")
	}

	// Initialize SMS service
	smsService := services.NewSMSService(cfg.SMS.URL, cfg.SMS.APIKey, cfg.SMS.SenderID)

	// Initialize OTP service
	var otpService *services.OTPService
	if redisService != nil {
		otpService = services.NewOTPService(redisService, smsService)
		log.Println("OTP service initialized successfully")
	}

	// Initialize JWT manager (1 hour token expiration - matches Spring Boot)
	// Spring Boot: EXPIRATION_TIME_IN_MILLISECONDS = 1000L * 60L * 60L (1 hour)
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, 1*time.Hour)

	// Initialize S3 service (optional - will gracefully degrade if not configured)
	s3Service, err := services.NewS3Service(cfg)
	if err != nil {
		log.Printf("Warning: S3 service initialization failed: %v (Image uploads will use fallback)", err)
		s3Service = nil // Set to nil to use fallback in handlers
	} else {
		log.Println("S3 service initialized successfully")
	}

	// Initialize services
	userService := services.NewUserService(otpService)
	categoryService := services.NewCategoryService()
	productService := services.NewProductService(categoryService)
	couponService := services.NewCouponService()
	discountService := services.NewDiscountService(productService)
	deliveryCostService := services.NewDeliveryCostService()
	orderService := services.NewOrderService(userService, productService, discountService, deliveryCostService, couponService)
	couponUsageService := services.NewCouponUsageService(couponService, discountService)
	addressService := services.NewAddressService(userService)
	checkoutService := services.NewCheckoutService(otpService, userService, jwtManager)
	systemService := services.NewSystemService(redisService)
	affiliateService := services.NewAffiliateService()
	contactService := services.NewContactService()

	// Initialize handlers
	var authHandler *handlers.AuthHandler
	var userHandler *handlers.UserHandler
	var checkoutHandler *handlers.CheckoutHandler
	
	if otpService != nil {
		authHandler = handlers.NewAuthHandler(jwtManager, otpService, userService)
		userHandler = handlers.NewUserHandler(userService, otpService)
		checkoutHandler = handlers.NewCheckoutHandler(otpService, checkoutService)
	} else {
		log.Println("Warning: Auth, User, and Checkout handlers not initialized due to missing OTP service")
	}
	
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	productHandler := handlers.NewProductHandler(productService, s3Service)
	orderHandler := handlers.NewOrderHandler(orderService)
	couponHandler := handlers.NewCouponHandler(couponService)
	discountHandler := handlers.NewDiscountHandler(discountService)
	couponUsageHandler := handlers.NewCouponUsageHandler(couponUsageService)
	addressHandler := handlers.NewAddressHandler(addressService)
	deliveryCostHandler := handlers.NewDeliveryCostHandler(deliveryCostService)
	healthHandler := handlers.NewHealthHandler()
	systemHandler := handlers.NewSystemHandler(systemService)
	affiliateHandler := handlers.NewAffiliateHandler(affiliateService)
	contactHandler := handlers.NewContactHandler(contactService)

	// Setup Gin router
	router := gin.Default()

	// Apply CORS middleware globally
	router.Use(middleware.CORSMiddleware(cfg))

	// Health check endpoint
	router.GET("/health", healthHandler.Health)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public) - only if authHandler is available
		if authHandler != nil {
			authRoutes := v1.Group("/auth")
			{
				authRoutes.POST("/login", authHandler.Login)
			}
		}

		// User routes - only if userHandler is available
		if userHandler != nil {
			userRoutes := v1.Group("/user")
			{
				// Public routes
				userRoutes.GET("/reset-password-request", userHandler.ResetPasswordRequest)
				userRoutes.POST("/reset-password", userHandler.ResetPassword)

				// Protected routes (require authentication)
				userRoutes.GET("/profile", middleware.AuthMiddleware(jwtManager), userHandler.GetProfile)
				userRoutes.POST("/change-password", middleware.AuthMiddleware(jwtManager), userHandler.ChangePassword)
				userRoutes.PUT("/:id", middleware.AuthMiddleware(jwtManager), userHandler.UpdateUser)

				// Admin only routes
				userRoutes.GET("/:id", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), userHandler.GetUserByID)
				userRoutes.GET("", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), userHandler.GetAllUsers)
			}
		}

		// Category routes
		categoryRoutes := v1.Group("/category")
		{
			// Public routes
			categoryRoutes.GET("/all", categoryHandler.GetAllCategories)
			categoryRoutes.GET("/get/:id", categoryHandler.GetCategoryByID)

			// Admin only routes
			categoryRoutes.POST("/add", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), categoryHandler.CreateCategory)
			categoryRoutes.PUT("/update/:id", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), categoryHandler.UpdateCategory)
			categoryRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), categoryHandler.DeleteCategory)
		}

		// Product routes
		productRoutes := v1.Group("/product")
		{
			// Public routes
			productRoutes.GET("/all", productHandler.GetAllProducts)
			productRoutes.GET("/get/:id", productHandler.GetProductByID)
			productRoutes.GET("/get-by-category/:id", productHandler.GetProductsByCategory)
			productRoutes.GET("/search", productHandler.SearchProducts)

			// Admin only routes
			productRoutes.POST("/add", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), productHandler.CreateProduct)
			productRoutes.PUT("/update/:id", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), productHandler.UpdateProduct)
			productRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), productHandler.DeleteProduct)
			productRoutes.PUT("/upload-image/:id", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), productHandler.UploadProductImage)
		}

		// Order routes
		orderRoutes := v1.Group("/order")
		{
			// Authenticated routes
			orderRoutes.POST("", middleware.AuthMiddleware(jwtManager), orderHandler.PlaceOrder)
			orderRoutes.GET("/history", middleware.AuthMiddleware(jwtManager), orderHandler.GetOrderHistory)

			// Admin only routes
			orderRoutes.PUT("/:orderId", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), orderHandler.UpdateOrderStatus)
			orderRoutes.GET("/filter", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), orderHandler.FilterOrders)
		}

		// Coupon routes
		couponRoutes := v1.Group("/coupon")
		{
			// Admin only routes
			couponRoutes.GET("/:id", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), couponHandler.GetCouponByID)
			couponRoutes.GET("/all", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), couponHandler.GetAllCoupons)
			couponRoutes.POST("", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), couponHandler.CreateCoupon)
			couponRoutes.PUT("/:id", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), couponHandler.UpdateCoupon)
			couponRoutes.DELETE("/:id", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), couponHandler.DeleteCoupon)
		}

		// Coupon Usage routes
		couponUsageRoutes := v1.Group("/coupon-usage")
		{
			// Authenticated routes
			couponUsageRoutes.POST("", middleware.AuthMiddleware(jwtManager), couponUsageHandler.ApplyCoupon)
			couponUsageRoutes.DELETE("", middleware.AuthMiddleware(jwtManager), couponUsageHandler.RemoveCoupon)
		}

		// Discount routes
		discountRoutes := v1.Group("/discount")
		{
			// Admin only routes
			discountRoutes.GET("", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), discountHandler.GetAllDiscounts)
			discountRoutes.POST("", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), discountHandler.AddDiscounts)
			discountRoutes.PUT("", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), discountHandler.UpdateDiscounts)
			discountRoutes.DELETE("/:id", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), discountHandler.DeleteDiscount)
		}

		// Address routes
		addressRoutes := v1.Group("/address")
		{
			// Authenticated routes
			addressRoutes.POST("", middleware.AuthMiddleware(jwtManager), addressHandler.AddAddress)
			addressRoutes.PUT("/:id", middleware.AuthMiddleware(jwtManager), addressHandler.UpdateAddress)
			addressRoutes.GET("", middleware.AuthMiddleware(jwtManager), addressHandler.GetAddressByUser)
		}

		// Delivery Cost routes
		deliveryCostRoutes := v1.Group("/delivery-cost")
		{
			// Public routes
			deliveryCostRoutes.GET("/all", deliveryCostHandler.GetAllDeliveryCosts)
			deliveryCostRoutes.GET("/location", deliveryCostHandler.GetDeliveryCostByLocation)

			// Admin only routes
			deliveryCostRoutes.POST("", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), deliveryCostHandler.AddCost)
			deliveryCostRoutes.PUT("/:id", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), deliveryCostHandler.UpdateCost)
			deliveryCostRoutes.GET("/:id", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), deliveryCostHandler.GetDeliveryCostByID)
			deliveryCostRoutes.DELETE("/:id", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), deliveryCostHandler.DeleteDeliveryCost)
		}

		// Checkout routes - only if checkoutHandler is available
		if checkoutHandler != nil {
			checkoutRoutes := v1.Group("/checkout")
			{
				// Public routes
				checkoutRoutes.POST("/send-otp", checkoutHandler.SendCheckoutOTP)
				checkoutRoutes.POST("/verify-otp", checkoutHandler.VerifyCheckoutOTP)
			}
		}

		// System routes
		systemRoutes := v1.Group("/system")
		{
			// Public routes
			systemRoutes.GET("/wake-up", systemHandler.WakeUp)
			systemRoutes.GET("/ping", systemHandler.Ping)
		}

		// Affiliate routes
		affiliateRoutes := v1.Group("/affiliate")
		{
			// Public routes with rate limiting
			affiliateRoutes.POST("/submit", middleware.AffiliateRateLimitMiddleware(), affiliateHandler.SubmitAffiliateApplication)

			// Admin only routes
			affiliateRoutes.GET("/applications", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), affiliateHandler.GetAffiliateApplications)
			affiliateRoutes.GET("/applications/:id", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), affiliateHandler.GetAffiliateApplicationByID)
			affiliateRoutes.PUT("/applications/:id/status", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), affiliateHandler.UpdateAffiliateStatus)
			affiliateRoutes.DELETE("/applications/:id", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), affiliateHandler.DeleteAffiliateApplication)
		}

		// Contact routes
		contactRoutes := v1.Group("/contact")
		{
			// Public routes
			contactRoutes.POST("", contactHandler.SubmitContactForm)

			// Admin only routes
			contactRoutes.GET("", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), contactHandler.GetAllContactSubmissions)
			contactRoutes.GET("/:id", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), contactHandler.GetContactSubmissionByID)
			contactRoutes.DELETE("/:id", middleware.AuthMiddleware(jwtManager), middleware.RequireAdmin(), contactHandler.DeleteContactSubmission)
		}
	}

	// Start server
	serverAddr := ":" + cfg.Server.Port
	log.Printf("Starting server on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

