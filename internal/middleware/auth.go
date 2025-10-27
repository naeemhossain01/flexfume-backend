package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/auth"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
)

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header format must be Bearer {token}",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate the token
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			if err == auth.ErrExpiredToken {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Token has expired",
				})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid token",
				})
			}
			c.Abort()
			return
		}

		// Set user information in the context
		// Note: claims.Subject contains phoneNumber (matches Spring Boot)
		c.Set("phone_number", claims.Subject)
		c.Set("user_role", claims.Authorities)

		// Fetch user ID from database using phone number
		// This is needed for handlers that require user_id
		var user models.User
		if err := database.GetDB().Where("phone_number = ?", claims.Subject).First(&user).Error; err == nil {
			c.Set("user_id", user.ID)
		}

		// Continue to the next handler
		c.Next()
	}
}
