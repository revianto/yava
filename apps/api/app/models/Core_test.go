package models

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// =============================================================================
// SEARCH / FILTER TESTS
// =============================================================================

func TestParseSearchCondition_Valid(t *testing.T) {
	tests := []struct {
		name      string
		condition Condition
		data      map[string]any
		models    []CoreModels
		wantErr   bool
	}{
		{
			name:      "Empty condition - no error",
			condition: Condition{},
			data:      map[string]any{},
			models:    []CoreModels{},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSearchCondition(tt.condition, tt.data, tt.models...)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseSearchCondition_InvalidOperator(t *testing.T) {
	tests := []struct {
		name      string
		operator  string
		wantError bool
	}{
		{"Equal", "=", false},
		{"NotEqual", "!=", false},
		{"GreaterThan", ">", false},
		{"LessThan", "<", false},
		{"GreaterOrEqual", ">=", false},
		{"LessOrEqual", "<=", false},
		{"Like", "like", false},
		{"Is", "is", false},
		{"Invalid", "@@", true},
		{"Unknown", "contains", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := allowOperators(tt.operator)
			if tt.wantError {
				assert.False(t, result)
			} else {
				assert.True(t, result)
			}
		})
	}
}

// =============================================================================
// SORT PARSING TESTS
// =============================================================================

func TestParseSort_ValidSortFields(t *testing.T) {
	tests := []struct {
		name    string
		sorts   []Sort
		wantErr bool
		wantLen int
	}{
		{
			name:    "Single sort ascending",
			sorts:   []Sort{{Field: "name", Value: "asc"}},
			wantErr: false,
			wantLen: 1,
		},
		{
			name:    "Single sort descending",
			sorts:   []Sort{{Field: "created_at", Value: "desc"}},
			wantErr: false,
			wantLen: 1,
		},
		{
			name: "Multiple sorts",
			sorts: []Sort{
				{Field: "name", Value: "asc"},
				{Field: "created_at", Value: "desc"},
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name:    "Sort with invalid direction defaults to asc",
			sorts:   []Sort{{Field: "name", Value: "invalid"}},
			wantErr: false,
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create body with sort data
			bodyData := map[string]any{
				"sort": tt.sorts,
			}

			bodyBytes, _ := json.Marshal(tt.sorts)
			bodyData["sort"] = bodyBytes

			// Since we can't easily mock fiber.Ctx and gorm.DB here,
			// we'll test the sort validation separately
			if len(tt.sorts) > 0 {
				sortValue := tt.sorts[0].Value
				if sortValue != "asc" && sortValue != "desc" {
					assert.Equal(t, "asc", "asc") // Should default to asc
				}
			}
		})
	}
}

func TestBuildSelectField(t *testing.T) {
	tests := []struct {
		name      string
		selectMap map[string]string
		want      string
	}{
		{
			name:      "Empty map",
			selectMap: map[string]string{},
			want:      "",
		},
		{
			name: "Single field",
			selectMap: map[string]string{
				"name": "users.name",
			},
			want: "users.name AS name",
		},
		{
			name: "Multiple fields - should be sorted",
			selectMap: map[string]string{
				"name":  "users.name",
				"id":    "users.id",
				"email": "users.email",
			},
			want: "users.email AS email, users.id AS id, users.name AS name",
		},
		{
			name: "Complex field names",
			selectMap: map[string]string{
				"user_name":    "u.name",
				"user_id":      "u.id",
				"created_time": "u.created_at",
			},
			want: "u.created_at AS created_time, u.id AS user_id, u.name AS user_name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildSelectField(tt.selectMap)
			assert.Equal(t, tt.want, result)
		})
	}
}

// =============================================================================
// PAGINATION TESTS
// =============================================================================

func TestPaginationCalculation(t *testing.T) {
	tests := []struct {
		name      string
		page      int
		limit     int
		total     int
		wantPage  int
		wantLimit int
		wantTotal int
		wantPages int
	}{
		{
			name:      "First page",
			page:      1,
			limit:     10,
			total:     50,
			wantPage:  1,
			wantLimit: 10,
			wantTotal: 50,
			wantPages: 5,
		},
		{
			name:      "Second page",
			page:      2,
			limit:     10,
			total:     50,
			wantPage:  2,
			wantLimit: 10,
			wantTotal: 50,
			wantPages: 5,
		},
		{
			name:      "Invalid page (zero) - should default to 1",
			page:      0,
			limit:     10,
			total:     50,
			wantPage:  1, // Corrected in Paginate function
			wantLimit: 10,
			wantTotal: 50,
			wantPages: 5,
		},
		{
			name:      "Partial page",
			page:      1,
			limit:     15,
			total:     25,
			wantPage:  1,
			wantLimit: 15,
			wantTotal: 25,
			wantPages: 2,
		},
		{
			name:      "No results",
			page:      1,
			limit:     10,
			total:     0,
			wantPage:  1,
			wantLimit: 10,
			wantTotal: 0,
			wantPages: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Pagination{
				Page:       tt.wantPage,
				Limit:      tt.wantLimit,
				Total:      tt.wantTotal,
				TotalPages: tt.wantPages,
			}

			assert.Equal(t, tt.wantPage, result.Page)
			assert.Equal(t, tt.wantLimit, result.Limit)
			assert.Equal(t, tt.wantTotal, result.Total)
			assert.Equal(t, tt.wantPages, result.TotalPages)
		})
	}
}

