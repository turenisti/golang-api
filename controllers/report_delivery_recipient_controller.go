package controllers

import (
	"scheduling-report/services"
	"scheduling-report/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ReportDeliveryRecipientController struct {
	service  *services.ReportDeliveryRecipientService
	validate *validator.Validate
}

func NewReportDeliveryRecipientController() *ReportDeliveryRecipientController {
	return &ReportDeliveryRecipientController{
		service:  services.NewReportDeliveryRecipientService(),
		validate: validator.New(),
	}
}

func (ctrl *ReportDeliveryRecipientController) GetRecipients(c *fiber.Ctx) error {
	var isActive *bool
	if c.Query("is_active") != "" {
		val := c.Query("is_active") == "true"
		isActive = &val
	}

	recipients, err := ctrl.service.GetAll(isActive)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}

	return utils.SuccessResponse(c, recipients, "Recipients retrieved successfully")
}

func (ctrl *ReportDeliveryRecipientController) GetRecipientByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid recipient ID")
	}

	recipient, err := ctrl.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, 40403100, err.Error())
	}

	return utils.SuccessResponse(c, recipient, "Recipient retrieved successfully")
}

func (ctrl *ReportDeliveryRecipientController) GetRecipientsByDeliveryID(c *fiber.Ctx) error {
	deliveryID, err := strconv.Atoi(c.Params("delivery_id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid delivery ID")
	}

	recipients, err := ctrl.service.GetByDeliveryID(deliveryID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}

	return utils.SuccessResponse(c, recipients, "Recipients retrieved successfully")
}

func (ctrl *ReportDeliveryRecipientController) CreateRecipient(c *fiber.Ctx) error {
	var input services.CreateRecipientInput

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid request body")
	}

	if err := ctrl.validate.Struct(input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003102, err.Error())
	}

	recipient, err := ctrl.service.Create(input)
	if err != nil {
		if err.Error() == "delivery not found" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003103, err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"responseCode":    "20103100",
		"responseMessage": "Recipient created successfully",
		"data":            recipient,
	})
}

func (ctrl *ReportDeliveryRecipientController) UpdateRecipient(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid recipient ID")
	}

	var input services.UpdateRecipientInput

	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid request body")
	}

	if err := ctrl.validate.Struct(input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003102, err.Error())
	}

	recipient, err := ctrl.service.Update(id, input)
	if err != nil {
		if err.Error() == "recipient not found" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, 40403100, err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}

	return utils.SuccessResponse(c, recipient, "Recipient updated successfully")
}

func (ctrl *ReportDeliveryRecipientController) DeleteRecipient(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 40003101, "Invalid recipient ID")
	}

	if err := ctrl.service.Delete(id); err != nil {
		if err.Error() == "recipient not found" {
			return utils.ErrorResponse(c, fiber.StatusNotFound, 40403100, err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 40003199, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Recipient deleted successfully")
}
