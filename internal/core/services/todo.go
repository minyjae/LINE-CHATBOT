package services

import (
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
)

type todoService struct {
	repo repoPort.TodoRepository
}

var _ servicePort.TodoService = (*todoService)(nil)

func NewTodoServiceImpl(repo repoPort.TodoRepository) *todoService {
	return &todoService{repo: repo}
}

func (s *todoService) List(userID uint, limit, offset int) ([]*entities.Todo, error) {
	return s.repo.ListByUserID(userID, limit, offset)
}

func (s *todoService) ListPending(userID uint) ([]*entities.Todo, error) {
	return s.repo.ListPendingByUserID(userID)
}
