package repositories

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

type ConversationSessionRepository interface {
	Create(session *entities.ConversationSession) (*entities.ConversationSession, error)
	GetByUserIDAndKey(userID uint, sessionKey string) (*entities.ConversationSession, error)
	Upsert(session *entities.ConversationSession) (*entities.ConversationSession, error)
	DeleteExpired(now time.Time) error
}
