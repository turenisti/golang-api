package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"scheduling-report/services"
	"scheduling-report/utils"
)

type ReportConfigAuditController struct {
	service *services.ReportConfigAuditService
}

func NewReportConfigAuditController() *ReportConfigAuditController {
	return &ReportConfigAuditController{
		service: services.NewReportConfigAuditService(),
	}
}

// GetAudits handles GET /api/audits
// Query params: ?config_id=1&action=update&performed_by=admin
func (ctrl *ReportConfigAuditController) GetAudits(c *fiber.Ctx) error {
	var configID *int
	var action *string
	var performedBy *string

	// Parse query parameters
	if configIDStr := c.Query("config_id"); configIDStr != "" {
		if id, err := strconv.Atoi(configIDStr); err == nil {
			configID = &id
		}
	}

	if actionStr := c.Query("action"); actionStr != "" {
		action = &actionStr
	}

	if performedByStr := c.Query("performed_by"); performedByStr != "" {
		performedBy = &performedByStr
	}

	audits, err := ctrl.service.GetAll(configID, action, performedBy)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 1, err.Error())
	}

	return utils.SuccessResponse(c, audits, "Audit records retrieved successfully")
}

// GetAuditByID handles GET /api/audits/:id
func (ctrl *ReportConfigAuditController) GetAuditByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 1, "Invalid audit ID")
	}

	audit, err := ctrl.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, 0, "Audit record not found")
	}

	return utils.SuccessResponse(c, audit, "Audit record retrieved successfully")
}

// GetAuditsByConfigID handles GET /api/audits/config/:config_id
func (ctrl *ReportConfigAuditController) GetAuditsByConfigID(c *fiber.Ctx) error {
	configID, err := strconv.Atoi(c.Params("config_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 1, "Invalid config ID")
	}

	audits, err := ctrl.service.GetByConfigID(configID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 2, "Failed to retrieve audit records")
	}

	return utils.SuccessResponse(c, audits, "Config audit history retrieved successfully")
}

// GetRecentChanges handles GET /api/audits/recent?days=7
func (ctrl *ReportConfigAuditController) GetRecentChanges(c *fiber.Ctx) error {
	days := 7 // default
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	audits, err := ctrl.service.GetRecentChanges(days)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 2, "Failed to retrieve recent changes")
	}

	return utils.SuccessResponse(c, audits, "Recent changes retrieved successfully")
}
