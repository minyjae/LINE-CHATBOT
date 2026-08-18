package entities

import "time"

type Note struct {
	ID              uint      `json:"id"`
	UserID          uint      `json:"user_id"`
	Content         string    `json:"content"`
	Tags            []string  `json:"tags,omitempty"`
	SourceMessageID *uint     `json:"source_message_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
