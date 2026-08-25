package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	repoPort "minyjae/go-starter/internal/core/domain/ports/repositories"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	"minyjae/go-starter/types"
	calendarnlp "minyjae/go-starter/utils/assistant/calendar"
	moneynlp "minyjae/go-starter/utils/assistant/money"
	notenlp "minyjae/go-starter/utils/assistant/note"
	remindernlp "minyjae/go-starter/utils/assistant/reminder"
	sharednlp "minyjae/go-starter/utils/assistant/shared"
	summarynlp "minyjae/go-starter/utils/assistant/summary"
	todonlp "minyjae/go-starter/utils/assistant/todo"
)

// assistantService เป็นตัวกลาง route ข้อความผู้ใช้ไปยัง feature ที่เหมาะสม
// input: สร้างจาก NewAssistantServiceImpl พร้อม repository ของแต่ละ feature และ optional intent parser
// output: service ที่คืน AssistantMessageResult พร้อม intent/reply/data
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

// NewAssistantServiceImpl สร้าง assistant service implementation
// input: repository ของ intent/todo/expense/income/calendar/reminder/note และ optional AssistantIntentParser
// output: *assistantService ที่พร้อมถูกใช้ผ่าน AssistantService interface
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

// HandleTextMessage รับข้อความผู้ใช้แล้วเลือก flow ที่ต้องทำงาน
// input: AssistantMessageInput ที่มี user id, message log id, text, now และ timezone
// output: AssistantMessageResult ที่มี intent/reply/data หรือ error จาก handler/repository
func (s *assistantService) HandleTextMessage(input types.AssistantMessageInput) (*types.AssistantMessageResult, error) {
	text := strings.TrimSpace(input.Text)
	if text == "" {
		return s.saveUnknownIntent(input, "ยังไม่เห็นข้อความที่อยากให้เลขาช่วยเลยค่ะ ส่งมาได้เลย เดี๋ยวจัดการให้ค่ะ")
	}

	loc := loadLocation(input.Timezone)
	now := input.Now.In(loc)

	if target, period, ok := moneynlp.ParseReportRequest(text, now, loc); ok {
		switch target {
		case moneynlp.ReportTargetIncome:
			return s.handleIncomeReport(input, now, loc, period)
		case moneynlp.ReportTargetCashflow:
			return s.handleCashflowReport(input, now, loc, period)
		default:
			return s.handleExpenseReport(input, now, loc, period)
		}
	}
	if calendarnlp.IsCancelCommand(text) {
		return s.handleCancelCalendarEvent(input, text, now, loc)
	}
	if calendarnlp.IsListCommand(text) {
		return s.handleCalendarList(input, now, loc)
	}
	if calendarnlp.LooksLikeIntentWithoutTime(text) {
		return s.handleCreateNote(input, text, now)
	}

	if result, ok := s.handleWithIntentParser(input, text, now, loc); ok {
		return result, nil
	}

	if summarynlp.IsTomorrowCommand(text) {
		return s.handleTomorrowSummary(input, now, loc)
	}
	if moneynlp.IsExpenseSummary(text) {
		return s.handleExpenseSummary(input, now)
	}
	if todonlp.IsCommand(text) {
		return s.handleCreateTodo(input, text, now)
	}
	if remindernlp.IsCommand(text) {
		return s.handleCreateReminder(input, text, now, loc)
	}
	if notenlp.IsCommand(text) {
		return s.handleCreateNote(input, text, now)
	}
	if amount, ok := moneynlp.ExtractAmount(text); ok && moneynlp.LooksLikeIncome(text) {
		return s.handleCreateIncome(input, text, amount, now, loc)
	}
	if amount, ok := moneynlp.ExtractAmount(text); ok && moneynlp.LooksLikeExpense(text) {
		return s.handleCreateExpense(input, text, amount, now, loc)
	}
	if calendarnlp.LooksLikeEvent(text) {
		return s.handleCreateCalendarEvent(input, text, now, loc)
	}

	return s.saveUnknownIntent(input, "ตอนนี้เลขาคนนี้ถนัด todo, รายรับ, รายจ่าย, นัดหมาย, เตือนความจำ, note และสรุปพรุ่งนี้ค่ะ ส่งงานแนวนี้มาได้เลยค่ะ")
}

// handleCreateTodo สร้าง todo จากข้อความที่ parser ตีความแล้วว่าเป็น todo
// input: AssistantMessageInput, text ข้อความดิบ, now เวลาปัจจุบัน
// output: AssistantMessageResult สำหรับ create_todo พร้อมข้อมูล todo หรือ error
func (s *assistantService) handleCreateTodo(input types.AssistantMessageInput, text string, now time.Time) (*types.AssistantMessageResult, error) {
	title := todonlp.CleanupTitle(text)
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

// handleCreateExpense สร้างรายจ่ายจากข้อความและจำนวนเงินที่ parse ได้
// input: AssistantMessageInput, text ข้อความดิบ, amount จำนวนเงิน, now เวลาปัจจุบัน, loc timezone
// output: AssistantMessageResult สำหรับ create_expense พร้อมข้อมูล expense หรือ error
func (s *assistantService) handleCreateExpense(input types.AssistantMessageInput, text string, amount float64, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	category := moneynlp.InferExpenseCategory(text)
	description := moneynlp.CleanupExpenseDescription(text)
	spentAt := moneynlp.ParseEntryTime(text, now, loc)

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
		ReplyText: moneynlp.FormatCreateReply(expense.Description, expense.Amount, expense.SpentAt, loc, moneynlp.HasExplicitEntryTime(text)),
		Data:      expense,
	}, nil
}

