package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	"minyjae/go-starter/types"
)

func (s *assistantService) handleWithIntentParser(input types.AssistantMessageInput, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, bool) {
	if s.intentParser == nil {
		return nil, false
	}

	parsed, err := s.intentParser.Parse(types.IntentParseInput{
		Text:     text,
		Now:      now.Format(time.RFC3339),
		Timezone: loc.String(),
		Locale:   "th-TH",
	})
	if err != nil || parsed == nil || parsed.Confidence < 0.55 {
		return nil, false
	}

	result, err := s.handleParsedIntent(input, parsed, text, now, loc)
	if err != nil {
		return nil, false
	}
	return result, true
}

func (s *assistantService) handleParsedIntent(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	switch parsed.Intent {
	case "create_expense":
		return s.handleParsedCreateExpense(input, parsed, text, now, loc)
	case "expense_summary":
		return s.handleExpenseSummary(input, now)
	case "expense_report_daily", "expense_report_weekly", "expense_report_monthly":
		return s.handleParsedExpenseReport(input, parsed, text, now, loc)
	case "create_todo":
		return s.handleParsedCreateTodo(input, parsed, text, now, loc)
	case "tomorrow_summary":
		return s.handleTomorrowSummary(input, now, loc)
	case "create_reminder":
		return s.handleParsedCreateReminder(input, parsed, text, now, loc)
	case "create_note":
		return s.handleParsedCreateNote(input, parsed, text, now)
	case "create_calendar_event":
		return s.handleParsedCreateCalendarEvent(input, parsed, text, now, loc)
	case "cancel_calendar_event":
		return s.handleCancelCalendarEvent(input, text, now, loc)
	default:
		return s.saveParsedUnknownIntent(input, parsed, now)
	}
}

func (s *assistantService) handleParsedExpenseReport(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	period, ok := expenseReportPeriodFromParsed(parsed, now, loc)
	if !ok {
		period, ok = parseExpenseReportPeriod(text, now, loc)
	}
	if !ok {
		return nil, fmt.Errorf("parsed expense report period is missing")
	}
	return s.handleExpenseReport(input, now, loc, period)
}

