package repositories

import "minyjae/go-starter/internal/core/domain/entities"

type NoteRepository interface {
	Create(note *entities.Note) (*entities.Note, error)
	GetByID(id uint) (*entities.Note, error)
	ListByUserID(userID uint, limit, offset int) ([]*entities.Note, error)
	SearchByContent(userID uint, query string, limit int) ([]*entities.Note, error)
	Update(note *entities.Note) (*entities.Note, error)
	Delete(id uint) error
}
