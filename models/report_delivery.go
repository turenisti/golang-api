package models

import (
	"database/sql/driver"
	"encoding/json"
)

type DeliveryConfig map[string]interface{}

func (dc DeliveryConfig) Value() (driver.Value, error) {
	return json.Marshal(dc)
}

func (dc *DeliveryConfig) Scan(value interface{}) error {
	if value == nil {
		*dc = make(DeliveryConfig)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, dc)
}

type ReportDelivery struct {
	ID                   int            `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigID             int            `gorm:"not null;index;column:config_id" json:"config_id"`
	DeliveryName         string         `gorm:"size:200;not null;column:delivery_name" json:"delivery_name"`
	Method               string         `gorm:"type:enum('email','sftp','webhook','s3','file_share');not null;index;column:method" json:"method"`
	DeliveryConfig       DeliveryConfig `gorm:"type:json;not null;column:delivery_config" json:"delivery_config"`
	MaxRetry             int            `gorm:"not null;default:3;column:max_retry" json:"max_retry"`
	RetryIntervalMinutes int            `gorm:"not null;default:5;column:retry_interval_minutes" json:"retry_interval_minutes"`
	IsActive             bool           `gorm:"not null;default:1;index;column:is_active" json:"is_active"`
	CreatedAt        CustomTime      `gorm:"default:CURRENT_TIMESTAMP;column:created_at" json:"created_at"`
	UpdatedAt        CustomTime      `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:updated_at" json:"updated_at"`
	CreatedBy            string         `gorm:"size:100;not null;column:created_by" json:"created_by"`
	UpdatedBy            string         `gorm:"size:100;not null;column:updated_by" json:"updated_by"`
}

func (ReportDelivery) TableName() string {
	return "report_deliveries"
}
