package finance

import (
	"time"

	"minyjae/go-starter/utils"

	"github.com/gofiber/fiber/v2"
)

// Summary ดึงสรุปการเงินรวมรายรับ รายจ่าย และยอดสุทธิ
// input: Fiber context ที่มี userID ใน locals และ query period/date/month
// output: HTTP response ที่มี period, income, expense, net และ currency หรือ error response
func (h *financeController) Summary(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}

	period, err := ParseSummaryPeriod(c.Query("period"), c.Query("date"), c.Query("month"), time.Now(), time.Local)
	if err != nil {
		return h.Response.BadRequest(c, "Invalid period. Use period=day|week|month with date=YYYY-MM-DD or month=YYYY-MM", "INVALID_PERIOD")
	}

	incomeSummary, err := h.incomeService.SummaryByPeriod(userID, period.Start, period.End)
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to summarize incomes", err.Error(), "INCOME_SUMMARY_FAILED")
	}
	expenseSummary, err := h.expenseService.SummaryByPeriod(userID, period.Start, period.End)
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to summarize expenses", err.Error(), "EXPENSE_SUMMARY_FAILED")
	}

	return h.Response.Item(c, "Finance summary fetched", fiber.Map{
		"period":   period,
		"income":   incomeSummary.Total,
		"expense":  expenseSummary.Total,
		"net":      incomeSummary.Total - expenseSummary.Total,
		"currency": "THB",
	})
}
