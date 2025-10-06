package controllers

import (
	"scheduling-report/services"
	"scheduling-report/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ReportDeliveryController struct {
	service  *services.ReportDeliveryService
	validate *validator.Validate
}

func NewReportDeliveryController() *ReportDeliveryController {
	return &ReportDeliveryController{
		service:  services.NewReportDeliveryService(),
		validate: validator.New(),
	}
}

func (ctrl *ReportDeliveryController) GetDeliveries(c *fiber.Ctx) error {
	var isActive *bool
	if c.Query("is_active") != "" {
		val := c.Query("is_active") == "true"
		isActive = &val
	}

	deliveries, err := ctrl.service.GetAll(isActive)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}

	return utils.SuccessResponse(c, deliveries, "Deliveries retrieved successfully")
}

func (ctrl *ReportDeliveryController) GetDeliveryByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid delivery ID")
	}

	delivery, err := ctrl.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, 40403100, err.Error())
	}

	return utils.SuccessResponse(c, delivery, "Delivery retrieved successfully")
}

func (ctrl *ReportDeliveryController) GetDeliveriesByConfigID(c *fiber.Ctx) error {
	configID, err := strconv.Atoi(c.Params("config_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid config ID")
	}

	deliveries, err := ctrl.service.GetByConfigID(configID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}

	return utils.SuccessResponse(c, deliveries, "Deliveries retrieved successfully")
}

func (ctrl *ReportDeliveryController) CreateDelivery(c *fiber.Ctx) error {
	var input services.CreateDeliveryInput

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid request body")
	}

	input.CreatedBy = c.Get("X-User-ID", "system")
	ipAddr := c.IP()
	input.IPAddress = &ipAddr
	sessionID := c.Get("X-Session-ID", "")
	if sessionID != "" {
		input.SessionID = &sessionID
	}

	if err := ctrl.validate.Struct(input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003102, err.Error())
	}

	delivery, err := ctrl.service.Create(input)
	if err != nil {
		if err.Error() == "report config not found" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003103, err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"responseCode":    "20103100",
		"responseMessage": "Delivery created successfully",
		"data":            delivery,
	})
}

func (ctrl *ReportDeliveryController) UpdateDelivery(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid delivery ID")
	}

	var input services.UpdateDeliveryInput

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid request body")
	}

	input.UpdatedBy = c.Get("X-User-ID", "system")
	ipAddr := c.IP()
	input.IPAddress = &ipAddr
	sessionID := c.Get("X-Session-ID", "")
	if sessionID != "" {
		input.SessionID = &sessionID
	}

	if err := ctrl.validate.Struct(input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003102, err.Error())
	}

	delivery, err := ctrl.service.Update(id, input)
	if err != nil {
		if err.Error() == "delivery not found" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, 40403100, err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}

	return utils.SuccessResponse(c, delivery, "Delivery updated successfully")
}

func (ctrl *ReportDeliveryController) DeleteDelivery(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid delivery ID")
	}

	deletedBy := c.Query("deleted_by", "system")
	ipAddr := c.IP()
	sessionID := c.Query("session_id", "")
	var sessionIDPtr *string
	if sessionID != "" {
		sessionIDPtr = &sessionID
	}

	if err := ctrl.service.Delete(id, deletedBy, sessionIDPtr, &ipAddr); err != nil {
		if err.Error() == "delivery not found" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, 40403100, err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Delivery deleted successfully")
}
