package resources

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/helpers"
)

type ResponseExample struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func ExampleResource(c *fiber.Ctx, data any) any {
	return ToResource(c, data, ExampleSingleResource)
}

func ExampleSingleResource(c *fiber.Ctx, data any) ResponseExample {
	dataMap, ok := data.(map[string]any)
	if !ok {
		return ResponseExample{}
	}
	return ResponseExample{
		Id:   helpers.Conv(dataMap["id"]).String(),
		Name: helpers.Conv(dataMap["name"]).String(),
	}
}
