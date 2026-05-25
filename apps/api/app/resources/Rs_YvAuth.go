package resources

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/helpers"
)

type UserResponse struct {
	Id        string  `json:"id"`
	Email     string  `json:"email"`
	Name      string  `json:"name"`
	AvatarUrl *string `json:"avatar_url"`
}

func UserResource(c *fiber.Ctx, data any) any {
	return ToResource(c, data, UserSingleResource)
}

func UserSingleResource(c *fiber.Ctx, data any) UserResponse {
	m, _ := data.(map[string]any)
	var avatarUrl *string
	if v := helpers.Conv(m["avatar_url"]).String(); v != "" {
		avatarUrl = &v
	}
	return UserResponse{
		Id:        helpers.Conv(m["id"]).String(),
		Email:     helpers.Conv(m["email"]).String(),
		Name:      helpers.Conv(m["name"]).String(),
		AvatarUrl: avatarUrl,
	}
}
