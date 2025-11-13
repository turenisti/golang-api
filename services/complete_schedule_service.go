package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"scheduling-report/config"
	"scheduling-report/models"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type CompleteScheduleService struct{}

func NewCompleteScheduleService() *CompleteScheduleService {
	return &CompleteScheduleService{}
}

// maskSensitiveFields masks sensitive fields in delivery_config for security
func maskSensitiveFields(deliveryConfig models.DeliveryConfig, method string) models.DeliveryConfig {
	// Only mask for methods that might have sensitive data
	if method == "sftp" || method == "webhook" || method == "s3" || method == "email" {
		maskedConfig := make(models.DeliveryConfig)
		for k, v := range deliveryConfig {
			// Mask sensitive field names
			if k == "password" || k == "api_key" || k == "secret" || k == "secret_key" || k == "access_key" {
				maskedConfig[k] = "***MASKED***"
			} else {
				maskedConfig[k] = v
			}
		}
		return maskedConfig
	}
	return deliveryConfig
}

// CreateComplete creates a complete schedule with config, deliveries, and recipients in a single transaction
func (s *CompleteScheduleService) CreateComplete(req models.CompleteScheduleRequest) (*models.CompleteScheduleResponse, error) {
	db := config.DB
	var response *models.CompleteScheduleResponse

	// Start transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// Step 1: Parse and prepare config data
		var parameters models.Parameters
		if req.Configs.Parameters != nil {
			if err := json.Unmarshal(req.Configs.Parameters, &parameters); err != nil {
				return fmt.Errorf("invalid parameters JSON: %w", err)
			}
		}

		timeoutSeconds := 300
		maxRows := 10000
		if req.Configs.TimeoutSeconds != nil {
			timeoutSeconds = *req.Configs.TimeoutSeconds
		}
		if req.Configs.MaxRows != nil {
			maxRows = *req.Configs.MaxRows
		}

		configModel := models.ReportConfig{
			ReportName:     req.Configs.ReportName,
			ReportQuery:    req.Configs.ReportQuery,
			OutputFormat:   req.Configs.OutputFormat,
			DatasourceID:   req.Configs.DatasourceID,
			FileName:       req.Configs.FileName,
			Parameters:     parameters,
			TimeoutSeconds: timeoutSeconds,
			MaxRows:        maxRows,
			IsActive:       true,
			CreatedAt:      models.CustomTime{Time: now},
			UpdatedAt:      models.CustomTime{Time: now},
			CreatedBy:      req.CreatedBy,
			UpdatedBy:      req.CreatedBy,
			Version:        1,
		}

		if err := tx.Create(&configModel).Error; err != nil {
			return fmt.Errorf("failed to create config: %w", err)
		}

		// Create audit trail for config creation (store full config in after_value)
		afterValueJSON, _ := json.Marshal(configModel)
		afterValueStr := string(afterValueJSON)

		auditCreate := models.ReportConfigAudit{
			ConfigID:    &configModel.ID,
			Action:      "create",
			AfterValue:  &afterValueStr,
			PerformedBy: req.CreatedBy,
			PerformedAt: now,
		}
		if err := tx.Create(&auditCreate).Error; err != nil {
			return fmt.Errorf("failed to create audit trail: %w", err)
		}

		// Step 2: Calculate next_run_at from cron expression
		nextRunAt, err := s.calculateNextRun(req.CronExpression, req.Timezone)
		if err != nil {
			return fmt.Errorf("invalid cron expression: %w", err)
		}

		// Step 3: Create schedule
		scheduleModel := models.ReportSchedule{
			ConfigID:       configModel.ID,
			CronExpression: req.CronExpression,
			Timezone:       req.Timezone,
			IsActive:       req.IsActive,
			LastRunAt:      req.LastRunAt,
			NextRunAt:      nextRunAt,
			CreatedAt:      models.CustomTime{Time: now},
			UpdatedAt:      models.CustomTime{Time: now},
			CreatedBy:      req.CreatedBy,
			UpdatedBy:      req.CreatedBy,
		}

		if err := tx.Create(&scheduleModel).Error; err != nil {
			return fmt.Errorf("failed to create schedule: %w", err)
		}

		// Step 4: Create deliveries and recipients
		deliveryResponses := []models.DeliveryResponseNested{}
		for _, deliveryReq := range req.Configs.Deliveries {
			// Parse delivery config
			var deliveryConfig models.DeliveryConfig
			if deliveryReq.DeliveryConfig != nil {
				if err := json.Unmarshal(deliveryReq.DeliveryConfig, &deliveryConfig); err != nil {
					return fmt.Errorf("invalid delivery_config JSON: %w", err)
				}
			}

			// Set defaults for delivery
			maxRetry := 3
			retryInterval := 5
			isActive := true
			if deliveryReq.MaxRetry != nil {
				maxRetry = *deliveryReq.MaxRetry
			}
			if deliveryReq.RetryIntervalMinutes != nil {
				retryInterval = *deliveryReq.RetryIntervalMinutes
			}
			if deliveryReq.IsActive != nil {
				isActive = *deliveryReq.IsActive
			}

			// Create delivery
			deliveryModel := models.ReportDelivery{
				ConfigID:             configModel.ID,
				DeliveryName:         deliveryReq.DeliveryName,
				Method:               deliveryReq.Method,
				MaxRetry:             maxRetry,
				RetryIntervalMinutes: retryInterval,
				IsActive:             isActive,
				DeliveryConfig:       deliveryConfig,
				CreatedAt:            models.CustomTime{Time: now},
				UpdatedAt:            models.CustomTime{Time: now},
				CreatedBy:            req.CreatedBy,
				UpdatedBy:            req.CreatedBy,
			}

			if err := tx.Create(&deliveryModel).Error; err != nil {
				return fmt.Errorf("failed to create delivery '%s': %w", deliveryReq.DeliveryName, err)
			}

			// Create recipients for this delivery
			recipientResponses := []models.RecipientResponseNested{}
			for _, recipientReq := range deliveryReq.Recipients {
				// Set default active
				recipientActive := true
				if recipientReq.IsActive != nil {
					recipientActive = *recipientReq.IsActive
				}

				recipientModel := models.ReportDeliveryRecipient{
					DeliveryID:     deliveryModel.ID,
					RecipientValue: recipientReq.RecipientValue,
					IsActive:       recipientActive,
					CreatedAt:      models.CustomTime{Time: now},
					UpdatedAt:      models.CustomTime{Time: now},
				}

				if err := tx.Create(&recipientModel).Error; err != nil {
					return fmt.Errorf("failed to create recipient '%s': %w", recipientReq.RecipientValue, err)
				}

				recipientResponses = append(recipientResponses, models.RecipientResponseNested{
					ID:             recipientModel.ID,
					RecipientValue: recipientModel.RecipientValue,
					IsActive:       recipientModel.IsActive,
				})
			}

			// Marshal delivery config back to json.RawMessage for response
			// Mask sensitive fields (password, api_key, secret)
			maskedConfig := maskSensitiveFields(deliveryModel.DeliveryConfig, deliveryModel.Method)
			deliveryConfigJSON, _ := json.Marshal(maskedConfig)

			deliveryResponses = append(deliveryResponses, models.DeliveryResponseNested{
				ID:                   deliveryModel.ID,
				DeliveryName:         deliveryModel.DeliveryName,
				Method:               deliveryModel.Method,
				MaxRetry:             deliveryModel.MaxRetry,
				RetryIntervalMinutes: deliveryModel.RetryIntervalMinutes,
				IsActive:             deliveryModel.IsActive,
				DeliveryConfig:       deliveryConfigJSON,
				Recipients:           recipientResponses,
			})
		}

		// Build response
		// Marshal parameters back to json.RawMessage for response
		parametersJSON, _ := json.Marshal(configModel.Parameters)

		response = &models.CompleteScheduleResponse{
			ScheduleID:     scheduleModel.ID,
			ConfigID:       configModel.ID,
			CronExpression: scheduleModel.CronExpression,
			Timezone:       scheduleModel.Timezone,
			IsActive:       scheduleModel.IsActive,
			LastRunAt:      scheduleModel.LastRunAt,
			NextRunAt:      scheduleModel.NextRunAt,
			CreatedAt:      scheduleModel.CreatedAt,
			UpdatedAt:      scheduleModel.UpdatedAt,
			CreatedBy:      scheduleModel.CreatedBy,
			UpdatedBy:      scheduleModel.UpdatedBy,
			Config: models.ConfigResponseNested{
				ID:             configModel.ID,
				ReportName:     configModel.ReportName,
				ReportQuery:    configModel.ReportQuery,
				OutputFormat:   configModel.OutputFormat,
				DatasourceID:   configModel.DatasourceID,
				FileName:       configModel.FileName,
				Parameters:     parametersJSON,
				TimeoutSeconds: configModel.TimeoutSeconds,
				MaxRows:        configModel.MaxRows,
				IsActive:       configModel.IsActive,
				Version:        configModel.Version,
				Deliveries:     deliveryResponses,
			},
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateComplete updates a complete schedule with partial update support (Option B: Flexible Update)
func (s *CompleteScheduleService) UpdateComplete(scheduleID int, req models.CompleteScheduleUpdateRequest) (*models.CompleteScheduleResponse, error) {
	db := config.DB
	var response *models.CompleteScheduleResponse

	// Start transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// Step 1: Find existing schedule
		var schedule models.ReportSchedule
		if err := tx.First(&schedule, scheduleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("schedule not found")
			}
			return fmt.Errorf("failed to find schedule: %w", err)
		}

		configID := schedule.ConfigID

		// Step 2: Update schedule fields (if provided)
		scheduleUpdates := map[string]interface{}{}
		if req.CronExpression != nil {
			scheduleUpdates["cron_expression"] = *req.CronExpression
			// Recalculate next_run_at
			nextRunAt, err := s.calculateNextRun(*req.CronExpression, schedule.Timezone)
			if err != nil {
				return fmt.Errorf("invalid cron expression: %w", err)
			}
			scheduleUpdates["next_run_at"] = nextRunAt
		}
		if req.Timezone != nil {
			scheduleUpdates["timezone"] = *req.Timezone
			// Recalculate next_run_at with new timezone
			cronExpr := schedule.CronExpression
			if req.CronExpression != nil {
				cronExpr = *req.CronExpression
			}
			nextRunAt, err := s.calculateNextRun(cronExpr, *req.Timezone)
			if err != nil {
				return fmt.Errorf("invalid timezone: %w", err)
			}
			scheduleUpdates["next_run_at"] = nextRunAt
		}
		if req.IsActive != nil {
			scheduleUpdates["is_active"] = *req.IsActive
		}
		if req.LastRunAt != nil {
			scheduleUpdates["last_run_at"] = req.LastRunAt
		}
		scheduleUpdates["updated_at"] = now
		scheduleUpdates["updated_by"] = req.UpdatedBy

		if len(scheduleUpdates) > 0 {
			if err := tx.Model(&schedule).Updates(scheduleUpdates).Error; err != nil {
				return fmt.Errorf("failed to update schedule: %w", err)
			}
		}

		// Step 3: Update config fields (if provided)
		var config models.ReportConfig
		if err := tx.First(&config, configID).Error; err != nil {
			return fmt.Errorf("failed to find config: %w", err)
		}

		if req.Configs != nil {
			// Store before state for audit (full config JSON)
			beforeValueJSON, _ := json.Marshal(config)
			beforeValueStr := string(beforeValueJSON)

			configUpdates := map[string]interface{}{}
			configUpdates["report_name"] = req.Configs.ReportName
			configUpdates["report_query"] = req.Configs.ReportQuery
			configUpdates["output_format"] = req.Configs.OutputFormat
			configUpdates["datasource_id"] = req.Configs.DatasourceID
			if req.Configs.FileName != nil {
				configUpdates["file_name"] = req.Configs.FileName
			}
			if req.Configs.Parameters != nil {
				configUpdates["parameters"] = req.Configs.Parameters
			}
			if req.Configs.TimeoutSeconds != nil {
				configUpdates["timeout_seconds"] = req.Configs.TimeoutSeconds
			}
			if req.Configs.MaxRows != nil {
				configUpdates["max_rows"] = req.Configs.MaxRows
			}
			configUpdates["updated_at"] = now
			configUpdates["updated_by"] = req.UpdatedBy
			configUpdates["version"] = gorm.Expr("version + 1")

			if err := tx.Model(&config).Updates(configUpdates).Error; err != nil {
				return fmt.Errorf("failed to update config: %w", err)
			}

			// Reload config to get updated values for after state
			tx.First(&config, configID)

			// Store after state for audit (full config JSON)
			afterValueJSON, _ := json.Marshal(config)
			afterValueStr := string(afterValueJSON)

			// Create audit trail for config update
			auditUpdate := models.ReportConfigAudit{
				ConfigID:    &configID,
				Action:      "update",
				BeforeValue: &beforeValueStr,
				AfterValue:  &afterValueStr,
				PerformedBy: req.UpdatedBy,
				PerformedAt: now,
			}
			if err := tx.Create(&auditUpdate).Error; err != nil {
				return fmt.Errorf("failed to create audit trail: %w", err)
			}
		}

		// Step 4: Handle deliveries (create/update/deactivate) - Option B: Flexible
		deliveryResponses := []models.DeliveryResponseNested{}
		if req.Configs != nil && len(req.Configs.Deliveries) > 0 {
			// Get all existing delivery IDs from request
			requestedDeliveryIDs := map[int]bool{}
			for _, deliveryReq := range req.Configs.Deliveries {
				if deliveryReq.ID != nil {
					requestedDeliveryIDs[*deliveryReq.ID] = true
				}
			}

			// Deactivate deliveries not in request
			if err := tx.Model(&models.ReportDelivery{}).
				Where("config_id = ? AND id NOT IN ?", configID, getMapKeys(requestedDeliveryIDs)).
				Updates(map[string]interface{}{
					"is_active":  false,
					"updated_at": now,
					"updated_by": req.UpdatedBy,
				}).Error; err != nil {
				return fmt.Errorf("failed to deactivate removed deliveries: %w", err)
			}

			// Process each delivery in request
			for _, deliveryReq := range req.Configs.Deliveries {
				var deliveryModel models.ReportDelivery

				if deliveryReq.ID != nil {
					// UPDATE existing delivery
					if err := tx.First(&deliveryModel, *deliveryReq.ID).Error; err != nil {
						return fmt.Errorf("delivery id %d not found: %w", *deliveryReq.ID, err)
					}

					deliveryUpdates := map[string]interface{}{
						"delivery_name": deliveryReq.DeliveryName,
						"method":        deliveryReq.Method,
						"updated_at":    now,
						"updated_by":    req.UpdatedBy,
					}
					if deliveryReq.MaxRetry != nil {
						deliveryUpdates["max_retry"] = deliveryReq.MaxRetry
					}
					if deliveryReq.RetryIntervalMinutes != nil {
						deliveryUpdates["retry_interval_minutes"] = deliveryReq.RetryIntervalMinutes
					}
					if deliveryReq.IsActive != nil {
						deliveryUpdates["is_active"] = deliveryReq.IsActive
					}
					if deliveryReq.DeliveryConfig != nil {
						deliveryUpdates["delivery_config"] = deliveryReq.DeliveryConfig
					}

					if err := tx.Model(&deliveryModel).Updates(deliveryUpdates).Error; err != nil {
						return fmt.Errorf("failed to update delivery %d: %w", *deliveryReq.ID, err)
					}
				} else {
					// CREATE new delivery
					// Parse delivery config
					var deliveryConfig models.DeliveryConfig
					if deliveryReq.DeliveryConfig != nil {
						if err := json.Unmarshal(deliveryReq.DeliveryConfig, &deliveryConfig); err != nil {
							return fmt.Errorf("invalid delivery_config JSON: %w", err)
						}
					}

					// Set defaults
					maxRetry := 3
					retryInterval := 5
					isActive := true
					if deliveryReq.MaxRetry != nil {
						maxRetry = *deliveryReq.MaxRetry
					}
					if deliveryReq.RetryIntervalMinutes != nil {
						retryInterval = *deliveryReq.RetryIntervalMinutes
					}
					if deliveryReq.IsActive != nil {
						isActive = *deliveryReq.IsActive
					}

					deliveryModel = models.ReportDelivery{
						ConfigID:             configID,
						DeliveryName:         deliveryReq.DeliveryName,
						Method:               deliveryReq.Method,
						MaxRetry:             maxRetry,
						RetryIntervalMinutes: retryInterval,
						IsActive:             isActive,
						DeliveryConfig:       deliveryConfig,
						CreatedAt:            models.CustomTime{Time: now},
						UpdatedAt:            models.CustomTime{Time: now},
						CreatedBy:            req.UpdatedBy,
						UpdatedBy:            req.UpdatedBy,
					}

					if err := tx.Create(&deliveryModel).Error; err != nil {
						return fmt.Errorf("failed to create delivery: %w", err)
					}
				}

				// Handle recipients for this delivery
				recipientResponses := []models.RecipientResponseNested{}
				if len(deliveryReq.Recipients) > 0 {
					// Get requested recipient IDs
					requestedRecipientIDs := map[int]bool{}
					for _, recipientReq := range deliveryReq.Recipients {
						if recipientReq.ID != nil {
							requestedRecipientIDs[*recipientReq.ID] = true
						}
					}

					// Hard delete recipients not in request
					if err := tx.Where("delivery_id = ? AND id NOT IN ?", deliveryModel.ID, getMapKeys(requestedRecipientIDs)).
						Delete(&models.ReportDeliveryRecipient{}).Error; err != nil {
						return fmt.Errorf("failed to delete removed recipients: %w", err)
					}

					// Process each recipient
					for _, recipientReq := range deliveryReq.Recipients {
						var recipientModel models.ReportDeliveryRecipient

						if recipientReq.ID != nil {
							// UPDATE existing recipient
							if err := tx.First(&recipientModel, *recipientReq.ID).Error; err != nil {
								return fmt.Errorf("recipient id %d not found: %w", *recipientReq.ID, err)
							}

							recipientUpdates := map[string]interface{}{
								"recipient_value": recipientReq.RecipientValue,
								"updated_at":      now,
							}
							if recipientReq.IsActive != nil {
								recipientUpdates["is_active"] = recipientReq.IsActive
							}

							if err := tx.Model(&recipientModel).Updates(recipientUpdates).Error; err != nil {
								return fmt.Errorf("failed to update recipient: %w", err)
							}
						} else {
							// CREATE new recipient
							recipientActive := true
							if recipientReq.IsActive != nil {
								recipientActive = *recipientReq.IsActive
							}

							recipientModel = models.ReportDeliveryRecipient{
								DeliveryID:     deliveryModel.ID,
								RecipientValue: recipientReq.RecipientValue,
								IsActive:       recipientActive,
								CreatedAt:      models.CustomTime{Time: now},
								UpdatedAt:      models.CustomTime{Time: now},
							}

							if err := tx.Create(&recipientModel).Error; err != nil {
								return fmt.Errorf("failed to create recipient: %w", err)
							}
						}

						recipientResponses = append(recipientResponses, models.RecipientResponseNested{
							ID:             recipientModel.ID,
							RecipientValue: recipientModel.RecipientValue,
							IsActive:       recipientModel.IsActive,
						})
					}
				}

				// Marshal delivery config back to json.RawMessage for response
				maskedConfig := maskSensitiveFields(deliveryModel.DeliveryConfig, deliveryModel.Method)
				deliveryConfigJSON, _ := json.Marshal(maskedConfig)

				deliveryResponses = append(deliveryResponses, models.DeliveryResponseNested{
					ID:                   deliveryModel.ID,
					DeliveryName:         deliveryModel.DeliveryName,
					Method:               deliveryModel.Method,
					MaxRetry:             deliveryModel.MaxRetry,
					RetryIntervalMinutes: deliveryModel.RetryIntervalMinutes,
					IsActive:             deliveryModel.IsActive,
					DeliveryConfig:       deliveryConfigJSON,
					Recipients:           recipientResponses,
				})
			}
		} else {
			// If no deliveries provided in update, fetch existing ones
			var existingDeliveries []models.ReportDelivery
			if err := tx.Where("config_id = ? AND is_active = ?", configID, true).Find(&existingDeliveries).Error; err != nil {
				return fmt.Errorf("failed to fetch deliveries: %w", err)
			}

			for _, delivery := range existingDeliveries {
				var recipients []models.ReportDeliveryRecipient
				if err := tx.Where("delivery_id = ? AND is_active = ?", delivery.ID, true).Find(&recipients).Error; err != nil {
					return fmt.Errorf("failed to fetch recipients: %w", err)
				}

				recipientResponses := []models.RecipientResponseNested{}
				for _, recipient := range recipients {
					recipientResponses = append(recipientResponses, models.RecipientResponseNested{
						ID:             recipient.ID,
						RecipientValue: recipient.RecipientValue,
						IsActive:       recipient.IsActive,
					})
				}

				// Marshal delivery config back to json.RawMessage for response
				maskedConfig := maskSensitiveFields(delivery.DeliveryConfig, delivery.Method)
				deliveryConfigJSON, _ := json.Marshal(maskedConfig)

				deliveryResponses = append(deliveryResponses, models.DeliveryResponseNested{
					ID:                   delivery.ID,
					DeliveryName:         delivery.DeliveryName,
					Method:               delivery.Method,
					MaxRetry:             delivery.MaxRetry,
					RetryIntervalMinutes: delivery.RetryIntervalMinutes,
					IsActive:             delivery.IsActive,
					DeliveryConfig:       deliveryConfigJSON,
					Recipients:           recipientResponses,
				})
			}
		}

		// Reload updated models
		tx.First(&schedule, scheduleID)
		tx.First(&config, configID)

		// Build response
		// Marshal parameters back to json.RawMessage for response
		parametersJSON, _ := json.Marshal(config.Parameters)

		response = &models.CompleteScheduleResponse{
			ScheduleID:     schedule.ID,
			ConfigID:       config.ID,
			CronExpression: schedule.CronExpression,
			Timezone:       schedule.Timezone,
			IsActive:       schedule.IsActive,
			LastRunAt:      schedule.LastRunAt,
			NextRunAt:      schedule.NextRunAt,
			CreatedAt:      schedule.CreatedAt,
			UpdatedAt:      schedule.UpdatedAt,
			CreatedBy:      schedule.CreatedBy,
			UpdatedBy:      schedule.UpdatedBy,
			Config: models.ConfigResponseNested{
				ID:             config.ID,
				ReportName:     config.ReportName,
				ReportQuery:    config.ReportQuery,
				OutputFormat:   config.OutputFormat,
				DatasourceID:   config.DatasourceID,
				FileName:       config.FileName,
				Parameters:     parametersJSON,
				TimeoutSeconds: config.TimeoutSeconds,
				MaxRows:        config.MaxRows,
				IsActive:       config.IsActive,
				Version:        config.Version,
				Deliveries:     deliveryResponses,
			},
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

// calculateNextRun calculates the next run time from cron expression and timezone
func (s *CompleteScheduleService) calculateNextRun(cronExpression string, timezone string) (*models.CustomTime, error) {
	// Parse timezone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone '%s': %w", timezone, err)
	}

	// Parse cron expression
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(cronExpression)
	if err != nil {
		return nil, fmt.Errorf("invalid cron expression '%s': %w", cronExpression, err)
	}

	// Calculate next run from now
	now := time.Now().In(loc)
	nextRun := schedule.Next(now)

	return &models.CustomTime{Time: nextRun}, nil
}

// Helper function to get map keys as slice (for NOT IN query)
func getMapKeys(m map[int]bool) []int {
	if len(m) == 0 {
		return []int{-1} // Return impossible ID to avoid SQL error with empty IN clause
	}
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
