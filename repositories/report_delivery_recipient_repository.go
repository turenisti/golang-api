package repository

import (
	"scheduling-report/config"
	"scheduling-report/models"

	"gorm.io/gorm"
)

type ReportDeliveryRecipientRepository struct {
	DB *gorm.DB
}

func NewReportDeliveryRecipientRepository() *ReportDeliveryRecipientRepository {
	return &ReportDeliveryRecipientRepository{
		DB: config.DB,
	}
}

// GetAll retrieves all active recipients
func (r *ReportDeliveryRecipientRepository) GetAll(isActive *bool) ([]models.ReportDeliveryRecipient, error) {
	var recipients []models.ReportDeliveryRecipient
	query := r.DB.Where("is_active = ?", true)

	if isActive != nil {
		query = r.DB.Where("is_active = ?", *isActive)
	}

	err := query.Order("created_at DESC").Find(&recipients).Error
	return recipients, err
}

// GetByID retrieves a recipient by ID
func (r *ReportDeliveryRecipientRepository) GetByID(id int) (*models.ReportDeliveryRecipient, error) {
	var recipient models.ReportDeliveryRecipient
	err := r.DB.Where("id = ?", id).First(&recipient).Error
	if err != nil {
		return nil, err
	}
	return &recipient, nil
}

// GetByDeliveryID retrieves all recipients for a delivery
func (r *ReportDeliveryRecipientRepository) GetByDeliveryID(deliveryID int) ([]models.ReportDeliveryRecipient, error) {
	var recipients []models.ReportDeliveryRecipient
	err := r.DB.Where("delivery_id = ? AND is_active = ?", deliveryID, true).Find(&recipients).Error
	return recipients, err
}

// Create inserts a new recipient
func (r *ReportDeliveryRecipientRepository) Create(recipient *models.ReportDeliveryRecipient) error {
	return r.DB.Create(recipient).Error
}

// Update updates an existing recipient
func (r *ReportDeliveryRecipientRepository) Update(recipient *models.ReportDeliveryRecipient) error {
	return r.DB.Save(recipient).Error
}

// Delete soft deletes a recipient
func (r *ReportDeliveryRecipientRepository) Delete(id int) error {
	return r.DB.Model(&models.ReportDeliveryRecipient{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

// CheckDeliveryExists verifies if a delivery exists
func (r *ReportDeliveryRecipientRepository) CheckDeliveryExists(deliveryID int) (bool, error) {
	var count int64
	err := r.DB.Model(&models.ReportDelivery{}).
		Where("id = ? AND is_active = ?", deliveryID, true).
		Count(&count).Error
	return count > 0, err
}
