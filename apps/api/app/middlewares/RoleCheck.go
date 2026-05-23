package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/exceptions"
	"github.com/revianto/yava/api/helpers"
)

// RequireRole middleware checks if user has one of the allowed roles
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userData, ok := c.Locals("userData").(fiber.Map)
		if !ok || userData == nil {
			locale := getLocaleFromCtx(c)
			return exceptions.ResponseErrorException(c, exceptions.AuthException(c, 401, helpers.TransError(locale, "unauthorized", nil)))
		}

		userRole := helpers.Conv(userData["role_code"]).String()
		for _, role := range roles {
			if role == userRole {
				return c.Next()
			}
		}

		locale := getLocaleFromCtx(c)
		return exceptions.ResponseErrorException(c, exceptions.AuthException(c, 403, helpers.TransError(locale, "forbidden", nil)))
	}
}
