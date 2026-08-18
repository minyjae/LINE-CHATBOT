package entities

import "time"

type User struct {
	ID           uint       `json:"id"`
	Email        *string    `json:"email,omitempty"`
	PasswordHash *string    `json:"-"`
	DisplayName  string     `json:"display_name"`
	Timezone     string     `json:"timezone"`
	Locale       string     `json:"locale"`
	Status       string     `json:"status"`
	LastSeenAt   *time.Time `json:"last_seen_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (u *User) IsActive() bool {
	return u != nil && u.Status == "active"
}
