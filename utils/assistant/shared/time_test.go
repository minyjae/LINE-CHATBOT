package shared

import (
	"testing"
	"time"
)

func TestParseNaturalTimeWithAfternoonMinute(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	now := time.Date(2026, time.August, 24, 12, 0, 0, 0, loc)

	got, ok := ParseNaturalTime("นัดดูรถ ตอน บ่าย2 30 นาที", now, loc, 9, 0)
	if !ok {
		t.Fatal("expected time")
	}
	want := time.Date(2026, time.August, 24, 14, 30, 0, 0, loc)
	if !got.Equal(want) {
		t.Fatalf("time = %s, want %s", got, want)
	}
}

func TestParseNaturalTimeWithThaiNoon(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	now := time.Date(2026, time.August, 24, 12, 0, 0, 0, loc)

	tests := []struct {
		text       string
		wantHour   int
		wantMinute int
	}{
		{text: "นัดดูรถที่ โตโยต้านครพิงค์ ตอนเที่ยง 30 นาที", wantHour: 12, wantMinute: 30},
		{text: "ลงตาราง ดูรถที่ โตโยต้า ตอนเที่ยงครึ่ง", wantHour: 12, wantMinute: 30},
		{text: "นัดดูรถที่ โจเส่ เที่ยงครึ่ง", wantHour: 12, wantMinute: 30},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			got, ok := ParseNaturalTime(tt.text, now, loc, 9, 0)
			if !ok {
				t.Fatal("expected time")
			}
			want := time.Date(2026, time.August, 24, tt.wantHour, tt.wantMinute, 0, 0, loc)
			if !got.Equal(want) {
				t.Fatalf("time = %s, want %s", got, want)
			}
		})
	}
}

func TestParseNaturalTimeWithCalendarDate(t *testing.T) {
	loc := time.FixedZone("Asia/Bangkok", 7*60*60)
	now := time.Date(2026, time.August, 24, 12, 0, 0, 0, loc)

	tests := []struct {
		text string
		want time.Time
	}{
		{
			text: "นัดประชุม 2026-08-30 10 โมง",
			want: time.Date(2026, time.August, 30, 10, 0, 0, 0, loc),
		},
		{
			text: "นัดประชุม 30/08/2026 10 โมง",
			want: time.Date(2026, time.August, 30, 10, 0, 0, 0, loc),
		},
		{
			text: "นัดประชุม 30 สิงหาคม 2026 10 โมง",
			want: time.Date(2026, time.August, 30, 10, 0, 0, 0, loc),
		},
		{
			text: "นัดประชุมวันที่ 30 10 โมง",
			want: time.Date(2026, time.August, 30, 10, 0, 0, 0, loc),
		},
		{
			text: "นัดดูจองรถ วันที่ 7 เดือน 9 ปี 2026 ตอน9โมง40",
			want: time.Date(2026, time.September, 7, 9, 40, 0, 0, loc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			got, ok := ParseNaturalTime(tt.text, now, loc, 9, 0)
			if !ok {
				t.Fatal("expected time")
			}
			if !got.Equal(tt.want) {
				t.Fatalf("time = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestExtractHourMinuteDetailedThaiTime(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		wantHour   int
		wantMinute int
	}{
		{name: "ten half", text: "ซื้อข้าวตอน 10 โมงครึ่ง", wantHour: 10, wantMinute: 30},
		{name: "dot clock", text: "ซื้อข้าวตอน 10.30น.", wantHour: 10, wantMinute: 30},
		{name: "tuum minutes", text: "ได้เงิน 4 ทุ่ม 23 นาที", wantHour: 22, wantMinute: 23},
		{name: "afternoon half", text: "บ่าย 2 ครึ่ง ซื้อกาแฟ", wantHour: 14, wantMinute: 30},
		{name: "evening", text: "กินข้าว 6 โมงเย็น", wantHour: 18, wantMinute: 0},
		{name: "colon clock", text: "ซื้อกาแฟ 14:30", wantHour: 14, wantMinute: 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hour, minute, ok := ExtractHourMinute(tt.text)
			if !ok {
				t.Fatal("expected time")
			}
			if hour != tt.wantHour || minute != tt.wantMinute {
				t.Fatalf("time = %02d:%02d, want %02d:%02d", hour, minute, tt.wantHour, tt.wantMinute)
			}
		})
	}
}
