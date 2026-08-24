package services

import (
	"strings"
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
)

type noteService struct {
	repo repoPort.NoteRepository
}

var _ servicePort.NoteService = (*noteService)(nil)

func NewNoteServiceImpl(repo repoPort.NoteRepository) *noteService {
	return &noteService{repo: repo}
}

func (s *noteService) Create(userID uint, note *entities.Note) (*entities.Note, error) {
	now := time.Now()
	note.UserID = userID
	note.CreatedAt = now
	note.UpdatedAt = now
	return s.repo.Create(note)
}

func (s *noteService) List(userID uint, limit, offset int) ([]*entities.Note, error) {
	return s.repo.ListByUserID(userID, limit, offset)
}

func (s *noteService) Search(userID uint, query string, limit int) ([]*entities.Note, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return s.repo.ListByUserID(userID, limit, 0)
	}
	return s.repo.SearchByContent(userID, query, limit)
}

func (s *noteService) Update(userID, id uint, note *entities.Note) (*entities.Note, error) {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if current.UserID != userID {
		return nil, servicePort.ErrForbidden
	}

	note.ID = id
	note.UserID = userID
	note.CreatedAt = current.CreatedAt
	note.UpdatedAt = time.Now()
	return s.repo.Update(note)
}

func (s *noteService) Delete(userID, id uint) error {
	current, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if current.UserID != userID {
		return servicePort.ErrForbidden
	}
	return s.repo.Delete(id)
}
