package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/revianto/yava/api/exceptions"
)

// AuthLimiter creates a rate limiter for auth endpoints (max 5 requests per minute)
func AuthLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return exceptions.ResponseErrorException(c, exceptions.ThrottleException(c, 429, "Tolong tunggu sebentar sebelum mencoba lagi."))
		},
	})
}
