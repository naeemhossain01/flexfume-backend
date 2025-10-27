package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
)

// AffiliateHandler handles affiliate-related requests
type AffiliateHandler struct {
	affiliateService *services.AffiliateService
}

// NewAffiliateHandler creates a new affiliate handler
func NewAffiliateHandler(affiliateService *services.AffiliateService) *AffiliateHandler {
	return &AffiliateHandler{
		affiliateService: affiliateService,
	}
}

// SubmitAffiliateApplication handles POST /api/v1/affiliate/submit
func (h *AffiliateHandler) SubmitAffiliateApplication(c *gin.Context) {
	var req models.SubmitAffiliateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Sanitize input to prevent XSS attacks
	req.Name = h.affiliateService.SanitizeInput(req.Name)
	req.About = h.affiliateService.SanitizeInput(req.About)
	req.Phone = h.affiliateService.SanitizeInput(req.Phone)
	req.Instagram = h.affiliateService.SanitizeInput(req.Instagram)
	req.Facebook = h.affiliateService.SanitizeInput(req.Facebook)
	req.YouTube = h.affiliateService.SanitizeInput(req.YouTube)
	req.LinkedIn = h.affiliateService.SanitizeInput(req.LinkedIn)
	req.OtherSocial = h.affiliateService.SanitizeInput(req.OtherSocial)

	// Submit application
	submission, err := h.affiliateService.SubmitAffiliateApplication(&req)
	if err != nil {
		switch err {
		case services.ErrAffiliateSubmissionRequired:
			c.JSON(http.StatusBadRequest, APIResponse{
				Error:   true,
				Message: "Name and about are required",
			})
		case services.ErrNoSocialMediaProvided:
			c.JSON(http.StatusBadRequest, APIResponse{
				Error:   true,
				Message: "At least one social media handle is required",
			})
		default:
			c.JSON(http.StatusInternalServerError, APIResponse{
				Error:   true,
				Message: "Failed to submit affiliate application",
			})
		}
		return
	}

	// Return success response
	c.JSON(http.StatusOK, APIResponse{
		Error:   false,
		Message: "Affiliate application submitted successfully",
		Response: map[string]interface{}{
			"id":          submission.ID,
			"status":      submission.Status,
			"submittedAt": submission.CreatedAt,
		},
	})
}

// GetAffiliateApplications handles GET /api/v1/affiliate/applications
func (h *AffiliateHandler) GetAffiliateApplications(c *gin.Context) {
	// Get query parameters
	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	// Parse pagination parameters
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Get applications
	result, err := h.affiliateService.GetAllAffiliateSubmissions(status, page, limit)
	if err != nil {
		switch err {
		case services.ErrInvalidAffiliateStatus:
			c.JSON(http.StatusBadRequest, APIResponse{
				Error:   true,
				Message: "Invalid status filter. Valid values are: PENDING, APPROVED, REJECTED",
			})
		default:
			c.JSON(http.StatusInternalServerError, APIResponse{
				Error:   true,
				Message: "Failed to retrieve affiliate applications",
			})
		}
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: result,
	})
}

// UpdateAffiliateStatus handles PUT /api/v1/affiliate/applications/{id}/status
func (h *AffiliateHandler) UpdateAffiliateStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "Affiliate application ID is required",
		})
		return
	}

	var req models.UpdateAffiliateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "Invalid request body: " + err.Error(),
		})
		return
	}

	// Update status
	submission, err := h.affiliateService.UpdateAffiliateStatus(id, &req)
	if err != nil {
		switch err {
		case services.ErrAffiliateSubmissionNotFound:
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: "Affiliate application not found",
			})
		case services.ErrInvalidAffiliateStatus:
			c.JSON(http.StatusBadRequest, APIResponse{
				Error:   true,
				Message: "Invalid status. Valid values are: PENDING, APPROVED, REJECTED",
			})
		default:
			c.JSON(http.StatusInternalServerError, APIResponse{
				Error:   true,
				Message: "Failed to update affiliate application status",
			})
		}
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "Affiliate application status updated successfully",
		Response: submission.ToResponse(),
	})
}

// GetAffiliateApplicationByID handles GET /api/v1/affiliate/applications/{id}
func (h *AffiliateHandler) GetAffiliateApplicationByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "Affiliate application ID is required",
		})
		return
	}

	submission, err := h.affiliateService.GetAffiliateSubmissionByID(id)
	if err != nil {
		switch err {
		case services.ErrAffiliateSubmissionNotFound:
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: "Affiliate application not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, APIResponse{
				Error:   true,
				Message: "Failed to retrieve affiliate application",
			})
		}
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: submission.ToResponse(),
	})
}

// DeleteAffiliateApplication handles DELETE /api/v1/affiliate/applications/{id}
func (h *AffiliateHandler) DeleteAffiliateApplication(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "Affiliate application ID is required",
		})
		return
	}

	err := h.affiliateService.DeleteAffiliateSubmission(id)
	if err != nil {
		switch err {
		case services.ErrAffiliateSubmissionNotFound:
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: "Affiliate application not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, APIResponse{
				Error:   true,
				Message: "Failed to delete affiliate application",
			})
		}
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:   false,
		Message: "Affiliate application deleted successfully",
	})
}
