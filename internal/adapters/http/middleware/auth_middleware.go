package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"minyjae/go-starter/utils"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		// fmt.Printf("In [middleware] : tokenStr = %v\n", tokenStr)
		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("display_name", claims.DisplayName)

		return c.Next()
	}
}
