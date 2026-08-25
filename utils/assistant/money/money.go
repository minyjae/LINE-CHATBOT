package money

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
	shared "minyjae/go-starter/utils/assistant/shared"
)

// ReportPeriod คือช่วงเวลาของรายงานการเงินที่ parse ได้จากข้อความ
// input: สร้างจาก ParseReportRequest/ParseExpenseReportPeriod
// output: struct ที่มี label, start/end สำหรับ query repository และ intent ที่ควรบันทึก
type ReportPeriod struct {
	Label  string
	Start  time.Time
	End    time.Time
	Intent string
}

// ReportTarget ระบุว่ารายงานที่ผู้ใช้ถามคือรายจ่าย รายรับ หรือ cashflow
// input: สร้างจาก detectMoneyReportTarget
// output: string enum สำหรับ route ไป handler ที่ถูกต้อง
type ReportTarget string

const (
	ReportTargetExpense  ReportTarget = "expense"
	ReportTargetIncome   ReportTarget = "income"
	ReportTargetCashflow ReportTarget = "cashflow"
)

// reportPeriodKind ระบุชนิดช่วงเวลาภายใน package money
// input: สร้างจาก parseReportPeriod
// output: string enum ที่ใช้ประกอบ intent เช่น expense_report_daily
type reportPeriodKind string

const (
	reportPeriodDaily   reportPeriodKind = "daily"
	reportPeriodWeekly  reportPeriodKind = "weekly"
	reportPeriodMonthly reportPeriodKind = "monthly"
)

// IsExpenseSummary ตรวจคำถามสรุปรายจ่ายแบบสั้นที่ผู้ใช้ถามด้วยคำว่าใช้เงิน
// input: text ข้อความดิบจากผู้ใช้
// output: true ถ้าข้อความคล้ายคำถามสรุปยอดใช้เงิน
func IsExpenseSummary(text string) bool {
	return strings.Contains(text, "ใช้เงิน") && (strings.Contains(text, "เดือน") || strings.Contains(text, "เท่าไหร่"))
}

// ParseExpenseReportPeriod parse เฉพาะช่วงเวลาของรายงานรายจ่าย
// input: text ข้อความดิบ, now เวลาปัจจุบัน, loc timezone
// output: ReportPeriod และ true เมื่อข้อความเป็น report ของรายจ่าย
func ParseExpenseReportPeriod(text string, now time.Time, loc *time.Location) (ReportPeriod, bool) {
	target, period, ok := ParseReportRequest(text, now, loc)
	if !ok || target != ReportTargetExpense {
		return ReportPeriod{}, false
	}
	return period, true
}

// ParseReportRequest parse คำถามรายงานการเงินเป็น target และช่วงเวลา
// input: text เช่น "สรุปรายจ่ายเดือนกันยายน", now เวลาปัจจุบัน, loc timezone
// output: ReportTarget, ReportPeriod, true เมื่อเป็นคำถาม report ที่ชัดเจน
func ParseReportRequest(text string, now time.Time, loc *time.Location) (ReportTarget, ReportPeriod, bool) {
	target, ok := detectMoneyReportTarget(text)
	if !ok || !shared.ContainsAny(text, moneyReportQueryWords()) {
		return "", ReportPeriod{}, false
	}
	if _, hasAmount := ExtractAmount(text); hasAmount && !shared.ContainsAny(text, strongMoneyReportQueryWords()) {
		return "", ReportPeriod{}, false
	}

	period, kind := parseReportPeriod(text, now, loc)
	period.Intent = fmt.Sprintf("%s_report_%s", target, kind)
	return target, period, true
}

// detectMoneyReportTarget ตรวจว่าคำถามรายงานต้องการรายรับ รายจ่าย หรือ cashflow
// input: text ข้อความดิบจากผู้ใช้
// output: ReportTarget และ true เมื่อพบ keyword การเงินที่รู้จัก
func detectMoneyReportTarget(text string) (ReportTarget, bool) {
	hasIncome := shared.ContainsAny(text, incomeWords())
	hasExpense := shared.ContainsAny(text, expenseWords())
	hasCashflow := shared.ContainsAny(text, cashflowWords()) || (hasIncome && hasExpense)

	switch {
	case hasCashflow:
		return ReportTargetCashflow, true
	case hasIncome:
		return ReportTargetIncome, true
	case hasExpense:
		return ReportTargetExpense, true
	default:
		return "", false
	}
}

