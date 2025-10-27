package middleware

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/config"
)

// CORSMiddleware handles CORS preflight requests and adds necessary headers
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		method := c.Request.Method
		
		// Set CORS headers based on environment and origin
		if cfg.CORS.Environment == "development" {
			// In development, allow all origins for easier testing
			c.Header("Access-Control-Allow-Origin", "*")
		} else if isOriginAllowed(origin, cfg.CORS.AllowedOrigins) {
			// In production/staging, only allow specific origins
			c.Header("Access-Control-Allow-Origin", origin)
		} else if len(cfg.CORS.AllowedOrigins) > 0 {
			// Fallback to first allowed origin if origin is not in the list
			c.Header("Access-Control-Allow-Origin", cfg.CORS.AllowedOrigins[0])
		} else {
			// No allowed origins configured, deny all
			c.Header("Access-Control-Allow-Origin", "null")
		}
		
		// Set CORS headers
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With, X-CSRF-Token")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400") // 24 hours
		
		// Handle preflight OPTIONS requests
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// isOriginAllowed checks if the given origin is in the allowed list
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}
	
	for _, allowed := range allowedOrigins {
		if strings.EqualFold(origin, allowed) {
			return true
		}
	}
	return false
}

// CORSMiddlewareWithRateLimit provides CORS with rate limiting for OPTIONS requests
func CORSMiddlewareWithRateLimit(cfg *config.Config) gin.HandlerFunc {
	rateLimiter := NewRateLimiter(100, time.Minute) // 100 requests per minute
	
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		method := c.Request.Method
		
		// Rate limit OPTIONS requests
		if method == "OPTIONS" {
			ip := c.ClientIP()
			if !rateLimiter.IsAllowed(ip) {
				log.Printf("CORS rate limit exceeded for IP: %s", ip)
				c.AbortWithStatus(http.StatusTooManyRequests)
				return
			}
			log.Printf("CORS preflight: %s %s from %s", method, c.Request.URL.Path, origin)
		}
		
		// Check if origin is allowed
		if isOriginAllowed(origin, cfg.CORS.AllowedOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
		} else if cfg.CORS.Environment == "development" {
			// In development, allow all origins for easier testing
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			// In production/staging, only allow specific origins
			c.Header("Access-Control-Allow-Origin", cfg.CORS.AllowedOrigins[0])
		}
		
		// Set CORS headers
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With, X-CSRF-Token")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400") // 24 hours
		
		// Handle preflight OPTIONS requests
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// RateLimiter provides rate limiting functionality
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// IsAllowed checks if the given IP is allowed to make a request
func (rl *RateLimiter) IsAllowed(ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	now := time.Now()
	cutoff := now.Add(-rl.window)
	
	// Clean old requests
	if times, exists := rl.requests[ip]; exists {
		var validTimes []time.Time
		for _, t := range times {
			if t.After(cutoff) {
				validTimes = append(validTimes, t)
			}
		}
		rl.requests[ip] = validTimes
	}
	
	// Check limit
	if len(rl.requests[ip]) >= rl.limit {
		return false
	}
	
	// Add current request
	rl.requests[ip] = append(rl.requests[ip], now)
	return true
}
