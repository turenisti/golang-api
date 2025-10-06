package services

import (
	"errors"
	"fmt"
	"scheduling-report/models"
	"scheduling-report/repositories"
)

type ReportConfigService struct {
	repo            *repository.ReportConfigRepository
	datasourceRepo  *repository.DatasourceRepository
	auditService    *ReportConfigAuditService
}

func NewReportConfigService() *ReportConfigService {
	return &ReportConfigService{
		repo:           repository.NewReportConfigRepository(),
		datasourceRepo: repository.NewDatasourceRepository(),
		auditService:   NewReportConfigAuditService(),
	}
}

// GetAll retrieves all report configs with optional filters
func (s *ReportConfigService) GetAll(isActive *bool, datasourceID *int) ([]models.ReportConfig, error) {
	return s.repo.GetAll(isActive, datasourceID)
}

// GetByID retrieves a single report config by ID
func (s *ReportConfigService) GetByID(id int) (*models.ReportConfig, error) {
	config, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("report config not found")
	}
	return config, nil
}

// CreateReportConfigInput defines input structure for creating report config
type CreateReportConfigInput struct {
	ReportName     string             `json:"report_name" validate:"required,min=3,max=200"`
	ReportQuery    string             `json:"report_query" validate:"required"`
	OutputFormat   string             `json:"output_format" validate:"required,oneof=csv xlsx pdf json"`
	DatasourceID   int                `json:"datasource_id" validate:"required"`
	Parameters     models.Parameters  `json:"parameters"`
	TimeoutSeconds int                `json:"timeout_seconds" validate:"min=1,max=3600"`
	MaxRows        int                `json:"max_rows" validate:"min=1,max=1000000"`
	CreatedBy      string             `json:"created_by" validate:"required"`
	IPAddress      *string            `json:"-"` // For audit
	SessionID      *string            `json:"-"` // For audit
}

// Create creates a new report config with audit logging
func (s *ReportConfigService) Create(input CreateReportConfigInput) (*models.ReportConfig, error) {
	// Validate datasource exists and is active
	datasource, err := s.datasourceRepo.GetByID(input.DatasourceID)
	if err != nil {
		return nil, errors.New("datasource not found")
	}
	if !datasource.IsActive {
		return nil, errors.New("datasource is not active")
	}

	// Check if name already exists
	exists, err := s.repo.CheckNameExists(input.ReportName, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("report config with name '%s' already exists", input.ReportName)
	}

	// Set defaults if not provided
	if input.TimeoutSeconds == 0 {
		input.TimeoutSeconds = 300
	}
	if input.MaxRows == 0 {
		input.MaxRows = 10000
	}

	config := &models.ReportConfig{
		ReportName:     input.ReportName,
		ReportQuery:    input.ReportQuery,
		OutputFormat:   input.OutputFormat,
		DatasourceID:   input.DatasourceID,
		Parameters:     input.Parameters,
		TimeoutSeconds: input.TimeoutSeconds,
		MaxRows:        input.MaxRows,
		IsActive:       true,
		CreatedBy:      input.CreatedBy,
		UpdatedBy:      input.CreatedBy,
		Version:        1,
	}

	if err := s.repo.Create(config); err != nil {
		return nil, err
	}

	// ✅ AUTO-CREATE AUDIT LOG
	s.auditService.CreateAuditLog(
		&config.ID,
		"create",
		nil,           // before_value (null for create)
		config,        // after_value (full config)
		input.CreatedBy,
		input.SessionID,
		input.IPAddress,
	)

	return config, nil
}

// UpdateReportConfigInput defines input structure for updating report config
type UpdateReportConfigInput struct {
	ReportName     string            `json:"report_name" validate:"required,min=3,max=200"`
	ReportQuery    string            `json:"report_query" validate:"required"`
	OutputFormat   string            `json:"output_format" validate:"required,oneof=csv xlsx pdf json"`
	DatasourceID   int               `json:"datasource_id" validate:"required"`
	Parameters     models.Parameters `json:"parameters"`
	TimeoutSeconds int               `json:"timeout_seconds" validate:"min=1,max=3600"`
	MaxRows        int               `json:"max_rows" validate:"min=1,max=1000000"`
	UpdatedBy      string            `json:"updated_by" validate:"required"`
	IPAddress      *string           `json:"-"` // For audit
	SessionID      *string           `json:"-"` // For audit
}

// Update updates an existing report config with audit logging
func (s *ReportConfigService) Update(id int, input UpdateReportConfigInput) (*models.ReportConfig, error) {
	// Get existing config (for audit before_value)
	existingConfig, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("report config not found")
	}

	// Validate datasource exists and is active
	datasource, err := s.datasourceRepo.GetByID(input.DatasourceID)
	if err != nil {
		return nil, errors.New("datasource not found")
	}
	if !datasource.IsActive {
		return nil, errors.New("datasource is not active")
	}

	// Check if new name conflicts with existing config
	if existingConfig.ReportName != input.ReportName {
		exists, err := s.repo.CheckNameExists(input.ReportName, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("report config with name '%s' already exists", input.ReportName)
		}
	}

	// Update fields
	existingConfig.ReportName = input.ReportName
	existingConfig.ReportQuery = input.ReportQuery
	existingConfig.OutputFormat = input.OutputFormat
	existingConfig.DatasourceID = input.DatasourceID
	existingConfig.Parameters = input.Parameters
	existingConfig.TimeoutSeconds = input.TimeoutSeconds
	existingConfig.MaxRows = input.MaxRows
	existingConfig.UpdatedBy = input.UpdatedBy

	if err := s.repo.Update(existingConfig); err != nil {
		return nil, err
	}

	// Reload to get incremented version
	updatedConfig, _ := s.repo.GetByID(id)

	// ✅ AUTO-CREATE AUDIT LOG
	s.auditService.CreateAuditLog(
		&updatedConfig.ID,
		"update",
		existingConfig,    // before_value (old config before update)
		updatedConfig,     // after_value (new config after update)
		input.UpdatedBy,
		input.SessionID,
		input.IPAddress,
	)

	return updatedConfig, nil
}

// Delete performs soft delete with audit logging
func (s *ReportConfigService) Delete(id int, deletedBy string, sessionID *string, ipAddress *string) error {
	// Get existing config for audit
	existingConfig, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("report config not found")
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// ✅ AUTO-CREATE AUDIT LOG
	s.auditService.CreateAuditLog(
		&id,
		"delete",
		existingConfig,    // before_value (config before deletion)
		nil,              // after_value (null after deletion)
		deletedBy,
		sessionID,
		ipAddress,
	)

	return nil
}

// ToggleActive activates or deactivates a config with audit logging
func (s *ReportConfigService) ToggleActive(id int, isActive bool, updatedBy string, sessionID *string, ipAddress *string) error {
	existingConfig, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("report config not found")
	}

	action := "deactivate"
	if isActive {
		action = "activate"
	}

	// Update is_active status
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// ✅ AUTO-CREATE AUDIT LOG
	s.auditService.CreateAuditLogWithFieldChange(
		&id,
		action,
		"is_active",
		fmt.Sprintf("%t", existingConfig.IsActive),
		fmt.Sprintf("%t", isActive),
		updatedBy,
		sessionID,
		ipAddress,
	)

	return nil
}