// parseReportPeriod แปลงคำบอกช่วงเวลาในคำถามรายงานเป็น start/end
// input: text ข้อความดิบ, now เวลาปัจจุบัน, loc timezone
// output: ReportPeriod และ reportPeriodKind โดย fallback เป็นรายวันของวันที่ปัจจุบัน
func parseReportPeriod(text string, now time.Time, loc *time.Location) (ReportPeriod, reportPeriodKind) {
	current := now.In(loc)
	referenceDate, hasReferenceDate := shared.ParseReferenceDate(text, current, loc)
	if strings.Contains(strings.ToLower(text), "week") || strings.Contains(text, "สัปดาห์") || strings.Contains(text, "อาทิตย์") {
		reference := current
		if hasReferenceDate {
			reference = referenceDate
		}
		start := shared.StartOfReportWeek(reference, loc)
		return ReportPeriod{
			Label: shared.FormatReportPeriodLabel("สัปดาห์", start, start.AddDate(0, 0, 7), loc),
			Start: start,
			End:   start.AddDate(0, 0, 7),
		}, reportPeriodWeekly
	}

	if strings.Contains(strings.ToLower(text), "month") || strings.Contains(text, "เดือน") || (shared.ContainsThaiMonthName(text) && !hasReferenceDate) {
		start := time.Date(current.Year(), current.Month(), 1, 0, 0, 0, 0, loc)
		if parsed, ok := shared.ParseReportMonth(text, current, loc); ok {
			start = parsed
		} else if hasReferenceDate {
			start = time.Date(referenceDate.Year(), referenceDate.Month(), 1, 0, 0, 0, 0, loc)
		}
		return ReportPeriod{
			Label: shared.FormatReportPeriodLabel("เดือน", start, start.AddDate(0, 1, 0), loc),
			Start: start,
			End:   start.AddDate(0, 1, 0),
		}, reportPeriodMonthly
	}

	reference := current
	if hasReferenceDate {
		reference = referenceDate
	}
	start := time.Date(reference.Year(), reference.Month(), reference.Day(), 0, 0, 0, 0, loc)
	return ReportPeriod{
		Label: shared.FormatReportPeriodLabel("วันที่", start, start.AddDate(0, 0, 1), loc),
		Start: start,
		End:   start.AddDate(0, 0, 1),
	}, reportPeriodDaily
}

// IsExpenseReportQuery ตรวจว่า text เป็นคำถามรายงานรายจ่ายหรือไม่
// input: text ข้อความดิบจากผู้ใช้
// output: true ถ้า ParseReportRequest สำเร็จและ target มีคำของรายจ่าย
func IsExpenseReportQuery(text string) bool {
	_, _, ok := ParseReportRequest(text, time.Now(), time.Local)
	return ok && shared.ContainsAny(text, expenseWords())
}

// IsIncomeReportQuery ตรวจว่า text เป็นคำถามรายงานรายรับหรือไม่
// input: text ข้อความดิบจากผู้ใช้
// output: true ถ้า ParseReportRequest สำเร็จและ target มีคำของรายรับ
func IsIncomeReportQuery(text string) bool {
	_, _, ok := ParseReportRequest(text, time.Now(), time.Local)
	return ok && shared.ContainsAny(text, incomeWords())
}

