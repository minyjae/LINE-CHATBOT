package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func UserIDFromLocals(c *fiber.Ctx) (uint, bool) {
	value := c.Locals("user_id")
	switch v := value.(type) {
	case uint:
		return v, true
	case int:
		if v <= 0 {
			return 0, false
		}
		return uint(v), true
	case float64:
		if v <= 0 {
			return 0, false
		}
		return uint(v), true
	case string:
		parsed, err := strconv.ParseUint(v, 10, 64)
		if err != nil || parsed == 0 {
			return 0, false
		}
		return uint(parsed), true
	default:
		return 0, false
	}
}
