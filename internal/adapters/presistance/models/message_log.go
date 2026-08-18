package models

import (
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	"minyjae/go-starter/utils"
)

type MessageLog struct {
	ID             uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint        `gorm:"not null;index" json:"user_id"`
	Source         string      `gorm:"type:varchar(32);not null;index" json:"source"`
	Direction      string      `gorm:"type:varchar(16);not null;index" json:"direction"`
	MessageType    string      `gorm:"type:varchar(32);not null;default:'text'" json:"message_type"`
	MessageText    string      `gorm:"type:text" json:"message_text"`
	RawPayload     utils.JSONB `gorm:"type:jsonb" json:"raw_payload,omitempty"`
	LineReplyToken *string     `gorm:"type:varchar(255)" json:"line_reply_token,omitempty"`
	CreatedAt      time.Time   `gorm:"autoCreateTime;index" json:"created_at"`
}

func (m *MessageLog) ToEntity() *entities.MessageLog {
	if m == nil {
		return nil
	}

	return &entities.MessageLog{
		ID:             m.ID,
		UserID:         m.UserID,
		Source:         m.Source,
		Direction:      m.Direction,
		MessageType:    m.MessageType,
		MessageText:    m.MessageText,
		RawPayload:     m.RawPayload.ToRawMessage(),
		LineReplyToken: m.LineReplyToken,
		CreatedAt:      m.CreatedAt,
	}
}

func MessageLogFromEntity(e *entities.MessageLog) *MessageLog {
	if e == nil {
		return nil
	}

	return &MessageLog{
		ID:             e.ID,
		UserID:         e.UserID,
		Source:         e.Source,
		Direction:      e.Direction,
		MessageType:    e.MessageType,
		MessageText:    e.MessageText,
		RawPayload:     utils.JSONB(e.RawPayload),
		LineReplyToken: e.LineReplyToken,
		CreatedAt:      e.CreatedAt,
	}
}
