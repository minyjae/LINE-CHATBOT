package expense

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"minyjae/go-starter/internal/core/domain/entities"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	"minyjae/go-starter/utils"
)

type expenseRequest struct {
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	SpentAt     time.Time `json:"spent_at"`
}

func (h *expenseController) Create(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}

	var req expenseRequest
	if err := c.BodyParser(&req); err != nil {
		return h.Response.BadRequest(c, "Invalid expense payload", "INVALID_EXPENSE_PAYLOAD")
	}
	if req.Amount <= 0 {
		return h.Response.BadRequest(c, "Amount must be greater than 0", "INVALID_AMOUNT")
	}

	expense, err := h.expenseService.Create(userID, &entities.Expense{
		Amount:      req.Amount,
		Currency:    req.Currency,
		Category:    req.Category,
		Description: req.Description,
		SpentAt:     req.SpentAt,
	})
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to create expense", err.Error(), "CREATE_EXPENSE_FAILED")
	}
	return h.Response.Created(c, "Expense created", expense)
}

func (h *expenseController) Update(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}
	id, ok := utils.UintParam(c, "id")
	if !ok {
		return h.Response.BadRequest(c, "Invalid expense id", "INVALID_EXPENSE_ID")
	}

	var req expenseRequest
	if err := c.BodyParser(&req); err != nil {
		return h.Response.BadRequest(c, "Invalid expense payload", "INVALID_EXPENSE_PAYLOAD")
	}
	if req.Amount <= 0 {
		return h.Response.BadRequest(c, "Amount must be greater than 0", "INVALID_AMOUNT")
	}

	expense, err := h.expenseService.Update(userID, id, &entities.Expense{
		Amount:      req.Amount,
		Currency:    req.Currency,
		Category:    req.Category,
		Description: req.Description,
		SpentAt:     req.SpentAt,
	})
	if err != nil {
		if errors.Is(err, servicePort.ErrForbidden) {
			return h.Response.Forbidden(c, "Forbidden")
		}
		return h.Response.InternalServerError(c, "Failed to update expense", err.Error(), "UPDATE_EXPENSE_FAILED")
	}
	return h.Response.Updated(c, "Expense updated", expense)
}

func (h *expenseController) Delete(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}
	id, ok := utils.UintParam(c, "id")
	if !ok {
		return h.Response.BadRequest(c, "Invalid expense id", "INVALID_EXPENSE_ID")
	}

	if err := h.expenseService.Delete(userID, id); err != nil {
		if errors.Is(err, servicePort.ErrForbidden) {
			return h.Response.Forbidden(c, "Forbidden")
		}
		return h.Response.InternalServerError(c, "Failed to delete expense", err.Error(), "DELETE_EXPENSE_FAILED")
	}
	return h.Response.Deleted(c, "Expense deleted")
}
