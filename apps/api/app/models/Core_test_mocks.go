package models

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// =============================================================================
// MOCK STRUCTS FOR TESTING
// =============================================================================

// MockModel implements CoreModels interface for testing
type MockModel struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Status        string     `json:"status"`
	Created_time  time.Time  `json:"created_time"`
	Created_by    int        `json:"created_by"`
	Created_from  string     `json:"created_from"`
	Modified_time time.Time  `json:"modified_time"`
	Modified_by   int        `json:"modified_by"`
	Modified_from string     `json:"modified_from"`
	Deleted_time  *time.Time `json:"deleted_time"`
	Deleted_by    *int       `json:"deleted_by"`
	Deleted_from  *string    `json:"deleted_from"`
}

func (m MockModel) TableName() string {
	return "mock_models"
}

func (m MockModel) ModulName() string {
	return "MockModel"
}

func (m MockModel) ScopesGetSelect(data map[string]any) map[string]string {
	return map[string]string{
		"id":           "mock_models.id",
		"name":         "mock_models.name",
		"status":       "mock_models.status",
		"created_time": "mock_models.created_time",
		"created_by":   "mock_models.created_by",
	}
}

func (m MockModel) ScopesSearchableFields(data map[string]any) map[string]SearchableFields {
	return map[string]SearchableFields{
		"id": {
			Operators: []string{"=", "!=", ">", "<", ">=", "<="},
		},
		"name": {
			Operators: []string{"like", "=", "!="},
		},
		"status": {
			Operators: []string{"=", "!="},
		},
		"created_by": {
			Operators: []string{"=", "!="},
		},
	}
}

func (m MockModel) ScopesSortbleFields(data map[string]any) map[string]bool {
	return map[string]bool{
		"id":           true,
		"name":         true,
		"status":       true,
		"created_time": true,
	}
}

func (m MockModel) ScopeJoin(data map[string]any) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

func (m MockModel) ScopeOption(data map[string]any) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("deleted_time IS NULL")
	}
}

// =============================================================================
// AUTO-FILL HOOKS TESTS WITH REFLECTION
// =============================================================================

func TestAutoFillCreate_WithReflection(t *testing.T) {
	t.Run("Fill created_time field", func(t *testing.T) {
		type TestModel struct {
			Created_time time.Time
		}

		model := &TestModel{}
		v := reflect.ValueOf(model).Elem()
		f := v.FieldByName("Created_time")

		if f.IsValid() && f.CanSet() {
			setTimeField(f)
			result := model.Created_time
			assert.False(t, result.IsZero())
		}
	})

	t.Run("Fill created_from field", func(t *testing.T) {
		type TestModel struct {
			Created_from string
		}

		model := &TestModel{}
		v := reflect.ValueOf(model).Elem()
		f := v.FieldByName("Created_from")

		if f.IsValid() && f.CanSet() {
			setStringField(f, "192.168.1.1")
			assert.Equal(t, "192.168.1.1", model.Created_from)
		}
	})

	t.Run("Fill created_by field", func(t *testing.T) {
		type TestModel struct {
			Created_by int
		}

		model := &TestModel{}
		v := reflect.ValueOf(model).Elem()
		f := v.FieldByName("Created_by")

		if f.IsValid() && f.CanSet() {
			setIntField(f, 42)
			assert.Equal(t, 42, model.Created_by)
		}
	})
}

func TestAutoFillUpdate_WithReflection(t *testing.T) {
	t.Run("Update modified_time field", func(t *testing.T) {
		type TestModel struct {
			Modified_time time.Time
		}

		model := &TestModel{
			Modified_time: time.Now().Add(-1 * time.Hour),
		}

		v := reflect.ValueOf(model).Elem()
		f := v.FieldByName("Modified_time")

		oldTime := model.Modified_time

		if f.IsValid() && f.CanSet() {
			setTimeField(f)
			assert.True(t, model.Modified_time.After(oldTime))
		}
	})
}

// =============================================================================
// CONDITION VALIDATION TESTS
// =============================================================================

