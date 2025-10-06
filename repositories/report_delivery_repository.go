package repository

import (
	"scheduling-report/config"
	"scheduling-report/models"

	"gorm.io/gorm"
)

type ReportDeliveryRepository struct {
	DB *gorm.DB
}

func NewReportDeliveryRepository() *ReportDeliveryRepository {
	return &ReportDeliveryRepository{
		DB: config.DB,
	}
}

// GetAll retrieves all active deliveries
func (r *ReportDeliveryRepository) GetAll(isActive *bool) ([]models.ReportDelivery, error) {
	var deliveries []models.ReportDelivery
	query := r.DB.Where("is_active = ?", true)

	if isActive != nil {
		query = r.DB.Where("is_active = ?", *isActive)
	}

	err := query.Order("created_at DESC").Find(&deliveries).Error
	return deliveries, err
}

// GetByID retrieves a delivery by ID
func (r *ReportDeliveryRepository) GetByID(id int) (*models.ReportDelivery, error) {
	var delivery models.ReportDelivery
	err := r.DB.Where("id = ?", id).First(&delivery).Error
	if err != nil {
		return nil, err
	}
	return &delivery, nil
}

// GetByConfigID retrieves all deliveries for a config
func (r *ReportDeliveryRepository) GetByConfigID(configID int) ([]models.ReportDelivery, error) {
	var deliveries []models.ReportDelivery
	err := r.DB.Where("config_id = ? AND is_active = ?", configID, true).Find(&deliveries).Error
	return deliveries, err
}

// Create inserts a new delivery
func (r *ReportDeliveryRepository) Create(delivery *models.ReportDelivery) error {
	return r.DB.Create(delivery).Error
}

// Update updates an existing delivery
func (r *ReportDeliveryRepository) Update(delivery *models.ReportDelivery) error {
	return r.DB.Save(delivery).Error
}

// Delete soft deletes a delivery
func (r *ReportDeliveryRepository) Delete(id int) error {
	return r.DB.Model(&models.ReportDelivery{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

// CheckConfigExists verifies if a report config exists
func (r *ReportDeliveryRepository) CheckConfigExists(configID int) (bool, error) {
	var count int64
	err := r.DB.Model(&models.ReportConfig{}).
		Where("id = ? AND is_active = ?", configID, true).
		Count(&count).Error
	return count > 0, err
}
