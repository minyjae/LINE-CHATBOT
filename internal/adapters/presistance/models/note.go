package models

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

type Note struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uint       `gorm:"not null;index" json:"user_id"`
	Content         string     `gorm:"type:text;not null" json:"content"`
	Tags            []string   `gorm:"serializer:json;type:jsonb" json:"tags,omitempty"`
	SourceMessageID *uint      `gorm:"index" json:"source_message_id,omitempty"`
	CreatedAt       time.Time  `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"index" json:"-"`
}

func (m *Note) ToEntity() *entities.Note {
	if m == nil {
		return nil
	}

	return &entities.Note{
		ID:              m.ID,
		UserID:          m.UserID,
		Content:         m.Content,
		Tags:            m.Tags,
		SourceMessageID: m.SourceMessageID,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func NoteFromEntity(e *entities.Note) *Note {
	if e == nil {
		return nil
	}

	return &Note{
		ID:              e.ID,
		UserID:          e.UserID,
		Content:         e.Content,
		Tags:            e.Tags,
		SourceMessageID: e.SourceMessageID,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}
