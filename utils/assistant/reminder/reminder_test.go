package reminder

import "testing"

func TestCleanupReminderTitleRemovesDateAndTime(t *testing.T) {
	got := CleanupTitle("เตือนดูจองรถ ตอน9โมง40 วันที่ 7 เดือน 9 ปี 2026")
	want := "ดูจองรถ"
	if got != want {
		t.Fatalf("title = %q, want %q", got, want)
	}
}
