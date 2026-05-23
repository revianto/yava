package repositories

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/models"
	"gorm.io/gorm"
)

// =============================================================================
// User REPOSITORY
// =============================================================================

// UserIndex returns paginated list of Users
func UserIndex(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (models.IndexData, any) {
	return models.GetIndexData(tx, data, c, locale, models.UserModel{})
}

// UserSingle returns single User by condition
func UserSingle(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, where func(db *gorm.DB) *gorm.DB) (map[string]any, any) {
	return models.GetSingleData(tx, data, c, locale, where, models.UserModel{})
}

// UserMultiple returns multiple Users by condition
func UserMultiple(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, where func(db *gorm.DB) *gorm.DB) ([]map[string]any, any) {
	return models.GetMultipleData(tx, data, c, locale, where, models.UserModel{})
}
