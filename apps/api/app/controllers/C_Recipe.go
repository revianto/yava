package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/resources"
	"github.com/revianto/yava/api/app/repositories"
	"github.com/revianto/yava/api/app/services"
	"github.com/revianto/yava/api/helpers"
)

func yvUserID(c *fiber.Ctx) int64 {
	id, _ := c.Locals("yv_user_id").(int64)
	return id
}

// GET /v1/recipes
func RecipeList(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	input := services.RecipeListInput{Page: page, Limit: limit}
	if vis := c.Query("visibility"); vis != "" {
		input.Visibility = &vis
	}
	if tid := c.Query("type_id"); tid != "" {
		if id, err := strconv.ParseInt(tid, 10, 64); err == nil {
			input.TypeId = &id
		}
	}
	if c.Query("mine") == "true" {
		uid := yvUserID(c)
		if uid > 0 {
			input.OwnerId = &uid
		}
	}

	recipes, total, svcErr := services.ListRecipes(getDB(c), input)
	if svcErr != nil {
		se := svcErr.(*services.ServiceError)
		return c.Status(se.Code).JSON(helpers.YvError(se.ErrCode, se.Message))
	}
	return c.JSON(helpers.YvList(resources.RecipeListResource(recipes), int64(page), int64(limit), total))
}

// GET /v1/recipes/:id
func RecipeShow(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.YvError("INVALID_ID", "ID tidak valid"))
	}
	recipe, svcErr := services.GetRecipe(getDB(c), id)
	if svcErr != nil {
		se := svcErr.(*services.ServiceError)
		return c.Status(se.Code).JSON(helpers.YvError(se.ErrCode, se.Message))
	}
	return c.JSON(helpers.YvSuccess(resources.RecipeResource(*recipe)))
}

// POST /v1/recipes
func RecipeCreate(c *fiber.Ctx) error {
	var body struct {
		TypeId      int64   `json:"type_id"`
		SubtypeId   *int64  `json:"subtype_id"`
		Name        string  `json:"name"`
		Description *string `json:"description"`
		Visibility  string  `json:"visibility"`
		ParamDose   *string `json:"param_dose"`
		ParamYield  *string `json:"param_yield"`
		ParamTemp   *string `json:"param_temp"`
		ParamGrind  *string `json:"param_grind"`
		ParamRatio  *string `json:"param_ratio"`
		Sessions    []struct {
			SortOrder   int     `json:"sort_order"`
			Name        string  `json:"name"`
			DurationSec int     `json:"duration_sec"`
			Note        *string `json:"note"`
		} `json:"sessions"`
		Notes []struct {
			SortOrder int    `json:"sort_order"`
			Content   string `json:"content"`
		} `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.YvError("INVALID_BODY", "Request body tidak valid"))
	}

	input := services.RecipeCreateInput{
		TypeId: body.TypeId, SubtypeId: body.SubtypeId,
		Name: body.Name, Description: body.Description, Visibility: body.Visibility,
		ParamDose: body.ParamDose, ParamYield: body.ParamYield,
		ParamTemp: body.ParamTemp, ParamGrind: body.ParamGrind, ParamRatio: body.ParamRatio,
	}
	for _, s := range body.Sessions {
		input.Sessions = append(input.Sessions, repositories.RecipeStepInput{
			SortOrder: s.SortOrder, Name: s.Name, DurationSec: s.DurationSec, Note: s.Note,
		})
	}
	for _, n := range body.Notes {
		input.Notes = append(input.Notes, repositories.RecipeNoteInput{
			SortOrder: n.SortOrder, Content: n.Content,
		})
	}

	recipe, svcErr := services.CreateRecipe(getDB(c), yvUserID(c), input)
	if svcErr != nil {
		se := svcErr.(*services.ServiceError)
		return c.Status(se.Code).JSON(helpers.YvError(se.ErrCode, se.Message))
	}
	return c.Status(fiber.StatusCreated).JSON(helpers.YvSuccess(resources.RecipeResource(*recipe)))
}

// PUT /v1/recipes/:id
func RecipeUpdate(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.YvError("INVALID_ID", "ID tidak valid"))
	}

	var body struct {
		TypeId      *int64  `json:"type_id"`
		SubtypeId   *int64  `json:"subtype_id"`
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Visibility  *string `json:"visibility"`
		ParamDose   *string `json:"param_dose"`
		ParamYield  *string `json:"param_yield"`
		ParamTemp   *string `json:"param_temp"`
		ParamGrind  *string `json:"param_grind"`
		ParamRatio  *string `json:"param_ratio"`
		Sessions    *[]struct {
			SortOrder   int     `json:"sort_order"`
			Name        string  `json:"name"`
			DurationSec int     `json:"duration_sec"`
			Note        *string `json:"note"`
		} `json:"sessions"`
		Notes *[]struct {
			SortOrder int    `json:"sort_order"`
			Content   string `json:"content"`
		} `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.YvError("INVALID_BODY", "Request body tidak valid"))
	}

	input := services.RecipeUpdateInput{
		TypeId: body.TypeId, SubtypeId: body.SubtypeId,
		Name: body.Name, Description: body.Description, Visibility: body.Visibility,
		ParamDose: body.ParamDose, ParamYield: body.ParamYield,
		ParamTemp: body.ParamTemp, ParamGrind: body.ParamGrind, ParamRatio: body.ParamRatio,
	}
	if body.Sessions != nil {
		sessions := make([]repositories.RecipeStepInput, len(*body.Sessions))
		for i, s := range *body.Sessions {
			sessions[i] = repositories.RecipeStepInput{SortOrder: s.SortOrder, Name: s.Name, DurationSec: s.DurationSec, Note: s.Note}
		}
		input.Sessions = &sessions
	}
	if body.Notes != nil {
		notes := make([]repositories.RecipeNoteInput, len(*body.Notes))
		for i, n := range *body.Notes {
			notes[i] = repositories.RecipeNoteInput{SortOrder: n.SortOrder, Content: n.Content}
		}
		input.Notes = &notes
	}

	recipe, svcErr := services.UpdateRecipe(getDB(c), id, yvUserID(c), input)
	if svcErr != nil {
		se := svcErr.(*services.ServiceError)
		return c.Status(se.Code).JSON(helpers.YvError(se.ErrCode, se.Message))
	}
	return c.JSON(helpers.YvSuccess(resources.RecipeResource(*recipe)))
}
