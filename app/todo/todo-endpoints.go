package todo

import (
	"fmt"

	"github.com/luridsnk/todo-go/app"
	"github.com/luridsnk/todo-go/common"
	"github.com/luridsnk/todo-go/store"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var subKey string = "sub"

func UseEndpoints(application *app.App, store *store.Store) {
	todoGroup := application.Group("api/v1/todo")
	{
		todoGroup.Get("/", getAllItems(store))
		todoGroup.Post("/", addNewItem(store))
		todoGroup.Put("/", updateItem(store))
		todoGroup.Delete("/:id", deleteItem(store))
	}
}

func getAllItems(store *store.Store) func(c *fiber.Ctx) error {

	type todoItemDto struct {
		Id          uuid.UUID `json:"id"`
		Description string    `json:"description"`
		IsDone      bool      `json:"isDone"`
	}

	return func(c *fiber.Ctx) error {
		id := common.ValueFromLocals[string](c, subKey)
		var todoItems []*todoItemDto
		rows, err := store.Query("select id, description, isDone from TodoItems where creatorId = $1", id)
		if err != nil {
			return err
		}

		for rows.Next() {
			var i todoItemDto
			err = rows.Scan(&i.Id, &i.Description, &i.IsDone)
			if err != nil {
				return err
			}
			todoItems = append(todoItems, &i)
		}

		return c.JSON(todoItems)
	}
}

func addNewItem(store *store.Store) func(c *fiber.Ctx) error {

	type todoItemDto struct {
		Description string `json:"description"`
		IsDone      bool   `json:"isDone"`
	}

	return func(c *fiber.Ctx) error {
		id := common.ValueFromLocals[string](c, subKey)
		i, err := common.ReadJson[todoItemDto](c.Body())
		if err != nil {
			return err
		}

		row, err := store.QueryRow("insert into TodoItems (creatorId, description, isDone) values ($1, $2, $3) returning id;", id, i.Description, i.IsDone)
		if err != nil {
			return err
		}

		var itemId string
		err = row.Scan(&itemId)
		if err != nil {
			return err
		}

		return c.SendString(itemId)
	}
}

func updateItem(store *store.Store) func(c *fiber.Ctx) error {

	type todoItemDto struct {
		Id          uuid.UUID `json:"id"`
		Description string    `json:"description"`
		IsDone      bool      `json:"isDone"`
	}

	return func(c *fiber.Ctx) error {
		id := common.ValueFromLocals[string](c, subKey)
		updatedItem, err := common.ReadJson[todoItemDto](c.Body())
		if err != nil {
			c.Status(400)
			return c.SendString("Given json object could not be parsed")
		}

		row, err := store.QueryRow(
			"select exists(select 1 from TodoItems where id=$1 and creatorId = $2)",
			updatedItem.Id,
			id)
		if err != nil {
			return err
		}
		var exists bool
		row.Scan(&exists)
		if !exists {
			c.Status(404)
			return c.SendString(fmt.Sprintf("No such value"))
		}

		executed, err := store.Execute(
			"update TodoItems SET description = $1, isDone = $2 WHERE id = $3;",
			updatedItem.Description,
			updatedItem.IsDone,
			updatedItem.Id)
		if err != nil || !executed {
			c.Status(400)
			return err
		}

		return c.SendStatus(204)
	}
}

func deleteItem(store *store.Store) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := common.ValueFromLocals[string](c, subKey)
		var itemId uuid.UUID
		itemId, err := uuid.Parse(c.Params("id"))
		if err != nil {
			c.Status(400)
			return c.SendString(fmt.Sprintf("Error occured: %s", err))
		}

		executed, err := store.Execute("delete from TodoItems where id = $1 and creatorId = $2", itemId, id)
		if err != nil || !executed {
			c.SendStatus(400)
			return err
		}

		return c.SendStatus(204)
	}
}
