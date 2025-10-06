package middlewares

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		traceID := c.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
			c.Locals("traceId", traceID)
		}

		log.Info().
			Str("traceId", traceID).
			Str("method", c.Method()).
			Str("path", c.OriginalURL()).
			Msg(fmt.Sprintf("Incoming request %s %s", c.Method(), c.Path()))

		err := c.Next()

		statusCode := c.Response().StatusCode()
		latency := time.Since(start)

		log.Info().
			Str("traceId", traceID).
			Str("method", c.Method()).
			Str("path", c.OriginalURL()).
			Int("statusCode", statusCode).
			Dur("processingTime", latency).
			Msg(fmt.Sprintf("Outgoing response %s %s %d", c.Method(), c.Path(), statusCode))

		return err
	}
}
