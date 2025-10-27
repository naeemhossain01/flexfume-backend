package services

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
)

// SystemService handles system health checks and wake-up operations
type SystemService struct {
	redisService *RedisService
}

// NewSystemService creates a new SystemService
func NewSystemService(redisService *RedisService) *SystemService {
	return &SystemService{
		redisService: redisService,
	}
}

// WakeUpResponse represents the response structure for wake-up checks
type WakeUpResponse struct {
	Status     string                 `json:"status"`
	Message    string                 `json:"message"`
	Timestamp  string                 `json:"timestamp"`
	WakeUpTime int64                  `json:"wakeUpTime"`
	Checks     map[string]interface{} `json:"checks"`
}

// CheckResult represents a single check result
type CheckResult struct {
	Status         string      `json:"status"`
	Message        string      `json:"message"`
	ResponseTimeMs int64       `json:"responseTimeMs,omitempty"`
	Details        interface{} `json:"details,omitempty"`
}

// PerformWakeUpChecks performs all system health checks
func (s *SystemService) PerformWakeUpChecks() WakeUpResponse {
	checks := make(map[string]interface{})

	// Check database connection
	checks["database"] = s.checkDatabase()

	// Check Redis connection
	checks["redis"] = s.checkRedis()

	// Check application health
	checks["application"] = s.checkApplication()

	// Determine overall status
	allHealthy := true
	for _, check := range checks {
		if checkMap, ok := check.(CheckResult); ok {
			if checkMap.Status != "UP" && checkMap.Status != "SKIP" {
				allHealthy = false
				break
			}
		}
	}

	status := "UP"
	if !allHealthy {
		status = "PARTIAL"
	}

	return WakeUpResponse{
		Status:     status,
		Message:    "Wake-up checks completed",
		Timestamp:  time.Now().Format(time.RFC3339),
		WakeUpTime: time.Now().UnixMilli(),
		Checks:     checks,
	}
}

// checkDatabase verifies database connectivity
func (s *SystemService) checkDatabase() CheckResult {
	startTime := time.Now()
	result := CheckResult{}

	db := database.DB
	if db == nil {
		result.Status = "DOWN"
		result.Message = "Database connection not initialized"
		result.ResponseTimeMs = time.Since(startTime).Milliseconds()
		return result
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Execute a simple query to test the connection
	var testResult int
	err := db.WithContext(ctx).Raw("SELECT 1 as test_connection").Scan(&testResult).Error
	
	responseTime := time.Since(startTime).Milliseconds()

	if err != nil {
		result.Status = "DOWN"
		result.Message = fmt.Sprintf("Database connection failed: %v", err)
		result.ResponseTimeMs = responseTime
		return result
	}

	if testResult == 1 {
		result.Status = "UP"
		result.Message = "Database connection successful"
		result.ResponseTimeMs = responseTime
	} else {
		result.Status = "DOWN"
		result.Message = "Database query failed"
		result.ResponseTimeMs = responseTime
	}

	return result
}

// checkRedis verifies Redis connectivity
func (s *SystemService) checkRedis() CheckResult {
	startTime := time.Now()
	result := CheckResult{}

	if s.redisService == nil {
		result.Status = "SKIP"
		result.Message = "Redis not configured or unavailable"
		result.ResponseTimeMs = 0
		return result
	}

	// Test Redis connection with a simple ping
	testKey := "wake_up_test"
	testValue := fmt.Sprintf("ping_%d", time.Now().UnixMilli())

	// Set a test value with 10 second expiration
	err := s.redisService.Set(testKey, testValue, 10*time.Second)
	if err != nil {
		result.Status = "DOWN"
		result.Message = fmt.Sprintf("Redis connection failed: %v", err)
		result.ResponseTimeMs = time.Since(startTime).Milliseconds()
		return result
	}

	// Try to retrieve the value
	retrievedValue, err := s.redisService.GetString(testKey)
	if err != nil {
		result.Status = "DOWN"
		result.Message = fmt.Sprintf("Redis read failed: %v", err)
		result.ResponseTimeMs = time.Since(startTime).Milliseconds()
		return result
	}

	// Verify the value matches (retrieved value will be JSON encoded string)
	if retrievedValue == "\""+testValue+"\"" || retrievedValue == testValue {
		result.Status = "UP"
		result.Message = "Redis connection successful"
		
		// Clean up test key
		_ = s.redisService.Delete(testKey)
	} else {
		result.Status = "DOWN"
		result.Message = "Redis value mismatch"
	}

	result.ResponseTimeMs = time.Since(startTime).Milliseconds()
	return result
}

// checkApplication checks application health and memory usage
func (s *SystemService) checkApplication() CheckResult {
	result := CheckResult{}

	// Get runtime memory statistics
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Calculate memory usage
	maxMemory := memStats.Sys
	totalMemory := memStats.TotalAlloc
	usedMemory := memStats.Alloc
	freeMemory := maxMemory - usedMemory

	memoryUsagePercent := float64(usedMemory) * 100.0 / float64(maxMemory)

	memory := map[string]interface{}{
		"maxMemoryMB":        maxMemory / (1024 * 1024),
		"totalMemoryMB":      totalMemory / (1024 * 1024),
		"usedMemoryMB":       usedMemory / (1024 * 1024),
		"freeMemoryMB":       freeMemory / (1024 * 1024),
		"memoryUsagePercent": fmt.Sprintf("%.2f", memoryUsagePercent),
	}

	details := map[string]interface{}{
		"memory":              memory,
		"availableProcessors": runtime.NumCPU(),
		"numGoroutines":       runtime.NumGoroutine(),
	}

	result.Status = "UP"
	result.Message = "Application is running"
	result.Details = details

	return result
}
