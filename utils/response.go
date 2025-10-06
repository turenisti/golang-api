package utils

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"scheduling-report/config"
)

type StandardResponse struct {
	Code    string      `json:"responseCode"`
	Message string      `json:"responseMessage"`
	Data    interface{} `json:"data,omitempty"`
}

func BuildCode(httpCode int, errorCode int) string {
	httpStr := strconv.Itoa(httpCode)
	serviceStr := config.Config.ServiceCode
	errorStr := fmt.Sprintf("%02d", errorCode)
	return fmt.Sprintf("%s%s%s", httpStr, serviceStr, errorStr)
}

func SuccessResponse(c *fiber.Ctx, data interface{}, message string) error {
	code := BuildCode(fiber.StatusOK, 0)
	return c.Status(fiber.StatusOK).JSON(StandardResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *fiber.Ctx, httpCode int, errorCode int, message string) error {
	code := BuildCode(httpCode, errorCode)
	return c.Status(httpCode).JSON(StandardResponse{
		Code:    code,
		Message: message,
	})
}
