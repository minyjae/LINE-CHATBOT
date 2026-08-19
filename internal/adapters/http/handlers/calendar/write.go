package calendar

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"minyjae/go-starter/internal/core/domain/entities"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	"minyjae/go-starter/utils"
)

type calendarEventRequest struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	StartAt     time.Time  `json:"start_at"`
	EndAt       *time.Time `json:"end_at"`
	Location    *string    `json:"location"`
	SyncStatus  string     `json:"sync_status"`
}

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
