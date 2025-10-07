package utils

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// CronValidation holds validation results for cron expressions
type CronValidation struct {
	Valid            bool     `json:"valid"`
	IntervalMinutes  int      `json:"interval_minutes"`
	ExecutionsPerDay float64  `json:"executions_per_day"`
	NextExecutions   []string `json:"next_executions"`
	Warnings         []string `json:"warnings"`
	Errors           []string `json:"errors"`
}

const (
	// MinimumIntervalMinutes prevents too-frequent schedules
	MinimumIntervalMinutes = 5
	// MaximumExecutionsPerDay limit (every 5 minutes = 288/day)
	MaximumExecutionsPerDay = 288
)

// ValidateCronExpression validates cron expression and calculates interval
func ValidateCronExpression(cronExpr string) CronValidation {
	result := CronValidation{
		Valid:    true,
		Warnings: []string{},
		Errors:   []string{},
	}

	// Parse cron expression
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(cronExpr)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Invalid cron syntax: %v", err))
		return result
	}

	// Calculate interval between executions
	now := time.Now()
	next1 := schedule.Next(now)
	next2 := schedule.Next(next1)
	next3 := schedule.Next(next2)
	next4 := schedule.Next(next3)
	next5 := schedule.Next(next4)

	interval := next2.Sub(next1)
	intervalMinutes := int(interval.Minutes())

	result.IntervalMinutes = intervalMinutes
	result.ExecutionsPerDay = float64(1440) / float64(intervalMinutes)

	// Get next 5 executions for preview
	result.NextExecutions = []string{
		next1.Format("2006-01-02 15:04:05"),
		next2.Format("2006-01-02 15:04:05"),
		next3.Format("2006-01-02 15:04:05"),
		next4.Format("2006-01-02 15:04:05"),
		next5.Format("2006-01-02 15:04:05"),
	}

	// Validation: Check minimum interval
	if intervalMinutes < MinimumIntervalMinutes {
		result.Valid = false
		result.Errors = append(result.Errors,
			fmt.Sprintf("Interval too short (%d minutes). Minimum allowed: %d minutes to prevent database overload",
				intervalMinutes, MinimumIntervalMinutes))
		return result
	}

	// Warning for high-frequency schedules
	if intervalMinutes < 15 {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("High frequency: %.0f executions per day. Please ensure your query is optimized and database can handle this load.",
				result.ExecutionsPerDay))
	}

	// Warning for very high frequency
	if intervalMinutes < 10 {
		result.Warnings = append(result.Warnings,
			"Very high frequency schedule detected. Monitor database performance closely.")
	}

	return result
}

// CalculateTimeRange calculates the time range for query based on interval
func CalculateTimeRange(lastRunAt *time.Time, cronExpr string, executionTime time.Time) map[string]interface{} {
	var startTime time.Time

	if lastRunAt != nil {
		// Use last_run_at for accurate time range
		startTime = *lastRunAt
	} else {
		// First run: calculate from cron expression
		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		schedule, err := parser.Parse(cronExpr)
		if err != nil {
			// Fallback to 24 hours ago
			startTime = executionTime.Add(-24 * time.Hour)
		} else {
			// Get previous scheduled time
			startTime = schedule.Next(executionTime.Add(-7 * 24 * time.Hour))
			for schedule.Next(startTime).Before(executionTime) {
				startTime = schedule.Next(startTime)
			}
		}
	}

	intervalHours := executionTime.Sub(startTime).Hours()

	return map[string]interface{}{
		"start_datetime":      startTime.Format("2006-01-02 15:04:05"),
		"end_datetime":        executionTime.Format("2006-01-02 15:04:05"),
		"start_date":          startTime.Format("2006-01-02"),
		"end_date":            executionTime.Format("2006-01-02"),
		"interval_hours":      fmt.Sprintf("%.2f", intervalHours),
		"calculation_method":  getCalculationMethod(lastRunAt),
		"yesterday":           executionTime.Add(-24 * time.Hour).Format("2006-01-02"),
		"last_week":           executionTime.Add(-7 * 24 * time.Hour).Format("2006-01-02"),
		"execution_time":      executionTime.Format("2006-01-02 15:04:05"),
	}
}

func getCalculationMethod(lastRunAt *time.Time) string {
	if lastRunAt != nil {
		return "last_run_at"
	}
	return "cron_detection"
}
