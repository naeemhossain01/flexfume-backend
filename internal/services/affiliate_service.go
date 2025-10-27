package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"gorm.io/gorm"
)

var (
	ErrAffiliateSubmissionNotFound = errors.New("affiliate submission not found")
	ErrInvalidAffiliateStatus      = errors.New("invalid affiliate status")
	ErrNoSocialMediaProvided        = errors.New("at least one social media handle is required")
	ErrAffiliateSubmissionRequired  = errors.New("name and about are required")
)

// AffiliateService handles affiliate submission-related business logic
type AffiliateService struct{}

// NewAffiliateService creates a new affiliate service
func NewAffiliateService() *AffiliateService {
	return &AffiliateService{}
}

// SubmitAffiliateApplication creates a new affiliate application
func (s *AffiliateService) SubmitAffiliateApplication(req *models.SubmitAffiliateRequest) (*models.AffiliateSubmission, error) {
	// Validate required fields
	if strings.TrimSpace(req.Name) == "" {
		return nil, ErrAffiliateSubmissionRequired
	}
	if strings.TrimSpace(req.About) == "" {
		return nil, ErrAffiliateSubmissionRequired
	}

	// Validate at least one social media handle is provided
	if !s.hasAtLeastOneSocialMedia(req) {
		return nil, ErrNoSocialMediaProvided
	}

	// Create affiliate submission
	submission := &models.AffiliateSubmission{
		Name:        strings.TrimSpace(req.Name),
		About:       strings.TrimSpace(req.About),
		Phone:       strings.TrimSpace(req.Phone),
		Instagram:   strings.TrimSpace(req.Instagram),
		Facebook:    strings.TrimSpace(req.Facebook),
		YouTube:     strings.TrimSpace(req.YouTube),
		LinkedIn:    strings.TrimSpace(req.LinkedIn),
		OtherSocial: strings.TrimSpace(req.OtherSocial),
		Status:      models.StatusPending,
		CreatedBy:   "SYSTEM", // Since this is a public submission
	}

	// Save to database
	if err := database.GetDB().Create(submission).Error; err != nil {
		return nil, fmt.Errorf("failed to create affiliate submission: %w", err)
	}

	return submission, nil
}

// GetAffiliateSubmissionByID retrieves an affiliate submission by ID
func (s *AffiliateService) GetAffiliateSubmissionByID(id string) (*models.AffiliateSubmission, error) {
	var submission models.AffiliateSubmission
	if err := database.GetDB().First(&submission, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrAffiliateSubmissionNotFound
		}
		return nil, err
	}
	return &submission, nil
}

// GetAllAffiliateSubmissions retrieves all affiliate submissions with pagination and filtering
func (s *AffiliateService) GetAllAffiliateSubmissions(status string, page, limit int) (*models.AffiliateSubmissionListResponse, error) {
	var submissions []models.AffiliateSubmission
	var total int64

	// Build query
	query := database.GetDB().Model(&models.AffiliateSubmission{})

	// Apply status filter if provided
	if status != "" {
		validStatuses := []string{"PENDING", "APPROVED", "REJECTED"}
		isValidStatus := false
		for _, validStatus := range validStatuses {
			if strings.EqualFold(status, validStatus) {
				query = query.Where("status = ?", strings.ToUpper(status))
				isValidStatus = true
				break
			}
		}
		if !isValidStatus {
			return nil, ErrInvalidAffiliateStatus
		}
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&submissions).Error; err != nil {
		return nil, err
	}

	// Convert to response format
	submissionResponses := make([]models.AffiliateSubmissionResponse, len(submissions))
	for i, submission := range submissions {
		submissionResponses[i] = submission.ToResponse()
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &models.AffiliateSubmissionListResponse{
		Submissions: submissionResponses,
		Total:       total,
		Page:        page,
		Limit:       limit,
		TotalPages:  totalPages,
	}, nil
}

// UpdateAffiliateStatus updates the status of an affiliate submission
func (s *AffiliateService) UpdateAffiliateStatus(id string, req *models.UpdateAffiliateStatusRequest) (*models.AffiliateSubmission, error) {
	// Get existing submission
	submission, err := s.GetAffiliateSubmissionByID(id)
	if err != nil {
		return nil, err
	}

	// Update status and admin notes
	submission.Status = req.Status
	submission.AdminNotes = strings.TrimSpace(req.AdminNotes)
	submission.UpdatedBy = "ADMIN" // This should be set by the handler with actual admin user ID

	// Save changes
	if err := database.GetDB().Save(submission).Error; err != nil {
		return nil, fmt.Errorf("failed to update affiliate submission: %w", err)
	}

	return submission, nil
}

// DeleteAffiliateSubmission soft deletes an affiliate submission
func (s *AffiliateService) DeleteAffiliateSubmission(id string) error {
	result := database.GetDB().Delete(&models.AffiliateSubmission{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrAffiliateSubmissionNotFound
	}
	return nil
}

// hasAtLeastOneSocialMedia checks if at least one social media handle is provided
func (s *AffiliateService) hasAtLeastOneSocialMedia(req *models.SubmitAffiliateRequest) bool {
	socialMediaHandles := []string{
		strings.TrimSpace(req.Instagram),
		strings.TrimSpace(req.Facebook),
		strings.TrimSpace(req.YouTube),
		strings.TrimSpace(req.LinkedIn),
		strings.TrimSpace(req.OtherSocial),
	}

	for _, handle := range socialMediaHandles {
		if handle != "" {
			return true
		}
	}
	return false
}

// SanitizeInput sanitizes text input to prevent XSS attacks
func (s *AffiliateService) SanitizeInput(input string) string {
	// Basic XSS prevention - remove potentially dangerous characters
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&#x27;")
	input = strings.ReplaceAll(input, "&", "&amp;")
	return strings.TrimSpace(input)
}

// SanitizeAffiliateSubmission sanitizes all text fields in an affiliate submission
func (s *AffiliateService) SanitizeAffiliateSubmission(submission *models.AffiliateSubmission) {
	submission.Name = s.SanitizeInput(submission.Name)
	submission.About = s.SanitizeInput(submission.About)
	submission.Phone = s.SanitizeInput(submission.Phone)
	submission.Instagram = s.SanitizeInput(submission.Instagram)
	submission.Facebook = s.SanitizeInput(submission.Facebook)
	submission.YouTube = s.SanitizeInput(submission.YouTube)
	submission.LinkedIn = s.SanitizeInput(submission.LinkedIn)
	submission.OtherSocial = s.SanitizeInput(submission.OtherSocial)
	submission.AdminNotes = s.SanitizeInput(submission.AdminNotes)
}
