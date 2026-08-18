package entities

import (
	"encoding/json"
	"time"
)

type MessageLog struct {
	ID             uint            `json:"id"`
	UserID         uint            `json:"user_id"`
	Source         string          `json:"source"`
	Direction      string          `json:"direction"`
	MessageType    string          `json:"message_type"`
	MessageText    string          `json:"message_text"`
	RawPayload     json.RawMessage `json:"raw_payload,omitempty"`
	LineReplyToken *string         `json:"line_reply_token,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
}
