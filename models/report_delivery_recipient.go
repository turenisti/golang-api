package models

import (
	"database/sql/driver"
	"encoding/json"
)

type RecipientConfig map[string]interface{}

func (rc RecipientConfig) Value() (driver.Value, error) {
	return json.Marshal(rc)
}

func (rc *RecipientConfig) Scan(value interface{}) error {
	if value == nil {
		*rc = make(RecipientConfig)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, rc)
}

type ReportDeliveryRecipient struct {
	ID              int             `gorm:"primaryKey;autoIncrement" json:"id"`
	DeliveryID      int             `gorm:"not null;index;column:delivery_id" json:"delivery_id"`
	RecipientType   string          `gorm:"size:20;not null;default:'email';column:recipient_type" json:"recipient_type"`
	RecipientValue  string          `gorm:"size:500;not null;column:recipient_value" json:"recipient_value"`
	RecipientConfig RecipientConfig `gorm:"type:json;column:recipient_config" json:"recipient_config"`
	IsActive        bool            `gorm:"not null;default:1;column:is_active" json:"is_active"`
	CreatedAt        CustomTime       `gorm:"default:CURRENT_TIMESTAMP;column:created_at" json:"created_at"`
	UpdatedAt        CustomTime       `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:updated_at" json:"updated_at"`
}

func (ReportDeliveryRecipient) TableName() string {
	return "report_delivery_recipients"
}
