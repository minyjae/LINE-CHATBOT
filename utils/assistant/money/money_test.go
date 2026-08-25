package money

import (
	"testing"
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

func TestParseMoneyReportRequestSpecificWeek(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	now := time.Date(2026, time.August, 24, 12, 0, 0, 0, loc)

	target, period, ok := ParseReportRequest("สรุปรายจ่ายสัปดาห์ของวันที่ 2026-08-10", now, loc)
	if !ok {
		t.Fatal("expected report request")
	}
	if target != ReportTargetExpense {
		t.Fatalf("target = %q, want %q", target, ReportTargetExpense)
	}
	wantStart := time.Date(2026, time.August, 10, 0, 0, 0, 0, loc)
	wantEnd := time.Date(2026, time.August, 17, 0, 0, 0, 0, loc)
	if !period.Start.Equal(wantStart) || !period.End.Equal(wantEnd) {
		t.Fatalf("period = %s..%s, want %s..%s", period.Start, period.End, wantStart, wantEnd)
	}
	if period.Intent != "expense_report_weekly" {
		t.Fatalf("intent = %q", period.Intent)
	}
}

func TestParseMoneyReportRequestThaiMonthBuddhistYear(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	now := time.Date(2026, time.August, 24, 12, 0, 0, 0, loc)

	target, period, ok := ParseReportRequest("สรุปรายรับเดือนสิงหาคม 2569", now, loc)
	if !ok {
		t.Fatal("expected report request")
	}
	if target != ReportTargetIncome {
		t.Fatalf("target = %q, want %q", target, ReportTargetIncome)
	}
	wantStart := time.Date(2026, time.August, 1, 0, 0, 0, 0, loc)
	wantEnd := time.Date(2026, time.September, 1, 0, 0, 0, 0, loc)
	if !period.Start.Equal(wantStart) || !period.End.Equal(wantEnd) {
		t.Fatalf("period = %s..%s, want %s..%s", period.Start, period.End, wantStart, wantEnd)
	}
	if period.Intent != "income_report_monthly" {
		t.Fatalf("intent = %q", period.Intent)
	}
}

func TestParseMoneyReportRequestCashflow(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	now := time.Date(2026, time.August, 24, 12, 0, 0, 0, loc)

	target, period, ok := ParseReportRequest("สรุปรายรับรายจ่ายวันนี้", now, loc)
	if !ok {
		t.Fatal("expected report request")
	}
	if target != ReportTargetCashflow {
		t.Fatalf("target = %q, want %q", target, ReportTargetCashflow)
	}
	if period.Intent != "cashflow_report_daily" {
		t.Fatalf("intent = %q", period.Intent)
	}
}

func TestParseMoneyReportRequestThaiMonthDateIsDaily(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	now := time.Date(2026, time.August, 24, 12, 0, 0, 0, loc)

	_, period, ok := ParseReportRequest("สรุปรายจ่ายวันที่ 10 สิงหาคม", now, loc)
	if !ok {
		t.Fatal("expected report request")
	}
	wantStart := time.Date(2026, time.August, 10, 0, 0, 0, 0, loc)
	wantEnd := time.Date(2026, time.August, 11, 0, 0, 0, 0, loc)
	if !period.Start.Equal(wantStart) || !period.End.Equal(wantEnd) {
		t.Fatalf("period = %s..%s, want %s..%s", period.Start, period.End, wantStart, wantEnd)
	}
	if period.Intent != "expense_report_daily" {
		t.Fatalf("intent = %q", period.Intent)
	}
}

func TestExtractAmountSkipsLeadingDate(t *testing.T) {
	amount, ok := ExtractAmount("วันที่ 1 ได้เงิน 1000")
	if !ok {
		t.Fatal("expected amount")
	}
	if amount != 1000 {
		t.Fatalf("amount = %v, want 1000", amount)
	}
}

func TestParseMoneyEntryTimeUsesExplicitTime(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	now := time.Date(2026, time.August, 24, 12, 34, 56, 0, loc)

	got := ParseEntryTime("ซื้อข้าวขาหมู 60 บาทตอน 9 โมงเช้า", now, loc)
	want := time.Date(2026, time.August, 24, 9, 0, 0, 0, loc)
	if !got.Equal(want) {
		t.Fatalf("time = %s, want %s", got, want)
	}
}

func TestParseMoneyEntryTimeUsesDateAndExplicitTime(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	now := time.Date(2026, time.August, 24, 12, 34, 56, 0, loc)

	got := ParseEntryTime("วันที่ 10 สิงหาคม ได้เงิน 1000 ตอน 14:30", now, loc)
	want := time.Date(2026, time.August, 10, 14, 30, 0, 0, loc)
	if !got.Equal(want) {
		t.Fatalf("time = %s, want %s", got, want)
	}
}

func TestMergeParsedMoneyTimeAddsExplicitTimeToParsedDate(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	parsedAt := time.Date(2026, time.August, 10, 0, 0, 0, 0, loc)

	got := MergeParsedTime(parsedAt, "ได้เงิน 1000 วันที่ 10 สิงหาคม ตอน 9 โมงเช้า", loc)
	want := time.Date(2026, time.August, 10, 9, 0, 0, 0, loc)
	if !got.Equal(want) {
		t.Fatalf("time = %s, want %s", got, want)
	}
}

func TestParseMoneyEntryTimeUsesDetailedThaiTime(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	now := time.Date(2026, time.August, 24, 12, 34, 56, 0, loc)

	got := ParseEntryTime("ซื้อข้าวขาหมู 60 บาทตอน 4 ทุ่ม 23 นาที", now, loc)
	want := time.Date(2026, time.August, 24, 22, 23, 0, 0, loc)
	if !got.Equal(want) {
		t.Fatalf("time = %s, want %s", got, want)
	}
}

func TestCleanupExpenseDescriptionRemovesAmountAndTimePhrase(t *testing.T) {
	got := CleanupExpenseDescription("วันนี้ซื้อข้าวขาหมู 60 บาท ตอน 10 โมงครึ่ง")
	want := "ซื้อข้าวขาหมู"
	if got != want {
		t.Fatalf("description = %q, want %q", got, want)
	}
}

func TestFormatMoneyCreateReplyUsesNormalizedTime(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	occurredAt := time.Date(2026, time.August, 24, 10, 30, 0, 0, loc)

	got := FormatCreateReply("ซื้อข้าวขาหมู", 60, occurredAt, loc, true)
	want := "ลงบัญชี ซื้อข้าวขาหมู 60.00 บาท ตอน 10:30 น. ให้เรียบร้อยค่ะ"
	if got != want {
		t.Fatalf("reply = %q, want %q", got, want)
	}
}

func TestFormatCashflowReportReplyIncludesItems(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	expenses := []*entities.Expense{{
		Description: "ข้าวขาหมู",
		Amount:      60,
		SpentAt:     time.Date(2026, time.August, 24, 10, 30, 0, 0, loc),
	}}

	got := FormatCashflowReportReply("วันที่ 24 Aug 2026", nil, expenses, 0, 60, loc)
	want := "เลขาสรุปการเงินวันที่ 24 Aug 2026 ให้แล้วค่ะ\n- รายรับ 0.00 บาท\n- รายจ่าย 60.00 บาท\n  - ข้าวขาหมู 60.00 บาท\n- สุทธิ -60.00 บาทค่ะ"
	if got != want {
		t.Fatalf("reply = %q, want %q", got, want)
	}
}
