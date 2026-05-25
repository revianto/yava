package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/models"
	"github.com/revianto/yava/api/app/repositories"
	"github.com/revianto/yava/api/exceptions"
	"github.com/revianto/yava/api/helpers"
	"gorm.io/gorm"
)

type ValidateRecipeCreate struct {
	Name       string `json:"name" validate:"required"`
	TypeId     int64  `json:"type_id" validate:"required,gt=0"`
	Visibility string `json:"visibility" validate:"required,oneof=private public group"`
}

type ValidateRecipeUpdate struct {
	Id string `json:"id" validate:"required"`
}

func RecipeList(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (any, any) {
	return repositories.RecipeIndex(tx, data, c, locale)
}

func RecipeShow(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id any) (any, any) {
	rid := helpers.Conv(id).Int64()
	if rid <= 0 {
		return nil, exceptions.ErrorException(c, fiber.StatusBadRequest, "ID tidak valid")
	}
	return repositories.RecipeSingle(tx, data, c, locale, func(db *gorm.DB) *gorm.DB {
		return db.Where("yv_recipe.id = ?", rid)
	})
}

func RecipeCreate(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (any, any) {
	if _, err := models.Validate(c, data, new(ValidateRecipeCreate), locale); err != nil {
		return nil, err
	}
	uid, _ := c.Locals("yv_user_id").(int64)
	data["owner_id"] = uid
	if helpers.Conv(data["visibility"]).String() == "" {
		data["visibility"] = "private"
	}
	return repositories.RecipeCreate(tx, data, c, locale)
}

func RecipeUpdate(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id any) (any, any) {
	rid := helpers.Conv(id).Int64()
	if rid <= 0 {
		return nil, exceptions.ErrorException(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	recipeMap, err := repositories.RecipeSingle(tx, data, c, locale, func(db *gorm.DB) *gorm.DB {
		return db.Where("yv_recipe.id = ?", rid)
	})
	if err != nil {
		return nil, err
	}

	if helpers.Conv(recipeMap["is_default"]).Bool() {
		return nil, exceptions.ErrorException(c, fiber.StatusForbidden, "Resep default tidak bisa diubah")
	}
	uid, _ := c.Locals("yv_user_id").(int64)
	ownerIdRaw := helpers.Conv(recipeMap["owner_id"]).Int64()
	if ownerIdRaw != uid {
		return nil, exceptions.ErrorException(c, fiber.StatusForbidden, "Kamu bukan pemilik resep ini")
	}

	if v, ok := data["visibility"].(string); ok && v != "" {
		if v != "private" && v != "public" && v != "group" {
			return nil, exceptions.ErrorException(c, fiber.StatusUnprocessableEntity, "Visibility tidak valid")
		}
	}

	return repositories.RecipeUpdate(tx, data, c, locale, rid, uid)
}

func recipeOwnerCheck(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id int64) (map[string]any, int64, any) {
	recipeMap, err := repositories.RecipeSingle(tx, data, c, locale, func(db *gorm.DB) *gorm.DB {
		return db.Where("yv_recipe.id = ?", id)
	})
	if err != nil {
		return nil, 0, err
	}
	if helpers.Conv(recipeMap["is_default"]).Bool() {
		return nil, 0, exceptions.ErrorException(c, fiber.StatusForbidden, "Resep default tidak bisa diubah")
	}
	uid, _ := c.Locals("yv_user_id").(int64)
	if helpers.Conv(recipeMap["owner_id"]).Int64() != uid {
		return nil, 0, exceptions.ErrorException(c, fiber.StatusForbidden, "Kamu bukan pemilik resep ini")
	}
	return recipeMap, uid, nil
}

func RecipeArchive(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id any) (any, any) {
	rid := helpers.Conv(id).Int64()
	if rid <= 0 {
		return nil, exceptions.ErrorException(c, fiber.StatusBadRequest, "ID tidak valid")
	}
	_, uid, err := recipeOwnerCheck(tx, data, c, locale, rid)
	if err != nil {
		return nil, err
	}
	return repositories.RecipeSetArchived(tx, data, c, locale, rid, uid, true)
}

func RecipeRestore(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id any) (any, any) {
	rid := helpers.Conv(id).Int64()
	if rid <= 0 {
		return nil, exceptions.ErrorException(c, fiber.StatusBadRequest, "ID tidak valid")
	}
	_, uid, err := recipeOwnerCheck(tx, data, c, locale, rid)
	if err != nil {
		return nil, err
	}
	return repositories.RecipeSetArchived(tx, data, c, locale, rid, uid, false)
}

func RecipeDuplicate(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id any) (any, any) {
	rid := helpers.Conv(id).Int64()
	if rid <= 0 {
		return nil, exceptions.ErrorException(c, fiber.StatusBadRequest, "ID tidak valid")
	}
	uid, _ := c.Locals("yv_user_id").(int64)
	return repositories.RecipeDuplicate(tx, data, c, locale, rid, uid)
}
