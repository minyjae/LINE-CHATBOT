package services

import "minyjae/go-starter/types"

type AssistantService interface {
	HandleTextMessage(input types.AssistantMessageInput) (*types.AssistantMessageResult, error)
}
