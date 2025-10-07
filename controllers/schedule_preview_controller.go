package controllers

import (
	"fmt"
	"regexp"
	"scheduling-report/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
)

type SchedulePreviewController struct{}

func NewSchedulePreviewController() *SchedulePreviewController {
	return &SchedulePreviewController{}
}

type PreviewRequest struct {
	ReportQuery    string `json:"report_query" validate:"required"`
	CronExpression string `json:"cron_expression" validate:"required"`
}

type ExecutionPreview struct {
	ExecutionTime string            `json:"execution_time"`
	TimeRange     map[string]string `json:"time_range"`
	ExampleQuery  string            `json:"example_query"`
}

// PreviewScheduleExecution previews how schedule will execute
func (ctrl *SchedulePreviewController) PreviewScheduleExecution(c *fiber.Ctx) error {
	var input PreviewRequest

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid request body")
	}

	// Validate cron expression
	cronValidation := utils.ValidateCronExpression(input.CronExpression)
	if !cronValidation.Valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"responseCode":    "40003103",
			"responseMessage": "Invalid cron expression",
			"errors":          cronValidation.Errors,
		})
	}

	// Parse cron
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(input.CronExpression)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003103, "Failed to parse cron expression")
	}

	// Generate preview for next 3 executions
	previews := []ExecutionPreview{}
	currentTime := time.Now()
	var lastRunAt *time.Time

	for i := 0; i < 3; i++ {
		nextRun := schedule.Next(currentTime)

		// Calculate time range
		timeRange := utils.CalculateTimeRange(lastRunAt, input.CronExpression, nextRun)

		// Replace template variables in query
		exampleQuery := replaceTemplateVariables(input.ReportQuery, timeRange)

		previews = append(previews, ExecutionPreview{
			ExecutionTime: nextRun.Format("2006-01-02 15:04:05"),
			TimeRange: map[string]string{
				"start": timeRange["start_datetime"].(string),
				"end":   timeRange["end_datetime"].(string),
			},
			ExampleQuery: exampleQuery,
		})

		// Update for next iteration
		currentTime = nextRun
		lastRunAt = &nextRun
	}

	return c.JSON(fiber.Map{
		"responseCode":    "20003100",
		"responseMessage": "Schedule preview generated successfully",
		"data": fiber.Map{
			"cron_validation":  cronValidation,
			"next_executions":  previews,
		},
	})
}

// replaceTemplateVariables replaces {{variable}} placeholders in query
func replaceTemplateVariables(query string, variables map[string]interface{}) string {
	result := query

	// Find all {{variable}} patterns
	re := regexp.MustCompile(`\{\{(\w+)\}\}`)
	matches := re.FindAllStringSubmatch(query, -1)

	for _, match := range matches {
		placeholder := match[0] // {{variable}}
		varName := match[1]     // variable

		if value, exists := variables[varName]; exists {
			result = regexp.MustCompile(regexp.QuoteMeta(placeholder)).
				ReplaceAllString(result, fmt.Sprintf("'%v'", value))
		}
	}

	return result
}
