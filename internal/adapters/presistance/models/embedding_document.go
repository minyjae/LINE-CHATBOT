package models

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	"minyjae/go-starter/utils"
)

type EmbeddingDocument struct {
	ID         uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint        `gorm:"not null;index" json:"user_id"`
	SourceType string      `gorm:"type:varchar(64);not null;index" json:"source_type"`
	SourceID   *uint       `gorm:"index" json:"source_id,omitempty"`
	Content    string      `gorm:"type:text;not null" json:"content"`
	Embedding  []float32   `gorm:"serializer:json;type:jsonb" json:"embedding,omitempty"`
	Metadata   utils.JSONB `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt  time.Time   `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt  time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

func (m *EmbeddingDocument) ToEntity() *entities.EmbeddingDocument {
	if m == nil {
		return nil
	}

	return &entities.EmbeddingDocument{
		ID:         m.ID,
		UserID:     m.UserID,
		SourceType: m.SourceType,
		SourceID:   m.SourceID,
		Content:    m.Content,
		Embedding:  m.Embedding,
		Metadata:   m.Metadata.ToRawMessage(),
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func EmbeddingDocumentFromEntity(e *entities.EmbeddingDocument) *EmbeddingDocument {
	if e == nil {
		return nil
	}

	return &EmbeddingDocument{
		ID:         e.ID,
		UserID:     e.UserID,
		SourceType: e.SourceType,
		SourceID:   e.SourceID,
		Content:    e.Content,
		Embedding:  e.Embedding,
		Metadata:   utils.JSONB(e.Metadata),
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}
