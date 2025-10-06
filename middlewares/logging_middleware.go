package middlewares

import (
	"fmt"
	"scheduling-report/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Get or generate trace ID
		traceID := c.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Locals("traceId", traceID)

		// Parse request body
		reqBody := c.Body()
		reqContentType := c.Get("Content-Type")

		// Parse request body to interface
		reqBodyAsInterface := utils.ParseBodyToJSON(reqBody, reqContentType)
		// Mask sensitive data
		reqBodyMasked := utils.MaskSensitiveData(reqBodyAsInterface)

		// Format headers
		reqHeaders := utils.FormatHeaders(c.GetReqHeaders())

		// Parse query parameters
		queryString := string(c.Request().URI().QueryString())
		queryParams := utils.ParseQueryParams(queryString)

		// Build log fields
		logFields := map[string]interface{}{
			"headers": reqHeaders,
		}

		// Add body only if it exists
		if reqBodyMasked != nil {
			logFields["body"] = reqBodyMasked
		}

		// Add query params only if they exist (for GET requests)
		if queryParams != nil && len(queryParams) > 0 {
			logFields["queryParams"] = queryParams
		}

		// Log incoming request
		log.Info().
			Str("traceId", traceID).
			Str("method", c.Method()).
			Str("path", c.OriginalURL()).
			Fields(logFields).
			Msg(fmt.Sprintf("Incoming request %s %s", c.Method(), c.Path()))

		// Process request
		err := c.Next()

		// Parse response
		statusCode := c.Response().StatusCode()
		latency := time.Since(start)
		resBody := c.Response().Body()
		resContentType := string(c.Response().Header.ContentType())

		// Parse response body
		resBodyAsInterface := utils.ParseBodyToJSON(resBody, resContentType)
		// Mask sensitive data in response
		resBodyMasked := utils.MaskSensitiveData(resBodyAsInterface)

		// Format response headers
		resHeaders := utils.FormatHeaders(c.GetRespHeaders())

		// Build response log fields
		resLogFields := map[string]interface{}{
			"headers": resHeaders,
		}

		// Add response body only if it exists
		if resBodyMasked != nil {
			resLogFields["body"] = resBodyMasked
		}

		// Log outgoing response
		log.Info().
			Str("traceId", traceID).
			Str("method", c.Method()).
			Str("path", c.OriginalURL()).
			Int("statusCode", statusCode).
			Dur("processingTime", latency).
			Fields(resLogFields).
			Msg(fmt.Sprintf("Outgoing response %s %s %d", c.Method(), c.Path(), statusCode))

		return err
	}
}
