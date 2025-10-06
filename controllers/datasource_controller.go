package controllers

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"scheduling-report/services"
	"scheduling-report/utils"
)

type DatasourceController struct {
	service  *services.DatasourceService
	validate *validator.Validate
}

func NewDatasourceController() *DatasourceController {
	return &DatasourceController{
		service:  services.NewDatasourceService(),
		validate: validator.New(),
	}
}

// GetDatasources handles GET /api/datasources
func (ctrl *DatasourceController) GetDatasources(c *fiber.Ctx) error {
	var isActive *bool

	if c.Query("is_active") != "" {
		val := c.Query("is_active") == "true"
		isActive = &val
	}

	datasources, err := ctrl.service.GetAll(isActive)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, 1, err.Error())
	}

	return utils.SuccessResponse(c, datasources, "Datasources retrieved successfully")
}

// GetDatasourceByID handles GET /api/datasources/:id
func (ctrl *DatasourceController) GetDatasourceByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 1, "Invalid datasource ID")
	}

	datasource, err := ctrl.service.GetByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, 0, "Datasource not found")
	}

	return utils.SuccessResponse(c, datasource, "Datasource retrieved successfully")
}

// CreateDatasource handles POST /api/datasources
func (ctrl *DatasourceController) CreateDatasource(c *fiber.Ctx) error {
	var input services.CreateDatasourceInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 1, "Invalid request body")
	}

	// Set created_by from context (default to "system" for now)
	input.CreatedBy = c.Get("X-User-ID", "system")

	if err := ctrl.validate.Struct(input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 2, err.Error())
	}

	datasource, err := ctrl.service.Create(input)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 3, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"responseCode":    utils.BuildCode(fiber.StatusCreated, 0),
		"responseMessage": "Datasource created successfully",
		"data":            datasource,
	})
}

// UpdateDatasource handles PUT /api/datasources/:id
func (ctrl *DatasourceController) UpdateDatasource(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 1, "Invalid datasource ID")
	}

	var input services.UpdateDatasourceInput
	if err := c.BodyParser(&input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 1, "Invalid request body")
	}

	// Set updated_by from context
	input.UpdatedBy = c.Get("X-User-ID", "system")

	if err := ctrl.validate.Struct(input); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 2, err.Error())
	}

	datasource, err := ctrl.service.Update(id, input)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 3, err.Error())
	}

	return utils.SuccessResponse(c, datasource, "Datasource updated successfully")
}

// DeleteDatasource handles DELETE /api/datasources/:id
func (ctrl *DatasourceController) DeleteDatasource(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 1, "Invalid datasource ID")
	}

	if err := ctrl.service.Delete(id); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, 3, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Datasource deleted successfully")
}
