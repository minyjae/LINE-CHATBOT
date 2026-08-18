package models

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

type User struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Email        *string    `gorm:"type:varchar(255);uniqueIndex" json:"email,omitempty"`
	PasswordHash *string    `gorm:"type:varchar(255)" json:"-"`
	DisplayName  string     `gorm:"type:varchar(255);not null" json:"display_name"`
	Timezone     string     `gorm:"type:varchar(64);not null;default:'Asia/Bangkok'" json:"timezone"`
	Locale       string     `gorm:"type:varchar(16);not null;default:'th-TH'" json:"locale"`
	Status       string     `gorm:"type:varchar(32);not null;default:'active';index" json:"status"`
	LastSeenAt   *time.Time `json:"last_seen_at,omitempty"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"-"`
}

func (m *User) ToEntity() *entities.User {
	if m == nil {
		return nil
	}

	return &entities.User{
		ID:           m.ID,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		DisplayName:  m.DisplayName,
		Timezone:     m.Timezone,
		Locale:       m.Locale,
		Status:       m.Status,
		LastSeenAt:   m.LastSeenAt,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func FromEntity(e *entities.User) *User {
	if e == nil {
		return nil
	}

	return &User{
		ID:           e.ID,
		Email:        e.Email,
		PasswordHash: e.PasswordHash,
		DisplayName:  e.DisplayName,
		Timezone:     e.Timezone,
		Locale:       e.Locale,
		Status:       e.Status,
		LastSeenAt:   e.LastSeenAt,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}
