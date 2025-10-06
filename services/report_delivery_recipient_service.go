package services

import (
	"errors"
	"scheduling-report/models"
	"scheduling-report/repositories"
)

type ReportDeliveryRecipientService struct {
	repo *repository.ReportDeliveryRecipientRepository
}

func NewReportDeliveryRecipientService() *ReportDeliveryRecipientService {
	return &ReportDeliveryRecipientService{
		repo: repository.NewReportDeliveryRecipientRepository(),
	}
}

type CreateRecipientInput struct {
	DeliveryID      int                    `json:"delivery_id" validate:"required"`
	RecipientType   string                 `json:"recipient_type" validate:"required"`
	RecipientValue  string                 `json:"recipient_value" validate:"required"`
	RecipientConfig map[string]interface{} `json:"recipient_config"`
}

type UpdateRecipientInput struct {
	RecipientType   string                 `json:"recipient_type" validate:"required"`
	RecipientValue  string                 `json:"recipient_value" validate:"required"`
	RecipientConfig map[string]interface{} `json:"recipient_config"`
}

func (s *ReportDeliveryRecipientService) GetAll(isActive *bool) ([]models.ReportDeliveryRecipient, error) {
	return s.repo.GetAll(isActive)
}

func (s *ReportDeliveryRecipientService) GetByID(id int) (*models.ReportDeliveryRecipient, error) {
	recipient, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("recipient not found")
	}
	return recipient, nil
}

func (s *ReportDeliveryRecipientService) GetByDeliveryID(deliveryID int) ([]models.ReportDeliveryRecipient, error) {
	return s.repo.GetByDeliveryID(deliveryID)
}

func (s *ReportDeliveryRecipientService) Create(input CreateRecipientInput) (*models.ReportDeliveryRecipient, error) {
	exists, err := s.repo.CheckDeliveryExists(input.DeliveryID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("delivery not found")
	}

	recipient := &models.ReportDeliveryRecipient{
		DeliveryID:      input.DeliveryID,
		RecipientType:   input.RecipientType,
		RecipientValue:  input.RecipientValue,
		RecipientConfig: input.RecipientConfig,
		IsActive:        true,
	}

	if err := s.repo.Create(recipient); err != nil {
		return nil, err
	}

	return recipient, nil
}

func (s *ReportDeliveryRecipientService) Update(id int, input UpdateRecipientInput) (*models.ReportDeliveryRecipient, error) {
	existingRecipient, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("recipient not found")
	}

	existingRecipient.RecipientType = input.RecipientType
	existingRecipient.RecipientValue = input.RecipientValue
	existingRecipient.RecipientConfig = input.RecipientConfig

	if err := s.repo.Update(existingRecipient); err != nil {
		return nil, err
	}

	return existingRecipient, nil
}

func (s *ReportDeliveryRecipientService) Delete(id int) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("recipient not found")
	}

	return s.repo.Delete(id)
}
