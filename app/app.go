package app

import (
	"github.com/gofiber/fiber/v2"
)

type App struct {
	*fiber.App
}

func New() *App {
	return &App{
		App: fiber.New(),
	}
}
