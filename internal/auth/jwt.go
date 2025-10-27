package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidToken is returned when the token is invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken is returned when the token has expired
	ErrExpiredToken = errors.New("token has expired")
)

// JWTClaims represents the claims in the JWT
// Matches Spring Boot JwtUtils structure:
// - subject: phoneNumber (username)
// - authorities: role (USER/ADMIN)
type JWTClaims struct {
	Authorities string `json:"authorities"` // Role claim (matches Spring Boot)
	jwt.RegisteredClaims
}

// JWTManager handles JWT operations
type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

// GenerateToken generates a new JWT token for a user
// Matches Spring Boot JwtUtils.generateToken(String username, String authorities)
// - subject: phoneNumber (username)
// - authorities: role
func (manager *JWTManager) GenerateToken(userID, phoneNumber, role string) (string, error) {
	claims := JWTClaims{
		Authorities: role, // Role as "authorities" claim (matches Spring Boot)
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   phoneNumber,                                                 // phoneNumber as subject (matches Spring Boot)
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(manager.tokenDuration)), // Expiration time
			IssuedAt:  jwt.NewNumericDate(time.Now()),                            // Issued at
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

// ValidateToken validates a JWT token and returns the claims
func (manager *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(manager.secretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Check if token is expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}
