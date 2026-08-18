package entities

import "time"

type Todo struct {
	ID              uint       `json:"id"`
	UserID          uint       `json:"user_id"`
	Title           string     `json:"title"`
	Description     *string    `json:"description,omitempty"`
	Status          string     `json:"status"`
	DueAt           *time.Time `json:"due_at,omitempty"`
	Priority        string     `json:"priority"`
	SourceMessageID *uint      `json:"source_message_id,omitempty"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (t *Todo) IsCompleted() bool {
	return t != nil && t.Status == "completed"
}
