package models

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

type Income struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uint       `gorm:"not null;index" json:"user_id"`
	Amount          float64    `gorm:"type:numeric(14,2);not null" json:"amount"`
	Currency        string     `gorm:"type:varchar(3);not null;default:'THB';index" json:"currency"`
	Category        string     `gorm:"type:varchar(64);not null;default:'uncategorized';index" json:"category"`
	Description     string     `gorm:"type:text" json:"description"`
	ReceivedAt      time.Time  `gorm:"not null;index" json:"received_at"`
	SourceMessageID *uint      `gorm:"index" json:"source_message_id,omitempty"`
	CreatedAt       time.Time  `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"index" json:"-"`
}

func (m *Income) ToEntity() *entities.Income {
	if m == nil {
		return nil
	}

	return &entities.Income{
		ID:              m.ID,
		UserID:          m.UserID,
		Amount:          m.Amount,
		Currency:        m.Currency,
		Category:        m.Category,
		Description:     m.Description,
		ReceivedAt:      m.ReceivedAt,
		SourceMessageID: m.SourceMessageID,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func IncomeFromEntity(e *entities.Income) *Income {
	if e == nil {
		return nil
	}

	return &Income{
		ID:              e.ID,
		UserID:          e.UserID,
		Amount:          e.Amount,
		Currency:        e.Currency,
		Category:        e.Category,
		Description:     e.Description,
		ReceivedAt:      e.ReceivedAt,
		SourceMessageID: e.SourceMessageID,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}
