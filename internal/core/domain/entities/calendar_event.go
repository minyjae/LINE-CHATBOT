package entities

import "time"

type CalendarEvent struct {
	ID              uint       `json:"id"`
	UserID          uint       `json:"user_id"`
	Title           string     `json:"title"`
	Description     *string    `json:"description,omitempty"`
	StartAt         time.Time  `json:"start_at"`
	EndAt           *time.Time `json:"end_at,omitempty"`
	Location        *string    `json:"location,omitempty"`
	GoogleEventID   *string    `json:"google_event_id,omitempty"`
	SyncStatus      string     `json:"sync_status"`
	SourceMessageID *uint      `json:"source_message_id,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
