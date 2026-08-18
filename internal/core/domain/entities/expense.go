package entities

import "time"

type Expense struct {
	ID              uint      `json:"id"`
	UserID          uint      `json:"user_id"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	Category        string    `json:"category"`
	Description     string    `json:"description"`
	SpentAt         time.Time `json:"spent_at"`
	SourceMessageID *uint     `json:"source_message_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
