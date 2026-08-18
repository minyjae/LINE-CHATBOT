package models

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	"minyjae/go-starter/utils"
)

type ConversationSession struct {
	ID         uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint        `gorm:"not null;uniqueIndex:idx_user_session_key" json:"user_id"`
	SessionKey string      `gorm:"type:varchar(128);not null;uniqueIndex:idx_user_session_key" json:"session_key"`
	Context    utils.JSONB `gorm:"type:jsonb" json:"context,omitempty"`
	ExpiresAt  time.Time   `gorm:"not null;index" json:"expires_at"`
	CreatedAt  time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

func (m *ConversationSession) ToEntity() *entities.ConversationSession {
	if m == nil {
		return nil
	}

	return &entities.ConversationSession{
		ID:         m.ID,
		UserID:     m.UserID,
		SessionKey: m.SessionKey,
		Context:    m.Context.ToRawMessage(),
		ExpiresAt:  m.ExpiresAt,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func ConversationSessionFromEntity(e *entities.ConversationSession) *ConversationSession {
	if e == nil {
		return nil
	}

	return &ConversationSession{
		ID:         e.ID,
		UserID:     e.UserID,
		SessionKey: e.SessionKey,
		Context:    utils.JSONB(e.Context),
		ExpiresAt:  e.ExpiresAt,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}
