package models

import "gorm.io/gorm"

type ExampleModel struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Created_time string `json:"created_time"`
	Created_by   int    `json:"created_by"`
	Created_from string `json:"created_from"`
}

func (ExampleModel) TableName() string { return "tr_examples" }
func (ExampleModel) ModulName() string { return "Example" }

func (s ExampleModel) ScopesGetSelect(data map[string]any) map[string]string {
	return map[string]string{
		"id":           "tr_examples.id",
		"name":         "tr_examples.name",
		"created_time": "tr_examples.created_time",
	}
}

func (s ExampleModel) ScopesSearchableFields(data map[string]any) map[string]SearchableFields {
	return map[string]SearchableFields{
		"id":   {Operators: []string{"=", "!="}},
		"name": {Operators: []string{"=", "like"}},
	}
}

func (s ExampleModel) ScopesSortbleFields(data map[string]any) map[string]bool {
	return map[string]bool{"id": true, "name": true}
}

func (s ExampleModel) ScopeJoin(data map[string]any) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB { return tx }
}

func (s ExampleModel) ScopeOption(data map[string]any) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB { return tx }
}

func (s *ExampleModel) BeforeCreate(tx *gorm.DB) error { return AutoFillCreate(s, tx) }
func (s *ExampleModel) BeforeUpdate(tx *gorm.DB) error { return AutoFillUpdate(s, tx) }
func (s *ExampleModel) BeforeDelete(tx *gorm.DB) error { return AutoFillDelete(s, tx) }
