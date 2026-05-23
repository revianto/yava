package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/controllers"
)

func ExampleRoutes(r fiber.Router) {
	example := r.Group("example")
	example.Get("/", controllers.ExampleList)
	example.Get("/:id", controllers.ExampleShow)
	example.Post("/", controllers.ExampleCreate)
	example.Put("/:id", controllers.ExampleUpdate)
	example.Delete("/:id", controllers.ExampleDelete)
}
