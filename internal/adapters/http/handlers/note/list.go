package note

import (
	"github.com/gofiber/fiber/v2"
	"minyjae/go-starter/utils"
)

// List ดึง note ของ user แบบแบ่งหน้า
// input: Fiber context ที่มี userID ใน locals และ query limit/offset
// output: HTTP response รายการ note หรือ error response
func (h *noteController) List(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}

	notes, err := h.noteService.List(userID, c.QueryInt("limit", 50), c.QueryInt("offset", 0))
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to list notes", err.Error(), "LIST_NOTES_FAILED")
	}
	return h.Response.Item(c, "Notes fetched", notes)
}

// Search ค้น note ของ user จาก query string
// input: Fiber context ที่มี userID ใน locals, query q และ limit
// output: HTTP response รายการ note ที่ตรงคำค้น หรือ error response
func (h *noteController) Search(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}

	notes, err := h.noteService.Search(userID, c.Query("q"), c.QueryInt("limit", 50))
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to search notes", err.Error(), "SEARCH_NOTES_FAILED")
	}
	return h.Response.Item(c, "Notes fetched", notes)
}
