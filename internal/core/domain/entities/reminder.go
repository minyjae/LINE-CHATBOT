package entities

import "time"

type Reminder struct {
	ID              uint       `json:"id"`
	UserID          uint       `json:"user_id"`
	Title           string     `json:"title"`
	RemindAt        time.Time  `json:"remind_at"`
	Status          string     `json:"status"`
	SentAt          *time.Time `json:"sent_at,omitempty"`
	SourceMessageID *uint      `json:"source_message_id,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (r *Reminder) IsPending() bool {
	return r != nil && r.Status == "pending"
}
