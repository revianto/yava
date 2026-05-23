package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/models"
	"github.com/revianto/yava/api/app/repositories"
	"github.com/revianto/yava/api/exceptions"
	"github.com/revianto/yava/api/helpers"
	"gorm.io/gorm"
)

// =============================================================================
// LIST & SHOW
// =============================================================================

func ExampleList(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (any, any) {
	return repositories.ExampleIndex(tx, data, c, locale)
}

func ExampleShow(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id any) (any, any) {
	return repositories.ExampleSingle(tx, data, c, locale, func(db *gorm.DB) *gorm.DB {
		return db.Where(models.ExampleModel{Id: helpers.Conv(id).Int()})
	})
}

// =============================================================================
// VALIDATION STRUCTS
// =============================================================================

type ValidateExampleCreate struct {
	Name string `json:"name" validate:"required"`
}

type ValidateExampleUpdate struct {
	Id   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type ValidateExampleDelete struct {
	Id string `json:"id" validate:"required"`
}

// =============================================================================
// CRUD OPERATIONS
// =============================================================================

func ExampleCreate(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (any, any) {
	if _, err := models.Validate(c, data, new(ValidateExampleCreate), locale); err != nil {
		return nil, err
	}

	m := models.ExampleModel{
		Name: helpers.Conv(data["name"]).String(),
	}

	if txErr := tx.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&m).Error
	}); txErr != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "failed to create "+m.ModulName())
	}

	return repositories.ExampleSingle(tx, data, c, locale, func(db *gorm.DB) *gorm.DB {
		return db.Where(models.ExampleModel{Id: m.Id})
	})
}

func ExampleUpdate(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, id any) (any, any) {
	data["id"] = id
	if _, err := models.Validate(c, data, new(ValidateExampleUpdate), locale); err != nil {
		return nil, err
	}

	if txErr := tx.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&models.ExampleModel{}).
			Where("id = ?", id).
			Update("name", helpers.Conv(data["name"]).String()).Error
	}); txErr != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "failed to update "+models.ExampleModel{}.ModulName())
	}

	return repositories.ExampleSingle(tx, data, c, locale, func(db *gorm.DB) *gorm.DB {
		return db.Where(models.ExampleModel{Id: helpers.Conv(id).Int()})
	})
}

func ExampleDelete(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (any, any) {
	if _, err := models.Validate(c, data, new(ValidateExampleDelete), locale); err != nil {
		return nil, err
	}

	m := models.ExampleModel{Id: helpers.Conv(data["id"]).Int()}

	if txErr := tx.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&m).Error
	}); txErr != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "failed to delete "+m.ModulName())
	}

	return nil, nil
}
