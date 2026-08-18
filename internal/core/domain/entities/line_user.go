package entities

import "time"

type LineUser struct {
	ID          uint       `json:"id"`
	UserID      uint       `json:"user_id"`
	LineUserID  string     `json:"line_user_id"`
	DisplayName *string    `json:"display_name,omitempty"`
	PictureURL  *string    `json:"picture_url,omitempty"`
	Status      string     `json:"status"`
	LastSeenAt  *time.Time `json:"last_seen_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
