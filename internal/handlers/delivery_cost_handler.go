package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
)

// DeliveryCostHandler handles delivery cost-related requests
type DeliveryCostHandler struct {
	deliveryCostService services.DeliveryCostServiceInterface
}

// NewDeliveryCostHandler creates a new delivery cost handler
func NewDeliveryCostHandler(deliveryCostService services.DeliveryCostServiceInterface) *DeliveryCostHandler {
	return &DeliveryCostHandler{
		deliveryCostService: deliveryCostService,
	}
}

// AddCost adds a new delivery cost configuration (Admin only)
// @Summary Add delivery cost
// @Description Add a new delivery cost configuration (Admin only)
// @Tags delivery-cost
// @Accept json
// @Produce json
// @Param deliveryCost body models.DeliveryCostRequest true "Delivery cost information"
// @Success 201 {object} map[string]interface{} "success response with delivery cost data"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/delivery-cost [post]
func (h *DeliveryCostHandler) AddCost(c *gin.Context) {
	var req models.DeliveryCostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	deliveryCost := &models.DeliveryCost{
		Location: req.Location,
		Service:  req.Service,
		Cost:     *req.Cost,
	}

	createdCost, err := h.deliveryCostService.AddCost(deliveryCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: createdCost.ToDeliveryCostInfo(),
	})
}

// UpdateCost updates an existing delivery cost configuration (Admin only)
// @Summary Update delivery cost
// @Description Update an existing delivery cost configuration (Admin only)
// @Tags delivery-cost
// @Accept json
// @Produce json
// @Param id path string true "Delivery Cost ID"
// @Param deliveryCost body models.DeliveryCostRequest true "Delivery cost information"
// @Success 200 {object} map[string]interface{} "success response with updated delivery cost data"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 404 {object} map[string]interface{} "delivery cost not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/delivery-cost/{id} [put]
func (h *DeliveryCostHandler) UpdateCost(c *gin.Context) {
	id := c.Param("id")

	var req models.DeliveryCostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	deliveryCost := &models.DeliveryCost{
		Location: req.Location,
		Service:  req.Service,
		Cost:     *req.Cost,
	}

	updatedCost, err := h.deliveryCostService.UpdateCost(id, deliveryCost)
	if err != nil {
		if err == services.ErrDeliveryCostNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: updatedCost.ToDeliveryCostInfo(),
	})
}

// GetDeliveryCostByID retrieves a delivery cost by ID (Admin only)
// @Summary Get delivery cost by ID
// @Description Get delivery cost details by ID (Admin only)
// @Tags delivery-cost
// @Produce json
// @Param id path string true "Delivery Cost ID"
// @Success 200 {object} map[string]interface{} "success response with delivery cost data"
// @Failure 404 {object} map[string]interface{} "delivery cost not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/delivery-cost/{id} [get]
func (h *DeliveryCostHandler) GetDeliveryCostByID(c *gin.Context) {
	id := c.Param("id")

	deliveryCost, err := h.deliveryCostService.GetDeliveryCostByID(id)
	if err != nil {
		if err == services.ErrDeliveryCostNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: deliveryCost.ToDeliveryCostInfo(),
	})
}

// GetAllDeliveryCosts retrieves all delivery cost configurations
// @Summary Get all delivery costs
// @Description Get all delivery cost configurations
// @Tags delivery-cost
// @Produce json
// @Success 200 {object} map[string]interface{} "success response with array of delivery costs"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/delivery-cost/all [get]
func (h *DeliveryCostHandler) GetAllDeliveryCosts(c *gin.Context) {
	deliveryCosts, err := h.deliveryCostService.GetAllDeliveryCosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	// Convert to info objects
	var deliveryCostInfos []models.DeliveryCostInfo
	for _, cost := range deliveryCosts {
		deliveryCostInfos = append(deliveryCostInfos, cost.ToDeliveryCostInfo())
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: deliveryCostInfos,
	})
}

// GetDeliveryCostByLocation retrieves delivery costs by location
// @Summary Get delivery costs by location
// @Description Get delivery costs for a specific location
// @Tags delivery-cost
// @Produce json
// @Param location query string true "Location name"
// @Success 200 {object} map[string]interface{} "success response with array of delivery costs"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/delivery-cost/location [get]
func (h *DeliveryCostHandler) GetDeliveryCostByLocation(c *gin.Context) {
	location := c.Query("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: "location parameter is required",
		})
		return
	}

	deliveryCosts, err := h.deliveryCostService.GetDeliveryCostByLocation(location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	// Convert to info objects
	var deliveryCostInfos []models.DeliveryCostInfo
	for _, cost := range deliveryCosts {
		deliveryCostInfos = append(deliveryCostInfos, cost.ToDeliveryCostInfo())
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: deliveryCostInfos,
	})
}

// DeleteDeliveryCost deletes a delivery cost configuration (Admin only)
// @Summary Delete delivery cost
// @Description Delete a delivery cost configuration (Admin only)
// @Tags delivery-cost
// @Produce json
// @Param id path string true "Delivery Cost ID"
// @Success 200 {object} map[string]interface{} "success response"
// @Failure 404 {object} map[string]interface{} "delivery cost not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/delivery-cost/{id} [delete]
func (h *DeliveryCostHandler) DeleteDeliveryCost(c *gin.Context) {
	id := c.Param("id")

	err := h.deliveryCostService.DeleteDeliveryCost(id)
	if err != nil {
		if err == services.ErrDeliveryCostNotFound {
			c.JSON(http.StatusNotFound, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: "Delivery cost deleted",
	})
}
