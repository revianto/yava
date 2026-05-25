package controllers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Helper untuk mendapatkan DB transaction dari context
func getDB(c *fiber.Ctx) *gorm.DB {
	return c.Locals("DB").(*gorm.DB)
}

// Helper untuk parse body data
func getBodyData(c *fiber.Ctx) fiber.Map {
	data := fiber.Map{}
	c.BodyParser(&data)
	return data
}

// Helper untuk mendapatkan locale, fallback ke "id"
func getLocale(c *fiber.Ctx) string {
	if locale, ok := c.Locals("locale").(string); ok && locale != "" {
		return locale
	}
	return "id"
}

// Helper untuk mendapatkan userData dari context (di-set oleh CheckUserToken middleware)
func getUserData(c *fiber.Ctx) fiber.Map {
	if userData, ok := c.Locals("userData").(fiber.Map); ok {
		return userData
	}
	return nil
}
