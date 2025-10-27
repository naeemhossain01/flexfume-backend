package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/services"
)

// AddressHandler handles address-related requests
type AddressHandler struct {
	addressService services.AddressServiceInterface
}

// NewAddressHandler creates a new address handler
func NewAddressHandler(addressService services.AddressServiceInterface) *AddressHandler {
	return &AddressHandler{
		addressService: addressService,
	}
}

// AddAddress adds a new address for a user
// @Summary Add a new address
// @Description Add a new address for authenticated user
// @Tags address
// @Accept json
// @Produce json
// @Param address body models.AddressRequest true "Address information"
// @Success 201 {object} map[string]interface{} "success response with address data"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 409 {object} map[string]interface{} "address already exists"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/address [post]
func (h *AddressHandler) AddAddress(c *gin.Context) {
	// Get user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Error:   true,
			Message: "User not authenticated",
		})
		return
	}

	var req models.AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	address := &models.Address{
		UserID:      userID.(string),
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		IsDefault:   req.IsDefault,
	}

	createdAddress, err := h.addressService.AddAddress(address)
	if err != nil {
		if err == services.ErrAddressAlreadyExists {
			c.JSON(http.StatusConflict, APIResponse{
				Error:   true,
				Message: err.Error(),
			})
			return
		}
		if err == services.ErrUserNotFound {
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

	c.JSON(http.StatusCreated, APIResponse{
		Error:    false,
		Message:  "SUCCESS",
		Response: createdAddress.ToAddressInfo(),
	})
}

// UpdateAddress updates an existing address
// @Summary Update an address
// @Description Update an existing address by ID for authenticated user
// @Tags address
// @Accept json
// @Produce json
// @Param id path string true "Address ID"
// @Param address body models.AddressRequest true "Address information"
// @Success 200 {object} map[string]interface{} "success response with updated address data"
// @Failure 400 {object} map[string]interface{} "bad request"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 404 {object} map[string]interface{} "address not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/address/{id} [put]
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	// Get user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Error:   true,
			Message: "User not authenticated",
		})
		return
	}

	addressID := c.Param("id")

	var req models.AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	address := &models.Address{
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		IsDefault:   req.IsDefault,
	}

	updatedAddress, err := h.addressService.UpdateAddress(addressID, address, userID.(string))
	if err != nil {
		if err == services.ErrAddressNotFound {
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
		Response: updatedAddress.ToAddressInfo(),
	})
}

// GetAddressByUser retrieves an address for authenticated user
// @Summary Get address for authenticated user
// @Description Get address for the authenticated user
// @Tags address
// @Produce json
// @Success 200 {object} map[string]interface{} "success response with address data"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 404 {object} map[string]interface{} "address not found"
// @Failure 500 {object} map[string]interface{} "internal server error"
// @Router /api/v1/address [get]
func (h *AddressHandler) GetAddressByUser(c *gin.Context) {
	// Get user ID from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Error:   true,
			Message: "User not authenticated",
		})
		return
	}

	address, err := h.addressService.GetAddressByUserID(userID.(string))
	if err != nil {
		if err == services.ErrAddressNotFound {
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
		Response: address.ToAddressInfo(),
	})
}
