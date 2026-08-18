package repositories

import (
	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

type assistantIntentRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.AssistantIntentRepository = (*assistantIntentRepositoryImpl)(nil)

func NewAssistantIntentRepository(db *gorm.DB) *assistantIntentRepositoryImpl {
	return &assistantIntentRepositoryImpl{db: db}
}

func (r *assistantIntentRepositoryImpl) Create(intent *entities.AssistantIntent) (*entities.AssistantIntent, error) {
	m := models.AssistantIntentFromEntity(intent)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *assistantIntentRepositoryImpl) GetByID(id uint) (*entities.AssistantIntent, error) {
	var m models.AssistantIntent
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *assistantIntentRepositoryImpl) ListByUserID(userID uint, limit, offset int) ([]*entities.AssistantIntent, error) {
	var rows []models.AssistantIntent
	if err := r.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(normalizeLimit(limit)).
		Offset(normalizeOffset(offset)).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.AssistantIntent, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result, nil
}

func (r *assistantIntentRepositoryImpl) Update(intent *entities.AssistantIntent) (*entities.AssistantIntent, error) {
	m := models.AssistantIntentFromEntity(intent)
	if err := r.db.Save(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *assistantIntentRepositoryImpl) UpdateStatus(id uint, status string, errorMessage *string) error {
	return r.db.Model(&models.AssistantIntent{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":        status,
			"error_message": errorMessage,
		}).Error
}
