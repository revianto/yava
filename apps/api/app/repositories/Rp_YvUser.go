package repositories

import (
	"github.com/revianto/yava/api/app/models"
	"gorm.io/gorm"
)

func YvUserFindByGoogleId(db *gorm.DB, googleId string) (*models.YvUser, error) {
	var user models.YvUser
	err := db.Where("google_id = ?", googleId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func YvUserFindById(db *gorm.DB, id int64) (*models.YvUser, error) {
	var user models.YvUser
	err := db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func YvUserUpsert(db *gorm.DB, googleId, email, name string, avatarURL *string) (*models.YvUser, error) {
	var user models.YvUser
	err := db.Where("google_id = ?", googleId).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		user = models.YvUser{
			GoogleId:  googleId,
			Email:     email,
			Name:      name,
			AvatarUrl: avatarURL,
		}
		if txErr := db.Transaction(func(tx *gorm.DB) error {
			return tx.Create(&user).Error
		}); txErr != nil {
			return nil, txErr
		}
		return &user, nil
	}
	if err != nil {
		return nil, err
	}
	// Update name and avatar if changed
	if txErr := db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&user).Updates(map[string]interface{}{
			"name":       name,
			"avatar_url": avatarURL,
		}).Error
	}); txErr != nil {
		return nil, txErr
	}
	return &user, nil
}
