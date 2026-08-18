package repositories

import (
	"time"

	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

type todoRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.TodoRepository = (*todoRepositoryImpl)(nil)

func NewTodoRepository(db *gorm.DB) *todoRepositoryImpl {
	return &todoRepositoryImpl{db: db}
}

func (r *todoRepositoryImpl) Create(todo *entities.Todo) (*entities.Todo, error) {
	m := models.TodoFromEntity(todo)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *todoRepositoryImpl) GetByID(id uint) (*entities.Todo, error) {
	var m models.Todo
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *todoRepositoryImpl) ListByUserID(userID uint, limit, offset int) ([]*entities.Todo, error) {
	var rows []models.Todo
	if err := r.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(normalizeLimit(limit)).
		Offset(normalizeOffset(offset)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return todosToEntities(rows), nil
}

func (r *todoRepositoryImpl) ListPendingByUserID(userID uint) ([]*entities.Todo, error) {
	var rows []models.Todo
	if err := r.db.Where("user_id = ? AND status = ?", userID, "pending").
		Order("due_at asc nulls last, created_at desc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return todosToEntities(rows), nil
}

func (r *todoRepositoryImpl) ListDueBetween(userID uint, start, end time.Time) ([]*entities.Todo, error) {
	var rows []models.Todo
	if err := r.db.Where("user_id = ? AND due_at >= ? AND due_at < ?", userID, start, end).
		Order("due_at asc").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return todosToEntities(rows), nil
}

func (r *todoRepositoryImpl) Update(todo *entities.Todo) (*entities.Todo, error) {
	m := models.TodoFromEntity(todo)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *todoRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Todo{}, id).Error
}

func todosToEntities(rows []models.Todo) []*entities.Todo {
	result := make([]*entities.Todo, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result
}
