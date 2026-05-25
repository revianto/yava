package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/models"
	"github.com/revianto/yava/api/app/repositories"
	"github.com/revianto/yava/api/exceptions"
	"github.com/revianto/yava/api/helpers"
	"gorm.io/gorm"
)

func RecipeTypeList(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (any, any) {
	return repositories.RecipeTypeIndex(tx, data, c, locale)
}

func RecipeSubtypeList(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, typeId any) (any, any) {
	id := helpers.Conv(typeId).Int64()
	if id <= 0 {
		return nil, exceptions.ErrorException(c, fiber.StatusBadRequest, "ID tidak valid")
	}
	_, err := repositories.RecipeTypeSingle(tx, data, c, locale, func(db *gorm.DB) *gorm.DB {
		return db.Where(models.RecipeType{Id: id})
	})
	if err != nil {
		return nil, err
	}
	return repositories.RecipeSubtypeMultiple(tx, data, c, locale, func(db *gorm.DB) *gorm.DB {
		return db.Where("type_id = ?", id)
	})
}
