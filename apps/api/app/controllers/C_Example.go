package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/resources"
	"github.com/revianto/yava/api/app/services"
	"github.com/revianto/yava/api/exceptions"
)

func ExampleList(c *fiber.Ctx) error {
	result, err := services.ExampleList(getDB(c), getBodyData(c), c, getLocale(c))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.ExampleResource(c, result))
}

func ExampleShow(c *fiber.Ctx) error {
	result, err := services.ExampleShow(getDB(c), getBodyData(c), c, getLocale(c), c.Params("id"))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.ExampleResource(c, result))
}

func ExampleCreate(c *fiber.Ctx) error {
	result, err := services.ExampleCreate(getDB(c), getBodyData(c), c, getLocale(c))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.ExampleResource(c, result))
}

func ExampleUpdate(c *fiber.Ctx) error {
	result, err := services.ExampleUpdate(getDB(c), getBodyData(c), c, getLocale(c), c.Params("id"))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.ExampleResource(c, result))
}

func ExampleDelete(c *fiber.Ctx) error {
	data := getBodyData(c)
	data["id"] = c.Params("id")
	_, err := services.ExampleDelete(getDB(c), data, c, getLocale(c))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.DeleteResource())
}
