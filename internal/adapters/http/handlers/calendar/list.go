package calendar

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"minyjae/go-starter/utils"
)

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

func parseDate(value string) (time.Time, error) {
	if value == "" {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local), nil
	}
	return time.ParseInLocation("2006-01-02", value, time.Local)
}
