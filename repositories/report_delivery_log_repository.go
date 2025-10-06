package repository

import (
	"scheduling-report/config"
	"scheduling-report/models"
	"gorm.io/gorm"
)

type ReportDeliveryLogRepository struct {
	DB *gorm.DB
}

func NewReportDeliveryLogRepository() *ReportDeliveryLogRepository {
	return &ReportDeliveryLogRepository{DB: config.DB}
}

func (r *ReportDeliveryLogRepository) GetAll(status *string, limit int) ([]models.ReportDeliveryLog, error) {
	var logs []models.ReportDeliveryLog
	query := r.DB.Order("sent_at DESC")
	
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&logs).Error
	return logs, err
}

func (r *ReportDeliveryLogRepository) GetByID(id int64) (*models.ReportDeliveryLog, error) {
	var log models.ReportDeliveryLog
	err := r.DB.Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *ReportDeliveryLogRepository) GetByExecutionID(executionID string) ([]models.ReportDeliveryLog, error) {
	var logs []models.ReportDeliveryLog
	err := r.DB.Where("execution_id = ?", executionID).Order("sent_at DESC").Find(&logs).Error
	return logs, err
}

func (r *ReportDeliveryLogRepository) GetByDeliveryID(deliveryID int, limit int) ([]models.ReportDeliveryLog, error) {
	var logs []models.ReportDeliveryLog
	query := r.DB.Where("delivery_id = ?", deliveryID).Order("sent_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&logs).Error
	return logs, err
}
