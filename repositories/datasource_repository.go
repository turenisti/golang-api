package repository

import (
	"scheduling-report/config"
	"scheduling-report/models"

	"gorm.io/gorm"
)

type DatasourceRepository struct {
	DB *gorm.DB
}

func NewDatasourceRepository() *DatasourceRepository {
	return &DatasourceRepository{DB: config.DB}
}

// GetAll retrieves all datasources with optional is_active filter
func (r *DatasourceRepository) GetAll(isActive *bool) ([]models.DataSource, error) {
	var datasources []models.DataSource
	query := r.DB

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Find(&datasources).Error
	return datasources, err
}

// GetByID retrieves a datasource by ID
func (r *DatasourceRepository) GetByID(id int) (*models.DataSource, error) {
	var datasource models.DataSource
	err := r.DB.First(&datasource, id).Error
	if err != nil {
		return nil, err
	}
	return &datasource, nil
}

// Create creates a new datasource
func (r *DatasourceRepository) Create(datasource *models.DataSource) error {
	return r.DB.Create(datasource).Error
}

// Update updates an existing datasource
func (r *DatasourceRepository) Update(datasource *models.DataSource) error {
	return r.DB.Save(datasource).Error
}

// Delete performs soft delete by setting is_active = false
func (r *DatasourceRepository) Delete(id int) error {
	return r.DB.Model(&models.DataSource{}).Where("id = ?", id).Update("is_active", false).Error
}

// CheckNameExists checks if a datasource name already exists (excluding given ID)
func (r *DatasourceRepository) CheckNameExists(name string, excludeID int) (bool, error) {
	var count int64
	query := r.DB.Model(&models.DataSource{}).Where("name = ?", name)

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}
