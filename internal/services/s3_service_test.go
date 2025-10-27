package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestS3Service_GeneratePresignedURL(t *testing.T) {
	// Note: This test requires AWS credentials and S3 bucket configuration
	// Skipping actual S3 operations in unit tests
	t.Skip("Skipping S3 service test - requires AWS credentials")

	service := NewS3Service()

	t.Run("Success - Generate presigned URL", func(t *testing.T) {
		fileName := "test-file.jpg"

		url, err := service.GeneratePresignedURL(fileName)

		assert.NoError(t, err)
		assert.NotEmpty(t, url)
	})
}

func TestS3Service_UploadFile(t *testing.T) {
	// Note: This test requires AWS credentials and S3 bucket configuration
	// Skipping actual S3 operations in unit tests
	t.Skip("Skipping S3 service test - requires AWS credentials")
}
