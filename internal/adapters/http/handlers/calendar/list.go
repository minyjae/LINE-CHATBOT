package calendar

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"minyjae/go-starter/utils"
)

// List ดึง calendar event ทั้งหมดของ user แบบแบ่งหน้า
// input: Fiber context ที่มี userID ใน locals และ query limit/offset
// output: HTTP response รายการ calendar event หรือ error response
func (h *calendarController) List(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}

	events, err := h.calendarService.List(userID, c.QueryInt("limit", 50), c.QueryInt("offset", 0))
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to list calendar events", err.Error(), "LIST_CALENDAR_EVENTS_FAILED")
	}
	return h.Response.Item(c, "Calendar events fetched", events)
}

// ListByDate ดึง calendar event ของ user เฉพาะวันที่ระบุ
// input: Fiber context ที่มี userID ใน locals และ query date รูปแบบ YYYY-MM-DD
// output: HTTP response รายการ calendar event ของวันนั้น หรือ error response
func (h *calendarController) ListByDate(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}

	date, err := parseDate(c.Query("date"))
	if err != nil {
		return h.Response.BadRequest(c, "Invalid date format. Use YYYY-MM-DD", "INVALID_DATE")
	}

	events, err := h.calendarService.ListByDate(userID, date, time.Local)
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to list calendar events by date", err.Error(), "LIST_CALENDAR_EVENTS_BY_DATE_FAILED")
	}
	return h.Response.Item(c, "Calendar events fetched", events)
}

// parseDate แปลง query date เป็น time.Time สำหรับ list calendar รายวัน
// input: value วันที่รูปแบบ YYYY-MM-DD หรือค่าว่าง
// output: time.Time ของวันที่นั้น หรือวันนี้เมื่อ value ว่าง และ error ถ้า format ไม่ถูกต้อง
func parseDate(value string) (time.Time, error) {
	if value == "" {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local), nil
	}
	return time.ParseInLocation("2006-01-02", value, time.Local)
}
