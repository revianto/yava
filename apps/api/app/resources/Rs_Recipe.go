package resources

import (
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/helpers"
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
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	AvatarUrl *string `json:"avatar_url"`
}

type RecipeResponse struct {
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description *string                `json:"description"`
	Type        *RecipeTypeResponse    `json:"type"`
	Subtype     *RecipeSubtypeResponse `json:"subtype"`
	Visibility  string                 `json:"visibility"`
	IsDefault   bool                   `json:"is_default"`
	IsArchived  bool                   `json:"is_archived"`
	SavesCount  int                    `json:"saves_count"`
	Params      RecipeParamsResponse   `json:"params"`
	Timeline    []TimelineItemResponse `json:"timeline"`
	Owner       *RecipeOwnerResponse   `json:"owner"`
}

func RecipeResource(c *fiber.Ctx, data any) any {
	return ToResource(c, data, RecipeSingleResource)
}

func RecipeSingleResource(c *fiber.Ctx, data any) RecipeResponse {
	m, ok := data.(map[string]any)
	if !ok {
		return RecipeResponse{}
	}

	resp := RecipeResponse{
		Id:         helpers.Conv(m["id"]).String(),
		Name:       helpers.Conv(m["name"]).String(),
		Visibility: helpers.Conv(m["visibility"]).String(),
		IsDefault:  helpers.Conv(m["is_default"]).Bool(),
		IsArchived: helpers.Conv(m["is_archived"]).Bool(),
		SavesCount: int(helpers.Conv(m["saves_count"]).Int()),
		Params: RecipeParamsResponse{
			Dose:  optStr(m["param_dose"]),
			Yield: optStr(m["param_yield"]),
			Temp:  optStr(m["param_temp"]),
			Grind: optStr(m["param_grind"]),
			Ratio: optStr(m["param_ratio"]),
		},
	}

	if v := helpers.Conv(m["description"]).String(); v != "" {
		resp.Description = &v
	}

	// Type
	if tm, ok := m["type"].(map[string]any); ok && tm != nil {
		t := RecipeTypeSingleResource(c, tm)
		resp.Type = &t
	}
	// Subtype
	if sm, ok := m["subtype"].(map[string]any); ok && sm != nil {
		s := RecipeSubtypeSingleResource(c, sm)
		resp.Subtype = &s
	}
	// Owner
	if om, ok := m["owner"].(map[string]any); ok && om != nil {
		o := RecipeOwnerResponse{
			Id:        helpers.Conv(om["id"]).String(),
			Name:      helpers.Conv(om["name"]).String(),
			AvatarUrl: optStr(om["avatar_url"]),
		}
		resp.Owner = &o
	}

	// Timeline: merge sessions + notes, sort by sort_order
	var timeline []TimelineItemResponse
	if sessions, ok := m["sessions"].([]any); ok {
		for _, si := range sessions {
			sm, _ := si.(map[string]any)
			name := helpers.Conv(sm["name"]).String()
			dur := int(helpers.Conv(sm["duration_sec"]).Int())
			item := TimelineItemResponse{
				Kind:        "session",
				SortOrder:   int(helpers.Conv(sm["sort_order"]).Int()),
				Name:        &name,
				DurationSec: &dur,
			}
			if n := helpers.Conv(sm["note"]).String(); n != "" {
				item.Note = &n
			}
			timeline = append(timeline, item)
		}
	}
	if notes, ok := m["notes"].([]any); ok {
		for _, ni := range notes {
			nm, _ := ni.(map[string]any)
			content := helpers.Conv(nm["content"]).String()
			timeline = append(timeline, TimelineItemResponse{
				Kind:      "note",
				SortOrder: int(helpers.Conv(nm["sort_order"]).Int()),
				Content:   &content,
			})
		}
	}
	sort.Slice(timeline, func(i, j int) bool { return timeline[i].SortOrder < timeline[j].SortOrder })
	resp.Timeline = timeline

	return resp
}

func optStr(v any) *string {
	s := helpers.Conv(v).String()
	if s == "" {
		return nil
	}
	return &s
}
