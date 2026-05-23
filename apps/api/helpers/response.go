package helpers

import "github.com/gofiber/fiber/v2"

// ResponseDeleted returns a standard delete response
func ResponseDeleted(c *fiber.Ctx, id any) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data has been deleted successfully",
		"id":      id,
	})
}
