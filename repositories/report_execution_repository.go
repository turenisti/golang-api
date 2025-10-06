package repository

import (
	"scheduling-report/config"
	"scheduling-report/models"
	"gorm.io/gorm"
)

type ReportExecutionRepository struct {
	DB *gorm.DB
}

func NewReportExecutionRepository() *ReportExecutionRepository {
	return &ReportExecutionRepository{DB: config.DB}
}

func (r *ReportExecutionRepository) GetAll(status *string, limit int) ([]models.ReportExecution, error) {
	var executions []models.ReportExecution
	query := r.DB.Order("started_at DESC")
	
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&executions).Error
	return executions, err
}

func (r *ReportExecutionRepository) GetByID(id string) (*models.ReportExecution, error) {
	var execution models.ReportExecution
	err := r.DB.Where("id = ?", id).First(&execution).Error
	if err != nil {
		return nil, err
	}
	return &execution, nil
}

func (r *ReportExecutionRepository) GetByConfigID(configID int, limit int) ([]models.ReportExecution, error) {
	var executions []models.ReportExecution
	query := r.DB.Where("config_id = ?", configID).Order("started_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&executions).Error
	return executions, err
}
