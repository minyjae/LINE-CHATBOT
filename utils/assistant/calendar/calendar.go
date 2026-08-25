package calendar

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	shared "minyjae/go-starter/utils/assistant/shared"
)

// CancelRequest คือผลลัพธ์จากการ parse คำสั่งยกเลิกนัด
// input: สร้างจาก ParseCancelRequest โดยใช้ข้อความผู้ใช้ เวลาปัจจุบัน และ timezone
// output: struct ที่บอก title/ช่วงวันที่/เวลา สำหรับนำไป filter นัดที่จะยกเลิก
type CancelRequest struct {
	Title   string
	Start   time.Time
	End     time.Time
	HasDate bool
	HasTime bool
	Hour    int
	Minute  int
}

// IsListCommand ตรวจว่า text เป็นคำสั่งขอดูนัดหรือรายการ calendar หรือไม่
// input: text ข้อความดิบจากผู้ใช้ เช่น "ดูนัดทั้งหมด", "มีนัดอะไรบ้าง"
// output: true ถ้าเป็นคำสั่ง list calendar
func IsListCommand(text string) bool {
	phrases := []string{
		"ดูนัดวันนี้",
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
	return shared.ContainsAny(text, phrases)
}

// LooksLikeEvent ตรวจว่า text น่าจะเป็นคำสั่งสร้าง calendar event หรือไม่
// input: text ข้อความดิบจากผู้ใช้
// output: true เมื่อไม่ใช่คำถาม/คำสั่งยกเลิก, มีเวลาชัดเจน, และมี prefix หรือ keyword ที่เหมือนการนัดหมาย
func LooksLikeEvent(text string) bool {
	if shared.IsQuestion(text) || IsCancelCommand(text) {
		return false
	}
	if !shared.HasExplicitTime(text) {
		return false
	}

	calendarPrefixes := []string{
		"นัด", "เพิ่มนัด", "ลงนัด", "บันทึกนัด", "สร้างนัด",
		"เพิ่มตาราง", "ลงตาราง", "บันทึกตาราง", "เพิ่ม calendar", "calendar",
	}
	if shared.HasPrefixAny(strings.TrimSpace(text), calendarPrefixes) {
		return true
	}

	eventWords := []string{"ประชุม", "นัด", "เจอ", "คุย", "โทรหา", "call", "meeting"}
	return shared.ContainsAny(text, eventWords)
}

// LooksLikeIntentWithoutTime ตรวจว่า text ดูเหมือนอยากสร้างนัด แต่ยังไม่ได้บอกเวลา
// input: text ข้อความดิบจากผู้ใช้ เช่น "นัดดูรถ"
// output: true ถ้าเป็น intent calendar แต่ไม่มี explicit time เพื่อให้ service ถามเวลาต่อ
func LooksLikeIntentWithoutTime(text string) bool {
	if shared.IsQuestion(text) || IsCancelCommand(text) || shared.HasExplicitTime(text) {
		return false
	}

	calendarPrefixes := []string{
		"นัด", "เพิ่มนัด", "ลงนัด", "บันทึกนัด", "สร้างนัด",
		"เพิ่มตาราง", "ลงตาราง", "บันทึกตาราง", "เพิ่ม calendar", "calendar",
	}
	eventWords := []string{"ประชุม", "นัด", "เจอ", "คุย", "โทรหา", "meeting"}
	return shared.HasPrefixAny(strings.TrimSpace(text), calendarPrefixes) || shared.ContainsAny(text, eventWords)
}

// CleanupTitle ตัดคำสั่ง วันที่ และเวลาออกจากข้อความสร้างนัด เพื่อเหลือ title
// input: text ข้อความดิบ เช่น "นัดดูรถ วันที่ 7 เดือน 9 ตอน 9 โมง 40"
// output: string title ที่สะอาด เช่น "ดูรถ"
func CleanupTitle(text string) string {
	cleaned := stripCalendarTimeClause(text)
	cleaned = shared.RemoveTimePhrases(cleaned)
	cleaned = shared.RemoveDatePhrases(cleaned)
	cleaned = shared.CleanupByRemoving(cleaned, []string{
		"เพิ่มนัด", "ลงนัด", "บันทึกนัด", "สร้างนัด", "นัด",
		"เพิ่มตาราง", "ลงตาราง", "บันทึกตาราง", "calendar", "Calendar",
	})
	cleaned = shared.CleanupByRemoving(cleaned, []string{"วันนี้", "พรุ่งนี้", "เมื่อวาน"})
	return shared.NormalizeDescriptionSpaces(cleaned)
}

// stripCalendarTimeClause ตัดข้อความตั้งแต่คำว่า "ตอน" หรือ "เวลา" เพื่อกันเวลาไหลเข้า title
// input: text ข้อความดิบหรือข้อความที่ยังมี phrase เวลา
// output: string ส่วนก่อน phrase เวลา หรือ text เดิมถ้าไม่พบ phrase เวลา
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

// ParseCancelRequest แปลงคำสั่งยกเลิกนัดเป็นเงื่อนไขสำหรับค้นหานัดที่จะลบ
// input: text ข้อความดิบ, now เวลาปัจจุบัน, loc timezone
// output: CancelRequest ที่มี title, ช่วงวันที่, และเวลา ถ้าผู้ใช้ระบุมา
func ParseCancelRequest(text string, now time.Time, loc *time.Location) CancelRequest {
	current := now.In(loc)
	title := CleanupCancelTitle(text)
	hour, minute, hasTime := shared.ExtractHourMinute(text)
	hasDate := strings.Contains(text, "วันนี้") || strings.Contains(text, "พรุ่งนี้")

	if strings.Contains(text, "พรุ่งนี้") {
		date := current.AddDate(0, 0, 1)
		start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
		return CancelRequest{
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
		return CancelRequest{
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
	return CancelRequest{
		Title:   title,
		Start:   start,
		End:     start.AddDate(0, 0, 30),
		HasDate: hasDate,
		HasTime: hasTime,
		Hour:    hour,
		Minute:  minute,
	}
}

// CleanupCancelTitle ตัดคำสั่งยกเลิกและคำเวลาออก เพื่อเหลือชื่อที่ใช้ match กับ event title
// input: text ข้อความดิบ เช่น "ยกเลิกนัด 10 โมง ประชุมทีม"
// output: string title/query เช่น "ประชุมทีม"
func CleanupCancelTitle(text string) string {
	cleaned := shared.CleanupByRemoving(text, []string{
		"ยกเลิกนัด", "ลบนัด", "ยกเลิกตาราง", "ลบตาราง",
		"cancel meeting", "cancel event", "cancel",
	})
	return shared.CleanupTimeWords(cleaned)
}

// FilterCancelCandidates กรอง event ที่ตรงกับเงื่อนไขยกเลิกนัด
// input: events รายการนัดในช่วงที่ query มา, request เงื่อนไขจาก ParseCancelRequest, loc timezone สำหรับเทียบเวลา
// output: []*CalendarEvent เฉพาะรายการที่เวลาและ title ตรงกับ request
func FilterCancelCandidates(events []*entities.CalendarEvent, request CancelRequest, loc *time.Location) []*entities.CalendarEvent {
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

// calendarTitleMatches เทียบ title ของ event กับ query แบบยืดหยุ่น
// input: title ชื่อนัดจริง, query ข้อความที่ผู้ใช้ระบุ
// output: true ถ้าตรงกันแบบเต็มคำ, contains กัน, หรือมีคำซ้อนทับกัน
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

// normalizeMatchText normalize ข้อความก่อนนำไปเทียบ title
// input: text ข้อความใด ๆ
// output: string ตัวพิมพ์เล็ก ตัดช่องว่างและ punctuation ที่ไม่จำเป็นออก
func normalizeMatchText(text string) string {
	lower := strings.ToLower(strings.TrimSpace(text))
	re := regexp.MustCompile(`[\s"'“”‘’.,!?;:()\[\]{}_-]+`)
	return re.ReplaceAllString(lower, "")
}

// FormatCancelClarification สร้าง reply เมื่อเจอหลาย event ที่น่าจะยกเลิกได้
// input: events รายการนัดที่ match หลายรายการ, loc timezone สำหรับแสดงเวลา
// output: string ข้อความให้ผู้ใช้ระบุชื่อหรือเวลาให้ชัดขึ้น
func FormatCancelClarification(events []*entities.CalendarEvent, loc *time.Location) string {
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

// FormatListReply สร้าง reply สำหรับแสดงรายการนัด
// input: events รายการนัดจาก repository, loc timezone สำหรับ format เวลา
// output: string ข้อความสรุปรายการนัด สูงสุด 10 รายการแรก พร้อมจำนวนที่เหลือ
func FormatListReply(events []*entities.CalendarEvent, loc *time.Location) string {
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

// IsCancelCommand ตรวจว่า text เป็นคำสั่งยกเลิก/ลบนัดหรือไม่
// input: text ข้อความดิบจากผู้ใช้
// output: true ถ้าพบคำสั่งยกเลิก calendar event
func IsCancelCommand(text string) bool {
	cancelWords := []string{"ยกเลิกนัด", "ลบนัด", "ยกเลิกตาราง", "ลบตาราง", "cancel meeting", "cancel event"}
	return shared.ContainsAny(text, cancelWords)
}
