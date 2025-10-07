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

// GetSchedulesWithDetails retrieves schedules with full config and delivery details
func (r *ReportScheduleRepository) GetSchedulesWithDetails(filters models.ScheduleDetailFilters) ([]models.ScheduleDetail, error) {
	var schedules []models.ReportSchedule
	query := r.DB.Model(&models.ReportSchedule{})

	// Schedule filters
	if filters.IsActive != nil {
		query = query.Where("report_schedules.is_active = ?", *filters.IsActive)
	}
	if filters.Timezone != "" {
		query = query.Where("report_schedules.timezone = ?", filters.Timezone)
	}
	if filters.ConfigID != nil {
		query = query.Where("report_schedules.config_id = ?", *filters.ConfigID)
	}
	if filters.CreatedBy != "" {
		query = query.Where("report_schedules.created_by = ?", filters.CreatedBy)
	}
	if filters.HasRun != nil {
		if *filters.HasRun {
			query = query.Where("report_schedules.last_run_at IS NOT NULL")
		} else {
			query = query.Where("report_schedules.last_run_at IS NULL")
		}
	}

	// Execute query to get schedules
	err := query.Order("report_schedules.created_at DESC").Find(&schedules).Error
	if err != nil {
		return nil, err
	}

	// Build detailed response
	var details []models.ScheduleDetail
	for _, schedule := range schedules {
		// Get config
		var config models.ReportConfig
		configQuery := r.DB.Where("id = ?", schedule.ConfigID)
		
		// Apply config filters
		if filters.ConfigIsActive != nil {
			configQuery = configQuery.Where("is_active = ?", *filters.ConfigIsActive)
		}
		if filters.DatasourceID != nil {
			configQuery = configQuery.Where("datasource_id = ?", *filters.DatasourceID)
		}
		if filters.OutputFormat != "" {
			configQuery = configQuery.Where("output_format = ?", filters.OutputFormat)
		}
		if filters.ConfigName != "" {
			configQuery = configQuery.Where("report_name LIKE ?", "%"+filters.ConfigName+"%")
		}

		err := configQuery.First(&config).Error
		if err != nil {
			// Skip schedules where config doesn't match filters
			continue
		}

		// Get deliveries
		var deliveries []models.ReportDelivery
		deliveryQuery := r.DB.Where("config_id = ?", config.ID)

		// Apply delivery filters
		if filters.DeliveryIsActive != nil {
			deliveryQuery = deliveryQuery.Where("is_active = ?", *filters.DeliveryIsActive)
		}
		if filters.DeliveryMethod != "" {
			deliveryQuery = deliveryQuery.Where("method = ?", filters.DeliveryMethod)
		}

		deliveryQuery.Find(&deliveries)

		// Build deliveries with recipients
		var deliveriesWithRecipients []models.DeliveryWithRecipients
		for _, delivery := range deliveries {
			// Get recipients for this delivery
			var recipients []models.ReportDeliveryRecipient
			r.DB.Where("delivery_id = ? AND is_active = ?", delivery.ID, true).Find(&recipients)

			// Build delivery with recipients
			deliveryWithRecipients := models.DeliveryWithRecipients{
				ID:                   delivery.ID,
				ConfigID:             delivery.ConfigID,
				DeliveryName:         delivery.DeliveryName,
				Method:               delivery.Method,
				DeliveryConfig:       delivery.DeliveryConfig,
				MaxRetry:             delivery.MaxRetry,
				Recipients:           recipients,
				RetryIntervalMinutes: delivery.RetryIntervalMinutes,
				IsActive:             delivery.IsActive,
				CreatedAt:            delivery.CreatedAt,
				UpdatedAt:            delivery.UpdatedAt,
				CreatedBy:            delivery.CreatedBy,
				UpdatedBy:            delivery.UpdatedBy,
			}

			deliveriesWithRecipients = append(deliveriesWithRecipients, deliveryWithRecipients)
		}

		// Build config with deliveries
		configWithDeliveries := models.ConfigWithDeliveries{
			ID:             config.ID,
			ReportName:     config.ReportName,
			ReportQuery:    config.ReportQuery,
			OutputFormat:   config.OutputFormat,
			DatasourceID:   config.DatasourceID,
			Parameters:     config.Parameters,
			TimeoutSeconds: config.TimeoutSeconds,
			MaxRows:        config.MaxRows,
			IsActive:       config.IsActive,
			CreatedAt:      config.CreatedAt,
			UpdatedAt:      config.UpdatedAt,
			CreatedBy:      config.CreatedBy,
			UpdatedBy:      config.UpdatedBy,
			Version:        config.Version,
			Deliveries:     deliveriesWithRecipients,
		}

		// Build schedule detail
		detail := models.ScheduleDetail{
			ID:             schedule.ID,
			Config:         &configWithDeliveries,
			CronExpression: schedule.CronExpression,
			Timezone:       schedule.Timezone,
			IsActive:       schedule.IsActive,
			LastRunAt:      schedule.LastRunAt,
			CreatedAt:      schedule.CreatedAt,
			UpdatedAt:      schedule.UpdatedAt,
			CreatedBy:      schedule.CreatedBy,
			UpdatedBy:      schedule.UpdatedBy,
		}

		details = append(details, detail)
	}

	return details, nil
}
