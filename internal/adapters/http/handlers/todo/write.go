package todo

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"minyjae/go-starter/internal/core/domain/entities"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	"minyjae/go-starter/utils"
)

type todoRequest struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	DueAt       *time.Time `json:"due_at"`
	Priority    string     `json:"priority"`
	CompletedAt *time.Time `json:"completed_at"`
}

func (h *todoController) Create(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}

	var req todoRequest
	if err := c.BodyParser(&req); err != nil {
		return h.Response.BadRequest(c, "Invalid todo payload", "INVALID_TODO_PAYLOAD")
	}
	if req.Title == "" {
		return h.Response.BadRequest(c, "Title is required", "TITLE_REQUIRED")
	}

	todo, err := h.todoService.Create(userID, &entities.Todo{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		DueAt:       req.DueAt,
		Priority:    req.Priority,
		CompletedAt: req.CompletedAt,
	})
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to create todo", err.Error(), "CREATE_TODO_FAILED")
	}
	return h.Response.Created(c, "Todo created", todo)
}

func (h *todoController) Update(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}
	id, ok := utils.UintParam(c, "id")
	if !ok {
		return h.Response.BadRequest(c, "Invalid todo id", "INVALID_TODO_ID")
	}

	var req todoRequest
	if err := c.BodyParser(&req); err != nil {
		return h.Response.BadRequest(c, "Invalid todo payload", "INVALID_TODO_PAYLOAD")
	}
	if req.Title == "" {
		return h.Response.BadRequest(c, "Title is required", "TITLE_REQUIRED")
	}

	todo, err := h.todoService.Update(userID, id, &entities.Todo{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		DueAt:       req.DueAt,
		Priority:    req.Priority,
		CompletedAt: req.CompletedAt,
	})
	if err != nil {
		if errors.Is(err, servicePort.ErrForbidden) {
			return h.Response.Forbidden(c, "Forbidden")
		}
		return h.Response.InternalServerError(c, "Failed to update todo", err.Error(), "UPDATE_TODO_FAILED")
	}
	return h.Response.Updated(c, "Todo updated", todo)
}

func (h *todoController) Delete(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}
	id, ok := utils.UintParam(c, "id")
	if !ok {
		return h.Response.BadRequest(c, "Invalid todo id", "INVALID_TODO_ID")
	}

	if err := h.todoService.Delete(userID, id); err != nil {
		if errors.Is(err, servicePort.ErrForbidden) {
			return h.Response.Forbidden(c, "Forbidden")
		}
		return h.Response.InternalServerError(c, "Failed to delete todo", err.Error(), "DELETE_TODO_FAILED")
	}
	return h.Response.Deleted(c, "Todo deleted")
}
