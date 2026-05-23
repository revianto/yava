package models

import (
	"time"

	"github.com/revianto/yava/api/helpers"
	"gorm.io/gorm"
)

// =============================================================================
// User MODEL
// =============================================================================

type UserModel struct {
	Id                    int64          `json:"id"`
	Name                  string         `json:"name"`
	Phone                 string         `json:"phone"`
	Password              string         `json:"password"`
	Role_id               int64          `json:"role_id"`
	Is_active             bool           `json:"is_active" gorm:"default:true"`
	Reset_code            *string        `json:"-"`
	Reset_code_expired_at *time.Time     `json:"-"`
	Created_time          time.Time      `json:"created_time"`
	Created_by            int            `json:"created_by"`
	Created_from          string         `json:"created_from"`
	Modified_time         time.Time      `json:"modified_time"`
	Modified_by           int            `json:"modified_by"`
	Modified_from         string         `json:"modified_from"`
	Deleted_time          gorm.DeletedAt `json:"-"`
	Deleted_by            int            `json:"-"`
	Deleted_from          string         `json:"-"`
}

// =============================================================================
// TABLE & MODULE NAME
// =============================================================================

func (UserModel) TableName() string { return "tr_users" }
func (UserModel) ModulName() string { return "User" }

// =============================================================================
// SCOPES - SELECT, SEARCH, SORT
// =============================================================================

func (s UserModel) ScopesGetSelect(data map[string]any) map[string]string {
	cols := map[string]string{}
	get := func(key, def string) string { return helpers.Conv(data).GetMapValueAsString(key, def) }

	if get("show_id", "1") == "1" {
		cols["id"] = `tr_users.id`
	}
	if get("show_name", "1") == "1" {
		cols["name"] = `tr_users.name`
	}
	if get("show_phone", "1") == "1" {
		cols["phone"] = `tr_users.phone`
	}
	if get("show_password", "1") == "1" {
		cols["password"] = `tr_users.password`
	}
	if get("show_role", "1") == "1" {
		cols["role"] = `(SELECT row_to_json(r) FROM (
			SELECT roles.id, roles.code, roles.name
			FROM tr_cd_roles roles
			WHERE roles.id = tr_users.role_id
				AND roles.deleted_time IS NULL
			LIMIT 1
		) r)`
		cols["role_code"] = `tr_cd_roles.code`
		cols["role_name"] = `tr_cd_roles.name`
	}
	if get("show_is_active", "1") == "1" {
		cols["is_active"] = `tr_users.is_active`
	}

	return cols
}

func (s UserModel) ScopesSearchableFields(data map[string]any) map[string]SearchableFields {
	return map[string]SearchableFields{
		"id":        {Operators: []string{"=", "!="}},
		"name":      {Operators: []string{"like", "=", "!="}},
		"phone":     {Operators: []string{"like", "=", "!="}},
		"role_id":   {Operators: []string{"=", "!="}},
		"role_code": {Operators: []string{"=", "!="}},
		"is_active": {Operators: []string{"=", "!="}},
	}
}

func (s UserModel) ScopesSortbleFields(data map[string]any) map[string]bool {
	return map[string]bool{"id": true, "name": true}
}

// =============================================================================
// SCOPES - JOIN & OPTIONS
// =============================================================================

func (s UserModel) ScopeJoin(data map[string]any) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Joins("LEFT JOIN tr_cd_roles ON tr_cd_roles.id = tr_users.role_id AND tr_cd_roles.deleted_time IS NULL")
	}
}

func (s UserModel) ScopeOption(data map[string]any) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("tr_users.deleted_time IS NULL")
	}
}

// =============================================================================
// GORM HOOKS
// =============================================================================

func (s *UserModel) BeforeCreate(tx *gorm.DB) error { return AutoFillCreate(s, tx) }
func (s *UserModel) BeforeUpdate(tx *gorm.DB) error { return AutoFillUpdate(s, tx) }
func (s *UserModel) BeforeDelete(tx *gorm.DB) error { return AutoFillDelete(s, tx) }
