package income

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"minyjae/go-starter/internal/core/domain/entities"
	servicePort "minyjae/go-starter/internal/core/domain/ports/services"
	"minyjae/go-starter/utils"
)

// incomeRequest คือ payload สำหรับสร้างหรือแก้ไขรายรับ
// input: JSON body จาก HTTP request
// output: struct ที่ controller แปลงต่อเป็น entities.Income
type incomeRequest struct {
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	ReceivedAt  time.Time `json:"received_at"`
}

// Create รับ HTTP request เพื่อสร้างรายรับ
// input: Fiber context ที่มี userID ใน locals และ JSON body incomeRequest
// output: HTTP created response พร้อม income ที่สร้าง หรือ error response
func (h *incomeController) Create(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}

	var req incomeRequest
	if err := c.BodyParser(&req); err != nil {
		return h.Response.BadRequest(c, "Invalid income payload", "INVALID_INCOME_PAYLOAD")
	}
	if req.Amount <= 0 {
		return h.Response.BadRequest(c, "Amount must be greater than 0", "INVALID_AMOUNT")
	}

	income, err := h.incomeService.Create(userID, &entities.Income{
		Amount:      req.Amount,
		Currency:    req.Currency,
		Category:    req.Category,
		Description: req.Description,
		ReceivedAt:  req.ReceivedAt,
	})
	if err != nil {
		return h.Response.InternalServerError(c, "Failed to create income", err.Error(), "CREATE_INCOME_FAILED")
	}
	return h.Response.Created(c, "Income created", income)
}

// Update รับ HTTP request เพื่อแก้ไขรายรับ
// input: Fiber context ที่มี userID ใน locals, path id, และ JSON body incomeRequest
// output: HTTP updated response พร้อม income ที่แก้ หรือ error response
func (h *incomeController) Update(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}
	id, ok := utils.UintParam(c, "id")
	if !ok {
		return h.Response.BadRequest(c, "Invalid income id", "INVALID_INCOME_ID")
	}

	var req incomeRequest
	if err := c.BodyParser(&req); err != nil {
		return h.Response.BadRequest(c, "Invalid income payload", "INVALID_INCOME_PAYLOAD")
	}
	if req.Amount <= 0 {
		return h.Response.BadRequest(c, "Amount must be greater than 0", "INVALID_AMOUNT")
	}

	income, err := h.incomeService.Update(userID, id, &entities.Income{
		Amount:      req.Amount,
		Currency:    req.Currency,
		Category:    req.Category,
		Description: req.Description,
		ReceivedAt:  req.ReceivedAt,
	})
	if err != nil {
		if errors.Is(err, servicePort.ErrForbidden) {
			return h.Response.Forbidden(c, "Forbidden")
		}
		return h.Response.InternalServerError(c, "Failed to update income", err.Error(), "UPDATE_INCOME_FAILED")
	}
	return h.Response.Updated(c, "Income updated", income)
}

// Delete รับ HTTP request เพื่อลบรายรับ
// input: Fiber context ที่มี userID ใน locals และ path id
// output: HTTP deleted response หรือ error response
func (h *incomeController) Delete(c *fiber.Ctx) error {
	userID, ok := utils.UserIDFromLocals(c)
	if !ok {
		return h.Response.Unauthorized(c, "Unauthorized", "UNAUTHORIZED")
	}
	id, ok := utils.UintParam(c, "id")
	if !ok {
		return h.Response.BadRequest(c, "Invalid income id", "INVALID_INCOME_ID")
	}

	if err := h.incomeService.Delete(userID, id); err != nil {
		if errors.Is(err, servicePort.ErrForbidden) {
			return h.Response.Forbidden(c, "Forbidden")
		}
		return h.Response.InternalServerError(c, "Failed to delete income", err.Error(), "DELETE_INCOME_FAILED")
	}
	return h.Response.Deleted(c, "Income deleted")
}
