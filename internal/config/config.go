package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
	Redis    RedisConfig
	SMS      SMSConfig
	AWS      AWSConfig
	CORS     CORSConfig
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

// SMSConfig holds SMS service configuration
type SMSConfig struct {
	URL      string
	APIKey   string
	SenderID string
}

// AWSConfig holds AWS S3 configuration
type AWSConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	BucketName      string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	Environment    string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "flexfume_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			TimeZone: getEnv("DB_TIMEZONE", "UTC"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
		},
		Server: ServerConfig{
			Port: getEnv("API_PORT", "8080"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		SMS: SMSConfig{
			URL:      getEnv("SMS_URL", ""),
			APIKey:   getEnv("SMS_API_KEY", ""),
			SenderID: getEnv("SMS_SENDER_ID", ""),
		},
		AWS: AWSConfig{
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			Region:          getEnv("AWS_REGION", "ap-southeast-1"),
			BucketName:      getEnv("AWS_S3_BUCKET_NAME", ""),
		},
		CORS: CORSConfig{
			AllowedOrigins: getCORSAllowedOrigins(),
			Environment:    getEnv("ENVIRONMENT", "development"),
		},
	}

	// Validate required fields
	if config.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return config, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getCORSAllowedOrigins returns allowed origins based on environment
func getCORSAllowedOrigins() []string {
	env := getEnv("ENVIRONMENT", "development")
	corsOrigins := getEnv("CORS_ALLOWED_ORIGINS", "")
	
	// If specific origins are provided via environment variable, use them
	if corsOrigins != "" {
		origins := strings.Split(corsOrigins, ",")
		var cleanOrigins []string
		for _, origin := range origins {
			cleanOrigins = append(cleanOrigins, strings.TrimSpace(origin))
		}
		return cleanOrigins
	}
	
	// Default origins based on environment
	switch env {
	case "production":
		return []string{
			"https://flexfume-frontend.vercel.app",
			"https://www.flexfume.com",
		}
	case "staging":
		return []string{
			"https://flexfume-staging.vercel.app",
			"http://localhost:3000",
		}
	default: // development
		return []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://127.0.0.1:3000",
		}
	}
}
