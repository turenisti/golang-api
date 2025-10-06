package repository

import (
	"scheduling-report/config"
	"scheduling-report/models"

	"gorm.io/gorm"
)

type DataSourceRepository struct {
	DB *gorm.DB
}

func NewDataSourceRepository() *DataSourceRepository {
	return &DataSourceRepository{DB: config.DB}
}

// GetAll retrieves all data sources
func (r *DataSourceRepository) GetAll(isActive *bool) ([]models.DataSource, error) {
	var dataSources []models.DataSource
	query := r.DB

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Order("created_at DESC").Find(&dataSources).Error
	return dataSources, err
}

// GetByID retrieves a single data source by ID
func (r *DataSourceRepository) GetByID(id string) (*models.DataSource, error) {
	var dataSource models.DataSource
	err := r.DB.First(&dataSource, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &dataSource, nil
}

// Create inserts a new data source
func (r *DataSourceRepository) Create(dataSource *models.DataSource) error {
	return r.DB.Create(dataSource).Error
}

// Update modifies an existing data source
func (r *DataSourceRepository) Update(dataSource *models.DataSource) error {
	return r.DB.Save(dataSource).Error
}

// SoftDelete sets is_active to false
func (r *DataSourceRepository) SoftDelete(id string) error {
	return r.DB.Model(&models.DataSource{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

// CheckUsage checks if datasource is used by any active report
func (r *DataSourceRepository) CheckUsage(id string) (int64, error) {
	var count int64
	err := r.DB.Model(&models.ReportConfig{}).
		Where("data_source_id = ? AND is_active = ?", id, true).
		Count(&count).Error
	return count, err
}
