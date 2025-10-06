package services

import (
	"encoding/json"
	"errors"
	"scheduling-report/models"
	"scheduling-report/repositories"
)

type ReportDeliveryService struct {
	repo         *repository.ReportDeliveryRepository
	auditService *ReportConfigAuditService
}

func NewReportDeliveryService() *ReportDeliveryService {
	return &ReportDeliveryService{
		repo:         repository.NewReportDeliveryRepository(),
		auditService: NewReportConfigAuditService(),
	}
}

type CreateDeliveryInput struct {
	ConfigID             int                    `json:"config_id" validate:"required"`
	DeliveryName         string                 `json:"delivery_name" validate:"required,min=3"`
	Method               string                 `json:"method" validate:"required,oneof=email sftp webhook s3 file_share"`
	DeliveryConfig       map[string]interface{} `json:"delivery_config" validate:"required"`
	MaxRetry             int                    `json:"max_retry" validate:"min=0,max=10"`
	RetryIntervalMinutes int                    `json:"retry_interval_minutes" validate:"min=1,max=60"`
	CreatedBy            string                 `json:"created_by"`
	SessionID            *string                `json:"session_id"`
	IPAddress            *string                `json:"ip_address"`
}

type UpdateDeliveryInput struct {
	DeliveryName         string                 `json:"delivery_name" validate:"required,min=3"`
	Method               string                 `json:"method" validate:"required,oneof=email sftp webhook s3 file_share"`
	DeliveryConfig       map[string]interface{} `json:"delivery_config" validate:"required"`
	MaxRetry             int                    `json:"max_retry" validate:"min=0,max=10"`
	RetryIntervalMinutes int                    `json:"retry_interval_minutes" validate:"min=1,max=60"`
	UpdatedBy            string                 `json:"updated_by"`
	SessionID            *string                `json:"session_id"`
	IPAddress            *string                `json:"ip_address"`
}

func (s *ReportDeliveryService) GetAll(isActive *bool) ([]models.ReportDelivery, error) {
	return s.repo.GetAll(isActive)
}

func (s *ReportDeliveryService) GetByID(id int) (*models.ReportDelivery, error) {
	delivery, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("delivery not found")
	}
	return delivery, nil
}

func (s *ReportDeliveryService) GetByConfigID(configID int) ([]models.ReportDelivery, error) {
	return s.repo.GetByConfigID(configID)
}

func (s *ReportDeliveryService) Create(input CreateDeliveryInput) (*models.ReportDelivery, error) {
	exists, err := s.repo.CheckConfigExists(input.ConfigID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("report config not found")
	}

	delivery := &models.ReportDelivery{
		ConfigID:             input.ConfigID,
		DeliveryName:         input.DeliveryName,
		Method:               input.Method,
		DeliveryConfig:       input.DeliveryConfig,
		MaxRetry:             input.MaxRetry,
		RetryIntervalMinutes: input.RetryIntervalMinutes,
		IsActive:             true,
		CreatedBy:            input.CreatedBy,
		UpdatedBy:            input.CreatedBy,
	}

	if err := s.repo.Create(delivery); err != nil {
		return nil, err
	}

	afterJSON, _ := json.Marshal(delivery)
	s.auditService.CreateAuditLog(&delivery.ConfigID, "create_delivery", nil, string(afterJSON), input.CreatedBy, input.SessionID, input.IPAddress)

	return delivery, nil
}

func (s *ReportDeliveryService) Update(id int, input UpdateDeliveryInput) (*models.ReportDelivery, error) {
	existingDelivery, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("delivery not found")
	}

	beforeJSON, _ := json.Marshal(existingDelivery)
	existingDelivery.DeliveryName = input.DeliveryName
	existingDelivery.Method = input.Method
	existingDelivery.DeliveryConfig = input.DeliveryConfig
	existingDelivery.MaxRetry = input.MaxRetry
	existingDelivery.RetryIntervalMinutes = input.RetryIntervalMinutes
	existingDelivery.UpdatedBy = input.UpdatedBy

	if err := s.repo.Update(existingDelivery); err != nil {
		return nil, err
	}

	afterJSON, _ := json.Marshal(existingDelivery)
	s.auditService.CreateAuditLog(&existingDelivery.ConfigID, "update_delivery", string(beforeJSON), string(afterJSON), input.UpdatedBy, input.SessionID, input.IPAddress)

	return existingDelivery, nil
}

func (s *ReportDeliveryService) Delete(id int, deletedBy string, sessionID *string, ipAddress *string) error {
	existingDelivery, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("delivery not found")
	}

	beforeJSON, _ := json.Marshal(existingDelivery)
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	s.auditService.CreateAuditLog(&existingDelivery.ConfigID, "delete_delivery", string(beforeJSON), nil, deletedBy, sessionID, ipAddress)
	return nil
}
