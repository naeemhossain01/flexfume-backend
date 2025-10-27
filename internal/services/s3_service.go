package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/config"
)

// S3Service handles file uploads to AWS S3
type S3Service struct {
	s3Client   *s3.S3
	bucketName string
	region     string
}

// NewS3Service creates a new S3 service instance
func NewS3Service(cfg *config.Config) (*S3Service, error) {
	if cfg.AWS.AccessKeyID == "" || cfg.AWS.SecretAccessKey == "" {
		return nil, fmt.Errorf("AWS credentials not configured")
	}

	if cfg.AWS.BucketName == "" {
		return nil, fmt.Errorf("AWS S3 bucket name not configured")
	}

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.AWS.Region),
		Credentials: credentials.NewStaticCredentials(
			cfg.AWS.AccessKeyID,
			cfg.AWS.SecretAccessKey,
			"",
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &S3Service{
		s3Client:   s3.New(sess),
		bucketName: cfg.AWS.BucketName,
		region:     cfg.AWS.Region,
	}, nil
}

// UploadFile uploads a file to S3 and returns the public URL
func (s *S3Service) UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// Generate unique filename
	filename := s.generateUniqueFilename(fileHeader.Filename)

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Determine content type
	contentType := s.getContentType(fileHeader.Filename)

	// Upload to S3
	_, err = s.s3Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(filename),
		Body:          bytes.NewReader(fileBytes),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(int64(len(fileBytes))),
		ACL:           aws.String("public-read"), // Make file publicly accessible
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Generate public URL
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, filename)
	return url, nil
}

// DeleteFile deletes a file from S3
func (s *S3Service) DeleteFile(fileURL string) error {
	// Extract filename from URL
	filename := s.extractFilenameFromURL(fileURL)
	if filename == "" {
		return fmt.Errorf("invalid file URL")
	}

	_, err := s.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(filename),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// generateUniqueFilename generates a unique filename using UUID
func (s *S3Service) generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	uniqueID := uuid.New().String()
	timestamp := time.Now().Unix()
	return fmt.Sprintf("products/%d-%s%s", timestamp, uniqueID, ext)
}

// extractFilenameFromURL extracts the filename from S3 URL
func (s *S3Service) extractFilenameFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return ""
	}
	// Return the last two parts (e.g., "products/filename.jpg")
	return strings.Join(parts[len(parts)-2:], "/")
}

// getContentType determines the content type based on file extension
func (s *S3Service) getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
	}

	if contentType, ok := contentTypes[ext]; ok {
		return contentType
	}
	return "application/octet-stream"
}

// ValidateImageFile validates if the file is a valid image
func (s *S3Service) ValidateImageFile(fileHeader *multipart.FileHeader) error {
	// Check file size (max 5MB)
	maxSize := int64(5 * 1024 * 1024) // 5MB
	if fileHeader.Size > maxSize {
		return fmt.Errorf("file size exceeds maximum allowed size of 5MB")
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg"}
	
	isValid := false
	for _, validExt := range validExtensions {
		if ext == validExt {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("invalid file type. Allowed types: jpg, jpeg, png, gif, webp, svg")
	}

	return nil
}
