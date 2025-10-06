package controllers

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"scheduling-report/services"
	"scheduling-report/utils"
)

type ReportConfigController struct {
	service  *services.ReportConfigService
	validate *validator.Validate
}

func NewReportConfigController() *ReportConfigController {
	return &ReportConfigController{
		service:  services.NewReportConfigService(),
		validate: validator.New(),
	}
}

// GetReportConfigs handles GET /api/report-configs
func (ctrl *ReportConfigController) GetReportConfigs(c *fiber.Ctx) error {
	var isActive *bool
	var datasourceID *int

	if c.Query("is_active") != "" {
		val := c.Query("is_active") == "true"
		isActive = &val
	}

	if c.Query("datasource_id") != "" {
		if id, err := strconv.Atoi(c.Query("datasource_id")); err == nil {
			datasourceID = &id
		}
	}

	configs, err := ctrl.service.GetAll(isActive, datasourceID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 1, err.Error())
	}

	return utils.SuccessResponse(c, configs, "Report configs retrieved successfully")
}

// GetReportConfigByID handles GET /api/report-configs/:id
func (ctrl *ReportConfigController) GetReportConfigByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 1, "Invalid report config ID")
	}

	config, err := ctrl.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, 0, "Report config not found")
	}

	return utils.SuccessResponse(c, config, "Report config retrieved successfully")
}

// CreateReportConfig handles POST /api/report-configs
func (ctrl *ReportConfigController) CreateReportConfig(c *fiber.Ctx) error {
	var input services.CreateReportConfigInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 1, "Invalid request body")
	}

	// Set user context for audit
	input.CreatedBy = c.Get("X-User-ID", "system")

	// Capture IP and session for audit
	ipAddr := c.IP()
	input.IPAddress = &ipAddr
	sessionID := c.Get("X-Session-ID", "")
	if sessionID != "" {
		input.SessionID = &sessionID
	}

	if err := ctrl.validate.Struct(input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 2, err.Error())
	}

	config, err := ctrl.service.Create(input)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 3, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"responseCode":    utils.BuildCode(fiber.StatusCreated, 0),
		"responseMessage": "Report config created successfully",
		"data":            config,
	})
}

// UpdateReportConfig handles PUT /api/report-configs/:id
func (ctrl *ReportConfigController) UpdateReportConfig(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 1, "Invalid report config ID")
	}

	var input services.UpdateReportConfigInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 1, "Invalid request body")
	}

	// Set user context for audit
	input.UpdatedBy = c.Get("X-User-ID", "system")

	// Capture IP and session for audit
	ipAddr := c.IP()
	input.IPAddress = &ipAddr
	sessionID := c.Get("X-Session-ID", "")
	if sessionID != "" {
		input.SessionID = &sessionID
	}

	if err := ctrl.validate.Struct(input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 2, err.Error())
	}

	config, err := ctrl.service.Update(id, input)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 3, err.Error())
	}

	return utils.SuccessResponse(c, config, "Report config updated successfully")
}

// DeleteReportConfig handles DELETE /api/report-configs/:id
func (ctrl *ReportConfigController) DeleteReportConfig(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 1, "Invalid report config ID")
	}

	// Get user context for audit
	deletedBy := c.Get("X-User-ID", "system")
	ipAddr := c.IP()
	sessionID := c.Get("X-Session-ID", "")
	var sessionPtr *string
	if sessionID != "" {
		sessionPtr = &sessionID
	}

	if err := ctrl.service.Delete(id, deletedBy, sessionPtr, &ipAddr); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 3, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Report config deleted successfully")
}
