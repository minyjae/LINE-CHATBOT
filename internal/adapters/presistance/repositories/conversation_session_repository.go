package repositories

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"minyjae/go-starter/internal/adapters/presistance/models"
	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
)

type conversationSessionRepositoryImpl struct {
	db *gorm.DB
}

var _ repoPort.ConversationSessionRepository = (*conversationSessionRepositoryImpl)(nil)

func NewConversationSessionRepository(db *gorm.DB) *conversationSessionRepositoryImpl {
	return &conversationSessionRepositoryImpl{db: db}
}

func (r *conversationSessionRepositoryImpl) Create(session *entities.ConversationSession) (*entities.ConversationSession, error) {
	m := models.ConversationSessionFromEntity(session)
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *conversationSessionRepositoryImpl) GetByUserIDAndKey(userID uint, sessionKey string) (*entities.ConversationSession, error) {
	var m models.ConversationSession
	if err := r.db.Where("user_id = ? AND session_key = ?", userID, sessionKey).First(&m).Error; err != nil {
		return nil, err
	}
	return m.ToEntity(), nil
}

func (r *conversationSessionRepositoryImpl) Upsert(session *entities.ConversationSession) (*entities.ConversationSession, error) {
	m := models.ConversationSessionFromEntity(session)
	if err := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "session_key"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"context",
			"expires_at",
			"updated_at",
		}),
	}).Create(m).Error; err != nil {
		return nil, err
	}

	return r.GetByUserIDAndKey(m.UserID, m.SessionKey)
}

func (r *conversationSessionRepositoryImpl) DeleteExpired(now time.Time) error {
	return r.db.Where("expires_at <= ?", now).Delete(&models.ConversationSession{}).Error
}
