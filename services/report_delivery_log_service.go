package services

import (
	"errors"
	"scheduling-report/models"
	"scheduling-report/repositories"
)

type ReportDeliveryLogService struct {
	repo *repository.ReportDeliveryLogRepository
}

func NewReportDeliveryLogService() *ReportDeliveryLogService {
	return &ReportDeliveryLogService{repo: repository.NewReportDeliveryLogRepository()}
}

func (s *ReportDeliveryLogService) GetAll(status *string, limit int) ([]models.ReportDeliveryLog, error) {
	return s.repo.GetAll(status, limit)
}

func (s *ReportDeliveryLogService) GetByID(id int64) (*models.ReportDeliveryLog, error) {
	log, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("delivery log not found")
	}
	return log, nil
}

func (s *ReportDeliveryLogService) GetByExecutionID(executionID string) ([]models.ReportDeliveryLog, error) {
	return s.repo.GetByExecutionID(executionID)
}

func (s *ReportDeliveryLogService) GetByDeliveryID(deliveryID int, limit int) ([]models.ReportDeliveryLog, error) {
	return s.repo.GetByDeliveryID(deliveryID, limit)
}
