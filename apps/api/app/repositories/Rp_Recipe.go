package repositories

import (
	"github.com/revianto/yava/api/app/models"
	"gorm.io/gorm"
)

type RecipeListParams struct {
	OwnerId    *int64
	Visibility *string
	TypeId     *int64
	Page       int
	Limit      int
}

func RecipeList(db *gorm.DB, p RecipeListParams) ([]models.Recipe, int64, error) {
	var recipes []models.Recipe
	var total int64

	q := db.Model(&models.Recipe{}).
		Preload("Type").Preload("Subtype").Preload("Owner").
		Where("yv_recipe.deleted_at IS NULL")

	if p.OwnerId != nil {
		q = q.Where("owner_id = ?", *p.OwnerId)
	}
	if p.Visibility != nil {
		q = q.Where("visibility = ?", *p.Visibility)
	}
	if p.TypeId != nil {
		q = q.Where("type_id = ?", *p.TypeId)
	}

	q.Count(&total)

	offset := (p.Page - 1) * p.Limit
	err := q.Order("created_at DESC").Offset(offset).Limit(p.Limit).Find(&recipes).Error
	return recipes, total, err
}

func RecipeById(db *gorm.DB, id int64) (*models.Recipe, error) {
	var recipe models.Recipe
	err := db.Preload("Type").Preload("Subtype").Preload("Owner").
		Preload("Sessions", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC") }).
		Preload("Notes", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order ASC") }).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&recipe).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

type RecipeCreateInput struct {
	OwnerId     int64
	TypeId      int64
	SubtypeId   *int64
	Name        string
	Description *string
	Visibility  string
	ParamDose   *string
	ParamYield  *string
	ParamTemp   *string
	ParamGrind  *string
	ParamRatio  *string
	Sessions    []RecipeStepInput
	Notes       []RecipeNoteInput
}

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

func RecipeCreate(db *gorm.DB, input RecipeCreateInput) (*models.Recipe, error) {
	recipe := models.Recipe{
		OwnerId:     &input.OwnerId,
		TypeId:      input.TypeId,
		SubtypeId:   input.SubtypeId,
		Name:        input.Name,
		Description: input.Description,
		Visibility:  input.Visibility,
		ParamDose:   input.ParamDose,
		ParamYield:  input.ParamYield,
		ParamTemp:   input.ParamTemp,
		ParamGrind:  input.ParamGrind,
		ParamRatio:  input.ParamRatio,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&recipe).Error; err != nil {
			return err
		}
		for _, s := range input.Sessions {
			sess := models.RecipeSession{
				RecipeId:    recipe.Id,
				SortOrder:   s.SortOrder,
				Name:        s.Name,
				DurationSec: s.DurationSec,
				Note:        s.Note,
			}
			if err := tx.Create(&sess).Error; err != nil {
				return err
			}
		}
		for _, n := range input.Notes {
			note := models.RecipeNote{
				RecipeId:  recipe.Id,
				SortOrder: n.SortOrder,
				Content:   n.Content,
			}
			if err := tx.Create(&note).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return RecipeById(db, recipe.Id)
}

type RecipeUpdateInput struct {
	TypeId      *int64
	SubtypeId   *int64
	Name        *string
	Description *string
	Visibility  *string
	ParamDose   *string
	ParamYield  *string
	ParamTemp   *string
	ParamGrind  *string
	ParamRatio  *string
	Sessions    *[]RecipeStepInput
	Notes       *[]RecipeNoteInput
}

func RecipeUpdate(db *gorm.DB, id int64, ownerId int64, input RecipeUpdateInput) (*models.Recipe, error) {
	err := db.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{}
		if input.TypeId != nil {
			updates["type_id"] = *input.TypeId
		}
		if input.SubtypeId != nil {
			updates["subtype_id"] = *input.SubtypeId
		}
		if input.Name != nil {
			updates["name"] = *input.Name
		}
		if input.Description != nil {
			updates["description"] = *input.Description
		}
		if input.Visibility != nil {
			updates["visibility"] = *input.Visibility
		}
		if input.ParamDose != nil {
			updates["param_dose"] = *input.ParamDose
		}
		if input.ParamYield != nil {
			updates["param_yield"] = *input.ParamYield
		}
		if input.ParamTemp != nil {
			updates["param_temp"] = *input.ParamTemp
		}
		if input.ParamGrind != nil {
			updates["param_grind"] = *input.ParamGrind
		}
		if input.ParamRatio != nil {
			updates["param_ratio"] = *input.ParamRatio
		}
		if len(updates) > 0 {
			if err := tx.Model(&models.Recipe{}).
				Where("id = ? AND owner_id = ? AND is_default = FALSE AND deleted_at IS NULL", id, ownerId).
				Updates(updates).Error; err != nil {
				return err
			}
		}
		if input.Sessions != nil {
			if err := tx.Where("recipe_id = ?", id).Delete(&models.RecipeSession{}).Error; err != nil {
				return err
			}
			for _, s := range *input.Sessions {
				sess := models.RecipeSession{RecipeId: id, SortOrder: s.SortOrder, Name: s.Name, DurationSec: s.DurationSec, Note: s.Note}
				if err := tx.Create(&sess).Error; err != nil {
					return err
				}
			}
		}
		if input.Notes != nil {
			if err := tx.Where("recipe_id = ?", id).Delete(&models.RecipeNote{}).Error; err != nil {
				return err
			}
			for _, n := range *input.Notes {
				note := models.RecipeNote{RecipeId: id, SortOrder: n.SortOrder, Content: n.Content}
				if err := tx.Create(&note).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return RecipeById(db, id)
}
