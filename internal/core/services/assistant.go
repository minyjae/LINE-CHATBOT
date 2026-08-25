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
	incomeRepo   repoPort.IncomeRepository
	calendarRepo repoPort.CalendarEventRepository
	reminderRepo repoPort.ReminderRepository
	noteRepo     repoPort.NoteRepository
	intentParser servicePort.AssistantIntentParser
}

var _ servicePort.AssistantService = (*assistantService)(nil)

func NewAssistantServiceImpl(
	intentRepo repoPort.AssistantIntentRepository,
	todoRepo repoPort.TodoRepository,
	expenseRepo repoPort.ExpenseRepository,
	incomeRepo repoPort.IncomeRepository,
	calendarRepo repoPort.CalendarEventRepository,
	reminderRepo repoPort.ReminderRepository,
	noteRepo repoPort.NoteRepository,
	intentParser servicePort.AssistantIntentParser,
) *assistantService {
	return &assistantService{
		intentRepo:   intentRepo,
		todoRepo:     todoRepo,
		expenseRepo:  expenseRepo,
		incomeRepo:   incomeRepo,
		calendarRepo: calendarRepo,
		reminderRepo: reminderRepo,
		noteRepo:     noteRepo,
		intentParser: intentParser,
	}
}

func (s *assistantService) HandleTextMessage(input types.AssistantMessageInput) (*types.AssistantMessageResult, error) {
	text := strings.TrimSpace(input.Text)
	if text == "" {
		return s.saveUnknownIntent(input, "ยังไม่เห็นข้อความที่อยากให้เลขาช่วยเลยค่ะ ส่งมาได้เลย เดี๋ยวจัดการให้ค่ะ")
	}

	loc := loadLocation(input.Timezone)
	now := input.Now.In(loc)

	if target, period, ok := parseMoneyReportRequest(text, now, loc); ok {
		switch target {
		case moneyReportTargetIncome:
			return s.handleIncomeReport(input, now, loc, period)
		case moneyReportTargetCashflow:
			return s.handleCashflowReport(input, now, loc, period)
		default:
			return s.handleExpenseReport(input, now, loc, period)
		}
	}
	if isCalendarCancelCommand(text) {
		return s.handleCancelCalendarEvent(input, text, now, loc)
	}
	if isCalendarListCommand(text) {
		return s.handleCalendarList(input, now, loc)
	}
	if looksLikeCalendarIntentWithoutTime(text) {
		return s.handleCreateNote(input, text, now)
	}

	if result, ok := s.handleWithIntentParser(input, text, now, loc); ok {
		return result, nil
	}

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
	if amount, ok := extractAmount(text); ok && looksLikeIncome(text) {
		return s.handleCreateIncome(input, text, amount, now, loc)
	}
	if amount, ok := extractAmount(text); ok && looksLikeExpense(text) {
		return s.handleCreateExpense(input, text, amount, now, loc)
	}
	if looksLikeCalendarEvent(text) {
		return s.handleCreateCalendarEvent(input, text, now, loc)
	}

	return s.saveUnknownIntent(input, "ตอนนี้เลขาคนนี้ถนัด todo, รายรับ, รายจ่าย, นัดหมาย, เตือนความจำ, note และสรุปพรุ่งนี้ค่ะ ส่งงานแนวนี้มาได้เลยค่ะ")
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
		ReplyText: fmt.Sprintf("จัด todo \"%s\" เข้าลิสต์ให้เรียบร้อยค่ะ", todo.Title),
		Data:      todo,
	}, nil
}

