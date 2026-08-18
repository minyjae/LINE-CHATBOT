package repositories

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

type TodoRepository interface {
	Create(todo *entities.Todo) (*entities.Todo, error)
	GetByID(id uint) (*entities.Todo, error)
	ListByUserID(userID uint, limit, offset int) ([]*entities.Todo, error)
	ListPendingByUserID(userID uint) ([]*entities.Todo, error)
	ListDueBetween(userID uint, start, end time.Time) ([]*entities.Todo, error)
	Update(todo *entities.Todo) (*entities.Todo, error)
	Delete(id uint) error
}