// handleCreateIncome สร้างรายรับจากข้อความและจำนวนเงินที่ parse ได้
// input: AssistantMessageInput, text ข้อความดิบ, amount จำนวนเงิน, now เวลาปัจจุบัน, loc timezone
// output: AssistantMessageResult สำหรับ create_income พร้อมข้อมูล income หรือ error
func (s *assistantService) handleCreateIncome(input types.AssistantMessageInput, text string, amount float64, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	category := moneynlp.InferIncomeCategory(text)
	description := moneynlp.CleanupIncomeDescription(text)
	receivedAt := moneynlp.ParseEntryTime(text, now, loc)

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
		ReplyText: moneynlp.FormatCreateReply(income.Description, income.Amount, income.ReceivedAt, loc, moneynlp.HasExplicitEntryTime(text)),
		Data:      income,
	}, nil
}

// handleExpenseSummary สรุปยอดรายจ่ายของเดือนปัจจุบัน
// input: AssistantMessageInput และ now เวลาปัจจุบัน
// output: AssistantMessageResult สำหรับ expense_summary พร้อมยอดรวม หรือ error
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

// handleExpenseReport สร้างรายงานรายจ่ายตามช่วงเวลาที่ parse ได้
// input: AssistantMessageInput, now เวลาปัจจุบัน, loc timezone, period ช่วงรายงาน
// output: AssistantMessageResult พร้อมรายการรายจ่าย ยอดรวม และช่วงเวลา หรือ error
func (s *assistantService) handleExpenseReport(input types.AssistantMessageInput, now time.Time, loc *time.Location, period moneynlp.ReportPeriod) (*types.AssistantMessageResult, error) {
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
		ReplyText: moneynlp.FormatExpenseReportReply(period.Label, expenses, total, loc),
		Data:      map[string]any{"expenses": expenses, "total": total, "start": period.Start, "end": period.End},
	}, nil
}

// handleIncomeReport สร้างรายงานรายรับตามช่วงเวลาที่ parse ได้
// input: AssistantMessageInput, now เวลาปัจจุบัน, loc timezone, period ช่วงรายงาน
// output: AssistantMessageResult พร้อมรายการรายรับ ยอดรวม และช่วงเวลา หรือ error
func (s *assistantService) handleIncomeReport(input types.AssistantMessageInput, now time.Time, loc *time.Location, period moneynlp.ReportPeriod) (*types.AssistantMessageResult, error) {
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
		ReplyText: moneynlp.FormatIncomeReportReply(period.Label, incomes, total, loc),
		Data:      map[string]any{"incomes": incomes, "total": total, "start": period.Start, "end": period.End},
	}, nil
}

