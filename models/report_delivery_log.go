package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type DeliveryDetails map[string]interface{}

func (dd DeliveryDetails) Value() (driver.Value, error) {
	return json.Marshal(dd)
}

func (dd *DeliveryDetails) Scan(value interface{}) error {
	if value == nil {
		*dd = make(DeliveryDetails)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, dd)
}

type ReportDeliveryLog struct {
	ID                int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigID          *int            `gorm:"index;column:config_id" json:"config_id"`
	DeliveryID        *int            `gorm:"index;column:delivery_id" json:"delivery_id"`
	ScheduleID        *int            `gorm:"column:schedule_id" json:"schedule_id"`
	ExecutionID       string          `gorm:"size:36;not null;index;column:execution_id" json:"execution_id"`
	Status            string          `gorm:"type:enum('pending','success','failed','retry');not null" json:"status"`
	SentAt            time.Time       `gorm:"default:CURRENT_TIMESTAMP;column:sent_at" json:"sent_at"`
	CompletedAt       *time.Time      `gorm:"column:completed_at" json:"completed_at"`
	RecipientCount    int             `gorm:"default:0;column:recipient_count" json:"recipient_count"`
	SuccessCount      int             `gorm:"default:0;column:success_count" json:"success_count"`
	FailureCount      int             `gorm:"default:0;column:failure_count" json:"failure_count"`
	RetryCount        int             `gorm:"not null;default:0;column:retry_count" json:"retry_count"`
	ErrorMessage      *string         `gorm:"type:text;column:error_message" json:"error_message"`
	DeliveryDetails   DeliveryDetails `gorm:"type:json;column:delivery_details" json:"delivery_details"`
	FileSizeBytes     *int64          `gorm:"column:file_size_bytes" json:"file_size_bytes"`
	ProcessingTimeMs  *int            `gorm:"column:processing_time_ms" json:"processing_time_ms"`
}

func (ReportDeliveryLog) TableName() string {
	return "report_delivery_logs"
}
