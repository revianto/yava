package repositories

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/models"
	"gorm.io/gorm"
)

func ExampleIndex(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (models.IndexData, any) {
	return models.GetIndexData(tx, data, c, locale, models.ExampleModel{})
}

func ExampleSingle(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, where func(db *gorm.DB) *gorm.DB) (map[string]any, any) {
	return models.GetSingleData(tx, data, c, locale, where, models.ExampleModel{})
}

func ExampleMultiple(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, where func(db *gorm.DB) *gorm.DB) ([]map[string]any, any) {
	return models.GetMultipleData(tx, data, c, locale, where, models.ExampleModel{})
}
