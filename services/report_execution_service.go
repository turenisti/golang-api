package services

import (
	"errors"
	"time"
	"scheduling-report/models"
	"scheduling-report/repositories"
	"github.com/google/uuid"
)

type ReportExecutionService struct {
	repo         *repository.ReportExecutionRepository
	configRepo   *repository.ReportConfigRepository
	scheduleRepo *repository.ReportScheduleRepository
}

func NewReportExecutionService() *ReportExecutionService {
	return &ReportExecutionService{
		repo:         repository.NewReportExecutionRepository(),
		configRepo:   repository.NewReportConfigRepository(),
		scheduleRepo: repository.NewReportScheduleRepository(),
	}
}

func (s *ReportExecutionService) GetAll(status *string, limit int) ([]models.ReportExecution, error) {
	return s.repo.GetAll(status, limit)
}

func (s *ReportExecutionService) GetByID(id string) (*models.ReportExecution, error) {
	execution, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("execution not found")
	}
	return execution, nil
}

func (s *ReportExecutionService) GetByConfigID(configID int, limit int) ([]models.ReportExecution, error) {
	return s.repo.GetByConfigID(configID, limit)
}

// ExecuteAsync creates a queued execution and sends it to Kafka
func (s *ReportExecutionService) ExecuteAsync(configID int, scheduleID *int, executedBy string) (*models.ReportExecution, error) {
	// 1. Validate config exists
	_, err := s.configRepo.GetByID(configID)
	if err != nil {
		return nil, errors.New("report config not found")
	}

	// 2. Validate schedule if provided
	if scheduleID != nil {
		schedule, err := s.scheduleRepo.GetByID(*scheduleID)
		if err != nil {
			return nil, errors.New("schedule not found")
		}
		if schedule.ConfigID != configID {
			return nil, errors.New("schedule does not belong to the specified config")
		}
	}

	// 3. Create execution record with status 'queued'
	executionID := uuid.New().String()
	now := time.Now()

	execution := &models.ReportExecution{
		ID:         executionID,
		ConfigID:   configID,
		ScheduleID: scheduleID,
		Status:     "queued",
		StartedAt:  now,
		ExecutedBy: executedBy,
	}

	if err := s.repo.Create(execution); err != nil {
		return nil, errors.New("failed to create execution record")
	}

	// 4. Produce message to Kafka
	kafkaProducer := GetKafkaProducer()
	if kafkaProducer == nil {
		return nil, errors.New("kafka producer not initialized")
	}

	executionReq := ExecutionRequest{
		ExecutionID: executionID,
		ConfigID:    configID,
		ScheduleID:  scheduleID,
		ExecutedBy:  executedBy,
		QueuedAt:    now.Format(time.RFC3339),
	}

	if err := kafkaProducer.ProduceExecutionRequest(executionReq); err != nil {
		return nil, errors.New("failed to produce message to kafka")
	}

	// 5. Return execution with queued status
	return execution, nil
}
