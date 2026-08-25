package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	"minyjae/go-starter/types"
	calendarnlp "minyjae/go-starter/utils/assistant/calendar"
	moneynlp "minyjae/go-starter/utils/assistant/money"
	notenlp "minyjae/go-starter/utils/assistant/note"
	remindernlp "minyjae/go-starter/utils/assistant/reminder"
	sharednlp "minyjae/go-starter/utils/assistant/shared"
	todonlp "minyjae/go-starter/utils/assistant/todo"
)

// handleWithIntentParser ลองใช้ AI intent parser ก่อน fallback ไป rule-based parser
// input: AssistantMessageInput, text ข้อความดิบ, now เวลาปัจจุบัน, loc timezone
// output: AssistantMessageResult และ true เมื่อ parser มั่นใจพอและ handler ทำงานสำเร็จ
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

// handleParsedIntent route intent ที่ AI parser ตีความแล้วไปยัง handler ที่ตรงกัน
// input: AssistantMessageInput, parsed intent/entities, text ข้อความดิบ, now, loc
// output: AssistantMessageResult จาก handler ที่เลือก หรือ error เมื่อ intent ใช้งานไม่ได้
func (s *assistantService) handleParsedIntent(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	switch parsed.Intent {
	case "create_expense":
		return s.handleParsedCreateExpense(input, parsed, text, now, loc)
	case "create_income":
		return s.handleParsedCreateIncome(input, parsed, text, now, loc)
	case "expense_summary":
		return s.handleExpenseSummary(input, now)
	case "expense_report_daily", "expense_report_weekly", "expense_report_monthly":
		return s.handleParsedExpenseReport(input, parsed, text, now, loc)
	case "income_report_daily", "income_report_weekly", "income_report_monthly":
		return s.handleParsedIncomeReport(input, parsed, text, now, loc)
	case "cashflow_report_daily", "cashflow_report_weekly", "cashflow_report_monthly":
		return s.handleParsedCashflowReport(input, parsed, text, now, loc)
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

// handleParsedExpenseReport สร้างรายงานรายจ่ายจากช่วงเวลาที่ AI parser ส่งมา
// input: AssistantMessageInput, parsed intent, text fallback, now, loc
// output: AssistantMessageResult ของ expense report หรือ error ถ้าหาช่วงเวลาไม่ได้
func (s *assistantService) handleParsedExpenseReport(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	period, ok := expenseReportPeriodFromParsed(parsed, now, loc)
	if !ok {
		period, ok = moneynlp.ParseExpenseReportPeriod(text, now, loc)
	}
	if !ok {
		return nil, fmt.Errorf("parsed expense report period is missing")
	}
	return s.handleExpenseReport(input, now, loc, period)
}

// handleParsedIncomeReport สร้างรายงานรายรับจากช่วงเวลาที่ AI parser ส่งมา
// input: AssistantMessageInput, parsed intent, text fallback, now, loc
// output: AssistantMessageResult ของ income report หรือ error ถ้าหาช่วงเวลาไม่ได้
func (s *assistantService) handleParsedIncomeReport(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	period, ok := expenseReportPeriodFromParsed(parsed, now, loc)
	if !ok {
		_, period, ok = moneynlp.ParseReportRequest(text, now, loc)
	}
	if !ok {
		return nil, fmt.Errorf("parsed income report period is missing")
	}
	return s.handleIncomeReport(input, now, loc, period)
}

// handleParsedCashflowReport สร้างรายงาน cashflow จากช่วงเวลาที่ AI parser ส่งมา
// input: AssistantMessageInput, parsed intent, text fallback, now, loc
// output: AssistantMessageResult ของ cashflow report หรือ error ถ้าหาช่วงเวลาไม่ได้
func (s *assistantService) handleParsedCashflowReport(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	period, ok := expenseReportPeriodFromParsed(parsed, now, loc)
	if !ok {
		_, period, ok = moneynlp.ParseReportRequest(text, now, loc)
	}
	if !ok {
		return nil, fmt.Errorf("parsed cashflow report period is missing")
	}
	return s.handleCashflowReport(input, now, loc, period)
}

// handleParsedCreateExpense สร้างรายจ่ายจาก entities ที่ AI parser ส่งมา
// input: AssistantMessageInput, parsed intent/entities, text fallback, now, loc
// output: AssistantMessageResult สำหรับ create_expense หรือ error ถ้า amount หาย/บันทึกไม่สำเร็จ
func (s *assistantService) handleParsedCreateExpense(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	if parsed.Entities.Amount <= 0 {
		return nil, fmt.Errorf("parsed expense amount is missing")
	}

	spentAt := moneynlp.ParseEntryTime(text, now, loc)
	if parsedAt, ok := parseOptionalIntentTime(parsed.Entities.SpentAt, loc); ok {
		spentAt = moneynlp.MergeParsedTime(parsedAt, text, loc)
	}
	description := moneynlp.CleanupExpenseDescription(firstNonEmpty(parsed.Entities.Description, parsed.Entities.Title, text))

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
		ReplyText: moneynlp.FormatCreateReply(expense.Description, expense.Amount, expense.SpentAt, loc, moneynlp.HasExplicitEntryTime(text)),
		Data:      expense,
	}, nil
}

// handleParsedCreateIncome สร้างรายรับจาก entities ที่ AI parser ส่งมา
// input: AssistantMessageInput, parsed intent/entities, text fallback, now, loc
// output: AssistantMessageResult สำหรับ create_income หรือ error ถ้า amount หาย/บันทึกไม่สำเร็จ
func (s *assistantService) handleParsedCreateIncome(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	if parsed.Entities.Amount <= 0 {
		return nil, fmt.Errorf("parsed income amount is missing")
	}

	receivedAt := moneynlp.ParseEntryTime(text, now, loc)
	if parsedAt, ok := parseOptionalIntentTime(parsed.Entities.ReceivedAt, loc); ok {
		receivedAt = moneynlp.MergeParsedTime(parsedAt, text, loc)
	} else if parsedAt, ok := parseOptionalIntentTime(parsed.Entities.SpentAt, loc); ok {
		receivedAt = moneynlp.MergeParsedTime(parsedAt, text, loc)
	}
	description := moneynlp.CleanupIncomeDescription(firstNonEmpty(parsed.Entities.Description, parsed.Entities.Title, text))

	income, err := s.incomeRepo.Create(&entities.Income{
		UserID:          input.UserID,
		Amount:          parsed.Entities.Amount,
		Currency:        defaultString(parsed.Entities.Currency, "THB"),
		Category:        defaultString(parsed.Entities.Category, "uncategorized"),
		Description:     defaultString(description, "income"),
		ReceivedAt:      receivedAt,
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
		ReplyText: moneynlp.FormatCreateReply(income.Description, income.Amount, income.ReceivedAt, loc, moneynlp.HasExplicitEntryTime(text)),
		Data:      income,
	}, nil
}

// handleParsedCreateTodo สร้าง todo จาก title/content ที่ AI parser ส่งมา
// input: AssistantMessageInput, parsed intent/entities, text fallback, now, loc
// output: AssistantMessageResult สำหรับ create_todo หรือ error ถ้า title หาย/บันทึกไม่สำเร็จ
func (s *assistantService) handleParsedCreateTodo(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	title := firstNonEmpty(parsed.Entities.Title, parsed.Entities.Content, todonlp.CleanupTitle(text))
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
		ReplyText: fmt.Sprintf("จัด todo \"%s\" เข้าลิสต์ให้เรียบร้อยค่ะ", todo.Title),
		Data:      todo,
	}, nil
}

// handleParsedCreateReminder สร้าง reminder จาก title/remind_at ที่ AI parser ส่งมา
// input: AssistantMessageInput, parsed intent/entities, text fallback, now, loc
// output: AssistantMessageResult สำหรับ create_reminder หรือ error ถ้า title/remind_at หาย
func (s *assistantService) handleParsedCreateReminder(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	title := firstNonEmpty(parsed.Entities.Title, parsed.Entities.Content, remindernlp.CleanupTitle(text))
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
		ReplyText: fmt.Sprintf("ตั้งเตือน \"%s\" เวลา %s ให้แล้วค่ะ เดี๋ยวถึงเวลาจะสะกิดให้ค่ะ", reminder.Title, reminder.RemindAt.Format("02 Jan 2006 15:04")),
		Data:      reminder,
	}, nil
}

