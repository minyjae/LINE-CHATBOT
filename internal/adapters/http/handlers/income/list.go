package income

import (
	"time"

	"github.com/gofiber/fiber/v2"
	financeHandler "minyjae/go-starter/internal/adapters/http/handlers/finance"
	"minyjae/go-starter/utils"
)

// List ดึงรายรับของ user แบบแบ่งหน้า
// input: Fiber context ที่มี userID ใน locals และ query limit/offset
// output: HTTP response รายการรายรับ หรือ error response
func (h *incomeController) List(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}

	incomes, err := h.incomeService.List(userID, c.QueryInt("limit", 50), c.QueryInt("offset", 0))
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to list incomes", err.Error(), "LIST_INCOMES_FAILED")
	}
	return h.Response.Item(c, "Incomes fetched", incomes)
}

// Summary ดึงสรุปรายรับตาม period/date/month จาก query string
// input: Fiber context ที่มี userID ใน locals และ query period/date/month
// output: HTTP response IncomeSummary หรือ error response
func (h *incomeController) Summary(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}

	period, err := financeHandler.ParseSummaryPeriod(c.Query("period"), c.Query("date"), c.Query("month"), time.Now(), time.Local)
	if err != nil {
		return h.Response.BadRequest(c, "Invalid period. Use period=day|week|month with date=YYYY-MM-DD or month=YYYY-MM", "INVALID_PERIOD")
	}

	summary, err := h.incomeService.SummaryByPeriod(userID, period.Start, period.End)
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to summarize incomes", err.Error(), "INCOME_SUMMARY_FAILED")
	}
	return h.Response.Item(c, "Income summary fetched", summary)
}
