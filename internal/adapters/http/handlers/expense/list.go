package expense

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"minyjae/go-starter/utils"
)

func (h *expenseController) List(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}

	expenses, err := h.expenseService.List(userID, c.QueryInt("limit", 50), c.QueryInt("offset", 0))
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to list expenses", err.Error(), "LIST_EXPENSES_FAILED")
	}
	return h.Response.Item(c, "Expenses fetched", expenses)
}

func (h *expenseController) Summary(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}

	year, month, err := parseMonth(c.Query("month"))
	if err != nil {
		return h.Response.BadRequest(c, "Invalid month format. Use YYYY-MM", "INVALID_MONTH")
	}

	summary, err := h.expenseService.SummaryByMonth(userID, year, month, time.Local)
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to summarize expenses", err.Error(), "EXPENSE_SUMMARY_FAILED")
	}
	return h.Response.Item(c, "Expense summary fetched", summary)
}

func parseMonth(value string) (int, time.Month, error) {
	if strings.TrimSpace(value) == "" {
		now := time.Now()
		return now.Year(), now.Month(), nil
	}

	parts := strings.Split(value, "-")
	if len(parts) != 2 {
		return 0, 0, fiber.ErrBadRequest
	}
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	monthNumber, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}
	if monthNumber < 1 || monthNumber > 12 {
		return 0, 0, fiber.ErrBadRequest
	}
	return year, time.Month(monthNumber), nil
}