// handleParsedCreateNote สร้าง note จาก content/description ที่ AI parser ส่งมา
// input: AssistantMessageInput, parsed intent/entities, text fallback, now
// output: AssistantMessageResult สำหรับ create_note หรือ error ถ้า content หาย/บันทึกไม่สำเร็จ
func (s *assistantService) handleParsedCreateNote(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, text string, now time.Time) (*types.AssistantMessageResult, error) {
	content := firstNonEmpty(parsed.Entities.Content, parsed.Entities.Description, notenlp.CleanupContent(text))
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
		ReplyText: "จดโน้ตใส่สมุดให้เรียบร้อยค่ะ",
		Data:      note,
	}, nil
}

// handleParsedCreateCalendarEvent สร้าง calendar event จาก title/start_at ที่ AI parser ส่งมา
// input: AssistantMessageInput, parsed intent/entities, text fallback, now, loc
// output: AssistantMessageResult สำหรับ create_calendar_event หรือ fallback เป็น note ถ้าข้อความไม่มีเวลา explicit
func (s *assistantService) handleParsedCreateCalendarEvent(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, text string, now time.Time, loc *time.Location) (*types.AssistantMessageResult, error) {
	if !sharednlp.HasExplicitTime(text) {
		return s.handleCreateNote(input, text, now)
	}

	title := firstNonEmpty(parsed.Entities.Title, calendarnlp.CleanupTitle(text))
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
		ReplyText: fmt.Sprintf("ลงนัด \"%s\" เวลา %s ให้เรียบร้อยค่ะ", event.Title, event.StartAt.Format("02 Jan 2006 15:04")),
		Data:      event,
	}, nil
}

