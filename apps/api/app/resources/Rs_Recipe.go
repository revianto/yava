package resources

import (
	"sort"

	"github.com/revianto/yava/api/app/models"
)

type RecipeParamsResponse struct {
	Dose  *string `json:"dose"`
	Yield *string `json:"yield"`
	Temp  *string `json:"temp"`
	Grind *string `json:"grind"`
	Ratio *string `json:"ratio"`
}

type TimelineItemResponse struct {
	Kind        string  `json:"kind"`
	SortOrder   int     `json:"sort_order"`
	Name        *string `json:"name,omitempty"`
	DurationSec *int    `json:"duration_sec,omitempty"`
	Note        *string `json:"note,omitempty"`
	Content     *string `json:"content,omitempty"`
}

type RecipeOwnerResponse struct {
	Id        int64   `json:"id"`
	Name      string  `json:"name"`
	AvatarUrl *string `json:"avatar_url"`
}

type RecipeResponse struct {
	Id          int64                `json:"id"`
	Name        string               `json:"name"`
	Description *string              `json:"description"`
	Type        *RecipeTypeResponse  `json:"type"`
	Subtype     *RecipeSubtypeResponse `json:"subtype"`
	Visibility  string               `json:"visibility"`
	IsDefault   bool                 `json:"is_default"`
	IsArchived  bool                 `json:"is_archived"`
	SavesCount  int                  `json:"saves_count"`
	Params      RecipeParamsResponse `json:"params"`
	Timeline    []TimelineItemResponse `json:"timeline"`
	Owner       *RecipeOwnerResponse `json:"owner"`
}

func RecipeResource(r models.Recipe) RecipeResponse {
	var typeResp *RecipeTypeResponse
	if r.Type != nil {
		t := RecipeTypeResource(*r.Type)
		typeResp = &t
	}
	var subtypeResp *RecipeSubtypeResponse
	if r.Subtype != nil {
		s := RecipeSubtypeResource(*r.Subtype)
		subtypeResp = &s
	}
	var ownerResp *RecipeOwnerResponse
	if r.Owner != nil {
		ownerResp = &RecipeOwnerResponse{Id: r.Owner.Id, Name: r.Owner.Name, AvatarUrl: r.Owner.AvatarUrl}
	}

	// Merge sessions + notes into unified timeline sorted by sort_order
	var timeline []TimelineItemResponse
	for _, s := range r.Sessions {
		kind := "session"
		timeline = append(timeline, TimelineItemResponse{
			Kind: kind, SortOrder: s.SortOrder,
			Name: &s.Name, DurationSec: &s.DurationSec, Note: s.Note,
		})
	}
	for _, n := range r.Notes {
		kind := "note"
		timeline = append(timeline, TimelineItemResponse{
			Kind: kind, SortOrder: n.SortOrder, Content: &n.Content,
		})
	}
	sort.Slice(timeline, func(i, j int) bool { return timeline[i].SortOrder < timeline[j].SortOrder })

	return RecipeResponse{
		Id: r.Id, Name: r.Name, Description: r.Description,
		Type: typeResp, Subtype: subtypeResp,
		Visibility: r.Visibility, IsDefault: r.IsDefault, IsArchived: r.IsArchived,
		SavesCount: r.SavesCount,
		Params: RecipeParamsResponse{
			Dose: r.ParamDose, Yield: r.ParamYield,
			Temp: r.ParamTemp, Grind: r.ParamGrind, Ratio: r.ParamRatio,
		},
		Timeline: timeline,
		Owner:    ownerResp,
	}
}

func RecipeListResource(recipes []models.Recipe) []RecipeResponse {
	out := make([]RecipeResponse, len(recipes))
	for i, r := range recipes {
		out[i] = RecipeResource(r)
	}
	return out
}
