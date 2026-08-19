package services

import (
	"strings"

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
