package models

import (
	"database/sql/driver"
	"encoding/json"
)

// Parameters stores report parameters as JSON
type Parameters map[string]interface{}

// Value implements driver.Valuer for JSON marshaling
func (p Parameters) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return json.Marshal(p)
}

// Scan implements sql.Scanner for JSON unmarshaling
func (p *Parameters) Scan(value interface{}) error {
	if value == nil {
		*p = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, p)
}

// ReportConfig matches report_configs table schema
type ReportConfig struct {
	ID             int        `gorm:"primaryKey;autoIncrement" json:"id"`
	ReportName     string     `gorm:"size:200;not null;index;column:report_name" json:"report_name"`
	ReportQuery    string     `gorm:"type:text;not null;column:report_query" json:"report_query"`
	OutputFormat   string     `gorm:"size:50;not null;default:'csv';column:output_format" json:"output_format"`
	DatasourceID   int        `gorm:"not null;index;column:datasource_id" json:"datasource_id"`
	Parameters     Parameters `gorm:"type:json;column:parameters" json:"parameters"`
	TimeoutSeconds int        `gorm:"default:300;column:timeout_seconds" json:"timeout_seconds"`
	MaxRows        int        `gorm:"default:10000;column:max_rows" json:"max_rows"`
	IsActive       bool       `gorm:"not null;default:1;index;column:is_active" json:"is_active"`
	CreatedAt        CustomTime  `gorm:"default:CURRENT_TIMESTAMP;index;column:created_at" json:"created_at"`
	UpdatedAt        CustomTime  `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:updated_at" json:"updated_at"`
	CreatedBy      string     `gorm:"size:100;not null;column:created_by" json:"created_by"`
	UpdatedBy      string     `gorm:"size:100;not null;column:updated_by" json:"updated_by"`
	Version        int        `gorm:"not null;default:1;column:version" json:"version"`
}

func (ReportConfig) TableName() string {
	return "report_configs"
}
