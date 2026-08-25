package todo

import (
	"minyjae/go-starter/utils"

	"github.com/gofiber/fiber/v2"
)

// List ดึง todo ของ user แบบแบ่งหน้า
// input: Fiber context ที่มี userID ใน locals และ query limit/offset
// output: HTTP response รายการ todo หรือ error response
func (h *todoController) List(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}

	todos, err := h.todoService.List(userID, c.QueryInt("limit", 50), c.QueryInt("offset", 0))
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to list todos", err.Error(), "LIST_TODOS_FAILED")
	}
	return h.Response.Item(c, "Todos fetched", todos)
}

// ListPending ดึง todo ของ user เฉพาะรายการที่ยัง pending
// input: Fiber context ที่มี userID ใน locals
// output: HTTP response รายการ pending todo หรือ error response
func (h *todoController) ListPending(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}

	todos, err := h.todoService.ListPending(userID)
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to list pending todos", err.Error(), "LIST_PENDING_TODOS_FAILED")
	}
	return h.Response.Item(c, "Pending todos fetched", todos)
}
