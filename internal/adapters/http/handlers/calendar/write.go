package calendar

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"minyjae/go-starter/internal/core/domain/entities"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	"minyjae/go-starter/utils"
)

// calendarEventRequest คือ payload สำหรับสร้างหรือแก้ไข calendar event
// input: JSON body จาก HTTP request
// output: struct ที่ controller แปลงต่อเป็น entities.CalendarEvent
type calendarEventRequest struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	StartAt     time.Time  `json:"start_at"`
	EndAt       *time.Time `json:"end_at"`
	Location    *string    `json:"location"`
	SyncStatus  string     `json:"sync_status"`
}

// Create รับ HTTP request เพื่อสร้าง calendar event
// input: Fiber context ที่มี userID ใน locals และ JSON body calendarEventRequest
// output: HTTP created response พร้อม event ที่สร้าง หรือ error response
func (h *calendarController) Create(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}

	var req calendarEventRequest
	if err := c.BodyParser(&req); err != nil {
		return h.Response.BadRequest(c, "Invalid calendar event payload", "INVALID_CALENDAR_EVENT_PAYLOAD")
	}
	if req.Title == "" || req.StartAt.IsZero() {
		return h.Response.BadRequest(c, "Title and start_at are required", "INVALID_CALENDAR_EVENT")
	}

	event, err := h.calendarService.Create(userID, &entities.CalendarEvent{
		Title:       req.Title,
		Description: req.Description,
		StartAt:     req.StartAt,
		EndAt:       req.EndAt,
		Location:    req.Location,
		SyncStatus:  req.SyncStatus,
	})
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to create calendar event", err.Error(), "CREATE_CALENDAR_EVENT_FAILED")
	}
	return h.Response.Created(c, "Calendar event created", event)
}

// Update รับ HTTP request เพื่อแก้ไข calendar event
// input: Fiber context ที่มี userID ใน locals, path id, และ JSON body calendarEventRequest
// output: HTTP updated response พร้อม event ที่แก้ หรือ error response
func (h *calendarController) Update(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}
	id, ok := utils.UintParam(c, "id")
	if !ok {
		return h.Response.BadRequest(c, "Invalid calendar event id", "INVALID_CALENDAR_EVENT_ID")
	}

	var req calendarEventRequest
	if err := c.BodyParser(&req); err != nil {
		return h.Response.BadRequest(c, "Invalid calendar event payload", "INVALID_CALENDAR_EVENT_PAYLOAD")
	}
	if req.Title == "" || req.StartAt.IsZero() {
		return h.Response.BadRequest(c, "Title and start_at are required", "INVALID_CALENDAR_EVENT")
	}

	event, err := h.calendarService.Update(userID, id, &entities.CalendarEvent{
		Title:       req.Title,
		Description: req.Description,
		StartAt:     req.StartAt,
		EndAt:       req.EndAt,
		Location:    req.Location,
		SyncStatus:  req.SyncStatus,
	})
	if err != nil {
		if errors.Is(err, servicePort.ErrForbidden) {
			return h.Response.Forbidden(c, "Forbidden")
		}
		return h.Response.InternalServerError(c, "Failed to update calendar event", err.Error(), "UPDATE_CALENDAR_EVENT_FAILED")
	}
	return h.Response.Updated(c, "Calendar event updated", event)
}

// Delete รับ HTTP request เพื่อลบ calendar event
// input: Fiber context ที่มี userID ใน locals และ path id
// output: HTTP deleted response หรือ error response
func (h *calendarController) Delete(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}
	id, ok := utils.UintParam(c, "id")
	if !ok {
		return h.Response.BadRequest(c, "Invalid calendar event id", "INVALID_CALENDAR_EVENT_ID")
	}

	if err := h.calendarService.Delete(userID, id); err != nil {
		if errors.Is(err, servicePort.ErrForbidden) {
			return h.Response.Forbidden(c, "Forbidden")
		}
		return h.Response.InternalServerError(c, "Failed to delete calendar event", err.Error(), "DELETE_CALENDAR_EVENT_FAILED")
	}
	return h.Response.Deleted(c, "Calendar event deleted")
}
