package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/repositories"
	"github.com/revianto/yava/api/exceptions"
	"github.com/revianto/yava/api/helpers"
	"gorm.io/gorm"
)

func AuthProcessGoogleUser(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (any, any) {
	googleId, _ := data["google_id"].(string)
	email, _ := data["email"].(string)
	if googleId == "" || email == "" {
		return nil, exceptions.ErrorException(c, fiber.StatusBadRequest, "Data pengguna Google tidak lengkap")
	}

	userMap, err := repositories.YvUserUpsert(tx, data, c, locale)
	if err != nil {
		return nil, err
	}

	avatarStr := helpers.Conv(userMap["avatar_url"]).String()
	token, tokenErr := helpers.YvCreateToken(
		helpers.Conv(userMap["id"]).Int64(),
		helpers.Conv(userMap["email"]).String(),
		helpers.Conv(userMap["name"]).String(),
		avatarStr,
	)
	if tokenErr != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusInternalServerError, "Gagal membuat token")
	}

	userMap["token"] = token
	return userMap, nil
}

func AuthMe(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (any, any) {
	uid, _ := c.Locals("yv_user_id").(int64)
	return repositories.YvUserSingle(tx, data, c, locale, func(db *gorm.DB) *gorm.DB {
		return db.Where("yv_user.id = ?", uid)
	})
}