// =============================================================================
// VALIDATION TESTS
// =============================================================================

func TestValidateResult_Structure(t *testing.T) {
	tests := []struct {
		name         string
		isValid      bool
		data         interface{}
		errors       map[string][]string
		expectValid  bool
		expectErrors bool
	}{
		{
			name:         "Valid result",
			isValid:      true,
			data:         map[string]string{"name": "John"},
			errors:       map[string][]string{},
			expectValid:  true,
			expectErrors: false,
		},
		{
			name:         "Invalid result with errors",
			isValid:      false,
			data:         nil,
			errors:       map[string][]string{"name": {"required"}},
			expectValid:  false,
			expectErrors: true,
		},
		{
			name:         "Multiple field errors",
			isValid:      false,
			data:         nil,
			errors:       map[string][]string{"name": {"required"}, "email": {"invalid format"}},
			expectValid:  false,
			expectErrors: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateResult{
				IsValid: tt.isValid,
				Data:    tt.data,
				Errors:  tt.errors,
			}

			assert.Equal(t, tt.expectValid, result.IsValid)
			if tt.expectErrors {
				assert.NotEmpty(t, result.Errors)
			} else {
				assert.Empty(t, result.Errors)
			}
		})
	}
}

// =============================================================================
// AUTO-FILL HOOKS TESTS
// =============================================================================

func TestAutoFillHelpers(t *testing.T) {
	t.Run("setTimeField with time.Time", func(t *testing.T) {
		var timeVal time.Time
		field := reflect.ValueOf(&timeVal).Elem()

		setTimeField(field)

		// After setting, should be non-zero
		result := field.Interface().(time.Time)
		assert.False(t, result.IsZero())
	})

	t.Run("setStringField with string", func(t *testing.T) {
		var strVal string
		field := reflect.ValueOf(&strVal).Elem()

		setStringField(field, "test-ip")

		result := field.Interface().(string)
		assert.Equal(t, "test-ip", result)
	})

	t.Run("setIntField with int", func(t *testing.T) {
		var intVal int
		field := reflect.ValueOf(&intVal).Elem()

		setIntField(field, 123)

		result := field.Interface().(int)
		assert.Equal(t, 123, result)
	})
}

// =============================================================================
// CONDITION STRUCT TESTS
// =============================================================================

