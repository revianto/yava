package repositories

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/models"
	"github.com/revianto/yava/api/exceptions"
	"gorm.io/gorm"
)

func YvUserSingle(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, where func(*gorm.DB) *gorm.DB) (map[string]any, any) {
	return models.GetSingleData(tx, data, c, locale, where, models.YvUser{})
}

func YvUserUpsert(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string) (map[string]any, any) {
	googleId, _ := data["google_id"].(string)
	email, _ := data["email"].(string)
	name, _ := data["name"].(string)
	avatarStr, _ := data["avatar_url"].(string)

	var avatarURL *string
	if avatarStr != "" {
		avatarURL = &avatarStr
	}

	var user models.YvUser
	err := tx.Where("google_id = ?", googleId).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		user = models.YvUser{GoogleId: googleId, Email: email, Name: name, AvatarUrl: avatarURL}
		if txErr := tx.Transaction(func(t *gorm.DB) error {
			return t.Create(&user).Error
		}); txErr != nil {
			return nil, exceptions.ErrorException(c, fiber.StatusInternalServerError, "Gagal menyimpan pengguna")
		}
	} else if err != nil {
		return nil, exceptions.ErrorException(c, fiber.StatusInternalServerError, "Gagal mengambil data pengguna")
	} else {
		if txErr := tx.Transaction(func(t *gorm.DB) error {
			return t.Model(&user).Updates(map[string]any{"name": name, "avatar_url": avatarURL}).Error
		}); txErr != nil {
			return nil, exceptions.ErrorException(c, fiber.StatusInternalServerError, "Gagal memperbarui pengguna")
		}
	}

	b, _ := json.Marshal(user)
	var m map[string]any
	json.Unmarshal(b, &m)
	return m, nil
}
