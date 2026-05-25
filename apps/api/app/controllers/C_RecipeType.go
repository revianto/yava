package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/resources"
	"github.com/revianto/yava/api/app/services"
	"github.com/revianto/yava/api/helpers"
)

// GET /v1/types
func RecipeTypeList(c *fiber.Ctx) error {
	types, err := services.GetRecipeTypes(getDB(c))
	if err != nil {
		se := err.(*services.ServiceError)
		return c.Status(se.Code).JSON(helpers.YvError(se.ErrCode, se.Message))
	}
	return c.JSON(helpers.YvSuccess(resources.RecipeTypeListResource(types)))
}

// GET /v1/types/:id/subtypes
func RecipeSubtypeList(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.YvError("INVALID_ID", "ID tidak valid"))
	}
	subs, svcErr := services.GetRecipeSubtypes(getDB(c), id)
	if svcErr != nil {
		se := svcErr.(*services.ServiceError)
		return c.Status(se.Code).JSON(helpers.YvError(se.ErrCode, se.Message))
	}
	return c.JSON(helpers.YvSuccess(resources.RecipeSubtypeListResource(subs)))
}
