package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/controllers"
	"github.com/revianto/yava/api/app/middlewares"
)

func YavaRoutes(app *fiber.App) {
	v1 := app.Group("/v1")

	// Auth — public
	auth := v1.Group("/auth")
	auth.Get("/google", controllers.AuthGoogle)
	auth.Get("/google/callback", controllers.AuthGoogleCallback)
	auth.Post("/logout", controllers.AuthLogout)
	auth.Get("/me", middlewares.YvAuth(), controllers.AuthMe)

	// Recipe types — public
	v1.Get("/types", controllers.RecipeTypeList)
	v1.Get("/types/:id/subtypes", controllers.RecipeSubtypeList)

	// Recipes — list & show are public, create & update require auth
	v1.Get("/recipes", controllers.RecipeList)
	v1.Get("/recipes/:id", controllers.RecipeShow)
	v1.Post("/recipes", middlewares.YvAuth(), controllers.RecipeCreate)
	v1.Put("/recipes/:id", middlewares.YvAuth(), controllers.RecipeUpdate)
}
