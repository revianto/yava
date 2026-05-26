package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/resources"
	"github.com/revianto/yava/api/app/services"
	"github.com/revianto/yava/api/exceptions"
)

// =============================================================================
// GROUP CRUD
// =============================================================================

// GET /v1/groups
func GroupList(c *fiber.Ctx) error {
	result, err := services.GroupList(getDB(c), getBodyData(c), c, getLocale(c))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.GroupResource(c, result))
}

// GET /v1/groups/:id
func GroupShow(c *fiber.Ctx) error {
	result, err := services.GroupShow(getDB(c), getBodyData(c), c, getLocale(c), c.Params("id"))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.GroupResource(c, result))
}

// POST /v1/groups
func GroupCreate(c *fiber.Ctx) error {
	result, err := services.GroupCreate(getDB(c), getBodyData(c), c, getLocale(c))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.Status(fiber.StatusCreated).JSON(resources.GroupResource(c, result))
}

// PUT /v1/groups/:id
func GroupUpdate(c *fiber.Ctx) error {
	result, err := services.GroupUpdate(getDB(c), getBodyData(c), c, getLocale(c), c.Params("id"))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.GroupResource(c, result))
}

// DELETE /v1/groups/:id
func GroupDelete(c *fiber.Ctx) error {
	_, err := services.GroupDelete(getDB(c), getBodyData(c), c, getLocale(c), c.Params("id"))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.DeleteResource())
}

// =============================================================================
// GROUP MEMBERS
// =============================================================================

// GET /v1/groups/:id/members
func GroupMemberList(c *fiber.Ctx) error {
	result, err := services.GroupMemberList(getDB(c), getBodyData(c), c, getLocale(c), c.Params("id"))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.GroupMemberResource(c, result))
}

// POST /v1/groups/:id/members
func GroupJoin(c *fiber.Ctx) error {
	result, err := services.GroupJoin(getDB(c), getBodyData(c), c, getLocale(c), c.Params("id"))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.Status(fiber.StatusCreated).JSON(resources.GroupMemberResource(c, result))
}

// DELETE /v1/groups/:id/members/:uid
func GroupMemberRemove(c *fiber.Ctx) error {
	_, err := services.GroupMemberRemove(getDB(c), getBodyData(c), c, getLocale(c), c.Params("id"), c.Params("uid"))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.DeleteResource())
}

// PATCH /v1/groups/:id/members/:uid/role
func GroupMemberSetRole(c *fiber.Ctx) error {
	result, err := services.GroupMemberSetRole(getDB(c), getBodyData(c), c, getLocale(c), c.Params("id"), c.Params("uid"))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.GroupMemberResource(c, result))
}

// =============================================================================
// GROUP RECIPES
// =============================================================================

// GET /v1/groups/:id/recipes
func GroupRecipeList(c *fiber.Ctx) error {
	result, err := services.GroupRecipeList(getDB(c), getBodyData(c), c, getLocale(c), c.Params("id"))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.GroupRecipeResource(c, result))
}

// GET /v1/groups/:id/recipes/pending
func GroupRecipePending(c *fiber.Ctx) error {
	result, err := services.GroupRecipePending(getDB(c), getBodyData(c), c, getLocale(c), c.Params("id"))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.GroupRecipeResource(c, result))
}

// POST /v1/groups/:id/recipes
func GroupRecipeSubmit(c *fiber.Ctx) error {
	result, err := services.GroupRecipeSubmit(getDB(c), getBodyData(c), c, getLocale(c), c.Params("id"))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.Status(fiber.StatusCreated).JSON(resources.GroupRecipeResource(c, result))
}

// PATCH /v1/groups/:id/recipes/:rid/approve
func GroupRecipeApprove(c *fiber.Ctx) error {
	result, err := services.GroupRecipeApprove(getDB(c), getBodyData(c), c, getLocale(c), c.Params("id"), c.Params("rid"))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.GroupRecipeResource(c, result))
}

// PATCH /v1/groups/:id/recipes/:rid/reject
func GroupRecipeReject(c *fiber.Ctx) error {
	result, err := services.GroupRecipeReject(getDB(c), getBodyData(c), c, getLocale(c), c.Params("id"), c.Params("rid"))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.GroupRecipeResource(c, result))
}

// DELETE /v1/groups/:id/recipes/:rid
func GroupRecipeRemove(c *fiber.Ctx) error {
	_, err := services.GroupRecipeRemove(getDB(c), getBodyData(c), c, getLocale(c), c.Params("id"), c.Params("rid"))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.DeleteResource())
}
