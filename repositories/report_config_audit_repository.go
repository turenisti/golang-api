package repository

import (
	"scheduling-report/config"
	"scheduling-report/models"
)

type ReportConfigAuditRepository struct{}

func NewReportConfigAuditRepository() *ReportConfigAuditRepository {
	return &ReportConfigAuditRepository{}
}

// GetAll retrieves all audit records with optional filters
func (r *ReportConfigAuditRepository) GetAll(configID *int, action *string, performedBy *string) ([]models.ReportConfigAudit, error) {
	var audits []models.ReportConfigAudit
	query := config.DB

	if configID != nil {
		query = query.Where("config_id = ?", *configID)
	}

	if action != nil {
		query = query.Where("action = ?", *action)
	}

	if performedBy != nil {
		query = query.Where("performed_by = ?", *performedBy)
	}

	err := query.Order("performed_at DESC").Find(&audits).Error
	return audits, err
}

// GetByID retrieves a single audit record by ID
func (r *ReportConfigAuditRepository) GetByID(id int64) (*models.ReportConfigAudit, error) {
	var audit models.ReportConfigAudit
	err := config.DB.First(&audit, id).Error
	if err != nil {
		return nil, err
	}
	return &audit, nil
}

// Create inserts a new audit record
func (r *ReportConfigAuditRepository) Create(audit *models.ReportConfigAudit) error {
	return config.DB.Create(audit).Error
}

// GetByConfigID retrieves all audit records for a specific config
func (r *ReportConfigAuditRepository) GetByConfigID(configID int) ([]models.ReportConfigAudit, error) {
	var audits []models.ReportConfigAudit
	err := config.DB.Where("config_id = ?", configID).
		Order("performed_at DESC").
		Find(&audits).Error
	return audits, err
}

// GetRecentChanges retrieves recent audit records (last N days)
func (r *ReportConfigAuditRepository) GetRecentChanges(days int) ([]models.ReportConfigAudit, error) {
	var audits []models.ReportConfigAudit
	err := config.DB.Where("performed_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", days).
		Order("performed_at DESC").
		Find(&audits).Error
	return audits, err
}
