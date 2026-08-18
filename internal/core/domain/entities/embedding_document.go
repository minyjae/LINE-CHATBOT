package entities

import (
	"encoding/json"
	"time"
)

type EmbeddingDocument struct {
	ID         uint            `json:"id"`
	UserID     uint            `json:"user_id"`
	SourceType string          `json:"source_type"`
	SourceID   *uint           `json:"source_id,omitempty"`
	Content    string          `json:"content"`
	Embedding  []float32       `json:"embedding,omitempty"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}