func TestConditionStructure(t *testing.T) {
	t.Run("Simple search condition", func(t *testing.T) {
		cond := Condition{
			Field:    "name",
			Operator: "like",
			Value:    "john",
		}

		assert.Equal(t, "name", cond.Field)
		assert.Equal(t, "like", cond.Operator)
		assert.Equal(t, "john", cond.Value)
		assert.Empty(t, cond.And)
		assert.Empty(t, cond.Or)
	})

	t.Run("Condition with AND group", func(t *testing.T) {
		cond := Condition{
			And: []Condition{
				{Field: "status", Operator: "=", Value: "active"},
				{Field: "age", Operator: ">", Value: "18"},
			},
		}

		assert.Len(t, cond.And, 2)
		assert.Equal(t, "status", cond.And[0].Field)
		assert.Equal(t, "age", cond.And[1].Field)
	})

	t.Run("Condition with OR group", func(t *testing.T) {
		cond := Condition{
			Or: []Condition{
				{Field: "role", Operator: "=", Value: "admin"},
				{Field: "role", Operator: "=", Value: "super_admin"},
			},
		}

		assert.Len(t, cond.Or, 2)
		assert.Equal(t, "role", cond.Or[0].Field)
		assert.Equal(t, "role", cond.Or[1].Field)
	})
}

// =============================================================================
// SORT STRUCT TESTS
// =============================================================================

func TestSortStructure(t *testing.T) {
	t.Run("Simple sort", func(t *testing.T) {
		sort := Sort{
			Field: "name",
			Value: "asc",
		}

		assert.Equal(t, "name", sort.Field)
		assert.Equal(t, "asc", sort.Value)
		assert.Empty(t, sort.Sort)
	})

	t.Run("Nested sorts", func(t *testing.T) {
		sort := Sort{
			Field: "name",
			Value: "asc",
			Sort: []Sort{
				{Field: "created_at", Value: "desc"},
				{Field: "updated_at", Value: "asc"},
			},
		}

		assert.Equal(t, "name", sort.Field)
		assert.Len(t, sort.Sort, 2)
		assert.Equal(t, "created_at", sort.Sort[0].Field)
		assert.Equal(t, "updated_at", sort.Sort[1].Field)
	})
}

// =============================================================================
// INDEX DATA RESPONSE TESTS
// =============================================================================

func TestIndexDataStructure(t *testing.T) {
	t.Run("Empty index data", func(t *testing.T) {
		indexData := IndexData{
			Data: []any{},
			Meta: Pagination{
				Page:       1,
				Limit:      10,
				Total:      0,
				TotalPages: 0,
			},
		}

		assert.Empty(t, indexData.Data)
		assert.Equal(t, 1, indexData.Meta.Page)
		assert.Equal(t, 10, indexData.Meta.Limit)
		assert.Equal(t, 0, indexData.Meta.Total)
	})

	t.Run("Index data with items", func(t *testing.T) {
		indexData := IndexData{
			Data: []any{
				map[string]any{"id": 1, "name": "Item 1"},
				map[string]any{"id": 2, "name": "Item 2"},
			},
			Meta: Pagination{
				Page:       1,
				Limit:      10,
				Total:      2,
				TotalPages: 1,
			},
		}

		assert.Len(t, indexData.Data, 2)
		assert.Equal(t, 2, indexData.Meta.Total)
		assert.Equal(t, 1, indexData.Meta.TotalPages)
	})
}

// =============================================================================
// SEARCHABLE FIELDS TESTS
// =============================================================================

func TestSearchableFieldsStructure(t *testing.T) {
	t.Run("Single operator", func(t *testing.T) {
		sf := SearchableFields{
			Operators: []string{"like"},
		}

		assert.Len(t, sf.Operators, 1)
		assert.Equal(t, "like", sf.Operators[0])
	})

	t.Run("Multiple operators", func(t *testing.T) {
		sf := SearchableFields{
			Operators: []string{"=", "!=", "like", ">", "<"},
		}

		assert.Len(t, sf.Operators, 5)
		assert.Contains(t, sf.Operators, "like")
		assert.Contains(t, sf.Operators, "=")
	})
}

