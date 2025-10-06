package controllers

import (
	"scheduling-report/services"
	"scheduling-report/utils"
	"strconv"
	"github.com/gofiber/fiber/v2"
)

type ReportDeliveryLogController struct {
	service *services.ReportDeliveryLogService
}

func NewReportDeliveryLogController() *ReportDeliveryLogController {
	return &ReportDeliveryLogController{service: services.NewReportDeliveryLogService()}
}

func (ctrl *ReportDeliveryLogController) GetDeliveryLogs(c *fiber.Ctx) error {
	var status *string
	if c.Query("status") != "" {
		s := c.Query("status")
		status = &s
	}
	limit, _ := strconv.Atoi(c.Query("limit", "100"))
	
	logs, err := ctrl.service.GetAll(status, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}
	return utils.SuccessResponse(c, logs, "Delivery logs retrieved successfully")
}

func (ctrl *ReportDeliveryLogController) GetDeliveryLogByID(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid log ID")
	}
	log, err := ctrl.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, 40403100, err.Error())
	}
	return utils.SuccessResponse(c, log, "Delivery log retrieved successfully")
}

func (ctrl *ReportDeliveryLogController) GetDeliveryLogsByExecutionID(c *fiber.Ctx) error {
	executionID := c.Params("execution_id")
	logs, err := ctrl.service.GetByExecutionID(executionID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}
	return utils.SuccessResponse(c, logs, "Delivery logs retrieved successfully")
}

func (ctrl *ReportDeliveryLogController) GetDeliveryLogsByDeliveryID(c *fiber.Ctx) error {
	deliveryID, err := strconv.Atoi(c.Params("delivery_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid delivery ID")
	}
	limit, _ := strconv.Atoi(c.Query("limit", "100"))
	
	logs, err := ctrl.service.GetByDeliveryID(deliveryID, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}
	return utils.SuccessResponse(c, logs, "Delivery logs retrieved successfully")
}
