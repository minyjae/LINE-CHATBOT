package services

import "minyjae/go-starter/internal/core/domain/entities"

// TodoService คือ contract business logic สำหรับ todo
// input: userID, todo entity หรือ id
// output: todo entity/list หรือ error
type TodoService interface {
	Create(userID uint, todo *entities.Todo) (*entities.Todo, error)
	List(userID uint, limit, offset int) ([]*entities.Todo, error)
	ListPending(userID uint) ([]*entities.Todo, error)
	Update(userID, id uint, todo *entities.Todo) (*entities.Todo, error)
	Delete(userID, id uint) error
}
