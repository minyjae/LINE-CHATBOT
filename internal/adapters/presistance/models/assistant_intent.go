package models

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	"minyjae/go-starter/utils"
)

type AssistantIntent struct {
	ID           uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint        `gorm:"not null;index" json:"user_id"`
	MessageLogID *uint       `gorm:"index" json:"message_log_id,omitempty"`
	Intent       string      `gorm:"type:varchar(64);not null;index" json:"intent"`
	Confidence   float64     `gorm:"type:numeric(5,4);not null;default:0" json:"confidence"`
	Entities     utils.JSONB `gorm:"type:jsonb" json:"entities,omitempty"`
	Status       string      `gorm:"type:varchar(32);not null;default:'pending';index" json:"status"`
	ErrorMessage *string     `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt    time.Time   `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

func (m *AssistantIntent) ToEntity() *entities.AssistantIntent {
	if m == nil {
		return nil
	}

	return &entities.AssistantIntent{
		ID:           m.ID,
		UserID:       m.UserID,
		MessageLogID: m.MessageLogID,
		Intent:       m.Intent,
		Confidence:   m.Confidence,
		Entities:     m.Entities.ToRawMessage(),
		Status:       m.Status,
		ErrorMessage: m.ErrorMessage,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func AssistantIntentFromEntity(e *entities.AssistantIntent) *AssistantIntent {
	if e == nil {
		return nil
	}

	return &AssistantIntent{
		ID:           e.ID,
		UserID:       e.UserID,
		MessageLogID: e.MessageLogID,
		Intent:       e.Intent,
		Confidence:   e.Confidence,
		Entities:     utils.JSONB(e.Entities),
		Status:       e.Status,
		ErrorMessage: e.ErrorMessage,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}
