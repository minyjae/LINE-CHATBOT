package services

import "minyjae/go-starter/types"

type AssistantIntentParser interface {
	Parse(input types.IntentParseInput) (*types.ParsedAssistantIntent, error)
}
