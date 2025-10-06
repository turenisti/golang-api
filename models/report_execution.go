package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type ExecutionContext map[string]interface{}

func (ec ExecutionContext) Value() (driver.Value, error) {
	return json.Marshal(ec)
}

func (ec *ExecutionContext) Scan(value interface{}) error {
	if value == nil {
		*ec = make(ExecutionContext)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, ec)
}

type ReportExecution struct {
	ID                   string           `gorm:"primaryKey;type:varchar(36)" json:"id"`
	ConfigID             int              `gorm:"not null;index;column:config_id" json:"config_id"`
	ScheduleID           *int             `gorm:"index;column:schedule_id" json:"schedule_id"`
	Status               string           `gorm:"type:enum('running','completed','failed','cancelled');not null;default:'running';index" json:"status"`
	StartedAt            time.Time        `gorm:"default:CURRENT_TIMESTAMP;index;column:started_at" json:"started_at"`
	CompletedAt          *time.Time       `gorm:"column:completed_at" json:"completed_at"`
	ExecutedBy           string           `gorm:"size:100;not null;index;column:executed_by" json:"executed_by"`
	ExecutionContext     ExecutionContext `gorm:"type:json;column:execution_context" json:"execution_context"`
	QueryExecutionTimeMs *int             `gorm:"column:query_execution_time_ms" json:"query_execution_time_ms"`
	RowsReturned         *int             `gorm:"column:rows_returned" json:"rows_returned"`
	FileGeneratedPath    *string          `gorm:"type:text;column:file_generated_path" json:"file_generated_path"`
	FileSizeBytes        *int64           `gorm:"column:file_size_bytes" json:"file_size_bytes"`
	ErrorMessage         *string          `gorm:"type:text;column:error_message" json:"error_message"`
}

func (ReportExecution) TableName() string {
	return "report_executions"
}
