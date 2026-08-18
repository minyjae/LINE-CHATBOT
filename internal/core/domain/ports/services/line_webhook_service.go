package services

import "minyjae/go-starter/types"

type LineWebhookService interface {
	HandleTextMessage(input types.LineTextMessageInput) (*types.LineTextMessageResult, error)
}
