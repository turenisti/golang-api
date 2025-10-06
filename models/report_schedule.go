package models


type ReportSchedule struct {
	ID             int        `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigID       int        `gorm:"not null;index;column:config_id" json:"config_id"`
	CronExpression string     `gorm:"size:100;not null;column:cron_expression" json:"cron_expression"`
	Timezone       string     `gorm:"size:50;default:'UTC';column:timezone" json:"timezone"`
	IsActive       bool       `gorm:"not null;default:1;index;column:is_active" json:"is_active"`
	LastRunAt      *CustomTime `gorm:"column:last_run_at" json:"last_run_at"`
	CreatedAt        CustomTime  `gorm:"default:CURRENT_TIMESTAMP;column:created_at" json:"created_at"`
	UpdatedAt        CustomTime  `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:updated_at" json:"updated_at"`
	CreatedBy      string     `gorm:"size:100;not null;column:created_by" json:"created_by"`
	UpdatedBy      string     `gorm:"size:100;not null;column:updated_by" json:"updated_by"`
}

func (ReportSchedule) TableName() string {
	return "report_schedules"
}
