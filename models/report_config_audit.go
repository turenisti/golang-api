package models

import "time"

type ReportConfigAudit struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigID      *int      `gorm:"index" json:"config_id"`
	Action        string    `gorm:"type:enum('create','update','delete','activate','deactivate');not null;index" json:"action"`
	FieldName     *string   `gorm:"size:100" json:"field_name"`
	BeforeValue   *string   `gorm:"type:text" json:"before_value"`
	AfterValue    *string   `gorm:"type:text" json:"after_value"`
	ChangeSummary *string   `gorm:"type:json" json:"change_summary"`
	PerformedBy   string    `gorm:"size:100;not null;index" json:"performed_by"`
	PerformedAt   time.Time `gorm:"index;default:CURRENT_TIMESTAMP" json:"performed_at"`
	SessionID     *string   `gorm:"size:100" json:"session_id"`
	IPAddress     *string   `gorm:"size:45" json:"ip_address"`
}

func (ReportConfigAudit) TableName() string {
	return "report_config_audits"
}
