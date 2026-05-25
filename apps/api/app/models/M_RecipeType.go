package models

import (
	"time"

	"gorm.io/gorm"
)

type RecipeType struct {
	Id        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Code      string    `json:"code" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"not null"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Subtypes []RecipeSubtype `json:"subtypes,omitempty" gorm:"foreignKey:TypeId"`
}

func (RecipeType) TableName() string { return "yv_cd_recipe_type" }
func (RecipeType) ModulName() string { return "RecipeType" }

func (s RecipeType) ScopesGetSelect(data map[string]any) map[string]string {
	return map[string]string{
		"id":         "yv_cd_recipe_type.id",
		"code":       "yv_cd_recipe_type.code",
		"name":       "yv_cd_recipe_type.name",
		"sort_order": "yv_cd_recipe_type.sort_order",
	}
}

func (s RecipeType) ScopesSearchableFields(data map[string]any) map[string]SearchableFields {
	return map[string]SearchableFields{
		"id":   {Operators: []string{"=", "!="}},
		"code": {Operators: []string{"=", "like"}},
	}
}

func (s RecipeType) ScopesSortbleFields(data map[string]any) map[string]bool {
	return map[string]bool{"id": true, "sort_order": true}
}

func (s RecipeType) ScopeJoin(data map[string]any) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB { return tx }
}

func (s RecipeType) ScopeOption(data map[string]any) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Order("sort_order ASC")
	}
}

type RecipeSubtype struct {
	Id        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TypeId    int64     `json:"type_id" gorm:"not null;index"`
	Code      string    `json:"code" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (RecipeSubtype) TableName() string { return "yv_cd_recipe_subtype" }
func (RecipeSubtype) ModulName() string { return "RecipeSubtype" }

func (s RecipeSubtype) ScopesGetSelect(data map[string]any) map[string]string {
	return map[string]string{
		"id":         "yv_cd_recipe_subtype.id",
		"type_id":    "yv_cd_recipe_subtype.type_id",
		"code":       "yv_cd_recipe_subtype.code",
		"name":       "yv_cd_recipe_subtype.name",
		"sort_order": "yv_cd_recipe_subtype.sort_order",
	}
}

func (s RecipeSubtype) ScopesSearchableFields(data map[string]any) map[string]SearchableFields {
	return map[string]SearchableFields{
		"id":      {Operators: []string{"=", "!="}},
		"type_id": {Operators: []string{"="}},
	}
}

func (s RecipeSubtype) ScopesSortbleFields(data map[string]any) map[string]bool {
	return map[string]bool{"id": true, "sort_order": true}
}

func (s RecipeSubtype) ScopeJoin(data map[string]any) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB { return tx }
}

func (s RecipeSubtype) ScopeOption(data map[string]any) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Order("sort_order ASC")
	}
}
