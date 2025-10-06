package controllers

import (
	"scheduling-report/services"
	"scheduling-report/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ReportScheduleController struct {
	service  *services.ReportScheduleService
	validate *validator.Validate
}

func NewReportScheduleController() *ReportScheduleController {
	return &ReportScheduleController{
		service:  services.NewReportScheduleService(),
		validate: validator.New(),
	}
}

// GetSchedules retrieves all schedules
func (ctrl *ReportScheduleController) GetSchedules(c *fiber.Ctx) error {
	var isActive *bool
	if c.Query("is_active") != "" {
		val := c.Query("is_active") == "true"
		isActive = &val
	}

	schedules, err := ctrl.service.GetAll(isActive)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}

	return utils.SuccessResponse(c, schedules, "Schedules retrieved successfully")
}

// GetScheduleByID retrieves a schedule by ID
func (ctrl *ReportScheduleController) GetScheduleByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid schedule ID")
	}

	schedule, err := ctrl.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, 40403100, err.Error())
	}

	return utils.SuccessResponse(c, schedule, "Schedule retrieved successfully")
}

// GetSchedulesByConfigID retrieves all schedules for a report config
func (ctrl *ReportScheduleController) GetSchedulesByConfigID(c *fiber.Ctx) error {
	configID, err := strconv.Atoi(c.Params("config_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid config ID")
	}

	schedules, err := ctrl.service.GetByConfigID(configID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}

	return utils.SuccessResponse(c, schedules, "Schedules retrieved successfully")
}

// CreateSchedule creates a new schedule
func (ctrl *ReportScheduleController) CreateSchedule(c *fiber.Ctx) error {
	var input services.CreateScheduleInput

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid request body")
	}

	// Get user ID from header
	input.CreatedBy = c.Get("X-User-ID", "system")

	// Capture IP and session for audit
	ipAddr := c.IP()
	input.IPAddress = &ipAddr
	sessionID := c.Get("X-Session-ID", "")
	if sessionID != "" {
		input.SessionID = &sessionID
	}

	if err := ctrl.validate.Struct(input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003102, err.Error())
	}

	schedule, err := ctrl.service.Create(input)
	if err != nil {
		if err.Error() == "invalid cron expression" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003102, err.Error())
		}
		if err.Error() == "report config not found" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003103, err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"responseCode":    "20103100",
		"responseMessage": "Schedule created successfully",
		"data":            schedule,
	})
}

// UpdateSchedule updates a schedule
func (ctrl *ReportScheduleController) UpdateSchedule(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid schedule ID")
	}

	var input services.UpdateScheduleInput

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid request body")
	}

	// Get user ID from header
	input.UpdatedBy = c.Get("X-User-ID", "system")

	// Capture IP and session for audit
	ipAddr := c.IP()
	input.IPAddress = &ipAddr
	sessionID := c.Get("X-Session-ID", "")
	if sessionID != "" {
		input.SessionID = &sessionID
	}

	if err := ctrl.validate.Struct(input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003102, err.Error())
	}

	schedule, err := ctrl.service.Update(id, input)
	if err != nil {
		if err.Error() == "invalid cron expression" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003102, err.Error())
		}
		if err.Error() == "schedule not found" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, 40403100, err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}

	return utils.SuccessResponse(c, schedule, "Schedule updated successfully")
}

// DeleteSchedule soft deletes a schedule
func (ctrl *ReportScheduleController) DeleteSchedule(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid schedule ID")
	}

	deletedBy := c.Query("deleted_by", "system")

	// Capture IP and session for audit
	ipAddr := c.IP()
	sessionID := c.Query("session_id", "")
	var sessionIDPtr *string
	if sessionID != "" {
		sessionIDPtr = &sessionID
	}

	if err := ctrl.service.Delete(id, deletedBy, sessionIDPtr, &ipAddr); err != nil {
		if err.Error() == "schedule not found" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, 40403100, err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Schedule deleted successfully")
}
