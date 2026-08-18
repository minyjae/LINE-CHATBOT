package models

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

type CalendarEvent struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uint       `gorm:"not null;index" json:"user_id"`
	Title           string     `gorm:"type:varchar(255);not null" json:"title"`
	Description     *string    `gorm:"type:text" json:"description,omitempty"`
	StartAt         time.Time  `gorm:"not null;index" json:"start_at"`
	EndAt           *time.Time `gorm:"index" json:"end_at,omitempty"`
	Location        *string    `gorm:"type:text" json:"location,omitempty"`
	GoogleEventID   *string    `gorm:"type:varchar(255);index" json:"google_event_id,omitempty"`
	SyncStatus      string     `gorm:"type:varchar(32);not null;default:'local';index" json:"sync_status"`
	SourceMessageID *uint      `gorm:"index" json:"source_message_id,omitempty"`
	CreatedAt       time.Time  `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"index" json:"-"`
}

func (m *CalendarEvent) ToEntity() *entities.CalendarEvent {
	if m == nil {
		return nil
	}

	return &entities.CalendarEvent{
		ID:              m.ID,
		UserID:          m.UserID,
		Title:           m.Title,
		Description:     m.Description,
		StartAt:         m.StartAt,
		EndAt:           m.EndAt,
		Location:        m.Location,
		GoogleEventID:   m.GoogleEventID,
		SyncStatus:      m.SyncStatus,
		SourceMessageID: m.SourceMessageID,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func CalendarEventFromEntity(e *entities.CalendarEvent) *CalendarEvent {
	if e == nil {
		return nil
	}

	return &CalendarEvent{
		ID:              e.ID,
		UserID:          e.UserID,
		Title:           e.Title,
		Description:     e.Description,
		StartAt:         e.StartAt,
		EndAt:           e.EndAt,
		Location:        e.Location,
		GoogleEventID:   e.GoogleEventID,
		SyncStatus:      e.SyncStatus,
		SourceMessageID: e.SourceMessageID,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}
