package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	"minyjae/go-starter/types"
)

type assistantService struct {
	intentRepo   repoPort.AssistantIntentRepository
	todoRepo     repoPort.TodoRepository
	expenseRepo  repoPort.ExpenseRepository
	calendarRepo repoPort.CalendarEventRepository
	reminderRepo repoPort.ReminderRepository
	noteRepo     repoPort.NoteRepository
}

var _ servicePort.AssistantService = (*assistantService)(nil)

func NewAssistantServiceImpl(
	intentRepo repoPort.AssistantIntentRepository,
	todoRepo repoPort.TodoRepository,
	expenseRepo repoPort.ExpenseRepository,
	calendarRepo repoPort.CalendarEventRepository,
	reminderRepo repoPort.ReminderRepository,
	noteRepo repoPort.NoteRepository,
) *assistantService {
	return &assistantService{
		intentRepo:   intentRepo,
		todoRepo:     todoRepo,
		expenseRepo:  expenseRepo,
		calendarRepo: calendarRepo,
		reminderRepo: reminderRepo,
		noteRepo:     noteRepo,
	}
}

func (s *assistantService) HandleTextMessage(input types.AssistantMessageInput) (*types.AssistantMessageResult, error) {
	text := strings.TrimSpace(input.Text)
	if text == "" {
		return s.saveUnknownIntent(input, "ยังไม่เห็นข้อความที่ต้องการให้ช่วยครับ")
	}

	loc := loadLocation(input.Timezone)
	now := input.Now.In(loc)

	if isTomorrowSummary(text) {
		return s.handleTomorrowSummary(input, now, loc)
	}
	if isExpenseSummary(text) {
		return s.handleExpenseSummary(input, now)
	}
	if isTodoCommand(text) {
		return s.handleCreateTodo(input, text, now)
	}
	if isReminderCommand(text) {
		return s.handleCreateReminder(input, text, now, loc)
	}
	if isNoteCommand(text) {
		return s.handleCreateNote(input, text, now)
	}
	if amount, ok := extractAmount(text); ok && looksLikeExpense(text) {
		return s.handleCreateExpense(input, text, amount, now)
	}
	if looksLikeCalendarEvent(text) {
		return s.handleCreateCalendarEvent(input, text, now, loc)
	}

	return s.saveUnknownIntent(input, "ตอนนี้ผมยังเข้าใจได้เฉพาะ todo, รายจ่าย, นัดหมาย, เตือนความจำ, note และสรุปพรุ่งนี้ครับ")
}

