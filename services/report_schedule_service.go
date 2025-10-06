package services

import (
	"encoding/json"
	"errors"
	"scheduling-report/models"
	"scheduling-report/repositories"

	"github.com/robfig/cron/v3"
)

type ReportScheduleService struct {
	repo         *repository.ReportScheduleRepository
	auditService *ReportConfigAuditService
}

func NewReportScheduleService() *ReportScheduleService {
	return &ReportScheduleService{
		repo:         repository.NewReportScheduleRepository(),
		auditService: NewReportConfigAuditService(),
	}
}

type CreateScheduleInput struct {
	ConfigID       int     `json:"config_id" validate:"required"`
	CronExpression string  `json:"cron_expression" validate:"required,min=9"`
	Timezone       string  `json:"timezone" validate:"required"`
	CreatedBy      string  `json:"created_by"`
	SessionID      *string `json:"session_id"`
	IPAddress      *string `json:"ip_address"`
}

type UpdateScheduleInput struct {
	CronExpression string  `json:"cron_expression" validate:"required,min=9"`
	Timezone       string  `json:"timezone" validate:"required"`
	UpdatedBy      string  `json:"updated_by"`
	SessionID      *string `json:"session_id"`
	IPAddress      *string `json:"ip_address"`
}

// GetAll retrieves all schedules
func (s *ReportScheduleService) GetAll(isActive *bool) ([]models.ReportSchedule, error) {
	return s.repo.GetAll(isActive)
}

// GetByID retrieves a schedule by ID
func (s *ReportScheduleService) GetByID(id int) (*models.ReportSchedule, error) {
	schedule, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("schedule not found")
	}
	return schedule, nil
}

// GetByConfigID retrieves all schedules for a config
func (s *ReportScheduleService) GetByConfigID(configID int) ([]models.ReportSchedule, error) {
	return s.repo.GetByConfigID(configID)
}

// ValidateCronExpression validates a cron expression
func (s *ReportScheduleService) ValidateCronExpression(expression string) error {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(expression)
	if err != nil {
		return errors.New("invalid cron expression")
	}
	return nil
}

// Create creates a new schedule with automatic audit logging
func (s *ReportScheduleService) Create(input CreateScheduleInput) (*models.ReportSchedule, error) {
	// Validate cron expression
	if err := s.ValidateCronExpression(input.CronExpression); err != nil {
		return nil, err
	}

	// Check if config exists
	exists, err := s.repo.CheckConfigExists(input.ConfigID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("report config not found")
	}

	schedule := &models.ReportSchedule{
		ConfigID:       input.ConfigID,
		CronExpression: input.CronExpression,
		Timezone:       input.Timezone,
		IsActive:       true,
		CreatedBy:      input.CreatedBy,
		UpdatedBy:      input.CreatedBy,
	}

	if err := s.repo.Create(schedule); err != nil {
		return nil, err
	}

	// Auto-create audit log
	afterJSON, _ := json.Marshal(schedule)
	afterValue := string(afterJSON)
	s.auditService.CreateAuditLog(
		&schedule.ConfigID,
		"create_schedule",
		nil,
		afterValue,
		input.CreatedBy,
		input.SessionID,
		input.IPAddress,
	)

	return schedule, nil
}

// Update updates a schedule with automatic audit logging
func (s *ReportScheduleService) Update(id int, input UpdateScheduleInput) (*models.ReportSchedule, error) {
	// Validate cron expression
	if err := s.ValidateCronExpression(input.CronExpression); err != nil {
		return nil, err
	}

	existingSchedule, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("schedule not found")
	}

	// Capture before state
	beforeJSON, _ := json.Marshal(existingSchedule)
	beforeValue := string(beforeJSON)

	// Update fields
	existingSchedule.CronExpression = input.CronExpression
	existingSchedule.Timezone = input.Timezone
	existingSchedule.UpdatedBy = input.UpdatedBy

	if err := s.repo.Update(existingSchedule); err != nil {
		return nil, err
	}

	// Capture after state
	afterJSON, _ := json.Marshal(existingSchedule)
	afterValue := string(afterJSON)

	// Auto-create audit log
	s.auditService.CreateAuditLog(
		&existingSchedule.ConfigID,
		"update_schedule",
		beforeValue,
		afterValue,
		input.UpdatedBy,
		input.SessionID,
		input.IPAddress,
	)

	return existingSchedule, nil
}

// Delete soft deletes a schedule with automatic audit logging
func (s *ReportScheduleService) Delete(id int, deletedBy string, sessionID *string, ipAddress *string) error {
	existingSchedule, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("schedule not found")
	}

	// Capture before state
	beforeJSON, _ := json.Marshal(existingSchedule)
	beforeValue := string(beforeJSON)

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// Auto-create audit log
	s.auditService.CreateAuditLog(
		&existingSchedule.ConfigID,
		"delete_schedule",
		beforeValue,
		nil,
		deletedBy,
		sessionID,
		ipAddress,
	)

	return nil
}
