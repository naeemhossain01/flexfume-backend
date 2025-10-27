package database

import (
	"fmt"
	"log"
	"time"

	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/config"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance
var DB *gorm.DB

// Connect establishes a connection to the PostgreSQL database
func Connect(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
		cfg.Database.TimeZone,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		PrepareStmt:            false, // Disable prepared statements to avoid "statement name already in use" error
		SkipDefaultTransaction: true,  // Improve performance by skipping default transactions
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)                  // Maximum number of idle connections
	sqlDB.SetMaxOpenConns(100)                 // Maximum number of open connections
	sqlDB.SetConnMaxLifetime(time.Hour * 1)    // Maximum lifetime of a connection (1 hour)

	log.Println("Database connection established successfully")

	// Run migrations
	if err := RunMigrations(cfg); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// AutoMigrate runs automatic migration for all models
func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Cart{},
		&models.Coupon{},
		&models.Discount{},
		&models.CouponUsage{},
		&models.Address{},
		&models.DeliveryCost{},
	)
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
