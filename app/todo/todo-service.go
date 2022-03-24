package todo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/luridsnk/todo-go/store"
)

type TodoItemService struct {
	store *store.Store
}

type GetItemsResult struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	IsDone      bool      `json:"isDone"`
}

func (todoSrv *TodoItemService) GetAllByUserId(id string) ([]*GetItemsResult, error) {
	var todoItems []*GetItemsResult
	rows, err := todoSrv.store.Query("select id, description, isDone from TodoItems where creatorId = $1", id)
	if err != nil {
		return nil, &fiber.Error{Message: err.Error(), Code: 500}
	}

	for rows.Next() {
		var i GetItemsResult
		err = rows.Scan(&i.Id, &i.Description, &i.IsDone)
		if err != nil {
			return nil, &fiber.Error{Message: err.Error(), Code: 400}
		}
		todoItems = append(todoItems, &i)
	}
	return todoItems, nil
}

type NewItemDto struct {
	Description string `json:"description"`
	IsDone      bool   `json:"isDone"`
}

func (todoSrv *TodoItemService) AddNewItem(userId string, newItem *NewItemDto) (string, error) {

	row, err := todoSrv.store.QueryRow(
		"insert into TodoItems (creatorId, description, isDone) values ($1, $2, $3) returning id;",
		userId,
		newItem.Description,
		newItem.IsDone)
	if err != nil {
		return "", &fiber.Error{Message: err.Error(), Code: 500}
	}

	var itemId string
	err = row.Scan(&itemId)
	if err != nil {
		return "", &fiber.Error{Message: err.Error(), Code: 400}
	}
	return itemId, nil
}

type UpdateItemDto struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	IsDone      bool      `json:"isDone"`
}

func (todoSrv *TodoItemService) UpdateItem(userId string, updItem *UpdateItemDto) error {
	row, err := todoSrv.store.QueryRow(
		"select exists(select 1 from TodoItems where id=$1 and creatorId = $2)",
		updItem.Id,
		userId)
	if err != nil {
		return &fiber.Error{Message: err.Error(), Code: 500}
	}
	var exists bool
	row.Scan(&exists)
	if !exists {
		return &fiber.Error{Message: "not found", Code: 404}
	}

	executed, err := todoSrv.store.Execute(
		"update TodoItems SET description = $1, isDone = $2 WHERE id = $3;",
		updItem.Description,
		updItem.IsDone,
		updItem.Id)
	if err != nil || !executed {
		return &fiber.Error{Message: "not found", Code: 400}
	}

	return nil
}

func (todoSrv *TodoItemService) DeleteItem(userId, itemId string) error {
	executed, err := todoSrv.store.Execute("delete from TodoItems where id = $1 and creatorId = $2",
		itemId,
		userId)
	if !executed {
		return &fiber.Error{Message: err.Error(), Code: 404}
	}

	return nil
}
