package models

import (
	"time"

	"gorm.io/gorm"
)

// AffiliateSubmissionStatus represents the status of an affiliate application
type AffiliateSubmissionStatus string

const (
	StatusPending  AffiliateSubmissionStatus = "PENDING"
	StatusApproved AffiliateSubmissionStatus = "APPROVED"
	StatusRejected AffiliateSubmissionStatus = "REJECTED"
)

// AffiliateSubmission represents an affiliate application in the system
type AffiliateSubmission struct {
	ID           string                     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string                     `gorm:"not null;size:255" json:"name"`
	About        string                     `gorm:"not null;type:text" json:"about"`
	Phone        string                     `gorm:"size:50" json:"phone"`
	Instagram    string                     `gorm:"size:255" json:"instagram"`
	Facebook     string                     `gorm:"size:255" json:"facebook"`
	YouTube      string                     `gorm:"column:youtube;size:255" json:"youtube"`
	LinkedIn     string                     `gorm:"column:linkedin;size:255" json:"linkedin"`
	OtherSocial  string                     `gorm:"size:255" json:"otherSocial"`
	Status       AffiliateSubmissionStatus `gorm:"type:varchar(50);default:'PENDING'" json:"status"`
	AdminNotes   string                     `gorm:"type:text" json:"adminNotes,omitempty"`
	CreatedAt    time.Time                  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time                  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt             `gorm:"column:deleted_at;index" json:"-"`
	CreatedBy    string                     `gorm:"column:created_by" json:"-"`
	UpdatedBy    string                     `gorm:"column:updated_by" json:"-"`
}

// TableName specifies the table name for the AffiliateSubmission model
func (AffiliateSubmission) TableName() string {
	return "affiliate_submissions"
}

// AffiliateSubmissionResponse represents the affiliate submission data returned in API responses
type AffiliateSubmissionResponse struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	About       string                     `json:"about"`
	Phone       string                     `json:"phone"`
	Instagram   string                     `json:"instagram"`
	Facebook    string                     `json:"facebook"`
	YouTube     string                     `json:"youtube"`
	LinkedIn    string                     `json:"linkedin"`
	OtherSocial string                     `json:"otherSocial"`
	Status      AffiliateSubmissionStatus  `json:"status"`
	AdminNotes  string                     `json:"adminNotes,omitempty"`
	CreatedAt   time.Time                  `json:"createdAt"`
	UpdatedAt   time.Time                  `json:"updatedAt"`
}

// ToResponse converts an AffiliateSubmission model to AffiliateSubmissionResponse
func (a *AffiliateSubmission) ToResponse() AffiliateSubmissionResponse {
	return AffiliateSubmissionResponse{
		ID:          a.ID,
		Name:        a.Name,
		About:       a.About,
		Phone:       a.Phone,
		Instagram:   a.Instagram,
		Facebook:    a.Facebook,
		YouTube:     a.YouTube,
		LinkedIn:    a.LinkedIn,
		OtherSocial: a.OtherSocial,
		Status:      a.Status,
		AdminNotes:  a.AdminNotes,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

// SubmitAffiliateRequest represents the request body for submitting an affiliate application
type SubmitAffiliateRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	About       string `json:"about" binding:"required,max=2000"`
	Phone       string `json:"phone" binding:"max=50"`
	Instagram   string `json:"instagram" binding:"max=255"`
	Facebook    string `json:"facebook" binding:"max=255"`
	YouTube     string `json:"youtube" binding:"max=255"`
	LinkedIn    string `json:"linkedin" binding:"max=255"`
	OtherSocial string `json:"otherSocial" binding:"max=255"`
}

// UpdateAffiliateStatusRequest represents the request body for updating affiliate status
type UpdateAffiliateStatusRequest struct {
	Status     AffiliateSubmissionStatus `json:"status" binding:"required,oneof=PENDING APPROVED REJECTED"`
	AdminNotes string                    `json:"notes" binding:"max=1000"`
}

// AffiliateSubmissionListResponse represents paginated affiliate submissions
type AffiliateSubmissionListResponse struct {
	Submissions []AffiliateSubmissionResponse `json:"submissions"`
	Total       int64                          `json:"total"`
	Page        int                            `json:"page"`
	Limit       int                            `json:"limit"`
	TotalPages  int                            `json:"totalPages"`
}
