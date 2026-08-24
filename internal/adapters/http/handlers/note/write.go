package note

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"minyjae/go-starter/internal/core/domain/entities"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	"minyjae/go-starter/utils"
)

type noteRequest struct {
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (h *noteController) Create(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}

	var req noteRequest
	if err := c.BodyParser(&req); err != nil {
		return h.Response.BadRequest(c, "Invalid note payload", "INVALID_NOTE_PAYLOAD")
	}
	if req.Content == "" {
		return h.Response.BadRequest(c, "Content is required", "CONTENT_REQUIRED")
	}

	note, err := h.noteService.Create(userID, &entities.Note{
		Content: req.Content,
		Tags:    req.Tags,
	})
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to create note", err.Error(), "CREATE_NOTE_FAILED")
	}
	return h.Response.Created(c, "Note created", note)
}

func (h *noteController) Update(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}
	id, ok := utils.UintParam(c, "id")
	if !ok {
		return h.Response.BadRequest(c, "Invalid note id", "INVALID_NOTE_ID")
	}

	var req noteRequest
	if err := c.BodyParser(&req); err != nil {
		return h.Response.BadRequest(c, "Invalid note payload", "INVALID_NOTE_PAYLOAD")
	}
	if req.Content == "" {
		return h.Response.BadRequest(c, "Content is required", "CONTENT_REQUIRED")
	}

	note, err := h.noteService.Update(userID, id, &entities.Note{
		Content: req.Content,
		Tags:    req.Tags,
	})
	if err != nil {
		if errors.Is(err, servicePort.ErrForbidden) {
			return h.Response.Forbidden(c, "Forbidden")
		}
		return h.Response.InternalServerError(c, "Failed to update note", err.Error(), "UPDATE_NOTE_FAILED")
	}
	return h.Response.Updated(c, "Note updated", note)
}

func (h *noteController) Delete(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}
	id, ok := utils.UintParam(c, "id")
	if !ok {
		return h.Response.BadRequest(c, "Invalid note id", "INVALID_NOTE_ID")
	}

	if err := h.noteService.Delete(userID, id); err != nil {
		if errors.Is(err, servicePort.ErrForbidden) {
			return h.Response.Forbidden(c, "Forbidden")
		}
		return h.Response.InternalServerError(c, "Failed to delete note", err.Error(), "DELETE_NOTE_FAILED")
	}
	return h.Response.Deleted(c, "Note deleted")
}
