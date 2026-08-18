package entities

import "time"

type IntegrationAccount struct {
	ID                    uint       `json:"id"`
	UserID                uint       `json:"user_id"`
	Provider              string     `json:"provider"`
	ProviderAccountID     *string    `json:"provider_account_id,omitempty"`
	AccessTokenEncrypted  string     `json:"-"`
	RefreshTokenEncrypted *string    `json:"-"`
	TokenExpiresAt        *time.Time `json:"token_expires_at,omitempty"`
	Scopes                []string   `json:"scopes,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}
