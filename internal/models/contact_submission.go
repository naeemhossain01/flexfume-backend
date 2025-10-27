package models

import (
	"time"

	"gorm.io/gorm"
)

// ContactSubmission represents a contact form submission in the system
type ContactSubmission struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"not null;size:255" json:"name"`
	Email     string    `gorm:"not null;size:255" json:"email"`
	Phone     string    `gorm:"size:50" json:"phone"`
	Subject   string    `gorm:"size:255" json:"subject"`
	Message   string    `gorm:"not null;type:text" json:"message"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
	CreatedBy string    `gorm:"column:created_by" json:"-"`
	UpdatedBy string    `gorm:"column:updated_by" json:"-"`
}

// TableName specifies the table name for the ContactSubmission model
func (ContactSubmission) TableName() string {
	return "contact_submissions"
}

// ContactSubmissionResponse represents the contact submission data returned in API responses
type ContactSubmissionResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToResponse converts a ContactSubmission model to ContactSubmissionResponse
func (c *ContactSubmission) ToResponse() ContactSubmissionResponse {
	return ContactSubmissionResponse{
		ID:        c.ID,
		Name:      c.Name,
		Email:     c.Email,
		Phone:     c.Phone,
		Subject:   c.Subject,
		Message:   c.Message,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// SubmitContactRequest represents the request body for submitting a contact form
type SubmitContactRequest struct {
	Name    string `json:"name" binding:"required,max=255"`
	Email   string `json:"email" binding:"required,email,max=255"`
	Phone   string `json:"phone" binding:"max=50"`
	Subject string `json:"subject" binding:"max=255"`
	Message string `json:"message" binding:"required"`
}
