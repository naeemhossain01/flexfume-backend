package services

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisService handles Redis operations
type RedisService struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisService creates a new Redis service
func NewRedisService(host, port, password string) (*RedisService, error) {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	}

	// Enable TLS for Upstash and other cloud Redis providers
	if host != "localhost" && host != "127.0.0.1" {
		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	client := redis.NewClient(options)

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisService{
		client: client,
		ctx:    ctx,
	}, nil
}

// Set stores a value in Redis with expiration
// Matches Spring Boot's RedisService which JSON-serializes all values
func (r *RedisService) Set(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(r.ctx, key, jsonData, expiration).Err()
}

// Get retrieves a value from Redis
func (r *RedisService) Get(key string, dest interface{}) error {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// GetString retrieves a string value from Redis
// Matches Spring Boot's RedisService which JSON-deserializes values
func (r *RedisService) GetString(key string) (string, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return "", err
	}

	// Unmarshal the JSON string to get the actual string value
	// Spring Boot stores "1574" as JSON string, so we need to unmarshal it
	var result string
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return "", err
	}

	return result, nil
}

// Delete removes a key from Redis
func (r *RedisService) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Exists checks if a key exists in Redis
func (r *RedisService) Exists(key string) (bool, error) {
	count, err := r.client.Exists(r.ctx, key).Result()
	return count > 0, err
}

// Close closes the Redis connection
func (r *RedisService) Close() error {
	return r.client.Close()
}

// Ensure RedisService implements RedisServiceInterface
var _ RedisServiceInterface = (*RedisService)(nil)
