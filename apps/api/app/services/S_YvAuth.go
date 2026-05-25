package services

import (
	"github.com/revianto/yava/api/app/models"
	"github.com/revianto/yava/api/app/repositories"
	"github.com/revianto/yava/api/helpers"
	"gorm.io/gorm"
)

type AuthResult struct {
	User  *models.YvUser
	Token string
}

func AuthProcessGoogleUser(db *gorm.DB, googleId, email, name string, avatarURL *string) (*AuthResult, error) {
	if googleId == "" || email == "" {
		return nil, &ServiceError{Code: 400, ErrCode: "INVALID_GOOGLE_USER", Message: "Data pengguna Google tidak lengkap"}
	}

	user, err := repositories.YvUserUpsert(db, googleId, email, name, avatarURL)
	if err != nil {
		return nil, &ServiceError{Code: 500, ErrCode: "USER_UPSERT_FAILED", Message: "Gagal menyimpan data pengguna"}
	}

	avatarStr := ""
	if user.AvatarUrl != nil {
		avatarStr = *user.AvatarUrl
	}
	token, err := helpers.YvCreateToken(user.Id, user.Email, user.Name, avatarStr)
	if err != nil {
		return nil, &ServiceError{Code: 500, ErrCode: "TOKEN_FAILED", Message: "Gagal membuat token"}
	}

	return &AuthResult{User: user, Token: token}, nil
}

func AuthGetMe(db *gorm.DB, userID int64) (*models.YvUser, error) {
	user, err := repositories.YvUserFindById(db, userID)
	if err != nil {
		return nil, &ServiceError{Code: 404, ErrCode: "USER_NOT_FOUND", Message: "Pengguna tidak ditemukan"}
	}
	return user, nil
}
