package services

import (
	"github.com/revianto/yava/api/app/models"
	"github.com/revianto/yava/api/app/repositories"
	"gorm.io/gorm"
)

func GetRecipeTypes(db *gorm.DB) ([]models.RecipeType, error) {
	types, err := repositories.RecipeTypeList(db)
	if err != nil {
		return nil, &ServiceError{Code: 500, ErrCode: "DB_ERROR", Message: err.Error()}
	}
	return types, nil
}

func GetRecipeSubtypes(db *gorm.DB, typeId int64) ([]models.RecipeSubtype, error) {
	if typeId <= 0 {
		return nil, &ServiceError{Code: 400, ErrCode: "INVALID_ID", Message: "ID tidak valid"}
	}
	_, err := repositories.RecipeTypeById(db, typeId)
	if err != nil {
		return nil, &ServiceError{Code: 404, ErrCode: "TYPE_NOT_FOUND", Message: "Jenis resep tidak ditemukan"}
	}
	subs, err := repositories.RecipeSubtypesByTypeId(db, typeId)
	if err != nil {
		return nil, &ServiceError{Code: 500, ErrCode: "DB_ERROR", Message: err.Error()}
	}
	return subs, nil
}
