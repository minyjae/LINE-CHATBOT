package services

import (
	"testing"
	"time"

	"minyjae/go-starter/internal/core/domain/entities"
)

func TestParseCalendarCancelRequest(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	now := time.Date(2026, 8, 21, 15, 30, 0, 0, loc)

	request := parseCalendarCancelRequest("ยกเลิกนัดพรุ่งนี้ 10 โมง ประชุมกับทีม", now, loc)

	if request.Title != "ประชุมกับทีม" {
		t.Fatalf("expected cleaned title, got %q", request.Title)
	}
	if !request.HasDate || !request.HasTime {
		t.Fatalf("expected date and time flags, got has_date=%v has_time=%v", request.HasDate, request.HasTime)
	}
	if request.Start.Day() != 22 || request.Hour != 10 || request.Minute != 0 {
		t.Fatalf("unexpected request window/time: start=%s hour=%d minute=%d", request.Start, request.Hour, request.Minute)
	}
}

func TestFilterCalendarCancelCandidates(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	event := &entities.CalendarEvent{
		ID:      1,
		Title:   "ประชุมกับทีม",
		StartAt: time.Date(2026, 8, 22, 10, 0, 0, 0, loc),
	}
	other := &entities.CalendarEvent{
		ID:      2,
		Title:   "คุยกับลูกค้า",
		StartAt: time.Date(2026, 8, 22, 10, 0, 0, 0, loc),
	}

	matches := filterCalendarCancelCandidates([]*entities.CalendarEvent{event, other}, calendarCancelRequest{
		Title:   "ประชุม",
		HasTime: true,
		Hour:    10,
		Minute:  0,
	}, loc)

	if len(matches) != 1 || matches[0].ID != event.ID {
		t.Fatalf("expected only event %d, got %#v", event.ID, matches)
	}
}
