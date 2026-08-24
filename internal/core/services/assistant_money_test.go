package services

import (
	"testing"
	"time"
)

func TestParseMoneyReportRequestSpecificWeek(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	now := time.Date(2026, time.August, 24, 12, 0, 0, 0, loc)

	target, period, ok := parseMoneyReportRequest("สรุปรายจ่ายสัปดาห์ของวันที่ 2026-08-10", now, loc)
	if !ok {
		t.Fatal("expected report request")
	}
	if target != moneyReportTargetExpense {
		t.Fatalf("target = %q, want %q", target, moneyReportTargetExpense)
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

	target, period, ok := parseMoneyReportRequest("สรุปรายรับเดือนสิงหาคม 2569", now, loc)
	if !ok {
		t.Fatal("expected report request")
	}
	if target != moneyReportTargetIncome {
		t.Fatalf("target = %q, want %q", target, moneyReportTargetIncome)
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

	target, period, ok := parseMoneyReportRequest("สรุปรายรับรายจ่ายวันนี้", now, loc)
	if !ok {
		t.Fatal("expected report request")
	}
	if target != moneyReportTargetCashflow {
		t.Fatalf("target = %q, want %q", target, moneyReportTargetCashflow)
	}
	if period.Intent != "cashflow_report_daily" {
		t.Fatalf("intent = %q", period.Intent)
	}
}

func TestParseMoneyReportRequestThaiMonthDateIsDaily(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	now := time.Date(2026, time.August, 24, 12, 0, 0, 0, loc)

	_, period, ok := parseMoneyReportRequest("สรุปรายจ่ายวันที่ 10 สิงหาคม", now, loc)
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
	amount, ok := extractAmount("วันที่ 1 ได้เงิน 1000")
	if !ok {
		t.Fatal("expected amount")
	}
	if amount != 1000 {
		t.Fatalf("amount = %v, want 1000", amount)
	}
}