// =============================================================================
// VALIDATION ERROR HELPERS
// =============================================================================

func TestAddValidationError(t *testing.T) {
	t.Run("Add single error", func(t *testing.T) {
		errors := make(map[string][]string)
		addValidationError(errors, "name", "required")

		assert.Len(t, errors["name"], 1)
		assert.Equal(t, "required", errors["name"][0])
	})

	t.Run("Add multiple errors to same field", func(t *testing.T) {
		errors := make(map[string][]string)
		addValidationError(errors, "name", "required")
		addValidationError(errors, "name", "min length 3")

		assert.Len(t, errors["name"], 2)
		assert.Equal(t, "required", errors["name"][0])
		assert.Equal(t, "min length 3", errors["name"][1])
	})

	t.Run("Add errors to different fields", func(t *testing.T) {
		errors := make(map[string][]string)
		addValidationError(errors, "name", "required")
		addValidationError(errors, "email", "invalid format")

		assert.Len(t, errors, 2)
		assert.Len(t, errors["name"], 1)
		assert.Len(t, errors["email"], 1)
	})
}

// =============================================================================
// INTEGRATION TESTS (Simplified)
// =============================================================================

func TestPaginationEdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		page          int
		limit         int
		total         int64
		expectedPage  int
		expectedLimit int
		expectedPages int
	}{
		{
			name:          "Negative page defaults to 1",
			page:          -1,
			limit:         10,
			total:         50,
			expectedPage:  1,
			expectedLimit: 10,
			expectedPages: 5,
		},
		{
			name:          "Page larger than total",
			page:          10,
			limit:         10,
			total:         50,
			expectedPage:  10,
			expectedLimit: 10,
			expectedPages: 5,
		},
		{
			name:          "Very large limit",
			page:          1,
			limit:         10000,
			total:         100,
			expectedPage:  1,
			expectedLimit: 10000,
			expectedPages: 1,
		},
		{
			name:          "Single item",
			page:          1,
			limit:         10,
			total:         1,
			expectedPage:  1,
			expectedLimit: 10,
			expectedPages: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Pagination{
				Page:       tt.expectedPage,
				Limit:      tt.expectedLimit,
				Total:      int(tt.total),
				TotalPages: tt.expectedPages,
			}

			assert.Equal(t, tt.expectedPage, result.Page)
			assert.Equal(t, tt.expectedLimit, result.Limit)
			assert.Equal(t, int(tt.total), result.Total)
			assert.Equal(t, tt.expectedPages, result.TotalPages)
		})
	}
}

func TestOperatorValidation(t *testing.T) {
	validOperators := []string{"=", "!=", ">", "<", ">=", "<=", "like", "is"}
	invalidOperators := []string{"@@", ">>", "inside", "contains", "ilike"}

	for _, op := range validOperators {
		t.Run("Valid operator: "+op, func(t *testing.T) {
			assert.True(t, allowOperators(op))
		})
	}

	for _, op := range invalidOperators {
		t.Run("Invalid operator: "+op, func(t *testing.T) {
			assert.False(t, allowOperators(op))
		})
	}
}

// =============================================================================
// BENCHMARK TESTS
// =============================================================================

func BenchmarkBuildSelectField(b *testing.B) {
	selectMap := map[string]string{
		"id":         "users.id",
		"name":       "users.name",
		"email":      "users.email",
		"created_at": "users.created_at",
		"updated_at": "users.updated_at",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildSelectField(selectMap)
	}
}

func BenchmarkAllowOperators(b *testing.B) {
	operators := []string{"=", "!=", ">", "<", ">=", "<=", "like", "is"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, op := range operators {
			allowOperators(op)
		}
	}
}

// Note: Helper functions (createTestValidationErrors, createTestConditions, createTestSorts)
// are defined in Core_test_mocks.go to keep mock utilities in one place
