package models

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

type LineUser struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint       `gorm:"not null;index" json:"user_id"`
	LineUserID  string     `gorm:"type:varchar(255);not null;uniqueIndex" json:"line_user_id"`
	DisplayName *string    `gorm:"type:varchar(255)" json:"display_name,omitempty"`
	PictureURL  *string    `gorm:"type:text" json:"picture_url,omitempty"`
	Status      string     `gorm:"type:varchar(32);not null;default:'active';index" json:"status"`
	LastSeenAt  *time.Time `json:"last_seen_at,omitempty"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"-"`
}

func (m *LineUser) ToEntity() *entities.LineUser {
	if m == nil {
		return nil
	}

	return &entities.LineUser{
		ID:          m.ID,
		UserID:      m.UserID,
		LineUserID:  m.LineUserID,
		DisplayName: m.DisplayName,
		PictureURL:  m.PictureURL,
		Status:      m.Status,
		LastSeenAt:  m.LastSeenAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func LineUserFromEntity(e *entities.LineUser) *LineUser {
	if e == nil {
		return nil
	}

	return &LineUser{
		ID:          e.ID,
		UserID:      e.UserID,
		LineUserID:  e.LineUserID,
		DisplayName: e.DisplayName,
		PictureURL:  e.PictureURL,
		Status:      e.Status,
		LastSeenAt:  e.LastSeenAt,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
