package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/helpers"
)

// YvAuth validates YAVA JWT from cookie or Authorization header.
func YvAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Cookies("yv_token")
		if tokenStr == "" {
			auth := c.Get("Authorization")
			tokenStr = strings.TrimPrefix(auth, "Bearer ")
		}
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.YvError("UNAUTHORIZED", "Token diperlukan"))
		}
		claims, err := helpers.YvParseToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.YvError("INVALID_TOKEN", "Token tidak valid atau sudah kedaluwarsa"))
		}
		c.Locals("yv_user_id", claims.UserID)
		c.Locals("yv_user_email", claims.Email)
		c.Locals("yv_user_name", claims.Name)
		c.Locals("yv_user_avatar", claims.AvatarURL)
		return c.Next()
	}
}

func YvUserID(c *fiber.Ctx) int64 {
	id, _ := c.Locals("yv_user_id").(int64)
	return id
}
