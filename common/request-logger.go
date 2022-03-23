package common

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RequestLoggingMiddleware(logger *AppLogger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// add request id
		requestId := uuid.NewString()
		c.Request().Header.Add("X-RequestId", requestId)
		// start request time logging
		startTime := time.Now()

		err := c.Next() // execution context

		elapsedMs := fmt.Sprintf("%d ms", time.Now().Sub(startTime).Milliseconds())
		var fields []interface{}
		fields = append(fields,
			"request_id", requestId,
			"method", c.Method(),
			"path", c.Path())
		statusCode := fmt.Sprintf("%d", c.Context().Response.StatusCode())
		if err != nil {
			fields = append(fields,
				"request_time", elapsedMs,
				"status_code", statusCode,
				"error", err)
			logger.Warnw("Request served with error", fields...)
			return err
		}

		fields = append(fields,
			"request_time", elapsedMs,
			"status_code", statusCode)
		logger.Infow("Request served", fields...)
		return err
	}
}
