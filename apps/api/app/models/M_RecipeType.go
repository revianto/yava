package models

import "time"

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
