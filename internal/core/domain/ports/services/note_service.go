package services

import "minyjae/go-starter/internal/core/domain/entities"

// NoteService คือ contract business logic สำหรับ note
// input: userID, note entity, id หรือ query
// output: note entity/list หรือ error
type NoteService interface {
	Create(userID uint, note *entities.Note) (*entities.Note, error)
	List(userID uint, limit, offset int) ([]*entities.Note, error)
	Search(userID uint, query string, limit int) ([]*entities.Note, error)
	Update(userID, id uint, note *entities.Note) (*entities.Note, error)
	Delete(userID, id uint) error
}
