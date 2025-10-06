package repository

import (
	"scheduling-report/config"
	"scheduling-report/models"

	"gorm.io/gorm"
)

type ReportConfigRepository struct {
	DB *gorm.DB
}

func NewReportConfigRepository() *ReportConfigRepository {
	return &ReportConfigRepository{DB: config.DB}
}

// GetAll retrieves all report configs with optional filters
func (r *ReportConfigRepository) GetAll(isActive *bool, datasourceID *int) ([]models.ReportConfig, error) {
	var configs []models.ReportConfig
	query := r.DB

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if datasourceID != nil {
		query = query.Where("datasource_id = ?", *datasourceID)
	}

	err := query.Order("created_at DESC").Find(&configs).Error
	return configs, err
}

// GetByID retrieves a report config by ID
func (r *ReportConfigRepository) GetByID(id int) (*models.ReportConfig, error) {
	var config models.ReportConfig
	err := r.DB.First(&config, id).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Create creates a new report config
func (r *ReportConfigRepository) Create(config *models.ReportConfig) error {
	return r.DB.Create(config).Error
}

// Update updates an existing report config
func (r *ReportConfigRepository) Update(config *models.ReportConfig) error {
	// Increment version on update
	return r.DB.Model(config).Updates(map[string]interface{}{
		"report_name":     config.ReportName,
		"report_query":    config.ReportQuery,
		"output_format":   config.OutputFormat,
		"datasource_id":   config.DatasourceID,
		"parameters":      config.Parameters,
		"timeout_seconds": config.TimeoutSeconds,
		"max_rows":        config.MaxRows,
		"updated_by":      config.UpdatedBy,
		"version":         gorm.Expr("version + 1"),
	}).Error
}

// Delete performs soft delete by setting is_active = false
func (r *ReportConfigRepository) Delete(id int) error {
	return r.DB.Model(&models.ReportConfig{}).Where("id = ?", id).Update("is_active", false).Error
}

// CheckNameExists checks if a report name already exists (excluding given ID)
func (r *ReportConfigRepository) CheckNameExists(name string, excludeID int) (bool, error) {
	var count int64
	query := r.DB.Model(&models.ReportConfig{}).Where("report_name = ?", name)

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

// GetByDatasourceID retrieves all configs for a specific datasource
func (r *ReportConfigRepository) GetByDatasourceID(datasourceID int) ([]models.ReportConfig, error) {
	var configs []models.ReportConfig
	err := r.DB.Where("datasource_id = ?", datasourceID).Find(&configs).Error
	return configs, err
}
