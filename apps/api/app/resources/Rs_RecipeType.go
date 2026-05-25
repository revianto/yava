package resources

import "github.com/revianto/yava/api/app/models"

type RecipeTypeResponse struct {
	Id        int64  `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
}

type RecipeSubtypeResponse struct {
	Id        int64  `json:"id"`
	TypeId    int64  `json:"type_id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
}

func RecipeTypeResource(t models.RecipeType) RecipeTypeResponse {
	return RecipeTypeResponse{Id: t.Id, Code: t.Code, Name: t.Name, SortOrder: t.SortOrder}
}

func RecipeTypeListResource(types []models.RecipeType) []RecipeTypeResponse {
	out := make([]RecipeTypeResponse, len(types))
	for i, t := range types {
		out[i] = RecipeTypeResource(t)
	}
	return out
}

func RecipeSubtypeResource(s models.RecipeSubtype) RecipeSubtypeResponse {
	return RecipeSubtypeResponse{Id: s.Id, TypeId: s.TypeId, Code: s.Code, Name: s.Name, SortOrder: s.SortOrder}
}

func RecipeSubtypeListResource(subs []models.RecipeSubtype) []RecipeSubtypeResponse {
	out := make([]RecipeSubtypeResponse, len(subs))
	for i, s := range subs {
		out[i] = RecipeSubtypeResource(s)
	}
	return out
}