func (s *assistantService) handleCreateTodo(input types.AssistantMessageInput, text string, now time.Time) (*types.AssistantMessageResult, error) {
	title := cleanupTodoTitle(text)
	if title == "" {
		title = text
	}

	todo, err := s.todoRepo.Create(&entities.Todo{
		UserID:          input.UserID,
		Title:           title,
		Status:          "pending",
		Priority:        "normal",
		SourceMessageID: input.MessageLogID,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	if err != nil {
		return nil, err
	}

	entitiesJSON := marshalEntities(map[string]any{"title": todo.Title})
	intent, err := s.intentRepo.Create(&entities.AssistantIntent{
		UserID:       input.UserID,
		MessageLogID: input.MessageLogID,
		Intent:       "create_todo",
		Confidence:   0.9,
		Entities:     entitiesJSON,
		Status:       "completed",
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return nil, err
	}

	return &types.AssistantMessageResult{
		Intent:    intent.Intent,
		ReplyText: fmt.Sprintf("เพิ่ม todo \"%s\" แล้วครับ", todo.Title),
		Data:      todo,
	}, nil
}

func (s *assistantService) handleCreateExpense(input types.AssistantMessageInput, text string, amount float64, now time.Time) (*types.AssistantMessageResult, error) {
	category := inferExpenseCategory(text)
	description := cleanupExpenseDescription(text)

	expense, err := s.expenseRepo.Create(&entities.Expense{
		UserID:          input.UserID,
		Amount:          amount,
		Currency:        "THB",
		Category:        category,
		Description:     description,
		SpentAt:         now,
		SourceMessageID: input.MessageLogID,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	if err != nil {
		return nil, err
	}

	entitiesJSON := marshalEntities(map[string]any{
		"amount":      expense.Amount,
		"currency":    expense.Currency,
		"category":    expense.Category,
		"description": expense.Description,
		"spent_at":    expense.SpentAt,
	})
	intent, err := s.intentRepo.Create(&entities.AssistantIntent{
		UserID:       input.UserID,
		MessageLogID: input.MessageLogID,
		Intent:       "create_expense",
		Confidence:   0.9,
		Entities:     entitiesJSON,
		Status:       "completed",
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return nil, err
	}

	return &types.AssistantMessageResult{
		Intent:    intent.Intent,
		ReplyText: fmt.Sprintf("บันทึก%s %.2f บาทแล้วครับ", expense.Description, expense.Amount),
		Data:      expense,
	}, nil
}

func (s *assistantService) handleExpenseSummary(input types.AssistantMessageInput, now time.Time) (*types.AssistantMessageResult, error) {
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	nextMonth := monthStart.AddDate(0, 1, 0)
	total, err := s.expenseRepo.SumBySpentAtBetween(input.UserID, monthStart, nextMonth)
	if err != nil {
		return nil, err
	}

	entitiesJSON := marshalEntities(map[string]any{
		"start": monthStart,
		"end":   nextMonth,
		"total": total,
	})
	intent, err := s.intentRepo.Create(&entities.AssistantIntent{
		UserID:       input.UserID,
		MessageLogID: input.MessageLogID,
		Intent:       "expense_summary",
		Confidence:   0.88,
		Entities:     entitiesJSON,
		Status:       "completed",
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return nil, err
	}

	return &types.AssistantMessageResult{
		Intent:    intent.Intent,
		ReplyText: fmt.Sprintf("เดือนนี้ใช้เงินไปแล้ว %.2f บาทครับ", total),
		Data:      map[string]any{"total": total, "start": monthStart, "end": nextMonth},
	}, nil
}

func (s *assistantService) handleCreateCalendarEvent(input types.AssistantMessageInput, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	startAt, ok := parseNaturalTime(text, now, loc, 9, 0)
	if !ok {
		startAt = now.Add(24 * time.Hour)
	}
	endAt := startAt.Add(time.Hour)
	title := cleanupCalendarTitle(text)
	if title == "" {
		title = text
	}

	event, err := s.calendarRepo.Create(&entities.CalendarEvent{
		UserID:          input.UserID,
		Title:           title,
		StartAt:         startAt,
		EndAt:           &endAt,
		SyncStatus:      "local",
		SourceMessageID: input.MessageLogID,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	if err != nil {
		return nil, err
	}

	entitiesJSON := marshalEntities(map[string]any{"title": event.Title, "start_at": event.StartAt, "end_at": event.EndAt})
	intent, err := s.intentRepo.Create(&entities.AssistantIntent{
		UserID:       input.UserID,
		MessageLogID: input.MessageLogID,
		Intent:       "create_calendar_event",
		Confidence:   0.82,
		Entities:     entitiesJSON,
		Status:       "completed",
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return nil, err
	}

	return &types.AssistantMessageResult{
		Intent:    intent.Intent,
		ReplyText: fmt.Sprintf("เพิ่มนัด \"%s\" เวลา %s แล้วครับ", event.Title, event.StartAt.Format("02 Jan 2006 15:04")),
		Data:      event,
	}, nil
}

func (s *assistantService) handleCreateReminder(input types.AssistantMessageInput, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	remindAt, ok := parseNaturalTime(text, now, loc, 9, 0)
	if !ok {
		remindAt = now.Add(24 * time.Hour)
	}
	title := cleanupReminderTitle(text)
	if title == "" {
		title = text
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

	entitiesJSON := marshalEntities(map[string]any{"title": reminder.Title, "remind_at": reminder.RemindAt})
	intent, err := s.intentRepo.Create(&entities.AssistantIntent{
		UserID:       input.UserID,
		MessageLogID: input.MessageLogID,
		Intent:       "create_reminder",
		Confidence:   0.86,
		Entities:     entitiesJSON,
		Status:       "completed",
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return nil, err
	}

	return &types.AssistantMessageResult{
		Intent:    intent.Intent,
		ReplyText: fmt.Sprintf("ตั้งเตือน \"%s\" เวลา %s แล้วครับ", reminder.Title, reminder.RemindAt.Format("02 Jan 2006 15:04")),
		Data:      reminder,
	}, nil
}

func (s *assistantService) handleCreateNote(input types.AssistantMessageInput, text string, now time.Time) (*types.AssistantMessageResult, error) {
	content := cleanupNoteContent(text)
	if content == "" {
		content = text
	}

	note, err := s.noteRepo.Create(&entities.Note{
		UserID:          input.UserID,
		Content:         content,
		SourceMessageID: input.MessageLogID,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	if err != nil {
		return nil, err
	}

	entitiesJSON := marshalEntities(map[string]any{"content": note.Content})
	intent, err := s.intentRepo.Create(&entities.AssistantIntent{
		UserID:       input.UserID,
		MessageLogID: input.MessageLogID,
		Intent:       "create_note",
		Confidence:   0.88,
		Entities:     entitiesJSON,
		Status:       "completed",
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return nil, err
	}

	return &types.AssistantMessageResult{
		Intent:    intent.Intent,
		ReplyText: "จดบันทึกไว้แล้วครับ",
		Data:      note,
	}, nil
}

func (s *assistantService) handleTomorrowSummary(input types.AssistantMessageInput, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	tomorrow := now.AddDate(0, 0, 1)
	start := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, loc)
	end := start.AddDate(0, 0, 1)

	events, err := s.calendarRepo.ListByStartBetween(input.UserID, start, end)
	if err != nil {
		return nil, err
	}
	todos, err := s.todoRepo.ListDueBetween(input.UserID, start, end)
	if err != nil {
		return nil, err
	}

	parts := []string{}
	for _, event := range events {
		parts = append(parts, fmt.Sprintf("%s %s", event.StartAt.In(loc).Format("15:04"), event.Title))
	}
	for _, todo := range todos {
		parts = append(parts, fmt.Sprintf("todo: %s", todo.Title))
	}

	reply := "พรุ่งนี้ยังไม่มีนัดหรือ todo ครับ"
	if len(parts) > 0 {
		reply = "พรุ่งนี้มี:\n- " + strings.Join(parts, "\n- ")
	}

	entitiesJSON := marshalEntities(map[string]any{"start": start, "end": end, "event_count": len(events), "todo_count": len(todos)})
	intent, err := s.intentRepo.Create(&entities.AssistantIntent{
		UserID:       input.UserID,
		MessageLogID: input.MessageLogID,
		Intent:       "tomorrow_summary",
		Confidence:   0.9,
		Entities:     entitiesJSON,
		Status:       "completed",
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return nil, err
	}

	return &types.AssistantMessageResult{
		Intent:    intent.Intent,
		ReplyText: reply,
		Data:      map[string]any{"events": events, "todos": todos},
	}, nil
}

func (s *assistantService) saveUnknownIntent(input types.AssistantMessageInput, reply string) (*types.AssistantMessageResult, error) {
	now := input.Now
	intent, err := s.intentRepo.Create(&entities.AssistantIntent{
		UserID:       input.UserID,
		MessageLogID: input.MessageLogID,
		Intent:       "unknown",
		Confidence:   0.2,
		Status:       "completed",
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return nil, err
	}

	return &types.AssistantMessageResult{Intent: intent.Intent, ReplyText: reply}, nil
}

func loadLocation(timezone string) *time.Location {
	if timezone == "" {
		timezone = "Asia/Bangkok"
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.FixedZone("Asia/Bangkok", 7*60*60)
	}
	return loc
}

func marshalEntities(value any) json.RawMessage {
	data, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	return data
}

func isTodoCommand(text string) bool {
	lower := strings.ToLower(text)
	return strings.Contains(lower, "todo") || strings.Contains(text, "ทูดู") || strings.Contains(text, "เพิ่มงาน")
}

func isExpenseSummary(text string) bool {
	return strings.Contains(text, "ใช้เงิน") && (strings.Contains(text, "เดือน") || strings.Contains(text, "เท่าไหร่"))
}

func isTomorrowSummary(text string) bool {
	return strings.Contains(text, "พรุ่งนี้") && (strings.Contains(text, "มีอะไร") || strings.Contains(text, "ทำอะไร"))
}

func isReminderCommand(text string) bool {
	return strings.Contains(text, "เตือน")
}

func isNoteCommand(text string) bool {
	lower := strings.ToLower(text)
	return strings.Contains(text, "จด") || strings.Contains(lower, "note")
}

func looksLikeExpense(text string) bool {
	keywords := []string{"บาท", "กิน", "ซื้อ", "จ่าย", "ค่า", "โอน", "กาแฟ", "ข้าว", "อาหาร"}
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

func looksLikeCalendarEvent(text string) bool {
	keywords := []string{"ประชุม", "นัด", "เจอ", "คุย", "พรุ่งนี้", "วันนี้", "โมง"}
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

func extractAmount(text string) (float64, bool) {
	re := regexp.MustCompile(`([0-9]+(?:\.[0-9]+)?)`)
	match := re.FindStringSubmatch(text)
	if len(match) < 2 {
		return 0, false
	}
	amount, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0, false
	}
	return amount, true
}

func inferExpenseCategory(text string) string {
	if strings.Contains(text, "ข้าว") || strings.Contains(text, "กิน") || strings.Contains(text, "อาหาร") || strings.Contains(text, "กาแฟ") {
		return "food"
	}
	if strings.Contains(text, "รถ") || strings.Contains(text, "แท็กซี่") || strings.Contains(text, "เดินทาง") {
		return "transport"
	}
	return "uncategorized"
}

func cleanupExpenseDescription(text string) string {
	re := regexp.MustCompile(`[0-9]+(?:\.[0-9]+)?\s*(บาท|บ|thb|THB)?`)
	cleaned := strings.TrimSpace(re.ReplaceAllString(text, ""))
	if cleaned == "" {
		return "รายจ่าย"
	}
	return cleaned
}

func cleanupTodoTitle(text string) string {
	replacers := []string{"เพิ่ม todo", "เพิ่ม Todo", "todo", "Todo", "เพิ่มทูดู", "ทูดู", "เพิ่มงาน"}
	return cleanupByRemoving(text, replacers)
}

func cleanupReminderTitle(text string) string {
	title := cleanupByRemoving(text, []string{"เตือน", "ช่วยเตือน", "แจ้งเตือน"})
	return cleanupTimeWords(title)
}

func cleanupNoteContent(text string) string {
	return cleanupByRemoving(text, []string{"จดไว้ว่า", "จดว่า", "จด", "note", "Note"})
}

func cleanupCalendarTitle(text string) string {
	return cleanupTimeWords(cleanupByRemoving(text, []string{"นัด", "เพิ่มนัด", "calendar", "Calendar"}))
}

func cleanupByRemoving(text string, tokens []string) string {
	result := text
	for _, token := range tokens {
		result = strings.ReplaceAll(result, token, "")
	}
	return strings.TrimSpace(result)
}

func cleanupTimeWords(text string) string {
	replacers := []string{"วันนี้", "พรุ่งนี้", "โมง", "ตอน", "เวลา"}
	result := text
	for _, token := range replacers {
		result = strings.ReplaceAll(result, token, "")
	}
	re := regexp.MustCompile(`\d{1,2}(:\d{2})?`)
	result = re.ReplaceAllString(result, "")
	return strings.TrimSpace(result)
}

func parseNaturalTime(text string, now time.Time, loc *time.Location, defaultHour, defaultMinute int) (time.Time, bool) {
	date := now.In(loc)
	if strings.Contains(text, "พรุ่งนี้") {
		date = date.AddDate(0, 0, 1)
	}

	hour, minute, hasTime := extractHourMinute(text)
	if !hasTime {
		hour = defaultHour
		minute = defaultMinute
	}

	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, loc), strings.Contains(text, "วันนี้") || strings.Contains(text, "พรุ่งนี้") || hasTime
}

func extractHourMinute(text string) (int, int, bool) {
	colon := regexp.MustCompile(`(\d{1,2})[:.](\d{2})`)
	if match := colon.FindStringSubmatch(text); len(match) == 3 {
		hour, _ := strconv.Atoi(match[1])
		minute, _ := strconv.Atoi(match[2])
		return normalizeThaiHour(text, hour), minute, true
	}

	clock := regexp.MustCompile(`(\d{1,2})\s*โมง`)
	if match := clock.FindStringSubmatch(text); len(match) == 2 {
		hour, _ := strconv.Atoi(match[1])
		return normalizeThaiHour(text, hour), 0, true
	}

	return 0, 0, false
}

func normalizeThaiHour(text string, hour int) int {
	if strings.Contains(text, "บ่าย") && hour < 12 {
		return hour + 12
	}
	if (strings.Contains(text, "เย็น") || strings.Contains(text, "ค่ำ")) && hour < 12 {
		return hour + 12
	}
	return hour
}
