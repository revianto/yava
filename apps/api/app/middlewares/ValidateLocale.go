package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

// ValidLocales defines allowed locales
var ValidLocales = map[string]bool{
	"en": true,
	"id": true,
}

// ValidateLocale middleware validates locale parameter, defaults to "id" if empty or invalid
func ValidateLocale(c *fiber.Ctx) error {
	locale := c.Params("locale")
	if !ValidLocales[locale] {
		locale = "id"
	}
	c.Locals("locale", locale)
	return c.Next()
}
