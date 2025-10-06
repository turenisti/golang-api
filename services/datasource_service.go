package services

import (
	"errors"
	"fmt"
	"scheduling-report/models"
	"scheduling-report/repositories"
)

type DatasourceService struct {
	repo *repository.DatasourceRepository
}

func NewDatasourceService() *DatasourceService {
	return &DatasourceService{
		repo: repository.NewDatasourceRepository(),
	}
}

// GetAll retrieves all datasources with optional filter
func (s *DatasourceService) GetAll(isActive *bool) ([]models.DataSource, error) {
	return s.repo.GetAll(isActive)
}

// GetByID retrieves a single datasource by ID
func (s *DatasourceService) GetByID(id int) (*models.DataSource, error) {
	datasource, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("datasource not found")
	}
	return datasource, nil
}

// CreateDatasourceInput defines input structure for creating datasource
type CreateDatasourceInput struct {
	Name             string                    `json:"name" validate:"required,min=3,max=100"`
	ConnectionURL    string                    `json:"connection_url" validate:"required"`
	DbType           string                    `json:"db_type" validate:"required,oneof=mysql postgresql oracle sqlserver mongodb bigquery snowflake"`
	ConnectionConfig models.ConnectionConfig   `json:"connection_config"`
	CreatedBy        string                    `json:"created_by" validate:"required"`
}

// Create creates a new datasource
func (s *DatasourceService) Create(input CreateDatasourceInput) (*models.DataSource, error) {
	// Check if name already exists
	exists, err := s.repo.CheckNameExists(input.Name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("datasource with name '%s' already exists", input.Name)
	}

	datasource := &models.DataSource{
		Name:             input.Name,
		ConnectionURL:    input.ConnectionURL,
		DbType:           input.DbType,
		ConnectionConfig: input.ConnectionConfig,
		IsActive:         true,
		CreatedBy:        input.CreatedBy,
		UpdatedBy:        input.CreatedBy,
	}

	if err := s.repo.Create(datasource); err != nil {
		return nil, err
	}

	return datasource, nil
}

// UpdateDatasourceInput defines input structure for updating datasource
type UpdateDatasourceInput struct {
	Name             string                  `json:"name" validate:"required,min=3,max=100"`
	ConnectionURL    string                  `json:"connection_url" validate:"required"`
	DbType           string                  `json:"db_type" validate:"required,oneof=mysql postgresql oracle sqlserver mongodb bigquery snowflake"`
	ConnectionConfig models.ConnectionConfig `json:"connection_config"`
	UpdatedBy        string                  `json:"updated_by" validate:"required"`
}

// Update updates an existing datasource
func (s *DatasourceService) Update(id int, input UpdateDatasourceInput) (*models.DataSource, error) {
	// Check if datasource exists
	datasource, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("datasource not found")
	}

	// Check if new name conflicts with existing datasource
	if datasource.Name != input.Name {
		exists, err := s.repo.CheckNameExists(input.Name, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("datasource with name '%s' already exists", input.Name)
		}
	}

	// Update fields
	datasource.Name = input.Name
	datasource.ConnectionURL = input.ConnectionURL
	datasource.DbType = input.DbType
	datasource.ConnectionConfig = input.ConnectionConfig
	datasource.UpdatedBy = input.UpdatedBy

	if err := s.repo.Update(datasource); err != nil {
		return nil, err
	}

	return datasource, nil
}

// Delete performs soft delete
func (s *DatasourceService) Delete(id int) error {
	// Check if datasource exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("datasource not found")
	}

	return s.repo.Delete(id)
}
