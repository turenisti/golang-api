package services

import (
	"errors"
	"scheduling-report/models"
	"scheduling-report/repositories"
)

type ReportExecutionService struct {
	repo *repository.ReportExecutionRepository
}

func NewReportExecutionService() *ReportExecutionService {
	return &ReportExecutionService{repo: repository.NewReportExecutionRepository()}
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
