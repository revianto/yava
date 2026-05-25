package resources

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/helpers"
)

type RecipeTypeResponse struct {
	Id   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type RecipeSubtypeResponse struct {
	Id     string `json:"id"`
	TypeId string `json:"type_id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
}

func RecipeTypeResource(c *fiber.Ctx, data any) any {
	return ToResource(c, data, RecipeTypeSingleResource)
}

func RecipeTypeSingleResource(c *fiber.Ctx, data any) RecipeTypeResponse {
	m, _ := data.(map[string]any)
	return RecipeTypeResponse{
		Id:   helpers.Conv(m["id"]).String(),
		Code: helpers.Conv(m["code"]).String(),
		Name: helpers.Conv(m["name"]).String(),
	}
}

func RecipeSubtypeResource(c *fiber.Ctx, data any) any {
	return ToResource(c, data, RecipeSubtypeSingleResource)
}

func RecipeSubtypeSingleResource(c *fiber.Ctx, data any) RecipeSubtypeResponse {
	m, _ := data.(map[string]any)
	return RecipeSubtypeResponse{
		Id:     helpers.Conv(m["id"]).String(),
		TypeId: helpers.Conv(m["type_id"]).String(),
		Code:   helpers.Conv(m["code"]).String(),
		Name:   helpers.Conv(m["name"]).String(),
	}
}
