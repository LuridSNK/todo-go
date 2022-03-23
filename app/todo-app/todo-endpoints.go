package todo

import (
	"fmt"
	"todo_app/app"
	"todo_app/common"
	"todo_app/store"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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
	return func(c *fiber.Ctx) error {
		id := common.GetUserFromContext(c)
		var todoItems []*TodoItem
		rows, err := store.Query("select * from TodoItems where creatorId = $1", id)
		if err != nil {
			return err
		}

		for rows.Next() {
			var i TodoItem
			err = rows.Scan(&i.Id, &i.Description, &i.CreatedAt, &i.IsDone, &i.CreatorId)
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

func addNewItem(store *store.Store) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := common.GetUserFromContext(c)
		i, err := common.ReadJson[TodoItem](c.Body())
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
	return func(c *fiber.Ctx) error {
		id := common.GetUserFromContext(c)
		updatedItem, err := common.ReadJson[TodoItem](c.Body())
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

func deleteItem(store *store.Store) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := common.GetUserFromContext(c)
		var itemId uuid.UUID
		itemId, err := uuid.Parse(c.Params("id"))
		if err != nil {
			c.Status(400)
			return c.SendString(fmt.Sprintf("Error occured: %s", err))
		}

		err = store.Execute("delete from TodoItems where id = $1 and creatorId = $2", itemId, id)
		if err != nil {
			return err
		}

		return c.SendStatus(204)
	}
}
