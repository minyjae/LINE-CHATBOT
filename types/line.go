package types

type LineWebhookRequest struct {
	Destination string             `json:"destination"`
	Events      []LineWebhookEvent `json:"events"`
}

type LineWebhookEvent struct {
	Type       string      `json:"type"`
	ReplyToken string      `json:"replyToken"`
	Source     LineSource  `json:"source"`
	Message    LineMessage `json:"message"`
}

type LineSource struct {
	Type   string `json:"type"`
	UserID string `json:"userId"`
}

type LineMessage struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Text string `json:"text"`
}
