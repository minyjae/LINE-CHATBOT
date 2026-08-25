package reminder

import (
	"github.com/gofiber/fiber/v2"
	"minyjae/go-starter/utils"
)

// List ดึง reminder ของ user แบบแบ่งหน้า
// input: Fiber context ที่มี userID ใน locals และ query limit/offset
// output: HTTP response รายการ reminder หรือ error response
func (h *reminderController) List(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}

	reminders, err := h.reminderService.List(userID, c.QueryInt("limit", 50), c.QueryInt("offset", 0))
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to list reminders", err.Error(), "LIST_REMINDERS_FAILED")
	}
	return h.Response.Item(c, "Reminders fetched", reminders)
}
