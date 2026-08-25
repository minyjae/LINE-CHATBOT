package expense

import (
	"time"

	"github.com/gofiber/fiber/v2"
	financeHandler "minyjae/go-starter/internal/adapters/http/handlers/finance"
	"minyjae/go-starter/utils"
)

// List ดึงรายจ่ายของ user แบบแบ่งหน้า
// input: Fiber context ที่มี userID ใน locals และ query limit/offset
// output: HTTP response รายการรายจ่าย หรือ error response
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

// Summary ดึงสรุปรายจ่ายตาม period/date/month จาก query string
// input: Fiber context ที่มี userID ใน locals และ query period/date/month
// output: HTTP response ExpenseSummary หรือ error response
func (h *expenseController) Summary(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}

	period, err := financeHandler.ParseSummaryPeriod(c.Query("period"), c.Query("date"), c.Query("month"), time.Now(), time.Local)
	if err != nil {
		return h.Response.BadRequest(c, "Invalid period. Use period=day|week|month with date=YYYY-MM-DD or month=YYYY-MM", "INVALID_PERIOD")
	}

	summary, err := h.expenseService.SummaryByPeriod(userID, period.Start, period.End)
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to summarize expenses", err.Error(), "EXPENSE_SUMMARY_FAILED")
	}
	return h.Response.Item(c, "Expense summary fetched", summary)
}
