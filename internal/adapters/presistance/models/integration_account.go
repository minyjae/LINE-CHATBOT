package models

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

type IntegrationAccount struct {
	ID                    uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID                uint       `gorm:"not null;uniqueIndex:idx_user_provider" json:"user_id"`
	Provider              string     `gorm:"type:varchar(64);not null;uniqueIndex:idx_user_provider" json:"provider"`
	ProviderAccountID     *string    `gorm:"type:varchar(255);index" json:"provider_account_id,omitempty"`
	AccessTokenEncrypted  string     `gorm:"type:text;not null" json:"-"`
	RefreshTokenEncrypted *string    `gorm:"type:text" json:"-"`
	TokenExpiresAt        *time.Time `gorm:"index" json:"token_expires_at,omitempty"`
	Scopes                []string   `gorm:"serializer:json;type:jsonb" json:"scopes,omitempty"`
	CreatedAt             time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt             *time.Time `gorm:"index" json:"-"`
}

func (m *IntegrationAccount) ToEntity() *entities.IntegrationAccount {
	if m == nil {
		return nil
	}

	return &entities.IntegrationAccount{
		ID:                    m.ID,
		UserID:                m.UserID,
		Provider:              m.Provider,
		ProviderAccountID:     m.ProviderAccountID,
		AccessTokenEncrypted:  m.AccessTokenEncrypted,
		RefreshTokenEncrypted: m.RefreshTokenEncrypted,
		TokenExpiresAt:        m.TokenExpiresAt,
		Scopes:                m.Scopes,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
	}
}

func IntegrationAccountFromEntity(e *entities.IntegrationAccount) *IntegrationAccount {
	if e == nil {
		return nil
	}

	return &IntegrationAccount{
		ID:                    e.ID,
		UserID:                e.UserID,
		Provider:              e.Provider,
		ProviderAccountID:     e.ProviderAccountID,
		AccessTokenEncrypted:  e.AccessTokenEncrypted,
		RefreshTokenEncrypted: e.RefreshTokenEncrypted,
		TokenExpiresAt:        e.TokenExpiresAt,
		Scopes:                e.Scopes,
		CreatedAt:             e.CreatedAt,
		UpdatedAt:             e.UpdatedAt,
	}
}