// LooksLikeExpense ตรวจว่า text น่าจะเป็นรายการรายจ่ายที่ต้องบันทึกหรือไม่
// input: text ข้อความดิบจากผู้ใช้
// output: true เมื่อพบ keyword รายจ่าย และไม่ชนกับรายรับ
func LooksLikeExpense(text string) bool {
	if LooksLikeIncome(text) {
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

// LooksLikeIncome ตรวจว่า text น่าจะเป็นรายการรายรับที่ต้องบันทึกหรือไม่
// input: text ข้อความดิบจากผู้ใช้
// output: true เมื่อพบ keyword สร้างรายรับ
func LooksLikeIncome(text string) bool {
	return shared.ContainsAny(text, incomeCreateWords())
}

// expenseWords คืน keyword ที่สื่อถึงรายจ่าย
// input: ไม่มี
// output: []string keyword สำหรับ detect/report รายจ่าย
func expenseWords() []string {
	return []string{"ใช้เงิน", "รายจ่าย", "ค่าใช้จ่าย", "จ่าย", "ซื้อ", "กิน", "อาหาร", "expense", "spent", "spend"}
}

// incomeWords คืน keyword ที่สื่อถึงรายรับ
// input: ไม่มี
// output: []string keyword สำหรับ detect/report รายรับ
func incomeWords() []string {
	return []string{"รายรับ", "รับเงิน", "ได้เงิน", "รายได้", "เงินเข้า", "income", "salary"}
}

// incomeCreateWords คืน keyword สำหรับตรวจการสร้างรายการรายรับ
// input: ไม่มี
// output: []string keyword ที่ใช้กับ LooksLikeIncome
func incomeCreateWords() []string {
	return []string{"รายรับ", "รับเงิน", "ได้เงิน", "ได้รับเงิน", "รายได้", "เงินเข้า", "เงินเดือน", "ขายได้", "income", "salary"}
}

// cashflowWords คืน keyword ที่สื่อถึงรายงานภาพรวมการเงิน
// input: ไม่มี
// output: []string keyword สำหรับ detect cashflow report
func cashflowWords() []string {
	return []string{"รายรับรายจ่าย", "รายรับและรายจ่าย", "รายรับ รายจ่าย", "เงินเข้าเงินออก", "สรุปการเงิน", "การเงิน", "สุทธิ", "คงเหลือ", "cashflow", "cash flow", "balance", "net"}
}

// moneyReportQueryWords คืน keyword ที่บอกว่าข้อความเป็นคำถามรายงาน
// input: ไม่มี
// output: []string keyword แบบกว้าง เช่น "สรุป", "รายการ", "วันนี้", "เดือน"
func moneyReportQueryWords() []string {
	return []string{"อะไร", "บ้าง", "เท่าไหร่", "เท่าไร", "รวม", "สรุป", "รายการ", "ดู", "เช็ค", "เช็ก", "ย้อนหลัง", "วันนี้", "เมื่อวาน", "สัปดาห์", "อาทิตย์", "เดือน", "summary", "report", "list", "total", "today", "yesterday", "week", "month"}
}

// strongMoneyReportQueryWords คืน keyword report ที่แรงพอจะไม่สับสนกับการบันทึกยอดเงิน
// input: ไม่มี
// output: []string keyword สำหรับกันเคสมีตัวเลขจำนวนเงินแต่จริง ๆ เป็นคำถาม report
func strongMoneyReportQueryWords() []string {
	return []string{"อะไร", "บ้าง", "เท่าไหร่", "เท่าไร", "รวม", "สรุป", "รายการ", "ดู", "เช็ค", "เช็ก", "ย้อนหลัง", "summary", "report", "list", "total"}
}

// ExtractAmount ดึงจำนวนเงินจากข้อความ
// input: text ข้อความดิบ เช่น "ซื้อกาแฟ 80 บาท" หรือ "วันที่ 7 กินข้าว 120"
// output: amount และ true เมื่อ parse ตัวเลขจำนวนเงินได้ โดยพยายามข้ามเลขวันที่
func ExtractAmount(text string) (float64, bool) {
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

// looksLikeDateNumber ตรวจว่าเลขตำแหน่งนี้น่าจะเป็นเลขวันที่ ไม่ใช่จำนวนเงิน
// input: text ข้อความเต็ม, start/end ตำแหน่ง index ของเลขที่พบ
// output: true ถ้าเลขอยู่หลัง "วันที่" หรืออยู่ใน pattern ที่มี / หรือ -
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

// ParseEntryTime parse วันที่/เวลาของรายการรับจ่าย
// input: text ข้อความดิบ, now เวลาปัจจุบัน, loc timezone
// output: time.Time ของรายการ โดยถ้าไม่มีวันที่/เวลาใช้ now เป็นค่า fallback
func ParseEntryTime(text string, now time.Time, loc *time.Location) time.Time {
	current := now.In(loc)
	hour, minute, hasTime := shared.ExtractHourMinute(text)
	if date, ok := shared.ParseReferenceDate(text, current, loc); ok {
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

// MergeParsedTime เติม explicit time จากข้อความลงใน parsedAt ถ้า parsedAt ยังเป็นเวลาเที่ยงคืน
// input: parsedAt วันที่ที่ parse มาก่อน, text ข้อความดิบ, loc timezone
// output: time.Time ที่รวมวันที่เดิมกับเวลาจาก text หรือคืน parsedAt เดิมถ้าไม่ต้องเติม
func MergeParsedTime(parsedAt time.Time, text string, loc *time.Location) time.Time {
	parsedAt = parsedAt.In(loc)
	hour, minute, hasTime := shared.ExtractHourMinute(text)
	if hasTime && parsedAt.Hour() == 0 && parsedAt.Minute() == 0 && parsedAt.Second() == 0 && parsedAt.Nanosecond() == 0 {
		return time.Date(parsedAt.Year(), parsedAt.Month(), parsedAt.Day(), hour, minute, 0, 0, loc)
	}
	return parsedAt
}

// HasExplicitEntryTime ตรวจว่ารายการรับจ่ายมีเวลาชัดเจนในข้อความหรือไม่
// input: text ข้อความดิบจากผู้ใช้
// output: true ถ้าเจอเวลาที่ parse เป็น hour/minute ได้
func HasExplicitEntryTime(text string) bool {
	_, _, ok := shared.ExtractHourMinute(text)
	return ok
}

// FormatCreateReply สร้าง reply หลังบันทึกรายการรับจ่าย
// input: description รายละเอียด, amount จำนวนเงิน, occurredAt เวลารายการ, loc timezone, includeTime กำหนดว่าจะแสดงเวลาหรือไม่
// output: string reply text สำหรับส่งกลับผู้ใช้
func FormatCreateReply(description string, amount float64, occurredAt time.Time, loc *time.Location, includeTime bool) string {
	description = shared.NormalizeDescriptionSpaces(description)
	if description == "" {
		description = "รายการ"
	}
	if includeTime {
		return fmt.Sprintf("ลงบัญชี %s %.2f บาท ตอน %s น. ให้เรียบร้อยค่ะ", description, amount, occurredAt.In(loc).Format("15:04"))
	}
	return fmt.Sprintf("ลงบัญชี %s %.2f บาทให้เรียบร้อยค่ะ", description, amount)
}

// InferExpenseCategory เดาหมวดหมู่รายจ่ายจาก keyword ในข้อความ
// input: text ข้อความดิบหรือ description
// output: string category เช่น "food", "transport" หรือ "uncategorized"
func InferExpenseCategory(text string) string {
	if strings.Contains(text, "ข้าว") || strings.Contains(text, "กิน") || strings.Contains(text, "อาหาร") || strings.Contains(text, "กาแฟ") {
		return "food"
	}
	if strings.Contains(text, "รถ") || strings.Contains(text, "แท็กซี่") || strings.Contains(text, "เดินทาง") {
		return "transport"
	}
	return "uncategorized"
}

// InferIncomeCategory เดาหมวดหมู่รายรับจาก keyword ในข้อความ
// input: text ข้อความดิบหรือ description
// output: string category เช่น "salary", "sales", "investment" หรือ "uncategorized"
func InferIncomeCategory(text string) string {
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

// CleanupExpenseDescription ตัดจำนวนเงิน วันที่ และเวลาออก เพื่อเหลือรายละเอียดรายจ่าย
// input: text ข้อความดิบจากผู้ใช้
// output: string description หรือ "รายจ่าย" เมื่อไม่เหลือข้อความ
func CleanupExpenseDescription(text string) string {
	return cleanupMoneyDescription(text, "รายจ่าย", nil)
}

// CleanupIncomeDescription ตัดคำสั่ง จำนวนเงิน วันที่ และเวลาออก เพื่อเหลือรายละเอียดรายรับ
// input: text ข้อความดิบจากผู้ใช้
// output: string description หรือ "รายรับ" เมื่อไม่เหลือข้อความ
func CleanupIncomeDescription(text string) string {
	return cleanupMoneyDescription(text, "รายรับ", []string{"บันทึก", "เพิ่ม", "รายรับ", "รับเงิน", "ได้เงิน", "ได้รับเงิน", "เงินเข้า", "income", "Income"})
}

// cleanupMoneyDescription ทำความสะอาด description กลางสำหรับรายรับ/รายจ่าย
// input: text ข้อความดิบ, fallback ค่าที่ใช้เมื่อ description ว่าง, removers token เฉพาะ feature ที่ต้องลบเพิ่ม
// output: string description ที่พร้อมบันทึก
func cleanupMoneyDescription(text, fallback string, removers []string) string {
	cleaned := strings.TrimSpace(text)
	cleaned = removeMoneyTimePhrases(cleaned)
	cleaned = removeMoneyDatePhrases(cleaned)

	amount := regexp.MustCompile(`[0-9]+(?:\.[0-9]+)?\s*(บาท|บ|thb|THB)?`)
	cleaned = amount.ReplaceAllString(cleaned, "")
	cleaned = shared.CleanupByRemoving(cleaned, removers)
	cleaned = shared.CleanupByRemoving(cleaned, []string{"วันนี้", "พรุ่งนี้", "เมื่อวาน", "ตอน", "เวลา"})
	cleaned = shared.NormalizeDescriptionSpaces(cleaned)

	if cleaned == "" {
		return fallback
	}
	return cleaned
}

// removeMoneyTimePhrases ลบ phrase เวลาออกจากข้อความการเงิน
// input: text ข้อความดิบหรือ description ที่ยังมีเวลา
// output: string ข้อความหลังตัดเวลาออก
func removeMoneyTimePhrases(text string) string {
	patterns := []string{
		`(?:ตอน|เวลา)?\s*\d{1,2}\s*โมง(?:\s*(?:เช้า|เย็น|ค่ำ))?(?:\s*(?:ครึ่ง|\d{1,2}(?:\s*นาที)?))?(?:\s*(?:เช้า|เย็น|ค่ำ))?`,
		`(?:ตอน|เวลา)?\s*\d{1,2}\s*ทุ่ม(?:\s*(?:ครึ่ง|\d{1,2}(?:\s*นาที)?))?`,
		`(?:ตอน|เวลา)?\s*บ่าย\s*\d{1,2}(?:\s*โมง)?(?:\s*(?:ครึ่ง|\d{1,2}(?:\s*นาที)?))?`,
		`(?:ตอน|เวลา)?\s*ตี\s*\d{1,2}(?:\s*(?:ครึ่ง|\d{1,2}(?:\s*นาที)?))?`,
		`(?:ตอน|เวลา)?\s*\d{1,2}[:.]\d{2}\s*(?:น\.?|นาฬิกา)?`,
		`(?:ตอน|เวลา)?\s*\d{1,2}\s*(?:น\.|นาฬิกา)(?:\s*(?:ครึ่ง|\d{1,2}(?:\s*นาที)?))?`,
	}
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		text = re.ReplaceAllString(text, "")
	}
	return text
}

// removeMoneyDatePhrases ลบ phrase วันที่ออกจากข้อความการเงิน
// input: text ข้อความดิบหรือ description ที่ยังมีวันที่
// output: string ข้อความหลังตัดวันที่ออก
func removeMoneyDatePhrases(text string) string {
	patterns := []string{
		`\b\d{4}-\d{1,2}-\d{1,2}\b`,
		`\b\d{1,2}[/-]\d{1,2}(?:[/-]\d{2,4})?\b`,
		`วันที่\s*\d{1,2}\s*เดือน\s*\d{1,2}(?:\s*ปี\s*\d{2,4})?`,
		`วันที่\s*\d{1,2}\s*(?:` + shared.ThaiMonthPattern() + `)?(?:\s*\d{2,4})?`,
	}
	for _, pattern := range patterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		text = re.ReplaceAllString(text, "")
	}
	return text
}

// FormatExpenseReportReply สร้าง reply สรุปรายจ่าย
// input: label ช่วงเวลา, expenses รายการรายจ่าย, total ยอดรวม, loc timezone สำหรับแสดงเวลา
// output: string รายงานรายจ่ายพร้อมรายการสูงสุด 10 รายการ
func FormatExpenseReportReply(label string, expenses []*entities.Expense, total float64, loc *time.Location) string {
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

// FormatIncomeReportReply สร้าง reply สรุปรายรับ
// input: label ช่วงเวลา, incomes รายการรายรับ, total ยอดรวม, loc timezone สำหรับแสดงเวลา
// output: string รายงานรายรับพร้อมรายการสูงสุด 10 รายการ
func FormatIncomeReportReply(label string, incomes []*entities.Income, total float64, loc *time.Location) string {
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

// FormatCashflowReportReply สร้าง reply สรุปเงินเข้า/เงินออก/สุทธิ
// input: label ช่วงเวลา, incomes/expenses รายการ, incomeTotal/expenseTotal ยอดรวม, loc timezone
// output: string รายงาน cashflow พร้อมยอดรายรับ รายจ่าย และสุทธิ
func FormatCashflowReportReply(label string, incomes []*entities.Income, expenses []*entities.Expense, incomeTotal, expenseTotal float64, loc *time.Location) string {
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
