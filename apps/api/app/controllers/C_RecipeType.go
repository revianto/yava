package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/resources"
	"github.com/revianto/yava/api/app/services"
	"github.com/revianto/yava/api/exceptions"
)

// GET /v1/types
func RecipeTypeList(c *fiber.Ctx) error {
	result, err := services.RecipeTypeList(getDB(c), getBodyData(c), c, getLocale(c))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.RecipeTypeResource(c, result))
}

// GET /v1/types/:id/subtypes
func RecipeSubtypeList(c *fiber.Ctx) error {
	result, err := services.RecipeSubtypeList(getDB(c), getBodyData(c), c, getLocale(c), c.Params("id"))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.RecipeSubtypeResource(c, result))
}
