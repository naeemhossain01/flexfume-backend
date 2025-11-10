package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/config"
)

// pathToFileURL converts a file path to a proper file:// URL
// This handles both Windows (C:\path) and Unix (/path) paths correctly
func pathToFileURL(path string) string {
	// Ensure we have an absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Printf("Warning: failed to get absolute path for %s: %v", path, err)
		absPath = path
	}
	
	// Convert to forward slashes
	absPath = filepath.ToSlash(absPath)
	
	// golang-migrate on Windows has issues with standard file:/// URLs
	// Use a simpler format that works on Windows: file://D:/path
	if runtime.GOOS == "windows" && len(absPath) >= 2 && absPath[1] == ':' {
		// Windows: file:// + D:/path = file://D:/path (NO leading slash before drive)
		urlStr := "file://" + absPath
		log.Printf("Generated file URL: %s", urlStr)
		return urlStr
	}
	
	// Unix: file:// + /path = file:///path (3 slashes total)
	urlStr := "file://" + absPath
	log.Printf("Generated file URL: %s", urlStr)
	return urlStr
}

// getMigrationsPath returns the absolute path to the migrations directory
func getMigrationsPath() (string, error) {
	// Try to get the path relative to the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	migrationsPath := filepath.Join(cwd, "migrations")
	
	// Check if migrations directory exists
	if _, err := os.Stat(migrationsPath); err == nil {
		log.Printf("Found migrations at: %s", migrationsPath)
		return pathToFileURL(migrationsPath), nil
	}

	// If not found in cwd, try relative to the executable
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	exeDir := filepath.Dir(exePath)
	migrationsPath = filepath.Join(exeDir, "migrations")
	
	if _, err := os.Stat(migrationsPath); err == nil {
		log.Printf("Found migrations at: %s", migrationsPath)
		return pathToFileURL(migrationsPath), nil
	}

	// If still not found, try parent directories (useful for development)
	migrationsPath = filepath.Join(exeDir, "..", "migrations")
	if _, err := os.Stat(migrationsPath); err == nil {
		absPath, _ := filepath.Abs(migrationsPath)
		log.Printf("Found migrations at: %s", absPath)
		return pathToFileURL(absPath), nil
	}

	// If still not found, try /app/migrations (Render native runtime)
	migrationsPath = "/app/migrations"
	if _, err := os.Stat(migrationsPath); err == nil {
		log.Printf("Found migrations at: %s", migrationsPath)
		return pathToFileURL(migrationsPath), nil
	}

	// Log debug information
	log.Printf("Current working directory: %s", cwd)
	log.Printf("Executable directory: %s", exeDir)
	log.Printf("Checked paths: %s, %s, %s, %s", 
		filepath.Join(cwd, "migrations"),
		filepath.Join(exeDir, "migrations"),
		filepath.Join(exeDir, "..", "migrations"),
		"/app/migrations")

	return "", fmt.Errorf("migrations directory not found in any expected location")
}

// RunMigrations runs database migrations
func RunMigrations(cfg *config.Config) error {
	// Create database connection string
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

	// Open database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Get migrations path
	migrationsPath, err := getMigrationsPath()
	if err != nil {
		return fmt.Errorf("failed to get migrations path: %w", err)
	}
	log.Printf("Using migrations URL: %s", migrationsPath)

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Run migrations
	log.Println("Running database migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply")
	} else {
		log.Println("Migrations applied successfully")
	}

	return nil
}

// RollbackMigration rolls back the last migration
func RollbackMigration(cfg *config.Config) error {
	// Create database connection string
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

	// Open database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Get migrations path
	migrationsPath, err := getMigrationsPath()
	if err != nil {
		return fmt.Errorf("failed to get migrations path: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Rollback one step
	log.Println("Rolling back last migration...")
	if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Println("Migration rolled back successfully")
	return nil
}

// MigrateDown rolls back all migrations
func MigrateDown(cfg *config.Config) error {
	// Create database connection string
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

	// Open database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Get migrations path
	migrationsPath, err := getMigrationsPath()
	if err != nil {
		return fmt.Errorf("failed to get migrations path: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Rollback all migrations
	log.Println("Rolling back all migrations...")
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback all migrations: %w", err)
	}

	log.Println("All migrations rolled back successfully")
	return nil
}
