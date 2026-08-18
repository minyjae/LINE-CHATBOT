package types

import "time"

type AssistantMessageInput struct {
	UserID       uint
	MessageLogID *uint
	Text         string
	Now          time.Time
	Timezone     string
}

type AssistantMessageResult struct {
	Intent    string      `json:"intent"`
	ReplyText string      `json:"reply_text"`
	Data      interface{} `json:"data,omitempty"`
}

type LineTextMessageInput struct {
	LineUserID string
	ReplyToken string
	Text       string
	Now        time.Time
}

type LineTextMessageResult struct {
	UserID       uint   `json:"user_id"`
	MessageLogID uint   `json:"message_log_id"`
	Intent       string `json:"intent"`
	ReplyText    string `json:"reply_text"`
}
