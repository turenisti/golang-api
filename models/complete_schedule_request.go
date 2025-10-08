package models

import "encoding/json"

// CompleteScheduleRequest represents the full schedule creation/update request
type CompleteScheduleRequest struct {
	CronExpression string                       `json:"cron_expression" validate:"required"`
	Timezone       string                       `json:"timezone" validate:"required"`
	IsActive       bool                         `json:"is_active"`
	LastRunAt      *CustomTime                  `json:"last_run_at"`
	CreatedBy      string                       `json:"created_by" validate:"required"`
	UpdatedBy      string                       `json:"updated_by"`
	Configs        ConfigWithDeliveriesRequest  `json:"configs" validate:"required"`
}

// CompleteScheduleUpdateRequest for partial updates
type CompleteScheduleUpdateRequest struct {
	CronExpression *string                       `json:"cron_expression"`
	Timezone       *string                       `json:"timezone"`
	IsActive       *bool                         `json:"is_active"`
	LastRunAt      *CustomTime                   `json:"last_run_at"`
	UpdatedBy      string                        `json:"updated_by" validate:"required"`
	Configs        *ConfigWithDeliveriesRequest  `json:"configs"`
}

// ConfigWithDeliveriesRequest represents report config with nested deliveries for create/update
type ConfigWithDeliveriesRequest struct {
	ReportName     string                        `json:"report_name" validate:"required"`
	ReportQuery    string                        `json:"report_query" validate:"required"`
	OutputFormat   string                        `json:"output_format" validate:"required,oneof=csv xlsx json"`
	DatasourceID   int                           `json:"datasource_id" validate:"required"`
	FileName       *string                       `json:"file_name"`
	Parameters     json.RawMessage               `json:"parameters"`
	TimeoutSeconds *int                          `json:"timeout_seconds"`
	MaxRows        *int                          `json:"max_rows"`
	Deliveries     []DeliveryWithRecipientsRequest `json:"deliveries" validate:"required,min=1"`
}

// DeliveryWithRecipientsRequest represents delivery method with nested recipients for create/update
type DeliveryWithRecipientsRequest struct {
	ID                   *int              `json:"id"` // For updates - if provided, update existing
	DeliveryName         string            `json:"delivery_name" validate:"required"`
	Method               string            `json:"method" validate:"required,oneof=email sftp webhook s3 file_share"`
	MaxRetry             *int              `json:"max_retry"`
	RetryIntervalMinutes *int              `json:"retry_interval_minutes"`
	IsActive             *bool             `json:"is_active"`
	DeliveryConfig       json.RawMessage   `json:"delivery_config"`
	Recipients           []RecipientRequest `json:"recipients" validate:"required,min=1"`
}

// RecipientRequest represents recipient in nested structure for create/update
type RecipientRequest struct {
	ID             *int   `json:"id"` // For updates - if provided, update existing
	RecipientValue string `json:"recipient_value" validate:"required"`
	IsActive       *bool  `json:"is_active"`
}

// CompleteScheduleResponse represents the full schedule response
type CompleteScheduleResponse struct {
	ScheduleID     int                      `json:"schedule_id"`
	ConfigID       int                      `json:"config_id"`
	CronExpression string                   `json:"cron_expression"`
	Timezone       string                   `json:"timezone"`
	IsActive       bool                     `json:"is_active"`
	LastRunAt      *CustomTime              `json:"last_run_at"`
	NextRunAt      *CustomTime              `json:"next_run_at"`
	CreatedAt      CustomTime               `json:"created_at"`
	UpdatedAt      CustomTime               `json:"updated_at"`
	CreatedBy      string                   `json:"created_by"`
	UpdatedBy      string                   `json:"updated_by"`
	Config         ConfigResponseNested     `json:"config"`
}

// ConfigResponseNested for response
type ConfigResponseNested struct {
	ID             int                      `json:"id"`
	ReportName     string                   `json:"report_name"`
	ReportQuery    string                   `json:"report_query"`
	OutputFormat   string                   `json:"output_format"`
	DatasourceID   int                      `json:"datasource_id"`
	FileName       *string                  `json:"file_name"`
	Parameters     json.RawMessage          `json:"parameters"`
	TimeoutSeconds int                      `json:"timeout_seconds"`
	MaxRows        int                      `json:"max_rows"`
	IsActive       bool                     `json:"is_active"`
	Version        int                      `json:"version"`
	Deliveries     []DeliveryResponseNested `json:"deliveries"`
}

// DeliveryResponseNested for response
type DeliveryResponseNested struct {
	ID                   int                     `json:"id"`
	DeliveryName         string                  `json:"delivery_name"`
	Method               string                  `json:"method"`
	MaxRetry             int                     `json:"max_retry"`
	RetryIntervalMinutes int                     `json:"retry_interval_minutes"`
	IsActive             bool                    `json:"is_active"`
	DeliveryConfig       json.RawMessage         `json:"delivery_config"`
	Recipients           []RecipientResponseNested `json:"recipients"`
}

// RecipientResponseNested for response
type RecipientResponseNested struct {
	ID             int    `json:"id"`
	RecipientValue string `json:"recipient_value"`
	IsActive       bool   `json:"is_active"`
}
