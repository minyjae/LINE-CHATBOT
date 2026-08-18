package repositories

import (
	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

type noteRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.NoteRepository = (*noteRepositoryImpl)(nil)

func NewNoteRepository(db *gorm.DB) *noteRepositoryImpl {
	return &noteRepositoryImpl{db: db}
}

func (r *noteRepositoryImpl) Create(note *entities.Note) (*entities.Note, error) {
	m := models.NoteFromEntity(note)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *noteRepositoryImpl) GetByID(id uint) (*entities.Note, error) {
	var m models.Note
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *noteRepositoryImpl) ListByUserID(userID uint, limit, offset int) ([]*entities.Note, error) {
	var rows []models.Note
	if err := r.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(normalizeLimit(limit)).
		Offset(normalizeOffset(offset)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return notesToEntities(rows), nil
}

func (r *noteRepositoryImpl) SearchByContent(userID uint, query string, limit int) ([]*entities.Note, error) {
	var rows []models.Note
	if err := r.db.Where("user_id = ? AND content ILIKE ?", userID, "%"+query+"%").
		Order("created_at desc").
		Limit(normalizeLimit(limit)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return notesToEntities(rows), nil
}

func (r *noteRepositoryImpl) Update(note *entities.Note) (*entities.Note, error) {
	m := models.NoteFromEntity(note)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *noteRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Note{}, id).Error
}

func notesToEntities(rows []models.Note) []*entities.Note {
	result := make([]*entities.Note, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result
}
