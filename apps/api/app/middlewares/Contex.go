package middlewares

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/models"
	"github.com/revianto/yava/api/database"
	"gorm.io/gorm"
)

func SetContexFiber(c *fiber.Ctx) error {
	ctx := context.WithValue(c.Context(), models.FiberCtxKey{}, c)
	c.Locals("DB", database.DB.Session(&gorm.Session{NewDB: true, Context: ctx}))
	return c.Next()
}
