package services

import "minyjae/go-starter/internal/core/domain/entities"

type NoteService interface {
	Create(userID uint, note *entities.Note) (*entities.Note, error)
	List(userID uint, limit, offset int) ([]*entities.Note, error)
	Search(userID uint, query string, limit int) ([]*entities.Note, error)
	Update(userID, id uint, note *entities.Note) (*entities.Note, error)
	Delete(userID, id uint) error
}
