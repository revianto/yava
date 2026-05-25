package repositories

import (
	"encoding/json"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/models"
	"github.com/revianto/yava/api/exceptions"
	"github.com/revianto/yava/api/helpers"
	"gorm.io/gorm"
)

func recipeToMap(r models.Recipe) map[string]any {
	b, _ := json.Marshal(r)
	var m map[string]any
	json.Unmarshal(b, &m)
	return m
}

func RecipeIndex(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (models.IndexData, any) {
	var recipes []models.Recipe
	var total int64

	q := tx.Model(&models.Recipe{}).
		Preload("Type").Preload("Subtype").Preload("Owner").
		Where("yv_recipe.deleted_at IS NULL")

	if v, _ := data["visibility"].(string); v != "" {
		q = q.Where("visibility = ?", v)
	}
	if tid := helpers.Conv(data["type_id"]).Int64(); tid > 0 {
		q = q.Where("type_id = ?", tid)
	}
	if mine, _ := data["mine"].(bool); mine {
		uid := helpers.Conv(data["owner_id"]).Int64()
		if uid > 0 {
			q = q.Where("owner_id = ?", uid)
		}
	}

	q.Count(&total)

	page := helpers.Conv(data["page"]).Default(c.Query("page", "1")).Int()
	limit := helpers.Conv(data["limit"]).Default(c.Query("limit", "20")).Int()
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	if err := q.Order("yv_recipe.created_at DESC").Offset(offset).Limit(limit).Find(&recipes).Error; err != nil {
		return models.IndexData{}, exceptions.ErrorException(c, fiber.StatusInternalServerError, "Gagal mengambil daftar resep")
	}

	items := make([]any, len(recipes))
	for i, r := range recipes {
		items[i] = recipeToMap(r)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return models.IndexData{
		Data: items,
		Meta: models.Pagination{Page: page, Limit: limit, Total: int(total), TotalPages: totalPages},
	}, nil
}

func RecipeSingle(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, where func(*gorm.DB) *gorm.DB) (map[string]any, any) {
	var recipe models.Recipe
	q := tx.Preload("Type").Preload("Subtype").Preload("Owner").
		Preload("Sessions", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC") }).
		Preload("Notes", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC") }).
		Where("yv_recipe.deleted_at IS NULL")
	if where != nil {
		q = where(q)
	}
	if err := q.First(&recipe).Error; err != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotFound, "Resep tidak ditemukan")
	}
	return recipeToMap(recipe), nil
}

// RecipeStepInput and RecipeNoteInput are used by service for create/update
type RecipeStepInput struct {
	SortOrder   int
	Name        string
	DurationSec int
	Note        *string
}

type RecipeNoteInput struct {
	SortOrder int
	Content   string
}

func RecipeCreate(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (map[string]any, any) {
	ownerId := helpers.Conv(data["owner_id"]).Int64()
	typeId := helpers.Conv(data["type_id"]).Int64()
	name, _ := data["name"].(string)
	description, _ := data["description"].(*string)
	visibility, _ := data["visibility"].(string)
	paramDose, _ := data["param_dose"].(*string)
	paramYield, _ := data["param_yield"].(*string)
	paramTemp, _ := data["param_temp"].(*string)
	paramGrind, _ := data["param_grind"].(*string)
	paramRatio, _ := data["param_ratio"].(*string)

	recipe := models.Recipe{
		OwnerId:     &ownerId,
		TypeId:      typeId,
		Name:        name,
		Description: description,
		Visibility:  visibility,
		ParamDose:   paramDose,
		ParamYield:  paramYield,
		ParamTemp:   paramTemp,
		ParamGrind:  paramGrind,
		ParamRatio:  paramRatio,
	}
	if sid := helpers.Conv(data["subtype_id"]).Int64(); sid > 0 {
		recipe.SubtypeId = &sid
	}

	sessions, _ := data["sessions"].([]RecipeStepInput)
	notes, _ := data["notes"].([]RecipeNoteInput)

	err := tx.Transaction(func(t *gorm.DB) error {
		if err := t.Create(&recipe).Error; err != nil {
			return err
		}
		for _, s := range sessions {
			sess := models.RecipeSession{RecipeId: recipe.Id, SortOrder: s.SortOrder, Name: s.Name, DurationSec: s.DurationSec, Note: s.Note}
			if err := t.Create(&sess).Error; err != nil {
				return err
			}
		}
		for _, n := range notes {
			note := models.RecipeNote{RecipeId: recipe.Id, SortOrder: n.SortOrder, Content: n.Content}
			if err := t.Create(&note).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "Gagal membuat resep")
	}
	return RecipeSingle(tx, data, c, locale, func(db *gorm.DB) *gorm.DB {
		return db.Where("yv_recipe.id = ?", recipe.Id)
	})
}

func RecipeUpdate(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id, ownerId int64) (map[string]any, any) {
	err := tx.Transaction(func(t *gorm.DB) error {
		updates := map[string]any{}
		if v := helpers.Conv(data["type_id"]).Int64(); v > 0 {
			updates["type_id"] = v
		}
		if v := helpers.Conv(data["subtype_id"]).Int64(); v > 0 {
			updates["subtype_id"] = v
		}
		if v, ok := data["name"].(string); ok && v != "" {
			updates["name"] = v
		}
		if v, ok := data["description"]; ok {
			updates["description"] = v
		}
		if v, ok := data["visibility"].(string); ok && v != "" {
			updates["visibility"] = v
		}
		for _, field := range []string{"param_dose", "param_yield", "param_temp", "param_grind", "param_ratio"} {
			if v, ok := data[field]; ok {
				updates[field] = v
			}
		}
		if len(updates) > 0 {
			if err := t.Model(&models.Recipe{}).
				Where("id = ? AND owner_id = ? AND is_default = FALSE AND deleted_at IS NULL", id, ownerId).
				Updates(updates).Error; err != nil {
				return err
			}
		}
		if sessions, ok := data["sessions"].([]RecipeStepInput); ok {
			if err := t.Where("recipe_id = ?", id).Delete(&models.RecipeSession{}).Error; err != nil {
				return err
			}
			for _, s := range sessions {
				sess := models.RecipeSession{RecipeId: id, SortOrder: s.SortOrder, Name: s.Name, DurationSec: s.DurationSec, Note: s.Note}
				if err := t.Create(&sess).Error; err != nil {
					return err
				}
			}
		}
		if notes, ok := data["notes"].([]RecipeNoteInput); ok {
			if err := t.Where("recipe_id = ?", id).Delete(&models.RecipeNote{}).Error; err != nil {
				return err
			}
			for _, n := range notes {
				note := models.RecipeNote{RecipeId: id, SortOrder: n.SortOrder, Content: n.Content}
				if err := t.Create(&note).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "Gagal memperbarui resep")
	}
	return RecipeSingle(tx, data, c, locale, func(db *gorm.DB) *gorm.DB {
		return db.Where("yv_recipe.id = ?", id)
	})
}

func RecipeSetArchived(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id, ownerId int64, archived bool) (map[string]any, any) {
	if txErr := tx.Transaction(func(t *gorm.DB) error {
		return t.Model(&models.Recipe{}).
			Where("id = ? AND owner_id = ? AND is_default = FALSE AND deleted_at IS NULL", id, ownerId).
			Update("is_archived", archived).Error
	}); txErr != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "Gagal memperbarui status resep")
	}
	return RecipeSingle(tx, data, c, locale, func(db *gorm.DB) *gorm.DB {
		return db.Where("yv_recipe.id = ?", id)
	})
}

func RecipeDuplicate(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id, ownerId int64) (map[string]any, any) {
	src, err := RecipeSingle(tx, data, c, locale, func(db *gorm.DB) *gorm.DB {
		return db.Where("yv_recipe.id = ?", id)
	})
	if err != nil {
		return nil, err
	}

	newRecipe := models.Recipe{
		OwnerId:    &ownerId,
		TypeId:     helpers.Conv(src["type_id"]).Int64(),
		Name:       helpers.Conv(src["name"]).String() + " (Salinan)",
		Visibility: "private",
		ParamDose:  optStrFromMap(src, "param_dose"),
		ParamYield: optStrFromMap(src, "param_yield"),
		ParamTemp:  optStrFromMap(src, "param_temp"),
		ParamGrind: optStrFromMap(src, "param_grind"),
		ParamRatio: optStrFromMap(src, "param_ratio"),
	}
	if desc := helpers.Conv(src["description"]).String(); desc != "" {
		newRecipe.Description = &desc
	}
	if sid := helpers.Conv(src["subtype_id"]).Int64(); sid > 0 {
		newRecipe.SubtypeId = &sid
	}

	sessions, _ := src["sessions"].([]any)
	notes, _ := src["notes"].([]any)

	if txErr := tx.Transaction(func(t *gorm.DB) error {
		if err := t.Create(&newRecipe).Error; err != nil {
			return err
		}
		for _, si := range sessions {
			sm, _ := si.(map[string]any)
			note := optStrFromMap(sm, "note")
			sess := models.RecipeSession{
				RecipeId:    newRecipe.Id,
				SortOrder:   int(helpers.Conv(sm["sort_order"]).Int()),
				Name:        helpers.Conv(sm["name"]).String(),
				DurationSec: int(helpers.Conv(sm["duration_sec"]).Int()),
				Note:        note,
			}
			if err := t.Create(&sess).Error; err != nil {
				return err
			}
		}
		for _, ni := range notes {
			nm, _ := ni.(map[string]any)
			note := models.RecipeNote{
				RecipeId:  newRecipe.Id,
				SortOrder: int(helpers.Conv(nm["sort_order"]).Int()),
				Content:   helpers.Conv(nm["content"]).String(),
			}
			if err := t.Create(&note).Error; err != nil {
				return err
			}
		}
		return nil
	}); txErr != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "Gagal menduplikasi resep")
	}

	return RecipeSingle(tx, data, c, locale, func(db *gorm.DB) *gorm.DB {
		return db.Where("yv_recipe.id = ?", newRecipe.Id)
	})
}

func optStrFromMap(m map[string]any, key string) *string {
	v := helpers.Conv(m[key]).String()
	if v == "" {
		return nil
	}
	return &v
}
