package resources

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/helpers"
)

// =============================================================================
// RESPONSE TYPES
// =============================================================================

type GroupResponse struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	AvatarUrl   *string `json:"avatar_url"`
	InviteCode  string  `json:"invite_code"`
	MemberCount int     `json:"member_count"`
}

type GroupMemberResponse struct {
	Id     string  `json:"id"`
	UserId string  `json:"user_id"`
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Avatar *string `json:"avatar_url"`
	Role   string  `json:"role"`
}

type GroupRecipeResponse struct {
	Id          string `json:"id"`
	GroupId     string `json:"group_id"`
	RecipeId    string `json:"recipe_id"`
	RecipeName  string `json:"recipe_name"`
	Status      string `json:"status"`
	SubmittedBy string `json:"submitted_by"`
}

// =============================================================================
// RESOURCE FUNCTIONS
// =============================================================================

func GroupResource(c *fiber.Ctx, data any) any {
	return ToResource(c, data, GroupSingleResource)
}

func GroupSingleResource(c *fiber.Ctx, data any) GroupResponse {
	m, ok := data.(map[string]any)
	if !ok {
		return GroupResponse{}
	}

	memberCount := 0
	if members, ok := m["members"].([]any); ok {
		memberCount = len(members)
	}

	var description *string
	if v := helpers.Conv(m["description"]).String(); v != "" {
		description = &v
	}
	var avatarUrl *string
	if v := helpers.Conv(m["avatar_url"]).String(); v != "" {
		avatarUrl = &v
	}

	return GroupResponse{
		Id:          helpers.Conv(m["id"]).String(),
		Name:        helpers.Conv(m["name"]).String(),
		Description: description,
		AvatarUrl:   avatarUrl,
		InviteCode:  helpers.Conv(m["invite_code"]).String(),
		MemberCount: memberCount,
	}
}

func GroupMemberResource(c *fiber.Ctx, data any) any {
	return ToResource(c, data, GroupMemberSingleResource)
}

func GroupMemberSingleResource(c *fiber.Ctx, data any) GroupMemberResponse {
	m, ok := data.(map[string]any)
	if !ok {
		return GroupMemberResponse{}
	}

	var name, email string
	var avatar *string
	if user, ok := m["user"].(map[string]any); ok {
		name = helpers.Conv(user["name"]).String()
		email = helpers.Conv(user["email"]).String()
		if v := helpers.Conv(user["avatar_url"]).String(); v != "" {
			avatar = &v
		}
	}

	return GroupMemberResponse{
		Id:     helpers.Conv(m["id"]).String(),
		UserId: helpers.Conv(m["user_id"]).String(),
		Name:   name,
		Email:  email,
		Avatar: avatar,
		Role:   helpers.Conv(m["role"]).String(),
	}
}

func GroupRecipeResource(c *fiber.Ctx, data any) any {
	return ToResource(c, data, GroupRecipeSingleResource)
}

func GroupRecipeSingleResource(c *fiber.Ctx, data any) GroupRecipeResponse {
	m, ok := data.(map[string]any)
	if !ok {
		return GroupRecipeResponse{}
	}

	var recipeName string
	if recipe, ok := m["recipe"].(map[string]any); ok {
		recipeName = helpers.Conv(recipe["name"]).String()
	}

	return GroupRecipeResponse{
		Id:          helpers.Conv(m["id"]).String(),
		GroupId:     helpers.Conv(m["group_id"]).String(),
		RecipeId:    helpers.Conv(m["recipe_id"]).String(),
		RecipeName:  recipeName,
		Status:      helpers.Conv(m["status"]).String(),
		SubmittedBy: helpers.Conv(m["submitted_by"]).String(),
	}
}
