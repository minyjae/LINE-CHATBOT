package reminder

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"minyjae/go-starter/internal/core/domain/entities"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	"minyjae/go-starter/utils"
)

type reminderRequest struct {
	Title    string     `json:"title"`
	RemindAt time.Time  `json:"remind_at"`
	Status   string     `json:"status"`
	SentAt   *time.Time `json:"sent_at"`
}

func (h *reminderController) Create(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}

	var req reminderRequest
	if err := c.BodyParser(&req); err != nil {
		return h.Response.BadRequest(c, "Invalid reminder payload", "INVALID_REMINDER_PAYLOAD")
	}
	if req.Title == "" || req.RemindAt.IsZero() {
		return h.Response.BadRequest(c, "Title and remind_at are required", "INVALID_REMINDER")
	}

	reminder, err := h.reminderService.Create(userID, &entities.Reminder{
		Title:    req.Title,
		RemindAt: req.RemindAt,
		Status:   req.Status,
		SentAt:   req.SentAt,
	})
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to create reminder", err.Error(), "CREATE_REMINDER_FAILED")
	}
	return h.Response.Created(c, "Reminder created", reminder)
}

func (h *reminderController) Update(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}
	id, ok := utils.UintParam(c, "id")
	if !ok {
		return h.Response.BadRequest(c, "Invalid reminder id", "INVALID_REMINDER_ID")
	}

	var req reminderRequest
	if err := c.BodyParser(&req); err != nil {
		return h.Response.BadRequest(c, "Invalid reminder payload", "INVALID_REMINDER_PAYLOAD")
	}
	if req.Title == "" || req.RemindAt.IsZero() {
		return h.Response.BadRequest(c, "Title and remind_at are required", "INVALID_REMINDER")
	}

	reminder, err := h.reminderService.Update(userID, id, &entities.Reminder{
		Title:    req.Title,
		RemindAt: req.RemindAt,
		Status:   req.Status,
		SentAt:   req.SentAt,
	})
	if err != nil {
		if errors.Is(err, servicePort.ErrForbidden) {
			return h.Response.Forbidden(c, "Forbidden")
		}
		return h.Response.InternalServerError(c, "Failed to update reminder", err.Error(), "UPDATE_REMINDER_FAILED")
	}
	return h.Response.Updated(c, "Reminder updated", reminder)
}

func (h *reminderController) Delete(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}
	id, ok := utils.UintParam(c, "id")
	if !ok {
		return h.Response.BadRequest(c, "Invalid reminder id", "INVALID_REMINDER_ID")
	}

	if err := h.reminderService.Delete(userID, id); err != nil {
		if errors.Is(err, servicePort.ErrForbidden) {
			return h.Response.Forbidden(c, "Forbidden")
		}
		return h.Response.InternalServerError(c, "Failed to delete reminder", err.Error(), "DELETE_REMINDER_FAILED")
	}
	return h.Response.Deleted(c, "Reminder deleted")
}
