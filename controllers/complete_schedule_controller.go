package controllers

import (
	"scheduling-report/models"
	"scheduling-report/services"
	"scheduling-report/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CompleteScheduleController struct {
	service *services.CompleteScheduleService
}

func NewCompleteScheduleController() *CompleteScheduleController {
	return &CompleteScheduleController{
		service: services.NewCompleteScheduleService(),
	}
}

// CreateComplete creates a complete schedule with config, deliveries, and recipients
// POST /api/schedules/complete
func (ctrl *CompleteScheduleController) CreateComplete(c *fiber.Ctx) error {
	var req models.CompleteScheduleRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40004001, "Invalid request body: "+err.Error())
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40004002, "Validation failed: "+err.Error())
	}

	// Create complete schedule
	response, err := ctrl.service.CreateComplete(req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40004099, "Failed to create schedule: "+err.Error())
	}

	return utils.SuccessResponse(c, response, "Complete schedule created successfully")
}

// UpdateComplete updates a complete schedule with partial update support
// PUT /api/schedules/complete/:id
func (ctrl *CompleteScheduleController) UpdateComplete(c *fiber.Ctx) error {
	// Get schedule ID from URL
	scheduleID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40004003, "Invalid schedule ID")
	}

	var req models.CompleteScheduleUpdateRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40004004, "Invalid request body: "+err.Error())
	}

	// Set updated_by from header if not provided
	if req.UpdatedBy == "" {
		req.UpdatedBy = c.Get("X-User-ID", "system")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40004005, "Validation failed: "+err.Error())
	}

	// Update complete schedule
	response, err := ctrl.service.UpdateComplete(scheduleID, req)
	if err != nil {
		if err.Error() == "schedule not found" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, 40404004, "Schedule not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40004099, "Failed to update schedule: "+err.Error())
	}

	return utils.SuccessResponse(c, response, "Complete schedule updated successfully")
}
