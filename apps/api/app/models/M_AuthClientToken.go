package models

import (
	"github.com/revianto/yava/api/helpers"
	"gorm.io/gorm"
)

// =============================================================================
// AuthClient MODEL
// =============================================================================

type AuthClientTokenModel struct {
	Id                   string `json:"id" gorm:"default:NULL"`
	Refresh_id           string `json:"refresh_id" gorm:"default:NULL"`
	Client_id            int    `json:"client_id" gorm:"default:"`
	User_id              int    `json:"user_id" gorm:"default:0"`
	Revoked              int    `json:"revoked" gorm:"default:0"`
	Expiary_time         string `json:"expiary_time" gorm:"default:NULL"`
	Refresh_expiary_time string `json:"refresh_expiary_time" gorm:"default:NULL"`
	Created_time         string `json:"created_time" gorm:"default:NULL"`
	Created_from         string `json:"created_from" gorm:"default:NULL"`
}

// =============================================================================
// TABLE & MODULE NAME
// =============================================================================

func (AuthClientTokenModel) TableName() string { return "tr_oauth_tokens" }
func (AuthClientTokenModel) ModulName() string { return "AuthClient Token" }

// =============================================================================
// SCOPES - SELECT, SEARCH, SORT
// =============================================================================

func (s AuthClientTokenModel) ScopesGetSelect(data map[string]any) map[string]string {
	cols := map[string]string{}
	get := func(key, def string) string { return helpers.Conv(data).GetMapValueAsString(key, def) }

	if get("show_id", "1") == "1" {
		cols["id"] = `tr_oauth_tokens.id`
	}
	if get("show_refresh_id", "1") == "1" {
		cols["refresh_id"] = `tr_oauth_tokens.refresh_id`
	}
	if get("show_client_id", "1") == "1" {
		cols["client_id"] = `tr_oauth_tokens.client_id`
	}
	if get("show_user_id", "1") == "1" {
		cols["user_id"] = `tr_oauth_tokens.user_id`
	}
	if get("show_revoked", "1") == "1" {
		cols["revoked"] = `tr_oauth_tokens.revoked`
	}
	if get("show_expiary_time", "1") == "1" {
		cols["expiary_time"] = `tr_oauth_tokens.expiary_time`
	}
	if get("show_refresh_expiary_time", "1") == "1" {
		cols["refresh_expiary_time"] = `tr_oauth_tokens.refresh_expiary_time`
	}

	return cols
}

func (s AuthClientTokenModel) ScopesSearchableFields(data map[string]any) map[string]SearchableFields {
	return map[string]SearchableFields{
		"id": {Operators: []string{"="}},
	}
}

func (s AuthClientTokenModel) ScopesSortbleFields(data map[string]any) map[string]bool {
	return map[string]bool{"id": true}
}

// =============================================================================
// SCOPES - JOIN & OPTIONS
// =============================================================================

func (s AuthClientTokenModel) ScopeJoin(data map[string]any) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB { return tx }
}

func (s AuthClientTokenModel) ScopeOption(data map[string]any) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx
	}
}

// =============================================================================
// GORM HOOKS
// =============================================================================

func (s *AuthClientTokenModel) BeforeCreate(tx *gorm.DB) error { return AutoFillCreate(s, tx) }
func (s *AuthClientTokenModel) BeforeUpdate(tx *gorm.DB) error { return AutoFillUpdate(s, tx) }
func (s *AuthClientTokenModel) BeforeDelete(tx *gorm.DB) error { return AutoFillDelete(s, tx) }
