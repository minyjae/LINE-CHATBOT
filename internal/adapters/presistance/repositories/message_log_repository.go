package repositories

import (
	"gorm.io/gorm"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

type messageLogRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.MessageLogRepository = (*messageLogRepositoryImpl)(nil)

func NewMessageLogRepository(db *gorm.DB) *messageLogRepositoryImpl {
	return &messageLogRepositoryImpl{db: db}
}

func (r *messageLogRepositoryImpl) Create(messageLog *entities.MessageLog) (*entities.MessageLog, error) {
	m := models.MessageLogFromEntity(messageLog)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *messageLogRepositoryImpl) GetByID(id uint) (*entities.MessageLog, error) {
	var m models.MessageLog
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *messageLogRepositoryImpl) ListByUserID(userID uint, limit, offset int) ([]*entities.MessageLog, error) {
	var rows []models.MessageLog
	if err := r.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(normalizeLimit(limit)).
		Offset(normalizeOffset(offset)).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.MessageLog, 0, len(rows))
	for i := range rows {
		result = append(result, rows[i].ToEntity())
	}
	return result, nil
}
