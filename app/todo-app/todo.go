package todo

import (
	"time"

	"github.com/google/uuid"
)

type TodoItem struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	IsDone      bool      `json:"isDone"`
	CreatorId   uuid.UUID `json:"creatorId"`
}
