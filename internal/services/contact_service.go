package services

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"gorm.io/gorm"
)

var (
	ErrContactSubmissionNotFound = errors.New("contact submission not found")
	ErrInvalidEmail              = errors.New("invalid email format")
	ErrInvalidPhone              = errors.New("invalid phone format")
	ErrContactSubmissionRequired  = errors.New("name, email, and message are required")
)

// ContactService handles contact form submission-related business logic
type ContactService struct{}

// NewContactService creates a new contact service
func NewContactService() *ContactService {
	return &ContactService{}
}

// SubmitContactForm creates a new contact form submission
func (s *ContactService) SubmitContactForm(req *models.SubmitContactRequest) (*models.ContactSubmission, error) {
	// Validate required fields
	if strings.TrimSpace(req.Name) == "" {
		return nil, ErrContactSubmissionRequired
	}
	if strings.TrimSpace(req.Email) == "" {
		return nil, ErrContactSubmissionRequired
	}
	if strings.TrimSpace(req.Message) == "" {
		return nil, ErrContactSubmissionRequired
	}

	// Validate email format
	if !s.isValidEmail(req.Email) {
		return nil, ErrInvalidEmail
	}

	// Validate phone format if provided
	if req.Phone != "" && !s.isValidPhone(req.Phone) {
		return nil, ErrInvalidPhone
	}

	// Create contact submission
	submission := &models.ContactSubmission{
		Name:      strings.TrimSpace(req.Name),
		Email:     strings.ToLower(strings.TrimSpace(req.Email)),
		Phone:     strings.TrimSpace(req.Phone),
		Subject:   strings.TrimSpace(req.Subject),
		Message:   strings.TrimSpace(req.Message),
		CreatedBy: "SYSTEM", // Since this is a public submission
	}

	// Save to database
	if err := database.GetDB().Create(submission).Error; err != nil {
		return nil, fmt.Errorf("failed to create contact submission: %w", err)
	}

	return submission, nil
}

// GetContactSubmissionByID retrieves a contact submission by ID
func (s *ContactService) GetContactSubmissionByID(id string) (*models.ContactSubmission, error) {
	var submission models.ContactSubmission
	if err := database.GetDB().First(&submission, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrContactSubmissionNotFound
		}
		return nil, err
	}
	return &submission, nil
}

// GetAllContactSubmissions retrieves all contact submissions with pagination
func (s *ContactService) GetAllContactSubmissions(page, limit int) ([]models.ContactSubmission, int64, error) {
	var submissions []models.ContactSubmission
	var total int64

	// Count total records
	if err := database.GetDB().Model(&models.ContactSubmission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := database.GetDB().Order("created_at DESC").Offset(offset).Limit(limit).Find(&submissions).Error; err != nil {
		return nil, 0, err
	}

	return submissions, total, nil
}

// DeleteContactSubmission soft deletes a contact submission
func (s *ContactService) DeleteContactSubmission(id string) error {
	result := database.GetDB().Delete(&models.ContactSubmission{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrContactSubmissionNotFound
	}
	return nil
}

// isValidEmail validates email format using regex
func (s *ContactService) isValidEmail(email string) bool {
	// Basic email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// isValidPhone validates phone format (basic validation)
func (s *ContactService) isValidPhone(phone string) bool {
	// Remove all non-digit characters for validation
	digitsOnly := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	// Phone should have 7-15 digits (international standard)
	return len(digitsOnly) >= 7 && len(digitsOnly) <= 15
}

// SanitizeInput sanitizes text input to prevent XSS attacks
func (s *ContactService) SanitizeInput(input string) string {
	// Basic XSS prevention - remove potentially dangerous characters
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&#x27;")
	input = strings.ReplaceAll(input, "&", "&amp;")
	return strings.TrimSpace(input)
}

// SanitizeContactSubmission sanitizes all text fields in a contact submission
func (s *ContactService) SanitizeContactSubmission(submission *models.ContactSubmission) {
	submission.Name = s.SanitizeInput(submission.Name)
	submission.Email = s.SanitizeInput(submission.Email)
	submission.Phone = s.SanitizeInput(submission.Phone)
	submission.Subject = s.SanitizeInput(submission.Subject)
	submission.Message = s.SanitizeInput(submission.Message)
}
