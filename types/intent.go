package types

import "encoding/json"

type IntentParseInput struct {
	Text     string `json:"text"`
	Now      string `json:"now"`
	Timezone string `json:"timezone"`
	Locale   string `json:"locale"`
}

type ParsedAssistantIntent struct {
	Intent     string          `json:"intent"`
	Confidence float64         `json:"confidence"`
	Entities   IntentEntities  `json:"entities"`
	Raw        json.RawMessage `json:"-"`
}

type IntentEntities struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Content     string   `json:"content,omitempty"`
	Amount      float64  `json:"amount,omitempty"`
	Currency    string   `json:"currency,omitempty"`
	Category    string   `json:"category,omitempty"`
	StartAt     string   `json:"start_at,omitempty"`
	EndAt       string   `json:"end_at,omitempty"`
	DueAt       string   `json:"due_at,omitempty"`
	RemindAt    string   `json:"remind_at,omitempty"`
	SpentAt     string   `json:"spent_at,omitempty"`
	ReceivedAt  string   `json:"received_at,omitempty"`
	Location    string   `json:"location,omitempty"`
	Priority    string   `json:"priority,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}
