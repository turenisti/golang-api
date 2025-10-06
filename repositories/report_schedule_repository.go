package repository

import (
	"scheduling-report/config"
	"scheduling-report/models"

	"gorm.io/gorm"
)

type ReportScheduleRepository struct {
	DB *gorm.DB
}

func NewReportScheduleRepository() *ReportScheduleRepository {
	return &ReportScheduleRepository{
		DB: config.DB,
	}
}

// GetAll retrieves all active report schedules with optional filtering
func (r *ReportScheduleRepository) GetAll(isActive *bool) ([]models.ReportSchedule, error) {
	var schedules []models.ReportSchedule
	query := r.DB.Where("is_active = ?", true)

	if isActive != nil {
		query = r.DB.Where("is_active = ?", *isActive)
	}

	err := query.Order("created_at DESC").Find(&schedules).Error
	return schedules, err
}

// GetByID retrieves a single report schedule by ID
func (r *ReportScheduleRepository) GetByID(id int) (*models.ReportSchedule, error) {
	var schedule models.ReportSchedule
	err := r.DB.Where("id = ?", id).First(&schedule).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

// GetByConfigID retrieves all schedules for a specific report config
func (r *ReportScheduleRepository) GetByConfigID(configID int) ([]models.ReportSchedule, error) {
	var schedules []models.ReportSchedule
	err := r.DB.Where("config_id = ? AND is_active = ?", configID, true).Find(&schedules).Error
	return schedules, err
}

// Create inserts a new report schedule
func (r *ReportScheduleRepository) Create(schedule *models.ReportSchedule) error {
	return r.DB.Create(schedule).Error
}

// Update updates an existing report schedule
func (r *ReportScheduleRepository) Update(schedule *models.ReportSchedule) error {
	return r.DB.Save(schedule).Error
}

// Delete soft deletes a report schedule by setting is_active to false
func (r *ReportScheduleRepository) Delete(id int) error {
	return r.DB.Model(&models.ReportSchedule{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

// CheckConfigExists verifies if a report config exists
func (r *ReportScheduleRepository) CheckConfigExists(configID int) (bool, error) {
	var count int64
	err := r.DB.Model(&models.ReportConfig{}).
		Where("id = ? AND is_active = ?", configID, true).
		Count(&count).Error
	return count > 0, err
}
