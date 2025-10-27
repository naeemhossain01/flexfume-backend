package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AffiliateRateLimitMiddleware creates a rate limiting middleware specifically for affiliate submissions
// This prevents spam submissions while allowing legitimate applications
func AffiliateRateLimitMiddleware() gin.HandlerFunc {
	// Allow 5 submissions per IP per hour to prevent spam
	rateLimiter := NewRateLimiter(5, time.Hour)
	
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		if !rateLimiter.IsAllowed(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   true,
				"message": "Too many affiliate applications submitted. Please try again later.",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}
