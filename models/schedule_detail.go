package models

// ScheduleDetail represents a schedule with full config and delivery details
type ScheduleDetail struct {
	ID             int                   `json:"id"`
	Config         *ConfigWithDeliveries `json:"configs"` // Named "configs" to match your requirement
	CronExpression string                `json:"cron_expression"`
	Timezone       string                `json:"timezone"`
	IsActive       bool                  `json:"is_active"`
	LastRunAt      *CustomTime           `json:"last_run_at"`
	CreatedAt      CustomTime            `json:"created_at"`
	UpdatedAt      CustomTime            `json:"updated_at"`
	CreatedBy      string                `json:"created_by"`
	UpdatedBy      string                `json:"updated_by"`
}

// ConfigWithDeliveries represents a report config with its deliveries
type ConfigWithDeliveries struct {
	ID             int                      `json:"id"`
	ReportName     string                   `json:"report_name"`
	ReportQuery    string                   `json:"report_query"`
	OutputFormat   string                   `json:"output_format"`
	DatasourceID   int                      `json:"datasource_id"`
	Parameters     Parameters               `json:"parameters"`
	TimeoutSeconds int                      `json:"timeout_seconds"`
	MaxRows        int                      `json:"max_rows"`
	IsActive       bool                     `json:"is_active"`
	CreatedAt      CustomTime               `json:"created_at"`
	UpdatedAt      CustomTime               `json:"updated_at"`
	CreatedBy      string                   `json:"created_by"`
	UpdatedBy      string                   `json:"updated_by"`
	Version        int                      `json:"version"`
	Deliveries     []DeliveryWithRecipients `json:"deliveries"`
}

// DeliveryWithRecipients represents a delivery with its recipients
type DeliveryWithRecipients struct {
	ID                   int                       `json:"id"`
	ConfigID             int                       `json:"config_id"`
	DeliveryName         string                    `json:"delivery_name"`
	Method               string                    `json:"method"`
	DeliveryConfig       DeliveryConfig            `json:"delivery_config"`
	MaxRetry             int                       `json:"max_retry"`
	Recipients           []ReportDeliveryRecipient `json:"recipients"`
	RetryIntervalMinutes int                       `json:"retry_interval_minutes"`
	IsActive             bool                      `json:"is_active"`
	CreatedAt            CustomTime                `json:"created_at"`
	UpdatedAt            CustomTime                `json:"updated_at"`
	CreatedBy            string                    `json:"created_by"`
	UpdatedBy            string                    `json:"updated_by"`
}

// ScheduleDetailFilters represents query filters for schedules/details endpoint
type ScheduleDetailFilters struct {
	IsActive           *bool   `query:"is_active"`             // Schedule active status
	Timezone           string  `query:"timezone"`              // Schedule timezone
	ConfigID           *int    `query:"config_id"`             // Specific config ID
	CreatedBy          string  `query:"created_by"`            // Schedule creator
	ConfigIsActive     *bool   `query:"config_is_active"`      // Config active status
	DatasourceID       *int    `query:"datasource_id"`         // Filter by datasource
	OutputFormat       string  `query:"output_format"`         // csv, excel, json, pdf
	ConfigName         string  `query:"config_name"`           // Partial match on report_name
	DeliveryIsActive   *bool   `query:"delivery_is_active"`    // Delivery active status
	DeliveryMethod     string  `query:"delivery_method"`       // email, sftp, webhook, s3, file_share
	HasRun             *bool   `query:"has_run"`               // Filter if last_run_at is not null
}
