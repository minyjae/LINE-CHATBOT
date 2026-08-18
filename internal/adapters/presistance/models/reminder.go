package models

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

type Reminder struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uint       `gorm:"not null;index" json:"user_id"`
	Title           string     `gorm:"type:varchar(255);not null" json:"title"`
	RemindAt        time.Time  `gorm:"not null;index" json:"remind_at"`
	Status          string     `gorm:"type:varchar(32);not null;default:'pending';index" json:"status"`
	SentAt          *time.Time `gorm:"index" json:"sent_at,omitempty"`
	SourceMessageID *uint      `gorm:"index" json:"source_message_id,omitempty"`
	CreatedAt       time.Time  `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"index" json:"-"`
}

func (m *Reminder) ToEntity() *entities.Reminder {
	if m == nil {
		return nil
	}

	return &entities.Reminder{
		ID:              m.ID,
		UserID:          m.UserID,
		Title:           m.Title,
		RemindAt:        m.RemindAt,
		Status:          m.Status,
		SentAt:          m.SentAt,
		SourceMessageID: m.SourceMessageID,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func ReminderFromEntity(e *entities.Reminder) *Reminder {
	if e == nil {
		return nil
	}

	return &Reminder{
		ID:              e.ID,
		UserID:          e.UserID,
		Title:           e.Title,
		RemindAt:        e.RemindAt,
		Status:          e.Status,
		SentAt:          e.SentAt,
		SourceMessageID: e.SourceMessageID,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}
