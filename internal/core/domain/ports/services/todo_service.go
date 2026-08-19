package services

import "minyjae/go-starter/internal/core/domain/entities"

type TodoService interface {
	List(userID uint, limit, offset int) ([]*entities.Todo, error)
	ListPending(userID uint) ([]*entities.Todo, error)
}
