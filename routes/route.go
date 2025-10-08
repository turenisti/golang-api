package routes

import (
	"scheduling-report/controllers"
	"scheduling-report/utils"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Initialize controllers
	datasourceCtrl := controllers.NewDatasourceController()
	reportConfigCtrl := controllers.NewReportConfigController()
	scheduleCtrl := controllers.NewReportScheduleController()
	previewCtrl := controllers.NewSchedulePreviewController()
	deliveryCtrl := controllers.NewReportDeliveryController()
	recipientCtrl := controllers.NewReportDeliveryRecipientController()
	executionCtrl := controllers.NewReportExecutionController()
	deliveryLogCtrl := controllers.NewReportDeliveryLogController()
	auditCtrl := controllers.NewReportConfigAuditController()

	// API routes
	api := app.Group("/api")

	// Datasources endpoints (Phase 1)
	api.Get("/datasources", datasourceCtrl.GetDatasources)
	api.Get("/datasources/:id", datasourceCtrl.GetDatasourceByID)
	api.Post("/datasources", datasourceCtrl.CreateDatasource)
	api.Put("/datasources/:id", datasourceCtrl.UpdateDatasource)
	api.Delete("/datasources/:id", datasourceCtrl.DeleteDatasource)

	// Report Configs endpoints (Phase 2)
	api.Get("/report-configs", reportConfigCtrl.GetReportConfigs)
	api.Get("/report-configs/:id", reportConfigCtrl.GetReportConfigByID)
	api.Post("/report-configs", reportConfigCtrl.CreateReportConfig)
	api.Put("/report-configs/:id", reportConfigCtrl.UpdateReportConfig)
	api.Delete("/report-configs/:id", reportConfigCtrl.DeleteReportConfig)

	// Schedules endpoints (Phase 3)
	api.Get("/schedules", scheduleCtrl.GetSchedules)
	api.Get("/schedules/details", scheduleCtrl.GetSchedulesWithDetails) // Schedule details with full config and deliveries - MUST be before :id
	api.Post("/schedules/validate-cron", scheduleCtrl.ValidateCronExpression) // Validate cron before creating schedule
	api.Post("/schedules/preview", previewCtrl.PreviewScheduleExecution) // Preview schedule execution with query examples
	api.Get("/schedules/:id", scheduleCtrl.GetScheduleByID)
	api.Get("/schedules/config/:config_id", scheduleCtrl.GetSchedulesByConfigID)
	api.Post("/schedules", scheduleCtrl.CreateSchedule)
	api.Put("/schedules/:id", scheduleCtrl.UpdateSchedule)
	api.Delete("/schedules/:id", scheduleCtrl.DeleteSchedule)

	// Deliveries endpoints (Phase 4)
	api.Get("/deliveries", deliveryCtrl.GetDeliveries)
	api.Get("/deliveries/:id", deliveryCtrl.GetDeliveryByID)
	api.Get("/deliveries/config/:config_id", deliveryCtrl.GetDeliveriesByConfigID)
	api.Post("/deliveries", deliveryCtrl.CreateDelivery)
	api.Put("/deliveries/:id", deliveryCtrl.UpdateDelivery)
	api.Delete("/deliveries/:id", deliveryCtrl.DeleteDelivery)

	// Recipients endpoints (Phase 4)
	api.Get("/recipients", recipientCtrl.GetRecipients)
	api.Get("/recipients/:id", recipientCtrl.GetRecipientByID)
	api.Get("/recipients/delivery/:delivery_id", recipientCtrl.GetRecipientsByDeliveryID)
	api.Post("/recipients", recipientCtrl.CreateRecipient)
	api.Put("/recipients/:id", recipientCtrl.UpdateRecipient)
	api.Delete("/recipients/:id", recipientCtrl.DeleteRecipient)

	// Executions endpoints (Phase 5 - read-only + async execution)
	api.Get("/executions", executionCtrl.GetExecutions)
	api.Get("/executions/execute-async", executionCtrl.ExecuteAsync) // NEW: Async execution via Kafka - MUST be before :id
	api.Get("/executions/:id", executionCtrl.GetExecutionByID)
	api.Get("/executions/config/:config_id", executionCtrl.GetExecutionsByConfigID)

	// Delivery Logs endpoints (Phase 5 - read-only)
	api.Get("/delivery-logs", deliveryLogCtrl.GetDeliveryLogs)
	api.Get("/delivery-logs/:id", deliveryLogCtrl.GetDeliveryLogByID)
	api.Get("/delivery-logs/execution/:execution_id", deliveryLogCtrl.GetDeliveryLogsByExecutionID)
	api.Get("/delivery-logs/delivery/:delivery_id", deliveryLogCtrl.GetDeliveryLogsByDeliveryID)

	// Audit Trail endpoints (read-only - audit logs created automatically)
	api.Get("/audits", auditCtrl.GetAudits)
	api.Get("/audits/:id", auditCtrl.GetAuditByID)
	api.Get("/audits/config/:config_id", auditCtrl.GetAuditsByConfigID)
	api.Get("/audits/recent", auditCtrl.GetRecentChanges)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "scheduling-report",
		})
	})

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return utils.ErrorResponse(c, fiber.StatusNotFound, 0, "Route not found")
	})
}
