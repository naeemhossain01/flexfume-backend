package services

import (
	"errors"
	"strings"

	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/database"
	"github.com/seamlance/client-flexfume-ecom-backend-go/internal/models"
	"gorm.io/gorm"
)

var (
	ErrDeliveryCostNotFound = errors.New("delivery cost not found")
)

// DeliveryCostService handles delivery cost-related business logic
type DeliveryCostService struct{}

// NewDeliveryCostService creates a new delivery cost service
func NewDeliveryCostService() *DeliveryCostService {
	return &DeliveryCostService{}
}

// AddCost creates a new delivery cost configuration
func (s *DeliveryCostService) AddCost(deliveryCost *models.DeliveryCost) (*models.DeliveryCost, error) {
	if err := database.GetDB().Create(deliveryCost).Error; err != nil {
		return nil, err
	}
	return deliveryCost, nil
}

// UpdateCost updates an existing delivery cost configuration
func (s *DeliveryCostService) UpdateCost(id string, deliveryCost *models.DeliveryCost) (*models.DeliveryCost, error) {
	// Get existing delivery cost
	existingCost, err := s.GetDeliveryCostByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if deliveryCost.Location != "" {
		existingCost.Location = deliveryCost.Location
	}
	if deliveryCost.Service != "" {
		existingCost.Service = deliveryCost.Service
	}
	if deliveryCost.Cost > 0 {
		existingCost.Cost = deliveryCost.Cost
	}

	// Save updated delivery cost
	if err := database.GetDB().Save(existingCost).Error; err != nil {
		return nil, err
	}

	return existingCost, nil
}

// GetDeliveryCostByID retrieves a delivery cost by ID
func (s *DeliveryCostService) GetDeliveryCostByID(id string) (*models.DeliveryCost, error) {
	var deliveryCost models.DeliveryCost
	if err := database.GetDB().First(&deliveryCost, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrDeliveryCostNotFound
		}
		return nil, err
	}
	return &deliveryCost, nil
}

// GetAllDeliveryCosts retrieves all delivery cost configurations
func (s *DeliveryCostService) GetAllDeliveryCosts() ([]models.DeliveryCost, error) {
	var deliveryCosts []models.DeliveryCost
	if err := database.GetDB().Find(&deliveryCosts).Error; err != nil {
		return nil, err
	}
	return deliveryCosts, nil
}

// GetDeliveryCostByLocation retrieves delivery costs by location (partial match)
func (s *DeliveryCostService) GetDeliveryCostByLocation(location string) ([]models.DeliveryCost, error) {
	var deliveryCosts []models.DeliveryCost
	
	// Use ILIKE for case-insensitive partial match (PostgreSQL)
	searchPattern := "%" + strings.ToLower(location) + "%"
	if err := database.GetDB().Where("LOWER(location) LIKE ?", searchPattern).Find(&deliveryCosts).Error; err != nil {
		return nil, err
	}
	
	return deliveryCosts, nil
}

// DeleteDeliveryCost deletes a delivery cost configuration
func (s *DeliveryCostService) DeleteDeliveryCost(id string) error {
	// Check if delivery cost exists
	if _, err := s.GetDeliveryCostByID(id); err != nil {
		return err
	}

	// Delete the delivery cost
	if err := database.GetDB().Delete(&models.DeliveryCost{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

// Ensure DeliveryCostService implements DeliveryCostServiceInterface
var _ DeliveryCostServiceInterface = (*DeliveryCostService)(nil)
