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

	// Recipes — list & show are public, write ops require auth
	v1.Get("/recipes", controllers.RecipeList)
	v1.Get("/recipes/:id", controllers.RecipeShow)
	v1.Post("/recipes", middlewares.YvAuth(), controllers.RecipeCreate)
	v1.Put("/recipes/:id", middlewares.YvAuth(), controllers.RecipeUpdate)
	v1.Patch("/recipes/:id/archive", middlewares.YvAuth(), controllers.RecipeArchive)
	v1.Patch("/recipes/:id/restore", middlewares.YvAuth(), controllers.RecipeRestore)
	v1.Post("/recipes/:id/duplicate", middlewares.YvAuth(), controllers.RecipeDuplicate)

	// Groups — all require auth
	groups := v1.Group("/groups", middlewares.YvAuth())
	groups.Get("/", controllers.GroupList)
	groups.Post("/", controllers.GroupCreate)
	groups.Get("/:id", controllers.GroupShow)
	groups.Put("/:id", controllers.GroupUpdate)
	groups.Delete("/:id", controllers.GroupDelete)
	groups.Get("/:id/members", controllers.GroupMemberList)
	groups.Post("/:id/members", controllers.GroupJoin)
	groups.Delete("/:id/members/:uid", controllers.GroupMemberRemove)
	groups.Patch("/:id/members/:uid/role", controllers.GroupMemberSetRole)
	groups.Get("/:id/recipes", controllers.GroupRecipeList)
	groups.Get("/:id/recipes/pending", controllers.GroupRecipePending)
	groups.Post("/:id/recipes", controllers.GroupRecipeSubmit)
	groups.Patch("/:id/recipes/:rid/approve", controllers.GroupRecipeApprove)
	groups.Patch("/:id/recipes/:rid/reject", controllers.GroupRecipeReject)
	groups.Delete("/:id/recipes/:rid", controllers.GroupRecipeRemove)
}