// handleCashflowReport สร้างรายงานรายรับ/รายจ่าย/สุทธิตามช่วงเวลาที่ parse ได้
// input: AssistantMessageInput, now เวลาปัจจุบัน, loc timezone, period ช่วงรายงาน
// output: AssistantMessageResult พร้อมรายการรายรับ รายจ่าย ยอดรวม และ net หรือ error
func (s *assistantService) handleCashflowReport(input types.AssistantMessageInput, now time.Time, loc *time.Location, period moneynlp.ReportPeriod) (*types.AssistantMessageResult, error) {
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
		ReplyText: moneynlp.FormatCashflowReportReply(period.Label, incomes, expenses, incomeTotal, expenseTotal, loc),
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

// handleCreateCalendarEvent สร้าง calendar event จากข้อความนัดหมาย
// input: AssistantMessageInput, text ข้อความดิบ, now เวลาปัจจุบัน, loc timezone
// output: AssistantMessageResult สำหรับ create_calendar_event พร้อมข้อมูล event หรือ error
func (s *assistantService) handleCreateCalendarEvent(input types.AssistantMessageInput, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	startAt, ok := sharednlp.ParseNaturalTime(text, now, loc, 9, 0)
	if !ok {
		startAt = now.Add(24 * time.Hour)
	}
	endAt := startAt.Add(time.Hour)
	title := calendarnlp.CleanupTitle(text)
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

// handleCancelCalendarEvent ยกเลิก calendar event จากข้อความที่ผู้ใช้ระบุ
// input: AssistantMessageInput, text ข้อความดิบ, now เวลาปัจจุบัน, loc timezone
// output: AssistantMessageResult สำหรับ cancel_calendar_event พร้อมสถานะ completed/not_found/needs_clarification หรือ error
func (s *assistantService) handleCancelCalendarEvent(input types.AssistantMessageInput, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	request := calendarnlp.ParseCancelRequest(text, now, loc)
	if request.Title == "" && !request.HasDate && !request.HasTime {
		return s.saveCalendarCancelIntent(input, now, request, "needs_clarification", "อยากให้เลขาถอนนัดไหนออกจากสมุดค่ะ ระบุชื่อ วัน หรือเวลาเพิ่มอีกนิด เช่น \"ยกเลิกนัดพรุ่งนี้ 10 โมง ประชุมกับทีม\" ค่ะ", nil)
	}

	candidates, err := s.calendarRepo.ListByStartBetween(input.UserID, request.Start, request.End)
	if err != nil {
		return nil, err
	}

	matches := calendarnlp.FilterCancelCandidates(candidates, request, loc)
	if len(matches) == 0 {
		return s.saveCalendarCancelIntent(input, now, request, "not_found", "เลขาค้นสมุดนัดแล้ว ยังไม่เจอนัดที่ตรงกับข้อความนี้ค่ะ ลองระบุวัน เวลา หรือชื่อนัดให้ชัดขึ้นอีกนิดค่ะ", nil)
	}
	if len(matches) > 1 {
		return s.saveCalendarCancelIntent(input, now, request, "needs_clarification", calendarnlp.FormatCancelClarification(matches, loc), matches)
	}

	event := matches[0]
	if err := s.calendarRepo.Delete(event.ID); err != nil {
		return nil, err
	}

	reply := fmt.Sprintf("ยกเลิกนัด \"%s\" เวลา %s ออกจากปฏิทินให้แล้วค่ะ", event.Title, event.StartAt.In(loc).Format("02 Jan 2006 15:04"))
	return s.saveCalendarCancelIntent(input, now, request, "completed", reply, matches)
}

// saveCalendarCancelIntent บันทึก intent ของการยกเลิกนัดและสร้าง result ส่งกลับ
// input: AssistantMessageInput, now, request ที่ parse ได้, status, reply, matches รายการนัดที่เกี่ยวข้อง
// output: AssistantMessageResult ที่มี reply/status/matches หรือ error จาก intent repository
func (s *assistantService) saveCalendarCancelIntent(input types.AssistantMessageInput, now time.Time, request calendarnlp.CancelRequest, status, reply string, matches []*entities.CalendarEvent) (*types.AssistantMessageResult, error) {
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

// handleCreateReminder สร้าง reminder จากข้อความเตือนความจำ
// input: AssistantMessageInput, text ข้อความดิบ, now เวลาปัจจุบัน, loc timezone
// output: AssistantMessageResult สำหรับ create_reminder พร้อมข้อมูล reminder หรือ error
func (s *assistantService) handleCreateReminder(input types.AssistantMessageInput, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	remindAt, ok := sharednlp.ParseNaturalTime(text, now, loc, 9, 0)
	if !ok {
		remindAt = now.Add(24 * time.Hour)
	}
	title := remindernlp.CleanupTitle(text)
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

// handleCreateNote สร้าง note จากข้อความที่ผู้ใช้ให้จด
// input: AssistantMessageInput, text ข้อความดิบ, now เวลาปัจจุบัน
// output: AssistantMessageResult สำหรับ create_note พร้อมข้อมูล note หรือ error
func (s *assistantService) handleCreateNote(input types.AssistantMessageInput, text string, now time.Time) (*types.AssistantMessageResult, error) {
	content := notenlp.CleanupContent(text)
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

// handleTomorrowSummary สรุปนัดและ todo ของวันพรุ่งนี้
// input: AssistantMessageInput, now เวลาปัจจุบัน, loc timezone
// output: AssistantMessageResult สำหรับ tomorrow_summary พร้อม events/todos หรือ error
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

// handleCalendarList ดึงรายการ calendar event ของ user เพื่อแสดงเป็นสรุปนัด
// input: AssistantMessageInput, now เวลาปัจจุบันสำหรับบันทึก intent, loc timezone สำหรับ format เวลา
// output: AssistantMessageResult สำหรับ calendar_list พร้อม events หรือ error
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
		ReplyText: calendarnlp.FormatListReply(events, loc),
		Data:      map[string]any{"events": events},
	}, nil
}

// saveUnknownIntent บันทึกข้อความที่ระบบยังตีความไม่ได้
// input: AssistantMessageInput และ reply ที่ต้องการตอบกลับ
// output: AssistantMessageResult intent unknown หรือ error จาก intent repository
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

// loadLocation โหลด timezone จากข้อความ config/input
// input: timezone เช่น "Asia/Bangkok" หรือค่าว่าง
// output: *time.Location โดย fallback เป็น Asia/Bangkok ถ้าโหลดไม่ได้
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

// marshalEntities แปลงข้อมูล entities ของ intent เป็น JSON
// input: value ข้อมูลใด ๆ ที่ต้องการเก็บใน AssistantIntent.Entities
// output: json.RawMessage หรือ nil ถ้า marshal ไม่สำเร็จ
func marshalEntities(value any) json.RawMessage {
	data, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	return data
}
