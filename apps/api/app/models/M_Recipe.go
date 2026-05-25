package models

import (
	"time"

	"gorm.io/gorm"
)

type Recipe struct {
	Id           int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	OwnerId      *int64         `json:"owner_id"`
	TypeId       int64          `json:"type_id" gorm:"not null;index"`
	SubtypeId    *int64         `json:"subtype_id"`
	Name         string         `json:"name" gorm:"not null"`
	Description  *string        `json:"description"`
	Visibility   string         `json:"visibility" gorm:"default:private;not null"`
	IsDefault    bool           `json:"is_default" gorm:"default:false;not null"`
	IsArchived   bool           `json:"is_archived" gorm:"default:false;not null"`
	SavesCount   int            `json:"saves_count" gorm:"default:0;not null"`
	ParamDose    *string        `json:"param_dose"`
	ParamYield   *string        `json:"param_yield"`
	ParamTemp    *string        `json:"param_temp"`
	ParamGrind   *string        `json:"param_grind"`
	ParamRatio   *string        `json:"param_ratio"`
	LastBrewedAt *time.Time     `json:"last_brewed_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	Type     *RecipeType    `json:"type,omitempty" gorm:"foreignKey:TypeId"`
	Subtype  *RecipeSubtype `json:"subtype,omitempty" gorm:"foreignKey:SubtypeId"`
	Owner    *YvUser        `json:"owner,omitempty" gorm:"foreignKey:OwnerId"`
	Sessions []RecipeSession `json:"sessions,omitempty" gorm:"foreignKey:RecipeId;references:Id"`
	Notes    []RecipeNote    `json:"notes,omitempty" gorm:"foreignKey:RecipeId;references:Id"`
}

func (Recipe) TableName() string { return "yv_recipe" }

type RecipeSession struct {
	Id          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	RecipeId    int64     `json:"recipe_id" gorm:"not null;index"`
	SortOrder   int       `json:"sort_order" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"`
	DurationSec int       `json:"duration_sec" gorm:"not null"`
	Note        *string   `json:"note"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (RecipeSession) TableName() string { return "yv_recipe_session" }

type RecipeNote struct {
	Id        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	RecipeId  int64     `json:"recipe_id" gorm:"not null;index"`
	SortOrder int       `json:"sort_order" gorm:"not null"`
	Content   string    `json:"content" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (RecipeNote) TableName() string { return "yv_recipe_note" }
