package repositories

import (
	"github.com/revianto/yava/api/app/models"
	"gorm.io/gorm"
)

func RecipeTypeList(db *gorm.DB) ([]models.RecipeType, error) {
	var types []models.RecipeType
	err := db.Order("sort_order ASC").Find(&types).Error
	return types, err
}

func RecipeTypeById(db *gorm.DB, id int64) (*models.RecipeType, error) {
	var t models.RecipeType
	err := db.First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func RecipeSubtypesByTypeId(db *gorm.DB, typeId int64) ([]models.RecipeSubtype, error) {
	var subs []models.RecipeSubtype
	err := db.Where("type_id = ?", typeId).Order("sort_order ASC").Find(&subs).Error
	return subs, err
}