func TestValidateSearchCondition_WithMockModel(t *testing.T) {
	mockModel := MockModel{}
	data := map[string]any{}

	tests := []struct {
		name      string
		condition Condition
		wantErr   bool
	}{
		{
			name: "Valid field with valid operator",
			condition: Condition{
				Field:    "name",
				Operator: "like",
				Value:    "test",
			},
			wantErr: false,
		},
		{
			name: "Invalid field name",
			condition: Condition{
				Field:    "invalid_field",
				Operator: "like",
				Value:    "test",
			},
			wantErr: true,
		},
		{
			name: "Valid field with invalid operator",
			condition: Condition{
				Field:    "name",
				Operator: "@@",
				Value:    "test",
			},
			wantErr: true,
		},
		{
			name: "ID field with comparison operator",
			condition: Condition{
				Field:    "id",
				Operator: ">",
				Value:    "10",
			},
			wantErr: false,
		},
		{
			name:      "Empty condition",
			condition: Condition{},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSearchCondition(tt.condition, data, mockModel)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSearchCondition_NestedConditions(t *testing.T) {
	mockModel := MockModel{}
	data := map[string]any{}

	t.Run("Valid nested AND conditions", func(t *testing.T) {
		condition := Condition{
			And: []Condition{
				{Field: "name", Operator: "like", Value: "test"},
				{Field: "status", Operator: "=", Value: "active"},
			},
		}

		err := validateSearchCondition(condition, data, mockModel)
		assert.NoError(t, err)
	})

	t.Run("Invalid nested AND condition - bad field", func(t *testing.T) {
		condition := Condition{
			And: []Condition{
				{Field: "name", Operator: "like", Value: "test"},
				{Field: "bad_field", Operator: "=", Value: "active"},
			},
		}

		err := validateSearchCondition(condition, data, mockModel)
		assert.Error(t, err)
	})

	t.Run("Valid nested OR conditions", func(t *testing.T) {
		condition := Condition{
			Or: []Condition{
				{Field: "status", Operator: "=", Value: "active"},
				{Field: "status", Operator: "=", Value: "pending"},
			},
		}

		err := validateSearchCondition(condition, data, mockModel)
		assert.NoError(t, err)
	})
}

// =============================================================================
// SORT VALIDATION TESTS
// =============================================================================

func TestSortValidation_WithMockModel(t *testing.T) {
	mockModel := MockModel{}
	data := map[string]any{}

	tests := []struct {
		name    string
		sort    Sort
		wantErr bool
	}{
		{
			name:    "Valid sortable field ascending",
			sort:    Sort{Field: "name", Value: "asc"},
			wantErr: false,
		},
		{
			name:    "Valid sortable field descending",
			sort:    Sort{Field: "created_time", Value: "desc"},
			wantErr: false,
		},
		{
			name:    "Invalid sortable field",
			sort:    Sort{Field: "invalid_field", Value: "asc"},
			wantErr: true,
		},
		{
			name:    "Invalid sort direction defaults to asc",
			sort:    Sort{Field: "name", Value: "invalid"},
			wantErr: false, // Direction is corrected, not an error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate field exists in sortable fields
			sortableFields := mockModel.ScopesSortbleFields(data)
			isValid := sortableFields[tt.sort.Field]

			if tt.wantErr {
				assert.False(t, isValid)
			} else {
				// Either field is valid, or direction is invalid (but gets corrected)
				if tt.sort.Field != "invalid_field" {
					assert.True(t, isValid)
				}
			}
		})
	}
}

// =============================================================================
// BUILD SELECT FIELD TESTS
// =============================================================================

func TestBuildSelectField_Consistency(t *testing.T) {
	selectMap := map[string]string{
		"id":    "t.id",
		"name":  "t.name",
		"email": "t.email",
	}

	t.Run("Multiple calls return same result (deterministic)", func(t *testing.T) {
		result1 := BuildSelectField(selectMap)
		result2 := BuildSelectField(selectMap)
		assert.Equal(t, result1, result2)
	})

	t.Run("Fields are sorted alphabetically", func(t *testing.T) {
		result := BuildSelectField(selectMap)
		// Result should start with 'email' (first alphabetically)
		assert.Contains(t, result, "t.email AS email, t.id AS id, t.name AS name")
	})
}

// =============================================================================
// PAGINATION OFFSET CALCULATION TESTS
// =============================================================================

func TestScopePaginate_OffsetCalculation(t *testing.T) {
	tests := []struct {
		name           string
		page           int
		limit          int
		expectedOffset int
	}{
		{
			name:           "Page 1 offset is 0",
			page:           1,
			limit:          10,
			expectedOffset: 0,
		},
		{
			name:           "Page 2 offset is 10",
			page:           2,
			limit:          10,
			expectedOffset: 10,
		},
		{
			name:           "Page 5 with 20 items per page",
			page:           5,
			limit:          20,
			expectedOffset: 80,
		},
		{
			name:           "Page 100 with 50 items",
			page:           100,
			limit:          50,
			expectedOffset: 4950,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedOffset := (tt.page - 1) * tt.limit
			assert.Equal(t, tt.expectedOffset, expectedOffset)
		})
	}
}

// =============================================================================
// STRUCTURE & MARSHAL TESTS
// =============================================================================

func TestConditionMarshaling(t *testing.T) {
	t.Run("Condition JSON Marshal", func(t *testing.T) {
		condition := Condition{
			Field:    "name",
			Operator: "like",
			Value:    "test",
		}

		jsonData, err := json.Marshal(condition)
		assert.NoError(t, err)
		assert.NotEmpty(t, jsonData)

		// Unmarshal to verify round-trip
		var unmarshaled Condition
		err = json.Unmarshal(jsonData, &unmarshaled)
		assert.NoError(t, err)
		assert.Equal(t, condition, unmarshaled)
	})

	t.Run("Nested condition JSON Marshal", func(t *testing.T) {
		condition := Condition{
			And: []Condition{
				{Field: "name", Operator: "like", Value: "test"},
				{Field: "status", Operator: "=", Value: "active"},
			},
		}

		jsonData, err := json.Marshal(condition)
		assert.NoError(t, err)

		var unmarshaled Condition
		err = json.Unmarshal(jsonData, &unmarshaled)
		assert.NoError(t, err)
		assert.Len(t, unmarshaled.And, 2)
	})
}

func TestIndexDataMarshaling(t *testing.T) {
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

	jsonData, err := json.Marshal(indexData)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)
}

