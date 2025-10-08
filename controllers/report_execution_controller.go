package controllers

import (
	"scheduling-report/services"
	"scheduling-report/utils"
	"strconv"
	"github.com/gofiber/fiber/v2"
)

type ReportExecutionController struct {
	service *services.ReportExecutionService
}

func NewReportExecutionController() *ReportExecutionController {
	return &ReportExecutionController{service: services.NewReportExecutionService()}
}

func (ctrl *ReportExecutionController) GetExecutions(c *fiber.Ctx) error {
	var status *string
	if c.Query("status") != "" {
		s := c.Query("status")
		status = &s
	}
	limit, _ := strconv.Atoi(c.Query("limit", "100"))
	
	executions, err := ctrl.service.GetAll(status, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}
	return utils.SuccessResponse(c, executions, "Executions retrieved successfully")
}

func (ctrl *ReportExecutionController) GetExecutionByID(c *fiber.Ctx) error {
	id := c.Params("id")
	execution, err := ctrl.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, 40403100, err.Error())
	}
	return utils.SuccessResponse(c, execution, "Execution retrieved successfully")
}

func (ctrl *ReportExecutionController) GetExecutionsByConfigID(c *fiber.Ctx) error {
	configID, err := strconv.Atoi(c.Params("config_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid config ID")
	}
	limit, _ := strconv.Atoi(c.Query("limit", "100"))

	executions, err := ctrl.service.GetByConfigID(configID, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}
	return utils.SuccessResponse(c, executions, "Executions retrieved successfully")
}

func (ctrl *ReportExecutionController) ExecuteAsync(c *fiber.Ctx) error {
	// Get config_id from query params
	configID, err := strconv.Atoi(c.Query("config_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003102, "Invalid or missing config_id parameter")
	}

	// Get optional schedule_id
	var scheduleID *int
	if scheduleIDStr := c.Query("schedule_id"); scheduleIDStr != "" {
		id, err := strconv.Atoi(scheduleIDStr)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003103, "Invalid schedule_id parameter")
		}
		scheduleID = &id
	}

	// Get executed_by from header or default to "system"
	executedBy := c.Get("X-User-ID", "system")

	// Execute async
	execution, err := ctrl.service.ExecuteAsync(configID, scheduleID, executedBy)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003104, err.Error())
	}

	return utils.SuccessResponse(c, execution, "Execution queued successfully")
}
