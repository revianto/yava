package services

import (
	"github.com/revianto/yava/api/app/models"
	"github.com/revianto/yava/api/app/repositories"
	"gorm.io/gorm"
)

type RecipeListInput struct {
	OwnerId    *int64
	Visibility *string
	TypeId     *int64
	Page       int
	Limit      int
}

func ListRecipes(db *gorm.DB, input RecipeListInput) ([]models.Recipe, int64, error) {
	if input.Page < 1 {
		input.Page = 1
	}
	if input.Limit < 1 || input.Limit > 100 {
		input.Limit = 20
	}
	recipes, total, err := repositories.RecipeList(db, repositories.RecipeListParams{
		OwnerId:    input.OwnerId,
		Visibility: input.Visibility,
		TypeId:     input.TypeId,
		Page:       input.Page,
		Limit:      input.Limit,
	})
	if err != nil {
		return nil, 0, &ServiceError{Code: 500, ErrCode: "DB_ERROR", Message: err.Error()}
	}
	return recipes, total, nil
}

func GetRecipe(db *gorm.DB, id int64) (*models.Recipe, error) {
	if id <= 0 {
		return nil, &ServiceError{Code: 400, ErrCode: "INVALID_ID", Message: "ID tidak valid"}
	}
	recipe, err := repositories.RecipeById(db, id)
	if err != nil {
		return nil, &ServiceError{Code: 404, ErrCode: "RECIPE_NOT_FOUND", Message: "Resep tidak ditemukan"}
	}
	return recipe, nil
}

type RecipeCreateInput struct {
	TypeId      int64
	SubtypeId   *int64
	Name        string
	Description *string
	Visibility  string
	ParamDose   *string
	ParamYield  *string
	ParamTemp   *string
	ParamGrind  *string
	ParamRatio  *string
	Sessions    []repositories.RecipeStepInput
	Notes       []repositories.RecipeNoteInput
}

func CreateRecipe(db *gorm.DB, userID int64, input RecipeCreateInput) (*models.Recipe, error) {
	if input.Name == "" {
		return nil, &ServiceError{Code: 422, ErrCode: "VALIDATION_ERROR", Message: "Nama resep wajib diisi"}
	}
	if input.TypeId <= 0 {
		return nil, &ServiceError{Code: 422, ErrCode: "VALIDATION_ERROR", Message: "Jenis resep wajib dipilih"}
	}
	if input.Visibility == "" {
		input.Visibility = "private"
	}
	if input.Visibility != "private" && input.Visibility != "public" && input.Visibility != "group" {
		return nil, &ServiceError{Code: 422, ErrCode: "VALIDATION_ERROR", Message: "Visibility tidak valid: private, public, atau group"}
	}

	recipe, err := repositories.RecipeCreate(db, repositories.RecipeCreateInput{
		OwnerId:     userID,
		TypeId:      input.TypeId,
		SubtypeId:   input.SubtypeId,
		Name:        input.Name,
		Description: input.Description,
		Visibility:  input.Visibility,
		ParamDose:   input.ParamDose,
		ParamYield:  input.ParamYield,
		ParamTemp:   input.ParamTemp,
		ParamGrind:  input.ParamGrind,
		ParamRatio:  input.ParamRatio,
		Sessions:    input.Sessions,
		Notes:       input.Notes,
	})
	if err != nil {
		return nil, &ServiceError{Code: 500, ErrCode: "CREATE_FAILED", Message: err.Error()}
	}
	return recipe, nil
}

type RecipeUpdateInput struct {
	TypeId      *int64
	SubtypeId   *int64
	Name        *string
	Description *string
	Visibility  *string
	ParamDose   *string
	ParamYield  *string
	ParamTemp   *string
	ParamGrind  *string
	ParamRatio  *string
	Sessions    *[]repositories.RecipeStepInput
	Notes       *[]repositories.RecipeNoteInput
}

func UpdateRecipe(db *gorm.DB, id, userID int64, input RecipeUpdateInput) (*models.Recipe, error) {
	recipe, err := repositories.RecipeById(db, id)
	if err != nil {
		return nil, &ServiceError{Code: 404, ErrCode: "RECIPE_NOT_FOUND", Message: "Resep tidak ditemukan"}
	}
	if recipe.IsDefault {
		return nil, &ServiceError{Code: 403, ErrCode: "FORBIDDEN", Message: "Resep default tidak bisa diubah"}
	}
	if recipe.OwnerId == nil || *recipe.OwnerId != userID {
		return nil, &ServiceError{Code: 403, ErrCode: "FORBIDDEN", Message: "Kamu bukan pemilik resep ini"}
	}
	if input.Visibility != nil {
		v := *input.Visibility
		if v != "private" && v != "public" && v != "group" {
			return nil, &ServiceError{Code: 422, ErrCode: "VALIDATION_ERROR", Message: "Visibility tidak valid"}
		}
	}

	updated, err := repositories.RecipeUpdate(db, id, userID, repositories.RecipeUpdateInput{
		TypeId: input.TypeId, SubtypeId: input.SubtypeId,
		Name: input.Name, Description: input.Description, Visibility: input.Visibility,
		ParamDose: input.ParamDose, ParamYield: input.ParamYield,
		ParamTemp: input.ParamTemp, ParamGrind: input.ParamGrind, ParamRatio: input.ParamRatio,
		Sessions: input.Sessions, Notes: input.Notes,
	})
	if err != nil {
		return nil, &ServiceError{Code: 500, ErrCode: "UPDATE_FAILED", Message: err.Error()}
	}
	return updated, nil
}
