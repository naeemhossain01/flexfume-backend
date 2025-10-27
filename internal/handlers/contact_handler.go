package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
)

// ContactHandler handles contact form-related requests
type ContactHandler struct {
	contactService *services.ContactService
}

// NewContactHandler creates a new contact handler
func NewContactHandler(contactService *services.ContactService) *ContactHandler {
	return &ContactHandler{
		contactService: contactService,
	}
}

// SubmitContactForm handles POST /api/v1/contact
func (h *ContactHandler) SubmitContactForm(c *gin.Context) {
	var req models.SubmitContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Sanitize input to prevent XSS attacks
	req.Name = h.contactService.SanitizeInput(req.Name)
	req.Email = h.contactService.SanitizeInput(req.Email)
	req.Phone = h.contactService.SanitizeInput(req.Phone)
	req.Subject = h.contactService.SanitizeInput(req.Subject)
	req.Message = h.contactService.SanitizeInput(req.Message)

	// Submit contact form
	submission, err := h.contactService.SubmitContactForm(&req)
	if err != nil {
		switch err {
		case services.ErrContactSubmissionRequired:
			c.JSON(http.StatusBadRequest, APIResponse{
				Error:   true,
				Message: "Name, email, and message are required",
			})
		case services.ErrInvalidEmail:
			c.JSON(http.StatusBadRequest, APIResponse{
				Error:   true,
				Message: "Invalid email format",
			})
		case services.ErrInvalidPhone:
			c.JSON(http.StatusBadRequest, APIResponse{
				Error:   true,
				Message: "Invalid phone format",
			})
		default:
			c.JSON(http.StatusInternalServerError, APIResponse{
				Error:   true,
				Message: "Failed to submit contact form",
			})
		}
		return
	}

	// Return success response (201 Created)
	c.JSON(http.StatusCreated, APIResponse{
		Error:   false,
		Message: "Contact form submitted successfully",
		Response: map[string]interface{}{
			"contactId":   submission.ID,
			"submittedAt": submission.CreatedAt,
		},
	})
}

// GetAllContactSubmissions handles GET /api/v1/contact (Admin only)
func (h *ContactHandler) GetAllContactSubmissions(c *gin.Context) {
	// Parse pagination parameters
	page := 1
	limit := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get contact submissions
	submissions, total, err := h.contactService.GetAllContactSubmissions(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: "Failed to retrieve contact submissions",
		})
		return
	}

	// Convert to response format
	submissionResponses := make([]models.ContactSubmissionResponse, len(submissions))
	for i, submission := range submissions {
		submissionResponses[i] = submission.ToResponse()
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:   false,
		Message: "Contact submissions retrieved successfully",
		Response: map[string]interface{}{
			"submissions": submissionResponses,
			"total":       total,
			"page":        page,
			"limit":       limit,
			"totalPages":  totalPages,
		},
	})
}

// GetContactSubmissionByID handles GET /api/v1/contact/:id (Admin only)
func (h *ContactHandler) GetContactSubmissionByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "Contact submission ID is required",
		})
		return
	}

	// Get contact submission
	submission, err := h.contactService.GetContactSubmissionByID(id)
	if err != nil {
		switch err {
		case services.ErrContactSubmissionNotFound:
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: "Contact submission not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, APIResponse{
				Error:   true,
				Message: "Failed to retrieve contact submission",
			})
		}
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:   false,
		Message: "Contact submission retrieved successfully",
		Response: submission.ToResponse(),
	})
}

// DeleteContactSubmission handles DELETE /api/v1/contact/:id (Admin only)
func (h *ContactHandler) DeleteContactSubmission(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "Contact submission ID is required",
		})
		return
	}

	// Delete contact submission
	err := h.contactService.DeleteContactSubmission(id)
	if err != nil {
		switch err {
		case services.ErrContactSubmissionNotFound:
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: "Contact submission not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, APIResponse{
				Error:   true,
				Message: "Failed to delete contact submission",
			})
		}
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:   false,
		Message: "Contact submission deleted successfully",
	})
}