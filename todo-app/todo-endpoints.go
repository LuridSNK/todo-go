package todo

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type App struct {
	*fiber.App
}

func New() *App {
	return &App{
		App: fiber.New(),
	}
}

func (app *App) UseTodoEndpoints(store map[uuid.UUID]TodoItem) {
	todoGroup := app.Group("api/v1/todo")
	{
		todoGroup.Get("/", GetAllItems(store))
		todoGroup.Post("/", AddNewItem(store))
		todoGroup.Put("/", UpdateItem(store))
		todoGroup.Delete("/:id", DeleteItem(store))
	}
}

func GetAllItems(store map[uuid.UUID]TodoItem) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		keys := make([]TodoItem, 0, len(store))
		for k := range store {
			keys = append(keys, store[k])
		}
		response, err := json.Marshal(keys)
		if err != nil {
			return c.SendString(fmt.Sprintf("Error: %s", err))
		}

		return c.SendString(string(response))
	}
}

func AddNewItem(store map[uuid.UUID]TodoItem) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var i TodoItem
		err := json.Unmarshal(c.Body(), &i)
		i.Id = uuid.New()
		i.CreatedAt = time.Now().UTC()
		if err != nil {
			return c.SendString(fmt.Sprintf("Error: %s", err))
		}

		store[i.Id] = i

		return c.SendString(i.Id.String())
	}
}

func UpdateItem(store map[uuid.UUID]TodoItem) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var updatedItem TodoItem
		if err := json.Unmarshal(c.Body(), &updatedItem); err != nil {
			c.Status(400)
			return c.SendString(fmt.Sprintf("No such value"))
		}

		if _, contains := store[updatedItem.Id]; !contains {
			c.Status(404)
			return c.SendString(fmt.Sprintf("No such value"))
		}

		store[updatedItem.Id] = updatedItem

		return c.SendStatus(204)
	}
}

func DeleteItem(store map[uuid.UUID]TodoItem) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var id uuid.UUID
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			c.Status(400)
			return c.SendString(fmt.Sprintf("Error occured: %s", err))
		}

		if _, contains := store[id]; !contains {
			c.Status(404)
			return c.SendString(fmt.Sprintf("No such value"))
		}
		delete(store, id)
		return c.SendStatus(204)
	}
}
