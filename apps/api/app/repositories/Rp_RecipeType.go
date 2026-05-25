package repositories

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/models"
	"gorm.io/gorm"
)

func RecipeTypeIndex(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (models.IndexData, any) {
	return models.GetIndexData(tx, data, c, locale, models.RecipeType{})
}

func RecipeTypeSingle(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, where func(*gorm.DB) *gorm.DB) (map[string]any, any) {
	return models.GetSingleData(tx, data, c, locale, where, models.RecipeType{})
}

func RecipeSubtypeMultiple(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, where func(*gorm.DB) *gorm.DB) ([]map[string]any, any) {
	return models.GetMultipleData(tx, data, c, locale, where, models.RecipeSubtype{})
}
