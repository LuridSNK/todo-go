package auth

import (
	"github.com/luridsnk/todo-go/common"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func JwtAuthMiddleware() func(c *fiber.Ctx) error {
	errorHandler := func(c *fiber.Ctx, err error) error {
		code := 401
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		c.SendStatus(code)
		return &fiber.Error{Code: code, Message: err.Error()}
	}

	jwtWare := jwtware.New(jwtware.Config{
		SigningKey:    []byte(secret),
		SigningMethod: jwtware.HS256,
		ContextKey:    common.TokenKey,
		ErrorHandler:  errorHandler,
	})

	return jwtWare
}
