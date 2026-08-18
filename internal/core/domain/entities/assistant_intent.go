package entities

import (
	"encoding/json"
	"time"
)

type AssistantIntent struct {
	ID           uint            `json:"id"`
	UserID       uint            `json:"user_id"`
	MessageLogID *uint           `json:"message_log_id,omitempty"`
	Intent       string          `json:"intent"`
	Confidence   float64         `json:"confidence"`
	Entities     json.RawMessage `json:"entities,omitempty"`
	Status       string          `json:"status"`
	ErrorMessage *string         `json:"error_message,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}