// saveParsedUnknownIntent บันทึก intent ที่ AI parser ส่งมาแต่ service ยังไม่รองรับ
// input: AssistantMessageInput, parsed intent/entities, now
// output: AssistantMessageResult intent unknown หรือ error จาก intent repository
func (s *assistantService) saveParsedUnknownIntent(input types.AssistantMessageInput, parsed *types.ParsedAssistantIntent, now time.Time) (*types.AssistantMessageResult, error) {
	if err := s.saveParsedIntent(input, parsed, now); err != nil {
		return nil, err
	}
	return &types.AssistantMessageResult{
		Intent:    "unknown",
		ReplyText: "เลขาขออภัยค่ะ ยังตีความคำสั่งนี้ไม่ออกค่ะ",
	}, nil
}

// saveParsedIntent บันทึก intent/entities/confidence ที่ AI parser ตีความได้
// input: AssistantMessageInput, parsed intent/entities, now
// output: error nil เมื่อบันทึก AssistantIntent สำเร็จ
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

// parseOptionalIntentTime parse เวลาจาก string entity ที่ AI parser ส่งมา
// input: value เวลาในรูปแบบ RFC3339, 2006-01-02T15:04:05, 2006-01-02T15:04 หรือ 2006-01-02; loc timezone
// output: time.Time ใน loc และ true เมื่อ parse ได้
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

// firstNonEmpty คืนค่าข้อความแรกที่ไม่ว่างหลัง trim
// input: values รายการข้อความที่ต้องการเลือก
// output: string ค่าแรกที่ไม่ว่าง หรือ "" ถ้าทุกค่าว่าง
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

// optionalString แปลง string ว่างให้เป็น nil pointer
// input: value ข้อความที่อาจว่าง
// output: *string เมื่อ value ไม่ว่าง หรือ nil เมื่อ value ว่าง
func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

// expenseReportPeriodFromParsed แปลง start_at/end_at จาก parsed intent เป็น ReportPeriod
// input: parsed intent ที่มี StartAt/EndAt, now เวลาปัจจุบัน, loc timezone
// output: ReportPeriod พร้อม label/intent หรือ false ถ้า start/end หายหรือช่วงเวลาไม่ถูกต้อง
func expenseReportPeriodFromParsed(parsed *types.ParsedAssistantIntent, now time.Time, loc *time.Location) (moneynlp.ReportPeriod, bool) {
	start, hasStart := parseOptionalIntentTime(parsed.Entities.StartAt, loc)
	end, hasEnd := parseOptionalIntentTime(parsed.Entities.EndAt, loc)
	if !hasStart || !hasEnd || !end.After(start) {
		return moneynlp.ReportPeriod{}, false
	}

	label := "ช่วงที่ขอ"
	switch parsed.Intent {
	case "expense_report_daily":
		label = sharednlp.FormatReportPeriodLabel("วันที่", start, end, loc)
	case "expense_report_weekly":
		label = sharednlp.FormatReportPeriodLabel("สัปดาห์", start, end, loc)
	case "expense_report_monthly":
		label = sharednlp.FormatReportPeriodLabel("เดือน", start, end, loc)
	case "income_report_daily", "cashflow_report_daily":
		label = sharednlp.FormatReportPeriodLabel("วันที่", start, end, loc)
	case "income_report_weekly", "cashflow_report_weekly":
		label = sharednlp.FormatReportPeriodLabel("สัปดาห์", start, end, loc)
	case "income_report_monthly", "cashflow_report_monthly":
		label = sharednlp.FormatReportPeriodLabel("เดือน", start, end, loc)
	}

	return moneynlp.ReportPeriod{
		Label:  label,
		Start:  start,
		End:    end,
		Intent: parsed.Intent,
	}, true
}
