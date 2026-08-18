package entities

import (
	"encoding/json"
	"time"
)

type ConversationSession struct {
	ID         uint            `json:"id"`
	UserID     uint            `json:"user_id"`
	SessionKey string          `json:"session_key"`
	Context    json.RawMessage `json:"context,omitempty"`
	ExpiresAt  time.Time       `json:"expires_at"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

func (s *ConversationSession) IsExpired(now time.Time) bool {
	return s != nil && !s.ExpiresAt.After(now)
}