func (s *assistantService) handleParsedCreateExpense(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	if parsed.Entities.Amount <= 0 {
		return nil, fmt.Errorf("parsed expense amount is missing")
	}

	spentAt := now
	if parsedAt, ok := parseOptionalIntentTime(parsed.Entities.SpentAt, loc); ok {
		spentAt = parsedAt
	}
	description := firstNonEmpty(parsed.Entities.Description, parsed.Entities.Title, cleanupExpenseDescription(text))

	expense, err := s.expenseRepo.Create(&entities.Expense{
		UserID:          input.UserID,
		Amount:          parsed.Entities.Amount,
		Currency:        defaultString(parsed.Entities.Currency, "THB"),
		Category:        defaultString(parsed.Entities.Category, "uncategorized"),
		Description:     defaultString(description, "expense"),
		SpentAt:         spentAt,
		SourceMessageID: input.MessageLogID,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	if err != nil {
		return nil, err
	}
	if err := s.saveParsedIntent(input, parsed, now); err != nil {
		return nil, err
	}

	return &types.AssistantMessageResult{
		Intent:    parsed.Intent,
		ReplyText: fmt.Sprintf("Recorded %s %.2f %s.", expense.Description, expense.Amount, expense.Currency),
		Data:      expense,
	}, nil
}

func (s *assistantService) handleParsedCreateTodo(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	title := firstNonEmpty(parsed.Entities.Title, parsed.Entities.Content, cleanupTodoTitle(text))
	if title == "" {
		return nil, fmt.Errorf("parsed todo title is missing")
	}

	var dueAt *time.Time
	if parsedDueAt, ok := parseOptionalIntentTime(parsed.Entities.DueAt, loc); ok {
		dueAt = &parsedDueAt
	}

	todo, err := s.todoRepo.Create(&entities.Todo{
		UserID:          input.UserID,
		Title:           title,
		Description:     optionalString(parsed.Entities.Description),
		Status:          "pending",
		DueAt:           dueAt,
		Priority:        defaultString(parsed.Entities.Priority, "normal"),
		SourceMessageID: input.MessageLogID,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	if err != nil {
		return nil, err
	}
	if err := s.saveParsedIntent(input, parsed, now); err != nil {
		return nil, err
	}

	return &types.AssistantMessageResult{
		Intent:    parsed.Intent,
		ReplyText: fmt.Sprintf("Added todo \"%s\".", todo.Title),
		Data:      todo,
	}, nil
}

func (s *assistantService) handleParsedCreateReminder(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	title := firstNonEmpty(parsed.Entities.Title, parsed.Entities.Content, cleanupReminderTitle(text))
	remindAt, ok := parseOptionalIntentTime(parsed.Entities.RemindAt, loc)
	if title == "" || !ok {
		return nil, fmt.Errorf("parsed reminder title or remind_at is missing")
	}

	reminder, err := s.reminderRepo.Create(&entities.Reminder{
		UserID:          input.UserID,
		Title:           title,
		RemindAt:        remindAt,
		Status:          "pending",
		SourceMessageID: input.MessageLogID,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	if err != nil {
		return nil, err
	}
	if err := s.saveParsedIntent(input, parsed, now); err != nil {
		return nil, err
	}

	return &types.AssistantMessageResult{
		Intent:    parsed.Intent,
		ReplyText: fmt.Sprintf("Reminder set for \"%s\" at %s.", reminder.Title, reminder.RemindAt.Format("02 Jan 2006 15:04")),
		Data:      reminder,
	}, nil
}

func (s *assistantService) handleParsedCreateNote(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, text string, now time.Time) (*types.AssistantMessageResult, error) {
	content := firstNonEmpty(parsed.Entities.Content, parsed.Entities.Description, cleanupNoteContent(text))
	if content == "" {
		return nil, fmt.Errorf("parsed note content is missing")
	}

	note, err := s.noteRepo.Create(&entities.Note{
		UserID:          input.UserID,
		Content:         content,
		Tags:            parsed.Entities.Tags,
		SourceMessageID: input.MessageLogID,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	if err != nil {
		return nil, err
	}
	if err := s.saveParsedIntent(input, parsed, now); err != nil {
		return nil, err
	}

	return &types.AssistantMessageResult{
		Intent:    parsed.Intent,
		ReplyText: "Saved note.",
		Data:      note,
	}, nil
}

func (s *assistantService) handleParsedCreateCalendarEvent(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	title := firstNonEmpty(parsed.Entities.Title, cleanupCalendarTitle(text))
	startAt, ok := parseOptionalIntentTime(parsed.Entities.StartAt, loc)
	if title == "" || !ok {
		return nil, fmt.Errorf("parsed calendar event title or start_at is missing")
	}

	var endAt *time.Time
	if parsedEndAt, ok := parseOptionalIntentTime(parsed.Entities.EndAt, loc); ok {
		endAt = &parsedEndAt
	}

	event, err := s.calendarRepo.Create(&entities.CalendarEvent{
		UserID:          input.UserID,
		Title:           title,
		Description:     optionalString(parsed.Entities.Description),
		StartAt:         startAt,
		EndAt:           endAt,
		Location:        optionalString(parsed.Entities.Location),
		SyncStatus:      "local",
		SourceMessageID: input.MessageLogID,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	if err != nil {
		return nil, err
	}
	if err := s.saveParsedIntent(input, parsed, now); err != nil {
		return nil, err
	}

	return &types.AssistantMessageResult{
		Intent:    parsed.Intent,
		ReplyText: fmt.Sprintf("Added event \"%s\" at %s.", event.Title, event.StartAt.Format("02 Jan 2006 15:04")),
		Data:      event,
	}, nil
}

func (s *assistantService) saveParsedUnknownIntent(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, now time.Time) (*types.AssistantMessageResult, error) {
	if err := s.saveParsedIntent(input, parsed, now); err != nil {
		return nil, err
	}
	return &types.AssistantMessageResult{
		Intent:    "unknown",
		ReplyText: "I could not understand that yet.",
	}, nil
}

func (s *assistantService) saveParsedIntent(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, now time.Time) error {
	entitiesJSON, err := json.Marshal(parsed.Entities)
	if err != nil {
		return err
	}

	_, err = s.intentRepo.Create(&entities.AssistantIntent{
		UserID:       input.UserID,
		MessageLogID: input.MessageLogID,
		Intent:       parsed.Intent,
		Confidence:   parsed.Confidence,
		Entities:     entitiesJSON,
		Status:       "completed",
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	return err
}

func parseOptionalIntentTime(value string, loc *time.Location) (time.Time, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, false
	}

	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return parsed.In(loc), true
	}
	if parsed, err := time.ParseInLocation("2006-01-02T15:04:05", value, loc); err == nil {
		return parsed, true
	}
	if parsed, err := time.ParseInLocation("2006-01-02T15:04", value, loc); err == nil {
		return parsed, true
	}
	if parsed, err := time.ParseInLocation("2006-01-02", value, loc); err == nil {
		return parsed, true
	}
	return time.Time{}, false
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func expenseReportPeriodFromParsed(parsed *types.ParsedAssistantIntent, now time.Time, loc *time.Location) (expenseReportPeriod, bool) {
	start, hasStart := parseOptionalIntentTime(parsed.Entities.StartAt, loc)
	end, hasEnd := parseOptionalIntentTime(parsed.Entities.EndAt, loc)
	if !hasStart || !hasEnd || !end.After(start) {
		return expenseReportPeriod{}, false
	}

	label := "ช่วงที่ขอ"
	switch parsed.Intent {
	case "expense_report_daily":
		label = "วันนี้"
	case "expense_report_weekly":
		label = "สัปดาห์นี้"
	case "expense_report_monthly":
		label = "เดือนนี้"
	}

	return expenseReportPeriod{
		Label:  label,
		Start:  start,
		End:    end,
		Intent: parsed.Intent,
	}, true
}
