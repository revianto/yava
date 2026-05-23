package middlewares

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/controllers"
	"github.com/revianto/yava/api/app/models"
	"gorm.io/gorm"
)

func SetContexFiber(c *fiber.Ctx) error {
	ctx := context.WithValue(c.Context(), models.FiberCtxKey{}, c)
	c.Locals("DB", controllers.DB.Session(&gorm.Session{NewDB: true, Context: ctx}))
	return c.Next()
}
