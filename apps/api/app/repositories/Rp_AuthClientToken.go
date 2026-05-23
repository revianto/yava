package repositories

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/models"
	"gorm.io/gorm"
)

// =============================================================================
// AuthClientToken REPOSITORY (UNION: AuthClientToken + Service)
// =============================================================================

// AuthClientTokenIndex returns paginated list of AuthClientTokens (AuthClientTokens + services combined)
func AuthClientTokenIndex(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (models.IndexData, any) {
	return models.GetIndexData(tx, data, c, locale, models.AuthClientTokenModel{})
}

// AuthClientTokenSingle returns single AuthClientToken by condition
func AuthClientTokenSingle(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, where func(db *gorm.DB) *gorm.DB) (map[string]any, any) {
	return models.GetSingleData(tx, data, c, locale, where, models.AuthClientTokenModel{})
}

// AuthClientTokenMultiple returns multiple AuthClientTokens by condition
func AuthClientTokenMultiple(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, where func(db *gorm.DB) *gorm.DB) ([]map[string]any, any) {
	return models.GetMultipleData(tx, data, c, locale, where, models.AuthClientTokenModel{})
}
