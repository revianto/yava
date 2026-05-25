package models

import (
	"time"

	"gorm.io/gorm"
)

type YvUser struct {
	Id        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	GoogleId  string    `json:"google_id" gorm:"uniqueIndex;not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"not null"`
	AvatarUrl *string   `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (YvUser) TableName() string { return "yv_user" }
func (YvUser) ModulName() string { return "YvUser" }

func (s YvUser) ScopesGetSelect(data map[string]any) map[string]string {
	return map[string]string{
		"id":         "yv_user.id",
		"email":      "yv_user.email",
		"name":       "yv_user.name",
		"avatar_url": "yv_user.avatar_url",
	}
}

func (s YvUser) ScopesSearchableFields(data map[string]any) map[string]SearchableFields {
	return map[string]SearchableFields{
		"id":    {Operators: []string{"=", "!="}},
		"email": {Operators: []string{"=", "like"}},
	}
}

func (s YvUser) ScopesSortbleFields(data map[string]any) map[string]bool {
	return map[string]bool{"id": true, "name": true}
}

func (s YvUser) ScopeJoin(data map[string]any) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB { return tx }
}

func (s YvUser) ScopeOption(data map[string]any) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB { return tx }
}
