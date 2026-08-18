package models

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

type Todo struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uint       `gorm:"not null;index" json:"user_id"`
	Title           string     `gorm:"type:varchar(255);not null" json:"title"`
	Description     *string    `gorm:"type:text" json:"description,omitempty"`
	Status          string     `gorm:"type:varchar(32);not null;default:'pending';index" json:"status"`
	DueAt           *time.Time `gorm:"index" json:"due_at,omitempty"`
	Priority        string     `gorm:"type:varchar(32);not null;default:'normal';index" json:"priority"`
	SourceMessageID *uint      `gorm:"index" json:"source_message_id,omitempty"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	CreatedAt       time.Time  `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"index" json:"-"`
}

func (m *Todo) ToEntity() *entities.Todo {
	if m == nil {
		return nil
	}

	return &entities.Todo{
		ID:              m.ID,
		UserID:          m.UserID,
		Title:           m.Title,
		Description:     m.Description,
		Status:          m.Status,
		DueAt:           m.DueAt,
		Priority:        m.Priority,
		SourceMessageID: m.SourceMessageID,
		CompletedAt:     m.CompletedAt,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func TodoFromEntity(e *entities.Todo) *Todo {
	if e == nil {
		return nil
	}

	return &Todo{
		ID:              e.ID,
		UserID:          e.UserID,
		Title:           e.Title,
		Description:     e.Description,
		Status:          e.Status,
		DueAt:           e.DueAt,
		Priority:        e.Priority,
		SourceMessageID: e.SourceMessageID,
		CompletedAt:     e.CompletedAt,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}
