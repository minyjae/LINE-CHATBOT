package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	"minyjae/go-starter/types"
)

var ErrParserNotConfigured = errors.New("openai intent parser is not configured")

type IntentParser struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

var _ servicePort.AssistantIntentParser = (*IntentParser)(nil)

func NewIntentParser(apiKey, model string) *IntentParser {
	if model == "" {
		model = "gpt-5"
	}
	return &IntentParser{
		apiKey: apiKey,
		model:  model,
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

func (p *IntentParser) Parse(input types.IntentParseInput) (*types.ParsedAssistantIntent, error) {
	if p.apiKey == "" {
		return nil, ErrParserNotConfigured
	}

	payload := map[string]any{
		"model": p.model,
		"input": []map[string]string{
			{
				"role":    "system",
				"content": "You classify Thai or English personal assistant messages. Return only JSON matching the schema. Use ISO 8601 timestamps with timezone when extracting dates. Supported intents: create_expense, create_income, expense_summary, expense_report_daily, expense_report_weekly, expense_report_monthly, income_report_daily, income_report_weekly, income_report_monthly, cashflow_report_daily, cashflow_report_weekly, cashflow_report_monthly, create_todo, tomorrow_summary, create_reminder, create_note, create_calendar_event, cancel_calendar_event, unknown. Use expense intents for money spent, income intents for money received, and cashflow intents when the user asks for both income and expense, net, balance, or cash flow. Calendar events must include an explicit time in the user message, such as 10 โมง, 16:00, or 9 น. Do not infer or invent a time for calendar events. If a message sounds like a meeting or schedule item but has no explicit time, classify it as create_note or unknown instead of create_calendar_event. Messages such as ยกเลิกนัด, ลบนัด, cancel meeting, or cancel event should be cancel_calendar_event. Questions about what was bought, eaten, spent, or expense lists should be expense_report_* instead of calendar.",
			},
			{
				"role":    "user",
				"content": fmt.Sprintf("now=%s timezone=%s locale=%s\nmessage=%s", input.Now, input.Timezone, input.Locale, input.Text),
			},
		},
		"text": map[string]any{
			"format": map[string]any{
				"type":        "json_schema",
				"name":        "assistant_intent",
				"description": "Intent and entities for a personal assistant message",
				"strict":      true,
				"schema":      intentSchema(),
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/responses", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response responsesAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("openai responses request failed with status %d: %s", resp.StatusCode, response.Error.Message)
	}

	output := response.OutputText()
	if output == "" {
		return nil, fmt.Errorf("openai response did not include output text")
	}

	var parsed types.ParsedAssistantIntent
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		return nil, err
	}
	parsed.Raw = json.RawMessage(output)
	return &parsed, nil
}

type responsesAPIResponse struct {
	Output []struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (r responsesAPIResponse) OutputText() string {
	for _, output := range r.Output {
		for _, content := range output.Content {
			if content.Type == "output_text" && content.Text != "" {
				return content.Text
			}
		}
	}
	return ""
}

func intentSchema() map[string]any {
	entityProperties := map[string]any{
		"title":       map[string]any{"type": "string"},
		"description": map[string]any{"type": "string"},
		"content":     map[string]any{"type": "string"},
		"amount":      map[string]any{"type": "number"},
		"currency":    map[string]any{"type": "string"},
		"category":    map[string]any{"type": "string"},
		"start_at":    map[string]any{"type": "string"},
		"end_at":      map[string]any{"type": "string"},
		"due_at":      map[string]any{"type": "string"},
		"remind_at":   map[string]any{"type": "string"},
		"spent_at":    map[string]any{"type": "string"},
		"received_at": map[string]any{"type": "string"},
		"location":    map[string]any{"type": "string"},
		"priority":    map[string]any{"type": "string"},
		"tags": map[string]any{
			"type":  "array",
			"items": map[string]any{"type": "string"},
		},
	}

	entityRequired := []string{
		"title",
		"description",
		"content",
		"amount",
		"currency",
		"category",
		"start_at",
		"end_at",
		"due_at",
		"remind_at",
		"spent_at",
		"received_at",
		"location",
		"priority",
		"tags",
	}

	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"intent": map[string]any{
				"type": "string",
				"enum": []string{
					"create_expense",
					"create_income",
					"expense_summary",
					"expense_report_daily",
					"expense_report_weekly",
					"expense_report_monthly",
					"income_report_daily",
					"income_report_weekly",
					"income_report_monthly",
					"cashflow_report_daily",
					"cashflow_report_weekly",
					"cashflow_report_monthly",
					"create_todo",
					"tomorrow_summary",
					"create_reminder",
					"create_note",
					"create_calendar_event",
					"cancel_calendar_event",
					"unknown",
				},
			},
			"confidence": map[string]any{
				"type":    "number",
				"minimum": 0,
				"maximum": 1,
			},
			"entities": map[string]any{
				"type":                 "object",
				"properties":           entityProperties,
				"required":             entityRequired,
				"additionalProperties": false,
			},
		},
		"required":             []string{"intent", "confidence", "entities"},
		"additionalProperties": false,
	}
}
