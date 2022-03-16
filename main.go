package main

import (
	"log"
	"todo_app/todo-app"

	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
)

var (
	port string
)

var storage = map[uuid.UUID]todo.TodoItem{}

func main() {
	port = ":5000"
	app := todo.New()

	app.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02T15:04:05",
		TimeZone:   "UTC",
	}))

	app.UseTodoEndpoints(storage)

	log.Fatal(app.Listen(port))
}
