package entities

import "time"

type Income struct {
	ID              uint      `json:"id"`
	UserID          uint      `json:"user_id"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	Category        string    `json:"category"`
	Description     string    `json:"description"`
	ReceivedAt      time.Time `json:"received_at"`
	SourceMessageID *uint     `json:"source_message_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
