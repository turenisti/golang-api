package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// ConnectionConfig stores additional connection parameters as JSON
type ConnectionConfig map[string]interface{}

// Value implements driver.Valuer for JSON marshaling
func (c ConnectionConfig) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	return json.Marshal(c)
}

// Scan implements sql.Scanner for JSON unmarshaling
func (c *ConnectionConfig) Scan(value interface{}) error {
	if value == nil {
		*c = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, c)
}

// DataSource matches report_datasources table schema
type DataSource struct {
	ID               int              `gorm:"primaryKey;autoIncrement" json:"id"`
	Name             string           `gorm:"size:100;not null;uniqueIndex" json:"name"`
	ConnectionURL    string           `gorm:"type:text;not null;column:connection_url" json:"connection_url"`
	DbType           string           `gorm:"type:enum('mysql','postgresql','oracle','sqlserver','mongodb','bigquery','snowflake');not null;index;column:db_type" json:"db_type"`
	ConnectionConfig ConnectionConfig `gorm:"type:json;column:connection_config" json:"connection_config"`
	IsActive         bool             `gorm:"not null;default:1;index;column:is_active" json:"is_active"`
	CreatedAt        time.Time        `gorm:"default:CURRENT_TIMESTAMP;column:created_at" json:"created_at"`
	UpdatedAt        time.Time        `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:updated_at" json:"updated_at"`
	CreatedBy        string           `gorm:"size:100;not null;column:created_by" json:"created_by"`
	UpdatedBy        string           `gorm:"size:100;not null;column:updated_by" json:"updated_by"`
}

func (DataSource) TableName() string {
	return "report_datasources"
}

// Legacy Datasource struct (for backward compatibility if needed)
type Datasource = DataSource
