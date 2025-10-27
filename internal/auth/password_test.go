package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	t.Run("Success - Hash password", func(t *testing.T) {
		password := "testPassword123"

		hash, err := HashPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
	})

	t.Run("Success - Different hashes for same password", func(t *testing.T) {
		password := "testPassword123"

		hash1, err1 := HashPassword(password)
		hash2, err2 := HashPassword(password)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, hash1, hash2) // bcrypt generates different salts
	})
}

func TestCheckPassword(t *testing.T) {
	t.Run("Success - Correct password", func(t *testing.T) {
		password := "testPassword123"
		hash, err := HashPassword(password)
		assert.NoError(t, err)

		result := CheckPassword(password, hash)

		assert.True(t, result)
	})

	t.Run("Error - Incorrect password", func(t *testing.T) {
		password := "testPassword123"
		wrongPassword := "wrongPassword456"
		hash, err := HashPassword(password)
		assert.NoError(t, err)

		result := CheckPassword(wrongPassword, hash)

		assert.False(t, result)
	})

	t.Run("Error - Invalid hash", func(t *testing.T) {
		password := "testPassword123"
		invalidHash := "not-a-valid-hash"

		result := CheckPassword(password, invalidHash)

		assert.False(t, result)
	})
}
