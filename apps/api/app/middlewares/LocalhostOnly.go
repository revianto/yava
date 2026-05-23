package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// LocalhostOnly middleware restricts access to localhost only
func LocalhostOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := c.IP()

		// Check if request is from localhost
		allowedIPs := []string{"127.0.0.1", "::1", "localhost"}

		isAllowed := false
		for _, allowedIP := range allowedIPs {
			if strings.HasPrefix(ip, allowedIP) {
				isAllowed = true
				break
			}
		}

		// Prevent bypass via Reverse Proxy (Nginx, etc.)
		// If these headers exist, the request originated from outside
		if c.Get("X-Forwarded-For") != "" || c.Get("X-Real-IP") != "" {
			isAllowed = false
		}

		if !isAllowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "Access denied",
				"message": "This endpoint is only accessible from localhost",
			})
		}

		return c.Next()
	}
}
