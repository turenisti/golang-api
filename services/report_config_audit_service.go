package services

import (
	"encoding/json"
	"errors"
	"scheduling-report/models"
	"scheduling-report/repositories"
	"time"
)

type ReportConfigAuditService struct {
	repo *repository.ReportConfigAuditRepository
}

func NewReportConfigAuditService() *ReportConfigAuditService {
	return &ReportConfigAuditService{
		repo: repository.NewReportConfigAuditRepository(),
	}
}

// GetAll retrieves all audit records with optional filters
func (s *ReportConfigAuditService) GetAll(configID *int, action *string, performedBy *string) ([]models.ReportConfigAudit, error) {
	// Validate action enum if provided
	if action != nil {
		validActions := map[string]bool{
			"create":     true,
			"update":     true,
			"delete":     true,
			"activate":   true,
			"deactivate": true,
		}
		if !validActions[*action] {
			return nil, errors.New("invalid action value. Must be one of: create, update, delete, activate, deactivate")
		}
	}

	return s.repo.GetAll(configID, action, performedBy)
}

// GetByID retrieves a single audit record
func (s *ReportConfigAuditService) GetByID(id int64) (*models.ReportConfigAudit, error) {
	audit, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("audit record not found")
	}
	return audit, nil
}

// CreateAuditLog creates an audit trail entry (used internally by other services)
func (s *ReportConfigAuditService) CreateAuditLog(
	configID *int,
	action string,
	beforeValue interface{},
	afterValue interface{},
	performedBy string,
	sessionID *string,
	ipAddress *string,
) error {
	// Convert before/after values to JSON strings
	var beforeJSON, afterJSON *string

	if beforeValue != nil {
		beforeBytes, err := json.Marshal(beforeValue)
		if err == nil {
			beforeStr := string(beforeBytes)
			beforeJSON = &beforeStr
		}
	}

	if afterValue != nil {
		afterBytes, err := json.Marshal(afterValue)
		if err == nil {
			afterStr := string(afterBytes)
			afterJSON = &afterStr
		}
	}

	audit := &models.ReportConfigAudit{
		ConfigID:    configID,
		Action:      action,
		BeforeValue: beforeJSON,
		AfterValue:  afterJSON,
		PerformedBy: performedBy,
		PerformedAt: time.Now(),
		SessionID:   sessionID,
		IPAddress:   ipAddress,
	}

	return s.repo.Create(audit)
}

// CreateAuditLogWithFieldChange creates audit log for specific field changes
func (s *ReportConfigAuditService) CreateAuditLogWithFieldChange(
	configID *int,
	action string,
	fieldName string,
	beforeValue string,
	afterValue string,
	performedBy string,
	sessionID *string,
	ipAddress *string,
) error {
	audit := &models.ReportConfigAudit{
		ConfigID:    configID,
		Action:      action,
		FieldName:   &fieldName,
		BeforeValue: &beforeValue,
		AfterValue:  &afterValue,
		PerformedBy: performedBy,
		PerformedAt: time.Now(),
		SessionID:   sessionID,
		IPAddress:   ipAddress,
	}

	return s.repo.Create(audit)
}

// GetByConfigID retrieves all audits for a specific config
func (s *ReportConfigAuditService) GetByConfigID(configID int) ([]models.ReportConfigAudit, error) {
	return s.repo.GetByConfigID(configID)
}

// GetRecentChanges retrieves recent audit records
func (s *ReportConfigAuditService) GetRecentChanges(days int) ([]models.ReportConfigAudit, error) {
	if days <= 0 {
		days = 7 // Default to 7 days
	}
	return s.repo.GetRecentChanges(days)
}
