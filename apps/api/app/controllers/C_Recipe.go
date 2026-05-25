package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/resources"
	"github.com/revianto/yava/api/app/services"
	"github.com/revianto/yava/api/exceptions"
)

// GET /v1/recipes
func RecipeList(c *fiber.Ctx) error {
	data := getBodyData(c)
	if v := c.Query("visibility"); v != "" {
		data["visibility"] = v
	}
	if v := c.Query("type_id"); v != "" {
		data["type_id"] = v
	}
	if c.Query("mine") == "true" {
		uid, _ := c.Locals("yv_user_id").(int64)
		data["mine"] = true
		data["owner_id"] = uid
	}
	result, err := services.RecipeList(getDB(c), data, c, getLocale(c))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.RecipeResource(c, result))
}

// GET /v1/recipes/:id
func RecipeShow(c *fiber.Ctx) error {
	result, err := services.RecipeShow(getDB(c), getBodyData(c), c, getLocale(c), c.Params("id"))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.RecipeResource(c, result))
}

// POST /v1/recipes
func RecipeCreate(c *fiber.Ctx) error {
	data := getBodyData(c)
	result, err := services.RecipeCreate(getDB(c), data, c, getLocale(c))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.Status(fiber.StatusCreated).JSON(resources.RecipeResource(c, result))
}

// PUT /v1/recipes/:id
func RecipeUpdate(c *fiber.Ctx) error {
	result, err := services.RecipeUpdate(getDB(c), getBodyData(c), c, getLocale(c), c.Params("id"))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.RecipeResource(c, result))
}