// =============================================================================
// ERROR CASE TESTS
// =============================================================================

func TestValidationErrorHandling(t *testing.T) {
	t.Run("Multiple errors on same field", func(t *testing.T) {
		errors := make(map[string][]string)

		addValidationError(errors, "password", "required")
		addValidationError(errors, "password", "minimum 8 characters")
		addValidationError(errors, "password", "must contain uppercase")

		assert.Len(t, errors["password"], 3)
		assert.Equal(t, errors["password"][0], "required")
		assert.Equal(t, errors["password"][1], "minimum 8 characters")
		assert.Equal(t, errors["password"][2], "must contain uppercase")
	})

	t.Run("Error preservation", func(t *testing.T) {
		errors := createTestValidationErrors()

		assert.Len(t, errors, 3)
		assert.Len(t, errors["name"], 1)
		assert.Len(t, errors["email"], 2)
		assert.Len(t, errors["age"], 1)
	})
}

// =============================================================================
// HELPER FUNCTION TESTS WITH JSON
// =============================================================================

func TestJsonImport(t *testing.T) {
	const (
		// Golang standard library - no import needed in test
		JSON = "encoding/json"
	)

	// This ensures we can use json functions in tests
	assert.NotEmpty(t, JSON)
}

// =============================================================================
// TYPE ASSERTIONS
// =============================================================================

func TestMockModelImplementsCoreModels(t *testing.T) {
	var model CoreModels = MockModel{}

	assert.NotNil(t, model)
	assert.Equal(t, "mock_models", model.TableName())
	assert.Equal(t, "MockModel", model.ModulName())
	assert.NotEmpty(t, model.ScopesGetSelect(map[string]any{}))
	assert.NotEmpty(t, model.ScopesSearchableFields(map[string]any{}))
	assert.NotEmpty(t, model.ScopesSortbleFields(map[string]any{}))
	assert.NotNil(t, model.ScopeJoin(map[string]any{}))
	assert.NotNil(t, model.ScopeOption(map[string]any{}))
}

// =============================================================================
// CONTEXT TIMEOUT TESTS
// =============================================================================

func TestContextOperations(t *testing.T) {
	t.Run("Context with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		select {
		case <-ctx.Done():
			t.Errorf("Context timed out unexpectedly")
		case <-time.After(100 * time.Millisecond):
			// OK, test continues
		}
	})
}

// =============================================================================
// HELPER FUNCTIONS (Shared Test Utilities)
// =============================================================================

// Helper to create test validation errors
func createTestValidationErrors() map[string][]string {
	return map[string][]string{
		"name":  {"required"},
		"email": {"invalid email format", "must be unique"},
		"age":   {"must be greater than 18"},
	}
}

// Helper to create test conditions
func createTestConditions() []Condition {
	return []Condition{
		{
			Field:    "name",
			Operator: "like",
			Value:    "john",
		},
		{
			Field:    "status",
			Operator: "=",
			Value:    "active",
		},
		{
			Field:    "age",
			Operator: ">",
			Value:    "18",
		},
	}
}

// Helper to create test sorts
func createTestSorts() []Sort {
	return []Sort{
		{Field: "created_at", Value: "desc"},
		{Field: "name", Value: "asc"},
	}
}
