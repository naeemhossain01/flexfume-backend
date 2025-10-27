package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTManager_GenerateToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key", 24*time.Hour)

	t.Run("Success - Generate valid token", func(t *testing.T) {
		userID := "test-user-123"
		phoneNumber := "+1234567890"
		role := "USER"

		token, err := manager.GenerateToken(userID, phoneNumber, role)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestJWTManager_ValidateToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key", 24*time.Hour)

	t.Run("Success - Validate valid token", func(t *testing.T) {
		userID := "test-user-123"
		phoneNumber := "+1234567890"
		role := "USER"

		token, err := manager.GenerateToken(userID, phoneNumber, role)
		assert.NoError(t, err)

		claims, err := manager.ValidateToken(token)

		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, phoneNumber, claims.PhoneNumber)
		assert.Equal(t, role, claims.Role)
	})

	t.Run("Error - Invalid token", func(t *testing.T) {
		claims, err := manager.ValidateToken("invalid-token")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("Error - Empty token", func(t *testing.T) {
		claims, err := manager.ValidateToken("")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("Error - Expired token", func(t *testing.T) {
		expiredManager := NewJWTManager("test-secret-key", -1*time.Hour)
		userID := "test-user-123"
		phoneNumber := "+1234567890"
		role := "USER"

		token, err := expiredManager.GenerateToken(userID, phoneNumber, role)
		assert.NoError(t, err)

		claims, err := manager.ValidateToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
		// The error message contains "expired" but may be wrapped
		assert.Contains(t, err.Error(), "expired")
	})

	t.Run("Error - Token signed with different key", func(t *testing.T) {
		differentManager := NewJWTManager("different-secret-key", 24*time.Hour)
		userID := "test-user-123"
		phoneNumber := "+1234567890"
		role := "USER"

		token, err := differentManager.GenerateToken(userID, phoneNumber, role)
		assert.NoError(t, err)

		claims, err := manager.ValidateToken(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}
