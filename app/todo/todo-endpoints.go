package todo

import (
	"github.com/luridsnk/todo-go/app"
	"github.com/luridsnk/todo-go/common"
	"github.com/luridsnk/todo-go/store"

	"github.com/gofiber/fiber/v2"
)

var subKey string = "sub"

func UseEndpoints(application *app.App, store *store.Store) {
	todoGroup := application.Group("api/v1/todo")
	{
		srv := &TodoItemService{store: store}
		todoGroup.Get("/", getAllItems(srv))
		todoGroup.Post("/", addNewItem(srv))
		todoGroup.Put("/", updateItem(srv))
		todoGroup.Delete("/:id", deleteItem(srv))
	}
}

func getAllItems(todoSrv *TodoItemService) func(c *fiber.Ctx) error {

	return func(c *fiber.Ctx) error {
		id := common.ValueFromLocals[string](c, subKey)
		todoItems, err := todoSrv.GetAllByUserId(id)
		if err != nil {
			return err
		}
		return c.JSON(todoItems)
	}
}

func addNewItem(todoSrv *TodoItemService) func(c *fiber.Ctx) error {

	return func(c *fiber.Ctx) error {
		userId := common.ValueFromLocals[string](c, subKey)
		i, err := common.ReadJson[NewItemDto](c.Body())
		if err != nil {
			return err
		}
		itemId, err := todoSrv.AddNewItem(userId, i)
		if err != nil {
			return err
		}

		return c.SendString(itemId)
	}
}

func updateItem(todoSrv *TodoItemService) func(c *fiber.Ctx) error {

	return func(c *fiber.Ctx) error {
		userId := common.ValueFromLocals[string](c, subKey)
		updatedItem, err := common.ReadJson[UpdateItemDto](c.Body())
		if err != nil {
			c.Status(400)
			return c.SendString("Given json object could not be parsed")
		}

		if err = todoSrv.UpdateItem(userId, updatedItem); err != nil {
			return err
		}

		return c.SendStatus(204)
	}
}

func deleteItem(todoSrv *TodoItemService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		userId := common.ValueFromLocals[string](c, subKey)
		itemId := c.Params("id")

		if err := todoSrv.DeleteItem(userId, itemId); err != nil {
			return err
		}

		return c.SendStatus(204)
	}
}