func (s *assistantService) handleCreateExpense(input types.AssistantMessageInput, text string, amount float64, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	category := inferExpenseCategory(text)
	description := cleanupExpenseDescription(text)
	spentAt := parseMoneyEntryTime(text, now, loc)

	expense, err := s.expenseRepo.Create(&entities.Expense{
		UserID:          input.UserID,
		Amount:          amount,
		Currency:        "THB",
		Category:        category,
		Description:     description,
		SpentAt:         spentAt,
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
		ReplyText: formatMoneyCreateReply(expense.Description, expense.Amount, expense.SpentAt, loc, hasExplicitMoneyEntryTime(text)),
		Data:      expense,
	}, nil
}

func (s *assistantService) handleCreateIncome(input types.AssistantMessageInput, text string, amount float64, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	category := inferIncomeCategory(text)
	description := cleanupIncomeDescription(text)
	receivedAt := parseMoneyEntryTime(text, now, loc)

	income, err := s.incomeRepo.Create(&entities.Income{
		UserID:          input.UserID,
		Amount:          amount,
		Currency:        "THB",
		Category:        category,
		Description:     description,
		ReceivedAt:      receivedAt,
		SourceMessageID: input.MessageLogID,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	if err != nil {
		return nil, err
	}

	entitiesJSON := marshalEntities(map[string]any{
		"amount":      income.Amount,
		"currency":    income.Currency,
		"category":    income.Category,
		"description": income.Description,
		"received_at": income.ReceivedAt,
	})
	intent, err := s.intentRepo.Create(&entities.AssistantIntent{
		UserID:       input.UserID,
		MessageLogID: input.MessageLogID,
		Intent:       "create_income",
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
		ReplyText: formatMoneyCreateReply(income.Description, income.Amount, income.ReceivedAt, loc, hasExplicitMoneyEntryTime(text)),
		Data:      income,
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
		ReplyText: fmt.Sprintf("เลขาเช็กให้แล้ว เดือนนี้ใช้เงินไป %.2f บาทค่ะ", total),
		Data:      map[string]any{"total": total, "start": monthStart, "end": nextMonth},
	}, nil
}

type expenseReportPeriod struct {
	Label  string
	Start  time.Time
	End    time.Time
	Intent string
}

func (s *assistantService) handleExpenseReport(input types.AssistantMessageInput, now time.Time, loc *time.Location, period expenseReportPeriod) (*types.AssistantMessageResult, error) {
	expenses, err := s.expenseRepo.ListBySpentAtBetween(input.UserID, period.Start, period.End)
	if err != nil {
		return nil, err
	}

	total := 0.0
	for _, expense := range expenses {
		total += expense.Amount
	}

	entitiesJSON := marshalEntities(map[string]any{
		"start": period.Start,
		"end":   period.End,
		"total": total,
		"count": len(expenses),
	})
	intent, err := s.intentRepo.Create(&entities.AssistantIntent{
		UserID:       input.UserID,
		MessageLogID: input.MessageLogID,
		Intent:       period.Intent,
		Confidence:   0.92,
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
		ReplyText: formatExpenseReportReply(period.Label, expenses, total, loc),
		Data:      map[string]any{"expenses": expenses, "total": total, "start": period.Start, "end": period.End},
	}, nil
}

func (s *assistantService) handleIncomeReport(input types.AssistantMessageInput, now time.Time, loc *time.Location, period expenseReportPeriod) (*types.AssistantMessageResult, error) {
	incomes, err := s.incomeRepo.ListByReceivedAtBetween(input.UserID, period.Start, period.End)
	if err != nil {
		return nil, err
	}

	total := 0.0
	for _, income := range incomes {
		total += income.Amount
	}

	entitiesJSON := marshalEntities(map[string]any{
		"start": period.Start,
		"end":   period.End,
		"total": total,
		"count": len(incomes),
	})
	intent, err := s.intentRepo.Create(&entities.AssistantIntent{
		UserID:       input.UserID,
		MessageLogID: input.MessageLogID,
		Intent:       period.Intent,
		Confidence:   0.92,
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
		ReplyText: formatIncomeReportReply(period.Label, incomes, total, loc),
		Data:      map[string]any{"incomes": incomes, "total": total, "start": period.Start, "end": period.End},
	}, nil
}

func (s *assistantService) handleCashflowReport(input types.AssistantMessageInput, now time.Time, loc *time.Location, period expenseReportPeriod) (*types.AssistantMessageResult, error) {
	incomes, err := s.incomeRepo.ListByReceivedAtBetween(input.UserID, period.Start, period.End)
	if err != nil {
		return nil, err
	}
	expenses, err := s.expenseRepo.ListBySpentAtBetween(input.UserID, period.Start, period.End)
	if err != nil {
		return nil, err
	}

	incomeTotal := 0.0
	for _, income := range incomes {
		incomeTotal += income.Amount
	}
	expenseTotal := 0.0
	for _, expense := range expenses {
		expenseTotal += expense.Amount
	}

	entitiesJSON := marshalEntities(map[string]any{
		"start":         period.Start,
		"end":           period.End,
		"income_total":  incomeTotal,
		"expense_total": expenseTotal,
		"net":           incomeTotal - expenseTotal,
		"income_count":  len(incomes),
		"expense_count": len(expenses),
	})
	intent, err := s.intentRepo.Create(&entities.AssistantIntent{
		UserID:       input.UserID,
		MessageLogID: input.MessageLogID,
		Intent:       period.Intent,
		Confidence:   0.92,
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
		ReplyText: formatCashflowReportReply(period.Label, incomes, expenses, incomeTotal, expenseTotal, loc),
		Data: map[string]any{
			"incomes":       incomes,
			"expenses":      expenses,
			"income_total":  incomeTotal,
			"expense_total": expenseTotal,
			"net":           incomeTotal - expenseTotal,
			"start":         period.Start,
			"end":           period.End,
		},
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
		ReplyText: fmt.Sprintf("ลงนัด \"%s\" เวลา %s ในปฏิทินให้เรียบร้อยค่ะ", event.Title, event.StartAt.Format("02 Jan 2006 15:04")),
		Data:      event,
	}, nil
}

type calendarCancelRequest struct {
	Title   string
	Start   time.Time
	End     time.Time
	HasDate bool
	HasTime bool
	Hour    int
	Minute  int
}

func (s *assistantService) handleCancelCalendarEvent(input types.AssistantMessageInput, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	request := parseCalendarCancelRequest(text, now, loc)
	if request.Title == "" && !request.HasDate && !request.HasTime {
		return s.saveCalendarCancelIntent(input, now, request, "needs_clarification", "อยากให้เลขาถอนนัดไหนออกจากสมุดค่ะ ระบุชื่อ วัน หรือเวลาเพิ่มอีกนิด เช่น \"ยกเลิกนัดพรุ่งนี้ 10 โมง ประชุมกับทีม\" ค่ะ", nil)
	}

	candidates, err := s.calendarRepo.ListByStartBetween(input.UserID, request.Start, request.End)
	if err != nil {
		return nil, err
	}

	matches := filterCalendarCancelCandidates(candidates, request, loc)
	if len(matches) == 0 {
		return s.saveCalendarCancelIntent(input, now, request, "not_found", "เลขาค้นสมุดนัดแล้ว ยังไม่เจอนัดที่ตรงกับข้อความนี้ค่ะ ลองระบุวัน เวลา หรือชื่อนัดให้ชัดขึ้นอีกนิดค่ะ", nil)
	}
	if len(matches) > 1 {
		return s.saveCalendarCancelIntent(input, now, request, "needs_clarification", formatCalendarCancelClarification(matches, loc), matches)
	}

	event := matches[0]
	if err := s.calendarRepo.Delete(event.ID); err != nil {
		return nil, err
	}

	reply := fmt.Sprintf("ยกเลิกนัด \"%s\" เวลา %s ออกจากปฏิทินให้แล้วค่ะ", event.Title, event.StartAt.In(loc).Format("02 Jan 2006 15:04"))
	return s.saveCalendarCancelIntent(input, now, request, "completed", reply, matches)
}

func (s *assistantService) saveCalendarCancelIntent(input types.AssistantMessageInput, now time.Time, request calendarCancelRequest, status, reply string, matches []*entities.CalendarEvent) (*types.AssistantMessageResult, error) {
	matchIDs := make([]uint, 0, len(matches))
	for _, match := range matches {
		matchIDs = append(matchIDs, match.ID)
	}

	entitiesJSON := marshalEntities(map[string]any{
		"title":     request.Title,
		"start":     request.Start,
		"end":       request.End,
		"has_date":  request.HasDate,
		"has_time":  request.HasTime,
		"match_ids": matchIDs,
	})
	intent, err := s.intentRepo.Create(&entities.AssistantIntent{
		UserID:       input.UserID,
		MessageLogID: input.MessageLogID,
		Intent:       "cancel_calendar_event",
		Confidence:   0.88,
		Entities:     entitiesJSON,
		Status:       status,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return nil, err
	}

	return &types.AssistantMessageResult{
		Intent:    intent.Intent,
		ReplyText: reply,
		Data:      map[string]any{"matches": matches, "status": status},
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
		ReplyText: fmt.Sprintf("ตั้งเตือน \"%s\" เวลา %s ให้แล้วค่ะ เดี๋ยวถึงเวลาจะสะกิดให้ค่ะ", reminder.Title, reminder.RemindAt.Format("02 Jan 2006 15:04")),
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
		ReplyText: "จดโน้ตใส่สมุดให้เรียบร้อยค่ะ",
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

	reply := "พรุ่งนี้สมุดยังโล่งค่ะ ยังไม่มีนัดหรือ todo ค่ะ"
	if len(parts) > 0 {
		reply = "พรุ่งนี้เลขาเช็กให้แล้ว มีรายการนี้ค่ะ:\n- " + strings.Join(parts, "\n- ") + "\nเดี๋ยวช่วยจำให้อีกแรงค่ะ"
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

func (s *assistantService) handleCalendarList(input types.AssistantMessageInput, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	events, err := s.calendarRepo.ListByUserID(input.UserID, 50, 0)
	if err != nil {
		return nil, err
	}

	entitiesJSON := marshalEntities(map[string]any{"event_count": len(events)})
	intent, err := s.intentRepo.Create(&entities.AssistantIntent{
		UserID:       input.UserID,
		MessageLogID: input.MessageLogID,
		Intent:       "calendar_list",
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
		ReplyText: formatCalendarListReply(events, loc),
		Data:      map[string]any{"events": events},
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

type moneyReportTarget string

const (
	moneyReportTargetExpense  moneyReportTarget = "expense"
	moneyReportTargetIncome   moneyReportTarget = "income"
	moneyReportTargetCashflow moneyReportTarget = "cashflow"
)

type reportPeriodKind string

const (
	reportPeriodDaily   reportPeriodKind = "daily"
	reportPeriodWeekly  reportPeriodKind = "weekly"
	reportPeriodMonthly reportPeriodKind = "monthly"
)

func parseExpenseReportPeriod(text string, now time.Time, loc *time.Location) (expenseReportPeriod, bool) {
	target, period, ok := parseMoneyReportRequest(text, now, loc)
	if !ok || target != moneyReportTargetExpense {
		return expenseReportPeriod{}, false
	}
	return period, true
}

func parseMoneyReportRequest(text string, now time.Time, loc *time.Location) (moneyReportTarget, expenseReportPeriod, bool) {
	target, ok := detectMoneyReportTarget(text)
	if !ok || !containsAny(text, moneyReportQueryWords()) {
		return "", expenseReportPeriod{}, false
	}
	if _, hasAmount := extractAmount(text); hasAmount && !containsAny(text, strongMoneyReportQueryWords()) {
		return "", expenseReportPeriod{}, false
	}

	period, kind := parseReportPeriod(text, now, loc)
	period.Intent = fmt.Sprintf("%s_report_%s", target, kind)
	return target, period, true
}

func detectMoneyReportTarget(text string) (moneyReportTarget, bool) {
	hasIncome := containsAny(text, incomeWords())
	hasExpense := containsAny(text, expenseWords())
	hasCashflow := containsAny(text, cashflowWords()) || (hasIncome && hasExpense)

	switch {
	case hasCashflow:
		return moneyReportTargetCashflow, true
	case hasIncome:
		return moneyReportTargetIncome, true
	case hasExpense:
		return moneyReportTargetExpense, true
	default:
		return "", false
	}
}

func parseReportPeriod(text string, now time.Time, loc *time.Location) (expenseReportPeriod, reportPeriodKind) {
	current := now.In(loc)
	referenceDate, hasReferenceDate := parseReportReferenceDate(text, current, loc)
	if strings.Contains(strings.ToLower(text), "week") || strings.Contains(text, "สัปดาห์") || strings.Contains(text, "อาทิตย์") {
		reference := current
		if hasReferenceDate {
			reference = referenceDate
		}
		start := startOfReportWeek(reference, loc)
		return expenseReportPeriod{
			Label: formatReportPeriodLabel("สัปดาห์", start, start.AddDate(0, 0, 7), loc),
			Start: start,
			End:   start.AddDate(0, 0, 7),
		}, reportPeriodWeekly
	}

	if strings.Contains(strings.ToLower(text), "month") || strings.Contains(text, "เดือน") || (containsThaiMonthName(text) && !hasReferenceDate) {
		start := time.Date(current.Year(), current.Month(), 1, 0, 0, 0, 0, loc)
		if parsed, ok := parseReportMonth(text, current, loc); ok {
			start = parsed
		} else if hasReferenceDate {
			start = time.Date(referenceDate.Year(), referenceDate.Month(), 1, 0, 0, 0, 0, loc)
		}
		return expenseReportPeriod{
			Label: formatReportPeriodLabel("เดือน", start, start.AddDate(0, 1, 0), loc),
			Start: start,
			End:   start.AddDate(0, 1, 0),
		}, reportPeriodMonthly
	}

	reference := current
	if hasReferenceDate {
		reference = referenceDate
	}
	start := time.Date(reference.Year(), reference.Month(), reference.Day(), 0, 0, 0, 0, loc)
	return expenseReportPeriod{
		Label: formatReportPeriodLabel("วันที่", start, start.AddDate(0, 0, 1), loc),
		Start: start,
		End:   start.AddDate(0, 0, 1),
	}, reportPeriodDaily
}

func isExpenseReportQuery(text string) bool {
	_, _, ok := parseMoneyReportRequest(text, time.Now(), time.Local)
	return ok && containsAny(text, expenseWords())
}

func isIncomeReportQuery(text string) bool {
	_, _, ok := parseMoneyReportRequest(text, time.Now(), time.Local)
	return ok && containsAny(text, incomeWords())
}

func isTomorrowSummary(text string) bool {
	return strings.Contains(text, "พรุ่งนี้") && (strings.Contains(text, "มีอะไร") || strings.Contains(text, "ทำอะไร"))
}

func isCalendarListCommand(text string) bool {
	phrases := []string{
		"ดูนัดทั้งหมด",
		"ดู calendar ทั้งหมด",
		"ดูปฏิทินทั้งหมด",
		"มีนัดอะไรบ้าง",
		"เช็คนัดทั้งหมด",
		"เช็กนัดทั้งหมด",
		"รายการนัดหมาย",
		"รายการนัด",
		"calendar ทั้งหมด",
		"นัดทั้งหมด",
	}
	return containsAny(text, phrases)
}

func isReminderCommand(text string) bool {
	return strings.Contains(text, "เตือน")
}

func isNoteCommand(text string) bool {
	lower := strings.ToLower(text)
	return strings.Contains(text, "จด") || strings.Contains(lower, "note")
}

func looksLikeExpense(text string) bool {
	if looksLikeIncome(text) {
		return false
	}
	keywords := []string{"บาท", "กิน", "ซื้อ", "จ่าย", "ค่า", "โอน", "กาแฟ", "ข้าว", "อาหาร"}
	for _, keyword := range keywords {
		if strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

func looksLikeIncome(text string) bool {
	return containsAny(text, incomeCreateWords())
}

func expenseWords() []string {
	return []string{"ใช้เงิน", "รายจ่าย", "ค่าใช้จ่าย", "จ่าย", "ซื้อ", "กิน", "อาหาร", "expense", "spent", "spend"}
}

func incomeWords() []string {
	return []string{"รายรับ", "รับเงิน", "ได้เงิน", "รายได้", "เงินเข้า", "income", "salary"}
}

func incomeCreateWords() []string {
	return []string{"รายรับ", "รับเงิน", "ได้เงิน", "ได้รับเงิน", "รายได้", "เงินเข้า", "เงินเดือน", "ขายได้", "income", "salary"}
}

func cashflowWords() []string {
	return []string{"รายรับรายจ่าย", "รายรับและรายจ่าย", "รายรับ รายจ่าย", "เงินเข้าเงินออก", "สรุปการเงิน", "การเงิน", "สุทธิ", "คงเหลือ", "cashflow", "cash flow", "balance", "net"}
}

func moneyReportQueryWords() []string {
	return []string{"อะไร", "บ้าง", "เท่าไหร่", "เท่าไร", "รวม", "สรุป", "รายการ", "ดู", "เช็ค", "เช็ก", "ย้อนหลัง", "วันนี้", "เมื่อวาน", "สัปดาห์", "อาทิตย์", "เดือน", "summary", "report", "list", "total", "today", "yesterday", "week", "month"}
}

func strongMoneyReportQueryWords() []string {
	return []string{"อะไร", "บ้าง", "เท่าไหร่", "เท่าไร", "รวม", "สรุป", "รายการ", "ดู", "เช็ค", "เช็ก", "ย้อนหลัง", "summary", "report", "list", "total"}
}

func looksLikeCalendarEvent(text string) bool {
	if isQuestion(text) || isCalendarCancelCommand(text) {
		return false
	}
	if !hasExplicitTime(text) {
		return false
	}

	calendarPrefixes := []string{
		"นัด", "เพิ่มนัด", "ลงนัด", "บันทึกนัด", "สร้างนัด",
		"เพิ่มตาราง", "ลงตาราง", "บันทึกตาราง", "เพิ่ม calendar", "calendar",
	}
	if hasPrefixAny(strings.TrimSpace(text), calendarPrefixes) {
		return true
	}

	eventWords := []string{"ประชุม", "นัด", "เจอ", "คุย", "โทรหา", "call", "meeting"}
	return containsAny(text, eventWords)
}

func looksLikeCalendarIntentWithoutTime(text string) bool {
	if isQuestion(text) || isCalendarCancelCommand(text) || hasExplicitTime(text) {
		return false
	}

	calendarPrefixes := []string{
		"นัด", "เพิ่มนัด", "ลงนัด", "บันทึกนัด", "สร้างนัด",
		"เพิ่มตาราง", "ลงตาราง", "บันทึกตาราง", "เพิ่ม calendar", "calendar",
	}
	eventWords := []string{"ประชุม", "นัด", "เจอ", "คุย", "โทรหา", "meeting"}
	return hasPrefixAny(strings.TrimSpace(text), calendarPrefixes) || containsAny(text, eventWords)
}

func hasExplicitTime(text string) bool {
	_, _, ok := extractHourMinute(text)
	return ok
}

func extractAmount(text string) (float64, bool) {
	withCurrency := regexp.MustCompile(`([0-9]+(?:\.[0-9]+)?)\s*(บาท|บ|thb|THB)`)
	if match := withCurrency.FindStringSubmatch(text); len(match) >= 2 {
		amount, err := strconv.ParseFloat(match[1], 64)
		if err == nil {
			return amount, true
		}
	}

	re := regexp.MustCompile(`[0-9]+(?:\.[0-9]+)?`)
	matches := re.FindAllStringIndex(text, -1)
	for _, match := range matches {
		raw := text[match[0]:match[1]]
		amount, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			continue
		}
		if amount > 31 && !looksLikeDateNumber(text, match[0], match[1]) {
			return amount, true
		}
	}
	if len(matches) == 0 {
		return 0, false
	}
	raw := text[matches[0][0]:matches[0][1]]
	amount, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, false
	}
	return amount, true
}

func looksLikeDateNumber(text string, start, end int) bool {
	before := text[:start]
	after := text[end:]
	if strings.HasSuffix(before, "วันที่ ") || strings.HasSuffix(before, "วันที่") {
		return true
	}
	if strings.HasPrefix(after, "/") || strings.HasPrefix(after, "-") {
		return true
	}
	if strings.HasSuffix(before, "/") || strings.HasSuffix(before, "-") {
		return true
	}
	return false
}

func parseMoneyEntryTime(text string, now time.Time, loc *time.Location) time.Time {
	current := now.In(loc)
	hour, minute, hasTime := extractHourMinute(text)
	if date, ok := parseReportReferenceDate(text, current, loc); ok {
		if hasTime {
			return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, loc)
		}
		return time.Date(date.Year(), date.Month(), date.Day(), current.Hour(), current.Minute(), current.Second(), current.Nanosecond(), loc)
	}
	if hasTime {
		return time.Date(current.Year(), current.Month(), current.Day(), hour, minute, 0, 0, loc)
	}
	return current
}

func mergeParsedMoneyTime(parsedAt time.Time, text string, loc *time.Location) time.Time {
	parsedAt = parsedAt.In(loc)
	hour, minute, hasTime := extractHourMinute(text)
	if hasTime && parsedAt.Hour() == 0 && parsedAt.Minute() == 0 && parsedAt.Second() == 0 && parsedAt.Nanosecond() == 0 {
		return time.Date(parsedAt.Year(), parsedAt.Month(), parsedAt.Day(), hour, minute, 0, 0, loc)
	}
	return parsedAt
}

func hasExplicitMoneyEntryTime(text string) bool {
	_, _, ok := extractHourMinute(text)
	return ok
}

func formatMoneyCreateReply(description string, amount float64, occurredAt time.Time, loc *time.Location, includeTime bool) string {
	description = normalizeDescriptionSpaces(description)
	if description == "" {
		description = "รายการ"
	}
	if includeTime {
		return fmt.Sprintf("ลงบัญชี %s %.2f บาท ตอน %s น. ให้เรียบร้อยค่ะ", description, amount, occurredAt.In(loc).Format("15:04"))
	}
	return fmt.Sprintf("ลงบัญชี %s %.2f บาทให้เรียบร้อยค่ะ", description, amount)
}

func parseReportReferenceDate(text string, current time.Time, loc *time.Location) (time.Time, bool) {
	lower := strings.ToLower(text)
	if strings.Contains(text, "วันนี้") || strings.Contains(lower, "today") {
		return time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, loc), true
	}
	if strings.Contains(text, "เมื่อวาน") || strings.Contains(lower, "yesterday") {
		date := current.AddDate(0, 0, -1)
		return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc), true
	}
	if strings.Contains(text, "พรุ่งนี้") || strings.Contains(lower, "tomorrow") {
		date := current.AddDate(0, 0, 1)
		return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc), true
	}

	if date, ok := parseISODateInText(text, loc); ok {
		return date, true
	}
	if date, ok := parseSlashDateInText(text, current, loc); ok {
		return date, true
	}
	if date, ok := parseThaiMonthDateInText(text, current, loc); ok {
		return date, true
	}
	if date, ok := parseDayOnlyInText(text, current, loc); ok {
		return date, true
	}
	return time.Time{}, false
}

func parseISODateInText(text string, loc *time.Location) (time.Time, bool) {
	re := regexp.MustCompile(`\b(\d{4})-(\d{1,2})-(\d{1,2})\b`)
	match := re.FindStringSubmatch(text)
	if len(match) != 4 {
		return time.Time{}, false
	}
	year, _ := strconv.Atoi(match[1])
	month, _ := strconv.Atoi(match[2])
	day, _ := strconv.Atoi(match[3])
	return buildDate(year, time.Month(month), day, loc)
}

func parseSlashDateInText(text string, current time.Time, loc *time.Location) (time.Time, bool) {
	re := regexp.MustCompile(`\b(\d{1,2})[/-](\d{1,2})(?:[/-](\d{2,4}))?\b`)
	match := re.FindStringSubmatch(text)
	if len(match) < 3 {
		return time.Time{}, false
	}
	day, _ := strconv.Atoi(match[1])
	month, _ := strconv.Atoi(match[2])
	year := current.Year()
	if len(match) >= 4 && match[3] != "" {
		year, _ = strconv.Atoi(match[3])
		year = normalizeYear(year)
	}
	return buildDate(year, time.Month(month), day, loc)
}

func parseThaiMonthDateInText(text string, current time.Time, loc *time.Location) (time.Time, bool) {
	monthPattern := thaiMonthPattern()
	re := regexp.MustCompile(`(?i)(\d{1,2})\s*(` + monthPattern + `)(?:\s*(\d{2,4}))?`)
	match := re.FindStringSubmatch(text)
	if len(match) < 3 {
		return time.Time{}, false
	}
	day, _ := strconv.Atoi(match[1])
	month, ok := thaiMonthNumber(match[2])
	if !ok {
		return time.Time{}, false
	}
	year := current.Year()
	if len(match) >= 4 && match[3] != "" {
		year, _ = strconv.Atoi(match[3])
		year = normalizeYear(year)
	}
	return buildDate(year, month, day, loc)
}

func parseDayOnlyInText(text string, current time.Time, loc *time.Location) (time.Time, bool) {
	re := regexp.MustCompile(`วันที่\s*(\d{1,2})`)
	match := re.FindStringSubmatch(text)
	if len(match) != 2 {
		return time.Time{}, false
	}
	day, _ := strconv.Atoi(match[1])
	return buildDate(current.Year(), current.Month(), day, loc)
}

func parseReportMonth(text string, current time.Time, loc *time.Location) (time.Time, bool) {
	if start, ok := parseISOMonthInText(text, loc); ok {
		return start, true
	}

	monthPattern := thaiMonthPattern()
	re := regexp.MustCompile(`(?i)(` + monthPattern + `)(?:\s*(\d{2,4}))?`)
	match := re.FindStringSubmatch(text)
	if len(match) >= 2 {
		month, ok := thaiMonthNumber(match[1])
		if ok {
			year := current.Year()
			if len(match) >= 3 && match[2] != "" {
				year, _ = strconv.Atoi(match[2])
				year = normalizeYear(year)
			}
			return time.Date(year, month, 1, 0, 0, 0, 0, loc), true
		}
	}

	numberMonth := regexp.MustCompile(`เดือน\s*(\d{1,2})(?:[/-](\d{2,4}))?`)
	match = numberMonth.FindStringSubmatch(text)
	if len(match) >= 2 {
		monthNumber, _ := strconv.Atoi(match[1])
		year := current.Year()
		if len(match) >= 3 && match[2] != "" {
			year, _ = strconv.Atoi(match[2])
			year = normalizeYear(year)
		}
		if monthNumber >= 1 && monthNumber <= 12 {
			return time.Date(year, time.Month(monthNumber), 1, 0, 0, 0, 0, loc), true
		}
	}

	return time.Time{}, false
}

func parseISOMonthInText(text string, loc *time.Location) (time.Time, bool) {
	re := regexp.MustCompile(`\b(\d{4})-(\d{1,2})\b`)
	match := re.FindStringSubmatch(text)
	if len(match) != 3 {
		return time.Time{}, false
	}
	year, _ := strconv.Atoi(match[1])
	month, _ := strconv.Atoi(match[2])
	if month < 1 || month > 12 {
		return time.Time{}, false
	}
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc), true
}

func buildDate(year int, month time.Month, day int, loc *time.Location) (time.Time, bool) {
	year = normalizeYear(year)
	if month < time.January || month > time.December || day < 1 || day > 31 {
		return time.Time{}, false
	}
	date := time.Date(year, month, day, 0, 0, 0, 0, loc)
	if date.Year() != year || date.Month() != month || date.Day() != day {
		return time.Time{}, false
	}
	return date, true
}

func normalizeYear(year int) int {
	if year < 100 {
		return 2000 + year
	}
	if year > 2400 {
		return year - 543
	}
	return year
}

func containsThaiMonthName(text string) bool {
	for name := range thaiMonthMap() {
		if strings.Contains(strings.ToLower(text), strings.ToLower(name)) {
			return true
		}
	}
	return false
}

func thaiMonthPattern() string {
	return strings.Join([]string{
		"มกราคม", "ม\\.ค\\.?", "jan", "january",
		"กุมภาพันธ์", "ก\\.พ\\.?", "feb", "february",
		"มีนาคม", "มี\\.ค\\.?", "mar", "march",
		"เมษายน", "เม\\.ย\\.?", "apr", "april",
		"พฤษภาคม", "พ\\.ค\\.?", "may",
		"มิถุนายน", "มิ\\.ย\\.?", "jun", "june",
		"กรกฎาคม", "ก\\.ค\\.?", "jul", "july",
		"สิงหาคม", "ส\\.ค\\.?", "aug", "august",
		"กันยายน", "ก\\.ย\\.?", "sep", "september",
		"ตุลาคม", "ต\\.ค\\.?", "oct", "october",
		"พฤศจิกายน", "พ\\.ย\\.?", "nov", "november",
		"ธันวาคม", "ธ\\.ค\\.?", "dec", "december",
	}, "|")
}

func thaiMonthNumber(value string) (time.Month, bool) {
	normalized := strings.ToLower(strings.TrimSpace(strings.TrimSuffix(value, ".")))
	normalized = strings.ReplaceAll(normalized, ".", "")
	month, ok := thaiMonthMap()[normalized]
	return month, ok
}

func thaiMonthMap() map[string]time.Month {
	return map[string]time.Month{
		"มกราคม": time.January, "มค": time.January, "jan": time.January, "january": time.January,
		"กุมภาพันธ์": time.February, "กพ": time.February, "feb": time.February, "february": time.February,
		"มีนาคม": time.March, "มีค": time.March, "mar": time.March, "march": time.March,
		"เมษายน": time.April, "เมย": time.April, "apr": time.April, "april": time.April,
		"พฤษภาคม": time.May, "พค": time.May, "may": time.May,
		"มิถุนายน": time.June, "มิย": time.June, "jun": time.June, "june": time.June,
		"กรกฎาคม": time.July, "กค": time.July, "jul": time.July, "july": time.July,
		"สิงหาคม": time.August, "สค": time.August, "aug": time.August, "august": time.August,
		"กันยายน": time.September, "กย": time.September, "sep": time.September, "september": time.September,
		"ตุลาคม": time.October, "ตค": time.October, "oct": time.October, "october": time.October,
		"พฤศจิกายน": time.November, "พย": time.November, "nov": time.November, "november": time.November,
		"ธันวาคม": time.December, "ธค": time.December, "dec": time.December, "december": time.December,
	}
}

func startOfReportWeek(date time.Time, loc *time.Location) time.Time {
	current := date.In(loc)
	daysSinceMonday := int(current.Weekday()) - int(time.Monday)
	if daysSinceMonday < 0 {
		daysSinceMonday += 7
	}
	startDate := current.AddDate(0, 0, -daysSinceMonday)
	return time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, loc)
}

func formatReportPeriodLabel(prefix string, start, end time.Time, loc *time.Location) string {
	switch prefix {
	case "วันที่":
		return "วันที่ " + start.In(loc).Format("02 Jan 2006")
	case "สัปดาห์":
		return fmt.Sprintf("สัปดาห์ %s-%s", start.In(loc).Format("02 Jan"), end.Add(-time.Nanosecond).In(loc).Format("02 Jan 2006"))
	case "เดือน":
		return "เดือน " + start.In(loc).Format("Jan 2006")
	default:
		return prefix
	}
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

func inferIncomeCategory(text string) string {
	lower := strings.ToLower(text)
	if strings.Contains(text, "เงินเดือน") || strings.Contains(lower, "salary") {
		return "salary"
	}
	if strings.Contains(text, "ขาย") {
		return "sales"
	}
	if strings.Contains(text, "ดอกเบี้ย") || strings.Contains(text, "ปันผล") || strings.Contains(lower, "dividend") {
		return "investment"
	}
	return "uncategorized"
}

func cleanupExpenseDescription(text string) string {
	return cleanupMoneyDescription(text, "รายจ่าย", nil)
}

func cleanupIncomeDescription(text string) string {
	return cleanupMoneyDescription(text, "รายรับ", []string{"บันทึก", "เพิ่ม", "รายรับ", "รับเงิน", "ได้เงิน", "ได้รับเงิน", "เงินเข้า", "income", "Income"})
}

func cleanupMoneyDescription(text, fallback string, removers []string) string {
	cleaned := strings.TrimSpace(text)
	cleaned = removeMoneyTimePhrases(cleaned)
	cleaned = removeMoneyDatePhrases(cleaned)

	amount := regexp.MustCompile(`[0-9]+(?:\.[0-9]+)?\s*(บาท|บ|thb|THB)?`)
	cleaned = amount.ReplaceAllString(cleaned, "")
	cleaned = cleanupByRemoving(cleaned, removers)
	cleaned = cleanupByRemoving(cleaned, []string{"วันนี้", "พรุ่งนี้", "เมื่อวาน", "ตอน", "เวลา"})
	cleaned = normalizeDescriptionSpaces(cleaned)

	if cleaned == "" {
		return fallback
	}
	return cleaned
}

func removeMoneyTimePhrases(text string) string {
	patterns := []string{
		`(?:ตอน|เวลา)?\s*\d{1,2}\s*โมง(?:\s*(?:เช้า|เย็น|ค่ำ))?(?:\s*(?:ครึ่ง|\d{1,2}\s*นาที))?(?:\s*(?:เช้า|เย็น|ค่ำ))?`,
		`(?:ตอน|เวลา)?\s*\d{1,2}\s*ทุ่ม(?:\s*(?:ครึ่ง|\d{1,2}\s*นาที))?`,
		`(?:ตอน|เวลา)?\s*บ่าย\s*\d{1,2}(?:\s*โมง)?(?:\s*(?:ครึ่ง|\d{1,2}\s*นาที))?`,
		`(?:ตอน|เวลา)?\s*ตี\s*\d{1,2}(?:\s*(?:ครึ่ง|\d{1,2}\s*นาที))?`,
		`(?:ตอน|เวลา)?\s*\d{1,2}[:.]\d{2}\s*(?:น\.?|นาฬิกา)?`,
		`(?:ตอน|เวลา)?\s*\d{1,2}\s*(?:น\.|นาฬิกา)(?:\s*(?:ครึ่ง|\d{1,2}\s*นาที))?`,
	}
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		text = re.ReplaceAllString(text, "")
	}
	return text
}

func removeMoneyDatePhrases(text string) string {
	patterns := []string{
		`\b\d{4}-\d{1,2}-\d{1,2}\b`,
		`\b\d{1,2}[/-]\d{1,2}(?:[/-]\d{2,4})?\b`,
		`วันที่\s*\d{1,2}\s*(?:` + thaiMonthPattern() + `)?(?:\s*\d{2,4})?`,
	}
	for _, pattern := range patterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		text = re.ReplaceAllString(text, "")
	}
	return text
}

func normalizeDescriptionSpaces(text string) string {
	text = strings.TrimSpace(text)
	spaces := regexp.MustCompile(`\s+`)
	text = spaces.ReplaceAllString(text, " ")
	punctuationSpaces := regexp.MustCompile(`\s+([,.;:!?])`)
	return punctuationSpaces.ReplaceAllString(text, "$1")
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
	cleaned := stripCalendarTimeClause(text)
	cleaned = removeCalendarTimePhrases(cleaned)
	cleaned = cleanupByRemoving(cleaned, []string{
		"เพิ่มนัด", "ลงนัด", "บันทึกนัด", "สร้างนัด", "นัด",
		"เพิ่มตาราง", "ลงตาราง", "บันทึกตาราง", "calendar", "Calendar",
	})
	cleaned = cleanupByRemoving(cleaned, []string{"วันนี้", "พรุ่งนี้", "เมื่อวาน"})
	return normalizeDescriptionSpaces(cleaned)
}

func stripCalendarTimeClause(text string) string {
	cutWords := []string{" ตอน ", " เวลา ", "ตอน", "เวลา"}
	cutAt := -1
	for _, word := range cutWords {
		if idx := strings.Index(text, word); idx >= 0 && (cutAt == -1 || idx < cutAt) {
			cutAt = idx
		}
	}
	if cutAt >= 0 {
		return strings.TrimSpace(text[:cutAt])
	}
	return text
}

func removeCalendarTimePhrases(text string) string {
	patterns := []string{
		`เที่ยง(?:\s*(?:ครึ่ง|\d{1,2}\s*นาที))?`,
		`\d{1,2}\s*โมง(?:\s*(?:เช้า|เย็น|ค่ำ))?(?:\s*(?:ครึ่ง|\d{1,2}\s*นาที))?(?:\s*(?:เช้า|เย็น|ค่ำ))?`,
		`\d{1,2}\s*ทุ่ม(?:\s*(?:ครึ่ง|\d{1,2}\s*นาที))?`,
		`บ่าย\s*\d{1,2}(?:\s*โมง)?(?:\s*(?:ครึ่ง|\d{1,2}\s*นาที))?`,
		`ตี\s*\d{1,2}(?:\s*(?:ครึ่ง|\d{1,2}\s*นาที))?`,
		`\d{1,2}[:.]\d{2}\s*(?:น\.?|นาฬิกา)?`,
		`\d{1,2}\s*(?:น\.|นาฬิกา)(?:\s*(?:ครึ่ง|\d{1,2}\s*นาที))?`,
	}
	for _, pattern := range patterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		text = re.ReplaceAllString(text, "")
	}
	return text
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

func parseCalendarCancelRequest(text string, now time.Time, loc *time.Location) calendarCancelRequest {
	current := now.In(loc)
	title := cleanupCalendarCancelTitle(text)
	hour, minute, hasTime := extractHourMinute(text)
	hasDate := strings.Contains(text, "วันนี้") || strings.Contains(text, "พรุ่งนี้")

	if strings.Contains(text, "พรุ่งนี้") {
		date := current.AddDate(0, 0, 1)
		start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
		return calendarCancelRequest{
			Title:   title,
			Start:   start,
			End:     start.AddDate(0, 0, 1),
			HasDate: true,
			HasTime: hasTime,
			Hour:    hour,
			Minute:  minute,
		}
	}

	if strings.Contains(text, "วันนี้") {
		start := time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, loc)
		return calendarCancelRequest{
			Title:   title,
			Start:   start,
			End:     start.AddDate(0, 0, 1),
			HasDate: true,
			HasTime: hasTime,
			Hour:    hour,
			Minute:  minute,
		}
	}

	start := time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, loc)
	return calendarCancelRequest{
		Title:   title,
		Start:   start,
		End:     start.AddDate(0, 0, 30),
		HasDate: hasDate,
		HasTime: hasTime,
		Hour:    hour,
		Minute:  minute,
	}
}

func extractHourMinute(text string) (int, int, bool) {
	if hour, minute, ok := extractThaiNoonTime(text); ok {
		return hour, minute, true
	}
	if hour, minute, ok := extractThaiTuumTime(text); ok {
		return hour, minute, true
	}
	if hour, minute, ok := extractThaiAfternoonTime(text); ok {
		return hour, minute, true
	}
	if hour, minute, ok := extractThaiClockTime(text); ok {
		return hour, minute, true
	}
	if hour, minute, ok := extractSeparatedHourMinute(text); ok {
		return hour, minute, true
	}
	if hour, minute, ok := extractThaiShortClockTime(text); ok {
		return hour, minute, true
	}
	if hour, minute, ok := extractThaiDawnTime(text); ok {
		return hour, minute, true
	}

	return 0, 0, false
}

func extractThaiNoonTime(text string) (int, int, bool) {
	re := regexp.MustCompile(`เที่ยง(?:\s*(ครึ่ง|\d{1,2}\s*นาที))?`)
	if match := re.FindStringSubmatch(text); len(match) == 2 {
		minute := parseThaiMinuteSuffix(match[1])
		return validHourMinute(12, minute)
	}
	return 0, 0, false
}

func extractSeparatedHourMinute(text string) (int, int, bool) {
	re := regexp.MustCompile(`(\d{1,2})[:.](\d{2})\s*(?:น\.?|นาฬิกา)?`)
	if match := re.FindStringSubmatch(text); len(match) == 3 {
		hour, _ := strconv.Atoi(match[1])
		minute, _ := strconv.Atoi(match[2])
		return validHourMinute(normalizeThaiHour(text, hour), minute)
	}
	return 0, 0, false
}

func extractThaiClockTime(text string) (int, int, bool) {
	re := regexp.MustCompile(`(\d{1,2})\s*โมง(?:\s*(เช้า|เย็น|ค่ำ))?(?:\s*(ครึ่ง|\d{1,2}\s*นาที))?(?:\s*(เช้า|เย็น|ค่ำ))?`)
	if match := re.FindStringSubmatch(text); len(match) == 5 {
		hour, _ := strconv.Atoi(match[1])
		minute := parseThaiMinuteSuffix(firstNonEmpty(match[3]))
		context := strings.TrimSpace(match[0] + " " + match[2] + " " + match[4])
		return validHourMinute(normalizeThaiHour(context, hour), minute)
	}
	return 0, 0, false
}

func extractThaiShortClockTime(text string) (int, int, bool) {
	re := regexp.MustCompile(`(\d{1,2})\s*(?:น\.|นาฬิกา)(?:\s*(ครึ่ง|\d{1,2}\s*นาที))?`)
	if match := re.FindStringSubmatch(text); len(match) == 3 {
		hour, _ := strconv.Atoi(match[1])
		minute := parseThaiMinuteSuffix(match[2])
		return validHourMinute(normalizeThaiHour(match[0], hour), minute)
	}
	return 0, 0, false
}

func extractThaiAfternoonTime(text string) (int, int, bool) {
	re := regexp.MustCompile(`บ่าย\s*(\d{1,2})(?:\s*โมง)?(?:\s*(ครึ่ง|\d{1,2}\s*นาที))?`)
	if match := re.FindStringSubmatch(text); len(match) == 3 {
		hour, _ := strconv.Atoi(match[1])
		minute := parseThaiMinuteSuffix(match[2])
		return validHourMinute(normalizeThaiHour(match[0], hour), minute)
	}
	return 0, 0, false
}

func extractThaiTuumTime(text string) (int, int, bool) {
	re := regexp.MustCompile(`(\d{1,2})\s*ทุ่ม(?:\s*(ครึ่ง|\d{1,2}\s*นาที))?`)
	if match := re.FindStringSubmatch(text); len(match) == 3 {
		tuum, _ := strconv.Atoi(match[1])
		minute := parseThaiMinuteSuffix(match[2])
		hour := 18 + tuum
		if hour == 24 {
			hour = 0
		}
		return validHourMinute(hour, minute)
	}
	return 0, 0, false
}

func extractThaiDawnTime(text string) (int, int, bool) {
	re := regexp.MustCompile(`ตี\s*(\d{1,2})(?:\s*(ครึ่ง|\d{1,2}\s*นาที))?`)
	if match := re.FindStringSubmatch(text); len(match) == 3 {
		hour, _ := strconv.Atoi(match[1])
		minute := parseThaiMinuteSuffix(match[2])
		return validHourMinute(hour, minute)
	}
	return 0, 0, false
}

func parseThaiMinuteSuffix(value string) int {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	if strings.Contains(value, "ครึ่ง") {
		return 30
	}
	re := regexp.MustCompile(`(\d{1,2})`)
	if match := re.FindStringSubmatch(value); len(match) == 2 {
		minute, _ := strconv.Atoi(match[1])
		return minute
	}
	return 0
}

func normalizeThaiHour(text string, hour int) int {
	if strings.Contains(text, "เที่ยงคืน") {
		return 0
	}
	if strings.Contains(text, "เที่ยง") && hour == 12 {
		return 12
	}
	if strings.Contains(text, "บ่าย") && hour < 12 {
		return hour + 12
	}
	if (strings.Contains(text, "เย็น") || strings.Contains(text, "ค่ำ")) && hour < 12 {
		return hour + 12
	}
	return hour
}

func validHourMinute(hour, minute int) (int, int, bool) {
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return 0, 0, false
	}
	return hour, minute, true
}

func cleanupCalendarCancelTitle(text string) string {
	cleaned := cleanupByRemoving(text, []string{
		"ยกเลิกนัด", "ลบนัด", "ยกเลิกตาราง", "ลบตาราง",
		"cancel meeting", "cancel event", "cancel",
	})
	return cleanupTimeWords(cleaned)
}

func filterCalendarCancelCandidates(events []*entities.CalendarEvent, request calendarCancelRequest, loc *time.Location) []*entities.CalendarEvent {
	matches := make([]*entities.CalendarEvent, 0, len(events))
	for _, event := range events {
		if request.HasTime {
			startAt := event.StartAt.In(loc)
			if startAt.Hour() != request.Hour || startAt.Minute() != request.Minute {
				continue
			}
		}
		if request.Title != "" && !calendarTitleMatches(event.Title, request.Title) {
			continue
		}
		matches = append(matches, event)
	}
	return matches
}

func calendarTitleMatches(title, query string) bool {
	normalizedTitle := normalizeMatchText(title)
	normalizedQuery := normalizeMatchText(query)
	if normalizedQuery == "" {
		return true
	}
	if normalizedTitle == normalizedQuery {
		return true
	}
	if strings.Contains(normalizedTitle, normalizedQuery) || strings.Contains(normalizedQuery, normalizedTitle) {
		return true
	}

	titleWords := strings.Fields(strings.ToLower(title))
	queryWords := strings.Fields(strings.ToLower(query))
	if len(titleWords) == 0 || len(queryWords) == 0 {
		return false
	}

	overlap := 0
	for _, queryWord := range queryWords {
		for _, titleWord := range titleWords {
			if queryWord == titleWord || strings.Contains(titleWord, queryWord) || strings.Contains(queryWord, titleWord) {
				overlap++
				break
			}
		}
	}
	return overlap > 0
}

func normalizeMatchText(text string) string {
	lower := strings.ToLower(strings.TrimSpace(text))
	re := regexp.MustCompile(`[\s"'“”‘’.,!?;:()\[\]{}_-]+`)
	return re.ReplaceAllString(lower, "")
}

func formatCalendarCancelClarification(events []*entities.CalendarEvent, loc *time.Location) string {
	limit := len(events)
	if limit > 5 {
		limit = 5
	}

	lines := []string{"เจอหลายนัดที่ใกล้เคียงกันค่ะ เลขาขอเวลา/ชื่อนัดให้ชัดขึ้นอีกนิด เช่น \"ยกเลิกนัด 10 โมง ประชุมกับทีม\" ค่ะ"}
	for i := 0; i < limit; i++ {
		event := events[i]
		lines = append(lines, fmt.Sprintf("- %s %s", event.StartAt.In(loc).Format("02 Jan 15:04"), event.Title))
	}
	if len(events) > limit {
		lines = append(lines, fmt.Sprintf("และอีก %d นัดค่ะ", len(events)-limit))
	}
	lines = append(lines, "เลือกนัดที่ใช่แล้วส่งมาได้เลยค่ะ")
	return strings.Join(lines, "\n")
}

func formatCalendarListReply(events []*entities.CalendarEvent, loc *time.Location) string {
	if len(events) == 0 {
		return "สมุดนัดยังโล่งค่ะ ยังไม่มีนัดที่บันทึกไว้ค่ะ"
	}

	limit := len(events)
	if limit > 10 {
		limit = 10
	}

	lines := []string{"เลขาเช็กสมุดนัดให้แล้ว มีรายการนัดหมายทั้งหมดค่ะ"}
	for i := 0; i < limit; i++ {
		event := events[i]
		lines = append(lines, fmt.Sprintf("- %s %s", event.StartAt.In(loc).Format("02 Jan 2006 15:04"), event.Title))
	}
	if len(events) > limit {
		lines = append(lines, fmt.Sprintf("และอีก %d นัดค่ะ", len(events)-limit))
	}
	lines = append(lines, "เช็กให้เรียบร้อยค่ะ")
	return strings.Join(lines, "\n")
}

func formatExpenseReportReply(label string, expenses []*entities.Expense, total float64, loc *time.Location) string {
	if len(expenses) == 0 {
		return fmt.Sprintf("%sยังไม่มีรายจ่ายที่บันทึกไว้ค่ะ", label)
	}

	lines := []string{fmt.Sprintf("เลขาสรุปรายจ่าย%s ให้แล้ว รวม %.2f บาทค่ะ", label, total)}
	limit := len(expenses)
	if limit > 10 {
		limit = 10
	}

	for i := limit - 1; i >= 0; i-- {
		expense := expenses[i]
		lines = append(lines, fmt.Sprintf("- %s %s %.2f บาท", expense.SpentAt.In(loc).Format("15:04"), expense.Description, expense.Amount))
	}
	if len(expenses) > limit {
		lines = append(lines, fmt.Sprintf("และอีก %d รายการค่ะ", len(expenses)-limit))
	}
	lines = append(lines, "สรุปให้เรียบร้อยค่ะ")

	return strings.Join(lines, "\n")
}

func formatIncomeReportReply(label string, incomes []*entities.Income, total float64, loc *time.Location) string {
	if len(incomes) == 0 {
		return fmt.Sprintf("%sยังไม่มีรายรับที่บันทึกไว้ค่ะ", label)
	}

	lines := []string{fmt.Sprintf("เลขาสรุปรายรับ%s ให้แล้ว รวม %.2f บาทค่ะ", label, total)}
	limit := len(incomes)
	if limit > 10 {
		limit = 10
	}

	for i := limit - 1; i >= 0; i-- {
		income := incomes[i]
		lines = append(lines, fmt.Sprintf("- %s %s %.2f บาท", income.ReceivedAt.In(loc).Format("15:04"), income.Description, income.Amount))
	}
	if len(incomes) > limit {
		lines = append(lines, fmt.Sprintf("และอีก %d รายการค่ะ", len(incomes)-limit))
	}
	lines = append(lines, "สรุปให้เรียบร้อยค่ะ")

	return strings.Join(lines, "\n")
}

func formatCashflowReportReply(label string, incomes []*entities.Income, expenses []*entities.Expense, incomeTotal, expenseTotal float64, loc *time.Location) string {
	net := incomeTotal - expenseTotal
	lines := []string{
		fmt.Sprintf("เลขาสรุปการเงิน%s ให้แล้วค่ะ", label),
		fmt.Sprintf("- รายรับ %.2f บาท", incomeTotal),
	}

	for i := len(incomes) - 1; i >= 0; i-- {
		income := incomes[i]
		lines = append(lines, fmt.Sprintf("  - %s %.2f บาท", income.Description, income.Amount))
	}

	lines = append(lines, fmt.Sprintf("- รายจ่าย %.2f บาท", expenseTotal))
	for i := len(expenses) - 1; i >= 0; i-- {
		expense := expenses[i]
		lines = append(lines, fmt.Sprintf("  - %s %.2f บาท", expense.Description, expense.Amount))
	}

	lines = append(lines, fmt.Sprintf("- สุทธิ %.2f บาทค่ะ", net))
	return strings.Join(lines, "\n")
}

func isQuestion(text string) bool {
	questionWords := []string{"อะไร", "บ้าง", "ไหม", "มั้ย", "เท่าไหร่", "เท่าไร", "หรือยัง", "ยังไง", "?", "？"}
	return containsAny(text, questionWords)
}

func isCalendarCancelCommand(text string) bool {
	cancelWords := []string{"ยกเลิกนัด", "ลบนัด", "ยกเลิกตาราง", "ลบตาราง", "cancel meeting", "cancel event"}
	return containsAny(text, cancelWords)
}

func containsAny(text string, keywords []string) bool {
	lower := strings.ToLower(text)
	for _, keyword := range keywords {
		if strings.Contains(lower, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func hasPrefixAny(text string, prefixes []string) bool {
	lower := strings.ToLower(text)
	for _, prefix := range prefixes {
		if strings.HasPrefix(lower, strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}
