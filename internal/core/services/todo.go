package services

import (
	"time"

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

func (s *todoService) Create(userID uint, todo *entities.Todo) (*entities.Todo, error) {
	now := time.Now()
	todo.UserID = userID
	todo.Status = defaultString(todo.Status, "pending")
	todo.Priority = defaultString(todo.Priority, "normal")
	todo.CreatedAt = now
	todo.UpdatedAt = now
	return s.repo.Create(todo)
}

func (s *todoService) List(userID uint, limit, offset int) ([]*entities.Todo, error) {
	return s.repo.ListByUserID(userID, limit, offset)
}

func (s *todoService) ListPending(userID uint) ([]*entities.Todo, error) {
	return s.repo.ListPendingByUserID(userID)
}

func (s *todoService) Update(userID, id uint, todo *entities.Todo) (*entities.Todo, error) {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if current.UserID != userID {
		return nil, servicePort.ErrForbidden
	}

	todo.ID = id
	todo.UserID = userID
	todo.CreatedAt = current.CreatedAt
	todo.UpdatedAt = time.Now()
	if todo.Status == "" {
		todo.Status = current.Status
	}
	if todo.Priority == "" {
		todo.Priority = current.Priority
	}
	return s.repo.Update(todo)
}

func (s *todoService) Delete(userID, id uint) error {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if current.UserID != userID {
		return servicePort.ErrForbidden
	}
	return s.repo.Delete(id)
}
