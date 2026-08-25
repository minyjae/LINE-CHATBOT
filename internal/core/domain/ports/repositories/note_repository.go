package repositories

import "minyjae/go-starter/internal/core/domain/entities"

// NoteRepository คือ contract สำหรับอ่าน/เขียนและค้น note
// input: domain entity, userID, id หรือ query content
// output: domain entity/list หรือ error จาก adapter ที่ implement จริง
type NoteRepository interface {
	Create(note *entities.Note) (*entities.Note, error)
	GetByID(id uint) (*entities.Note, error)
	ListByUserID(userID uint, limit, offset int) ([]*entities.Note, error)
	SearchByContent(userID uint, query string, limit int) ([]*entities.Note, error)
	Update(note *entities.Note) (*entities.Note, error)
	Delete(id uint) error
}
