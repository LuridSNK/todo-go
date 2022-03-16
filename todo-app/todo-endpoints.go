package todo

import (
	"fmt"
	"todo_app/common"
	"todo_app/store"

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

func (app *App) UseTodoEndpoints(store *store.Store) {
	todoGroup := app.Group("api/v1/todo")
	{
		todoGroup.Get("/", GetAllItems(store))
		todoGroup.Post("/", AddNewItem(store))
		todoGroup.Put("/", UpdateItem(store))
		todoGroup.Delete("/:id", DeleteItem(store))
	}
}

func GetAllItems(store *store.Store) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var todoItems []*TodoItem
		rows, err := store.Query("select * from TodoItems")
		if err != nil {
			return err
		}

		for rows.Next() {
			var i TodoItem
			err = rows.Scan(&i.Id, &i.Description, &i.CreatedAt, &i.IsDone)
			if err != nil {
				return err
			}
			todoItems = append(todoItems, &i)
		}

		response, err := common.ToJson(todoItems)
		if err != nil {
			return err
		}

		return c.SendString(string(response))
	}
}

func AddNewItem(store *store.Store) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		i, err := common.ReadJson[TodoItem](c.Body())
		if err != nil {
			return err
		}

		row, err := store.QueryRow("insert into TodoItems (description, isDone) values ($1, $2) returning id;", i.Description, i.IsDone)
		if err != nil {
			return err
		}

		var id uuid.UUID
		err = row.Scan(&id)
		if err != nil {
			return err
		}

		return c.SendString(id.String())
	}
}

func UpdateItem(store *store.Store) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		updatedItem, err := common.ReadJson[TodoItem](c.Body())
		if err != nil {
			c.Status(400)
			return c.SendString("Given json object could not be parsed")
		}

		row, err := store.QueryRow(
			"select exists(select 1 from TodoItems where id=$1)",
			updatedItem.Id)
		if err != nil {
			return err
		}
		var exists bool
		row.Scan(&exists)
		if !exists {
			c.Status(404)
			return c.SendString(fmt.Sprintf("No such value"))
		}

		err = store.Execute(
			"update TodoItems SET description = $1, isDone = $2 WHERE id = $3;",
			updatedItem.Description,
			updatedItem.IsDone,
			updatedItem.Id)
		if err != nil {
			return err
		}

		return c.SendStatus(204)
	}
}

func DeleteItem(store *store.Store) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var id uuid.UUID
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			c.Status(400)
			return c.SendString(fmt.Sprintf("Error occured: %s", err))
		}

		err = store.Execute("delete from TodoItems where id = $1", id)
		if err != nil {
			return err
		}

		return c.SendStatus(204)
	}
}
