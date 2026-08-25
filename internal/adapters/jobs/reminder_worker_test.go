package jobs

import (
	"testing"
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

func TestFormatReminderMessagePreAlert(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	sentAt := time.Date(2026, time.September, 7, 9, 30, 0, 0, loc)
	reminder := &entities.Reminder{
		Title:    "ดูจองรถ",
		RemindAt: time.Date(2026, time.September, 7, 9, 40, 0, 0, loc),
		Status:   "pending",
	}

	got := formatReminderMessage(reminder, sentAt)
	want := "อีก 10 นาทีจะถึงเวลา \"ดูจองรถ\" แล้วค่ะ"
	if got != want {
		t.Fatalf("message = %q, want %q", got, want)
	}
	if !shouldMarkReminderPreSent(reminder, sentAt) {
		t.Fatal("expected reminder to be marked pre_sent")
	}
}

func TestFormatReminderMessageDue(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	sentAt := time.Date(2026, time.September, 7, 9, 40, 0, 0, loc)
	reminder := &entities.Reminder{
		Title:    "ดูจองรถ",
		RemindAt: sentAt,
		Status:   "pre_sent",
	}

	got := formatReminderMessage(reminder, sentAt)
	want := "ถึงเวลาแล้วค่ะ เลขามาสะกิดเรื่อง \"ดูจองรถ\" ค่ะ"
	if got != want {
		t.Fatalf("message = %q, want %q", got, want)
	}
	if shouldMarkReminderPreSent(reminder, sentAt) {
		t.Fatal("did not expect due reminder to be marked pre_sent")
	}
}

