package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/exceptions"
	"github.com/revianto/yava/api/helpers"

	"gorm.io/gorm"
)

// FiberCtxKey is the Go context key for storing *fiber.Ctx in GORM's statement context.
// Using a typed key avoids collisions with other context values.
type FiberCtxKey struct{}

// FiberCtxFromDB retrieves *fiber.Ctx stored via context.WithValue in the GORM statement context.
func FiberCtxFromDB(tx *gorm.DB) *fiber.Ctx {
	if tx.Statement == nil || tx.Statement.Context == nil {
		return nil
	}
	ctx, _ := tx.Statement.Context.Value(FiberCtxKey{}).(*fiber.Ctx)
	return ctx
}

// =============================================================================
// INTERFACES & TYPES
// =============================================================================

// CoreModels interface defines the contract for models that support dynamic indexing
type CoreModels interface {
	ScopesGetSelect(map[string]any) map[string]string
	ScopesSearchableFields(map[string]any) map[string]SearchableFields
	ScopesSortbleFields(map[string]any) map[string]bool
	TableName() string
	ModulName() string
	ScopeJoin(map[string]any) func(*gorm.DB) *gorm.DB
	ScopeOption(map[string]any) func(*gorm.DB) *gorm.DB
}

// SearchableFields defines allowed operators for a searchable field
type SearchableFields struct {
	Operators []string
}

// IndexData represents paginated list response
type IndexData struct {
	Data []any      `json:"data"`
	Meta Pagination `json:"meta"`
}

// Pagination metadata for paginated responses
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Condition represents search condition (supports nested AND/OR)
type Condition struct {
	Field    string      `json:"field,omitempty"`
	Operator string      `json:"operator,omitempty"`
	Value    string      `json:"value,omitempty"`
	And      []Condition `json:"and,omitempty"`
	Or       []Condition `json:"or,omitempty"`
}

// Sort represents sort configuration (supports multiple sorts)
type Sort struct {
	Field string `json:"field,omitempty"`
	Value string `json:"value,omitempty"`
	Sort  []Sort `json:"sort,omitempty"`
}

// ValidateResult stores validation result
type ValidateResult struct {
	IsValid bool
	Data    interface{}
	Errors  map[string][]string
}

// =============================================================================
// MAIN DATA RETRIEVAL FUNCTIONS
// =============================================================================

// GetIndexData retrieves paginated list data with search, sort, and pagination
// Supports both single model and UNION queries (max 2 models)
func GetIndexData(tx *gorm.DB, data map[string]any, c *fiber.Ctx, locale string, models ...CoreModels) (IndexData, any) {
	returnData := []map[string]any{}

	// Parse Search
	search, errSearch := ParseSearch(tx, data, c, locale, models...)
	if errSearch != nil {
		return IndexData{}, exceptions.ErrorException(c, fiber.StatusBadRequest, errSearch.Error())
	}

	// Parse Sort
	sort, errSort := ParseSort(tx, data, c, locale, models...)
	if errSort != nil {
		return IndexData{}, exceptions.ErrorException(c, fiber.StatusBadRequest, errSort.Error())
	}

	// Parse Pagination
	paginate := Paginate(tx, data, c, locale, search, models...)

	var errSelect error

	if len(models) > 1 {
		// UNION STRATEGY
		if len(models) > 2 {
			return IndexData{}, exceptions.ErrorException(c, fiber.StatusBadRequest, "Union currently only supports 2 models.")
		}

		sql1 := tx.ToSQL(func(queryTx *gorm.DB) *gorm.DB {
			return queryTx.Model(models[0]).
				Select(BuildSelectField(models[0].ScopesGetSelect(data))).
				Scopes(
					models[0].ScopeJoin(data),
					models[0].ScopeOption(data),
					ScopeSearch(tx, data, c, locale, search, models[0]),
				).Find(&[]map[string]any{})
		})

		sql2 := tx.ToSQL(func(queryTx *gorm.DB) *gorm.DB {
			return queryTx.Model(models[1]).
				Select(BuildSelectField(models[1].ScopesGetSelect(data))).
				Scopes(
					models[1].ScopeJoin(data),
					models[1].ScopeOption(data),
					ScopeSearch(tx, data, c, locale, search, models[1]),
				).Find(&[]map[string]any{})
		})

		offset := (paginate.Page - 1) * paginate.Limit
		unionSQL := fmt.Sprintf("SELECT * FROM (%s UNION %s) AS u LIMIT %d OFFSET %d", sql1, sql2, paginate.Limit, offset)

		rows, err := tx.Raw(unionSQL).Rows()
		if err != nil {
			errSelect = err
		} else {
			defer rows.Close()
			cols, _ := rows.Columns()
			for rows.Next() {
				columns := make([]interface{}, len(cols))
				columnPointers := make([]interface{}, len(cols))
				for i := range columns {
					columnPointers[i] = &columns[i]
				}
				rows.Scan(columnPointers...)
				rowMap := make(map[string]any)
				for i, colName := range cols {
					rowMap[colName] = columns[i]
				}
				returnData = append(returnData, rowMap)
			}
		}

	} else {
		// SINGLE MODEL STRATEGY
		errSelect = BuildSelectQuery(tx, data, c, locale, models[0]).
			Scopes(
				ScopeSearch(tx, data, c, locale, search, models[0]),
				ScopeSort(tx, data, c, locale, sort, models[0]),
				ScopePaginate(c, paginate.Page, paginate.Limit),
			).Find(&returnData).Error
	}

	if errSelect != nil {
		return IndexData{}, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "Data not found or database error: "+errSelect.Error())
	}

	// Convert to []any
	convertedData := make([]any, len(returnData))
	for i, v := range returnData {
		convertedData[i] = v
	}

	return IndexData{
		Data: convertedData,
		Meta: paginate,
	}, nil
}

// GetSingleData retrieves single record with optional conditions
func GetSingleData(tx *gorm.DB, data map[string]any, c *fiber.Ctx, locale string, conditions func(*gorm.DB) *gorm.DB, models ...CoreModels) (map[string]any, any) {
	returnData := map[string]any{}
	var errSelect error
	data["show_single"] = 1

	search, errSearch := ParseSearch(tx, data, c, locale, models...)
	if errSearch != nil {
		return map[string]any{}, exceptions.ErrorException(c, fiber.StatusBadRequest, errSearch.Error())
	}

	if len(models) > 1 {
		// UNION STRATEGY
		if len(models) > 2 {
			return map[string]any{}, exceptions.ErrorException(c, fiber.StatusBadRequest, "Union currently only supports 2 models.")
		}

		query1 := BuildSelectQuery(tx, data, c, locale, models[0]).
			Scopes(ScopeSearch(tx, data, c, locale, search, models[0]))
		query2 := BuildSelectQuery(tx, data, c, locale, models[1]).
			Scopes(ScopeSearch(tx, data, c, locale, search, models[1]))

		unionQuery := tx.Raw("? UNION ?", query1, query2)
		q := tx.Table("(?) AS u", unionQuery)
		if conditions != nil {
			q = q.Scopes(conditions)
		}
		errSelect = q.Take(&returnData).Error

	} else {
		// SINGLE MODEL STRATEGY
		q := BuildSelectQuery(tx, data, c, locale, models[0]).
			Scopes(ScopeSearch(tx, data, c, locale, search, models[0]))
		if conditions != nil {
			q = q.Scopes(conditions)
		}
		errSelect = q.Take(&returnData).Error
	}

	if errSelect != nil {
		if errors.Is(errSelect, gorm.ErrRecordNotFound) {
			return map[string]any{}, exceptions.ErrorException(c, fiber.ErrBadRequest.Code, "Data "+models[0].ModulName()+" not found")
		}
		log.Printf("[DB ERROR] GetSingleData: %v", errSelect)
		return map[string]any{}, exceptions.ErrorException(c, fiber.ErrBadRequest.Code, "Terjadi kesalahan sistem, silakan coba lagi")
	}

	return returnData, nil
}

// GetMultipleData retrieves multiple records with optional conditions
func GetMultipleData(tx *gorm.DB, data map[string]any, c *fiber.Ctx, locale string, conditions func(*gorm.DB) *gorm.DB, models ...CoreModels) ([]map[string]any, any) {
	returnData := []map[string]any{}
	var errSelect error

	search, errSearch := ParseSearch(tx, data, c, locale, models...)
	if errSearch != nil {
		return []map[string]any{}, exceptions.ErrorException(c, fiber.StatusBadRequest, errSearch.Error())
	}

	if len(models) > 1 {
		// UNION STRATEGY
		if len(models) > 2 {
			return []map[string]any{}, exceptions.ErrorException(c, fiber.StatusBadRequest, "Union currently only supports 2 models.")
		}

		query1 := BuildSelectQuery(tx, data, c, locale, models[0]).
			Scopes(ScopeSearch(tx, data, c, locale, search, models[0]))
		query2 := BuildSelectQuery(tx, data, c, locale, models[1]).
			Scopes(ScopeSearch(tx, data, c, locale, search, models[1]))

		unionQuery := tx.Raw("? UNION ?", query1, query2)
		q := tx.Table("(?) AS u", unionQuery)
		if conditions != nil {
			q = q.Scopes(conditions)
		}
		errSelect = q.Find(&returnData).Error

	} else {
		// SINGLE MODEL STRATEGY
		q := BuildSelectQuery(tx, data, c, locale, models[0]).
			Scopes(ScopeSearch(tx, data, c, locale, search, models[0]))
		if conditions != nil {
			q = q.Scopes(conditions)
		}
		errSelect = q.Find(&returnData).Error
	}

	if errSelect != nil {
		if errors.Is(errSelect, gorm.ErrRecordNotFound) {
			return []map[string]any{}, nil
		}
		log.Printf("[DB ERROR] GetMultipleData: %v", errSelect)
		return []map[string]any{}, exceptions.ErrorException(c, fiber.StatusNotAcceptable, "Terjadi kesalahan sistem, silakan coba lagi")
	}

	return returnData, nil
}

// =============================================================================
// QUERY BUILDER FUNCTIONS
// =============================================================================

// BuildSelectField builds SELECT clause with sorted columns for consistent UNION
func BuildSelectField(selectMap map[string]string) string {
	keys := make([]string, 0, len(selectMap))
	for k := range selectMap {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	cols := []string{}
	for _, alias := range keys {
		cols = append(cols, fmt.Sprintf("%s AS %s", selectMap[alias], alias))
	}
	return strings.Join(cols, ", ")
}

// BuildSelectQuery builds base SELECT query with joins and options
func BuildSelectQuery(tx *gorm.DB, data map[string]any, c *fiber.Ctx, locale string, m CoreModels) *gorm.DB {
	return tx.Session(&gorm.Session{NewDB: true}).Model(m).
		Select(BuildSelectField(m.ScopesGetSelect(data))).
		Scopes(
			m.ScopeJoin(data),
			m.ScopeOption(data),
		)
}

// =============================================================================
// SCOPE FUNCTIONS
// =============================================================================

// ScopePaginate applies LIMIT and OFFSET
func ScopePaginate(c *fiber.Ctx, page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

// ScopeSearch applies WHERE conditions recursively
func ScopeSearch(tx *gorm.DB, data map[string]any, c *fiber.Ctx, locale string, s Condition, m ...CoreModels) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return applySearchRecursive(db, data, c, locale, s, m...)
	}
}

// ScopeSort applies ORDER BY clause
func ScopeSort(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, sort Sort, m ...CoreModels) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return applySortRecursive(tx, data, c, locale, sort, m...)
	}
}

// =============================================================================
// PAGINATION
// =============================================================================

// Paginate calculates pagination with accurate count including search conditions
func Paginate(tx *gorm.DB, data map[string]any, c *fiber.Ctx, locale string, search Condition, models ...CoreModels) Pagination {
	page := helpers.Conv(data["page"]).Default(c.Query("page", "1")).Int()
	limit := helpers.Conv(data["limit"]).Default(c.Query("limit", "10")).Int()
	if page <= 0 {
		page = 1
	}

	total := int64(0)

	if len(models) > 1 {
		// UNION count
		sql1 := tx.ToSQL(func(queryTx *gorm.DB) *gorm.DB {
			return queryTx.Model(models[0]).
				Select(BuildSelectField(models[0].ScopesGetSelect(data))).
				Scopes(
					models[0].ScopeJoin(data),
					models[0].ScopeOption(data),
					ScopeSearch(tx, data, c, locale, search, models[0]),
				).Find(&[]map[string]any{})
		})

		sql2 := tx.ToSQL(func(queryTx *gorm.DB) *gorm.DB {
			return queryTx.Model(models[1]).
				Select(BuildSelectField(models[1].ScopesGetSelect(data))).
				Scopes(
					models[1].ScopeJoin(data),
					models[1].ScopeOption(data),
					ScopeSearch(tx, data, c, locale, search, models[1]),
				).Find(&[]map[string]any{})
		})

		countSQL := fmt.Sprintf("SELECT COUNT(*) FROM (%s UNION %s) AS temp", sql1, sql2)
		tx.Raw(countSQL).Scan(&total)

	} else {
		// Single model count
		tx.Model(models[0]).Scopes(
			models[0].ScopeJoin(data),
			models[0].ScopeOption(data),
			ScopeSearch(tx, data, c, locale, search, models[0]),
		).Count(&total)
	}

	return Pagination{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}
}

// =============================================================================
// SORT FUNCTIONS
// =============================================================================

// ParseSort parses sort parameters from query params or JSON body
func ParseSort(tx *gorm.DB, data map[string]any, c *fiber.Ctx, locale string, models ...CoreModels) (Sort, error) {
	root := Sort{}

	// Try JSON Body first
	if data["sort"] != nil {
		arrSort := []Sort{}
		json.Unmarshal(helpers.Conv(data["sort"]).Byte(), &arrSort)

		for i, val := range arrSort {
			if strings.ToLower(val.Value) != "desc" && strings.ToLower(val.Value) != "asc" {
				arrSort[i].Value = "asc"
			}
			isValid := false
			for _, m := range models {
				if m.ScopesSortbleFields(map[string]any{})[val.Field] {
					isValid = true
					break
				}
			}
			if !isValid {
				return Sort{}, fmt.Errorf("sort field is not allowed: %s", val.Field)
			}
		}
		root.Sort = arrSort
		return root, nil
	}

	// Fallback to Query Params
	params := map[string][]string{}
	c.Context().QueryArgs().VisitAll(func(k, v []byte) {
		params[string(k)] = []string{string(v)}
	})

	tempMap := make(map[int]*Sort)

	validateField := func(f string) bool {
		for _, m := range models {
			if m.ScopesSortbleFields(map[string]any{})[f] {
				return true
			}
		}
		return false
	}

	for key, val := range params {
		if !strings.HasPrefix(key, "sort[") {
			continue
		}

		parts := strings.Split(key, "[")

		if len(parts) > 2 {
			// Indexed: sort[0][field]
			indexStr := strings.TrimSuffix(parts[1], "]")
			idx := helpers.Conv(indexStr).Int()
			prop := strings.TrimSuffix(parts[2], "]")

			if _, ok := tempMap[idx]; !ok {
				tempMap[idx] = &Sort{}
			}

			switch prop {
			case "field":
				if !validateField(val[0]) {
					return Sort{}, fmt.Errorf("sort field is not allowed: %s", val[0])
				}
				tempMap[idx].Field = val[0]
			case "value", "dir":
				v := strings.ToLower(val[0])
				if v != "asc" && v != "desc" {
					v = "asc"
				}
				tempMap[idx].Value = v
			}
		} else {
			// Root: sort[field]
			prop := strings.TrimSuffix(parts[1], "]")
			switch prop {
			case "field":
				if !validateField(val[0]) {
					return Sort{}, fmt.Errorf("sort field is not allowed: %s", val[0])
				}
				root.Field = val[0]
			case "value", "dir":
				v := strings.ToLower(val[0])
				if v != "asc" && v != "desc" {
					v = "asc"
				}
				root.Value = v
			}
		}
	}

	// Reconstruct sort array from map
	for i := 0; i < len(tempMap); i++ {
		if item, ok := tempMap[i]; ok {
			root.Sort = append(root.Sort, *item)
		} else {
			return root, fmt.Errorf("sort index missing: %d (must be sequential from 0)", i)
		}
	}

	return root, nil
}

// applySortRecursive applies sort conditions recursively
func applySortRecursive(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, s Sort, models ...CoreModels) *gorm.DB {
	if s.Field != "" {
		isValid := false
		for _, m := range models {
			if m.ScopesSortbleFields(map[string]any{})[s.Field] {
				isValid = true
				break
			}
		}

		if isValid {
			direction := "ASC"
			if strings.ToUpper(s.Value) == "DESC" {
				direction = "DESC"
			}

			sortColumn := s.Field
			if len(models) == 1 {
				sortColumn = models[0].ScopesGetSelect(data)[s.Field]
			}

			tx = tx.Order(fmt.Sprintf("%s %s", sortColumn, direction))
		}
	}

	for _, nestedSort := range s.Sort {
		tx = applySortRecursive(tx, data, c, locale, nestedSort, models...)
	}

	return tx
}

// =============================================================================
// SEARCH FUNCTIONS
// =============================================================================

// allowOperators checks if operator is allowed
func allowOperators(op string) bool {
	allowed := map[string]bool{
		"=":    true,
		"!=":   true,
		">":    true,
		"<":    true,
		">=":   true,
		"<=":   true,
		"like": true,
		"is":   true,
	}
	return allowed[op]
}

// ParseSearch parses search conditions from query params or JSON body
func ParseSearch(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, models ...CoreModels) (Condition, error) {
	var finalResult Condition

	var queryPayload struct {
		Search Condition `query:"search"`
	}

	// Try Query Params
	if err := c.QueryParser(&queryPayload); err == nil {
		if queryPayload.Search.Field != "" || len(queryPayload.Search.And) > 0 || len(queryPayload.Search.Or) > 0 {
			finalResult = queryPayload.Search
		}
	}

	// Fallback to Body
	if finalResult.Field == "" && len(finalResult.And) == 0 && len(finalResult.Or) == 0 {
		if len(c.Body()) > 0 {
			if err := c.BodyParser(&finalResult); err != nil {
				return Condition{}, fmt.Errorf("invalid json body format: %s", err.Error())
			}
		}
	}

	// Validate
	if err := validateSearchCondition(finalResult, data, models...); err != nil {
		return finalResult, err
	}

	return finalResult, nil
}

// applySearchRecursive applies search conditions recursively
func applySearchRecursive(tx *gorm.DB, data fiber.Map, c *fiber.Ctx, locale string, s Condition, models ...CoreModels) *gorm.DB {
	// Single field condition
	if s.Field != "" {
		colName := s.Field
		if len(models) == 1 {
			if realCol, ok := models[0].ScopesGetSelect(data)[s.Field]; ok {
				colName = realCol
			}
		}

		query := fmt.Sprintf("%s %s ?", colName, s.Operator)
		value := s.Value
		if s.Operator == "like" {
			value = "%" + s.Value + "%"
		}
		return tx.Where(query, value)
	}

	// AND Group
	if len(s.And) > 0 {
		for _, v := range s.And {
			tx = applySearchRecursive(tx, data, c, locale, v, models...)
		}
		return tx
	}

	// OR Group
	if len(s.Or) > 0 {
		if len(s.Or) > 0 {
			tx = applySearchRecursive(tx, data, c, locale, s.Or[0], models...)
		}
		for i := 1; i < len(s.Or); i++ {
			orQuery := applySearchRecursive(tx.Session(&gorm.Session{NewDB: true}), data, c, locale, s.Or[i], models...)
			tx = tx.Or(orQuery)
		}
		return tx
	}

	return tx
}

// validateSearchCondition validates search fields and operators
func validateSearchCondition(cond Condition, data fiber.Map, models ...CoreModels) error {
	if cond.Field == "" && len(cond.And) == 0 && len(cond.Or) == 0 {
		return nil
	}

	if cond.Field != "" {
		fieldFound := false
		operatorAllowed := false
		allowedOps := []string{}

		for _, m := range models {
			scopes := m.ScopesSearchableFields(data)
			if config, ok := scopes[cond.Field]; ok {
				fieldFound = true
				allowedOps = config.Operators
				if slices.Contains(config.Operators, cond.Operator) && allowOperators(cond.Operator) {
					operatorAllowed = true
				}
				if fieldFound {
					break
				}
			}
		}

		if !fieldFound {
			return fmt.Errorf("field '%s' is not allowed", cond.Field)
		}
		if !operatorAllowed {
			return fmt.Errorf("operator '%s' is not valid for field '%s' (allowed: %v)", cond.Operator, cond.Field, allowedOps)
		}
	}

	// Validate children
	for _, sub := range cond.And {
		if err := validateSearchCondition(sub, data, models...); err != nil {
			return err
		}
	}
	for _, sub := range cond.Or {
		if err := validateSearchCondition(sub, data, models...); err != nil {
			return err
		}
	}

	return nil
}

// =============================================================================
// VALIDATION
// =============================================================================

// Validate parses and validates data to struct
func Validate(c *fiber.Ctx, data interface{}, orig interface{}, locale string) (ValidateResult, interface{}) {
	result := ValidateResult{
		IsValid: false,
		Errors:  make(map[string][]string),
	}

	t := reflect.TypeOf(orig)
	if t.Kind() != reflect.Ptr {
		result.Errors["_system"] = []string{"System Error: orig must be a pointer"}
		return result, exceptions.ValidateMapException(c, 500, result.Errors)
	}

	newP := reflect.New(t.Elem()).Interface()

	// Convert data to []byte
	var byteData []byte
	switch d := data.(type) {
	case []byte:
		byteData = d
	case string:
		byteData = []byte(d)
	default:
		var err error
		byteData, err = json.Marshal(d)
		if err != nil {
			result.Errors["_system"] = []string{"Failed to marshal data: " + err.Error()}
			return result, exceptions.ValidateMapException(c, 400, result.Errors)
		}
	}

	// JSON Unmarshal
	if err := json.Unmarshal(byteData, newP); err != nil {
		switch te := err.(type) {
		case *json.SyntaxError:
			msg := helpers.Trans(locale, "validations.syntax_error", map[string]interface{}{"offset": te.Offset})
			if msg == "validations.syntax_error" {
				msg = fmt.Sprintf("Syntax error at offset %v", te.Offset)
			}
			addValidationError(result.Errors, "_json", msg)
		case *json.UnmarshalTypeError:
			msg := helpers.Trans(locale, "validations.invalid_type", map[string]interface{}{
				"attribute": te.Field,
				"expected":  te.Type.String(),
				"got":       te.Value,
			})
			if msg == "validations.invalid_type" {
				msg = fmt.Sprintf("Invalid type: expected %v, got %s", te.Type, te.Value)
			}
			addValidationError(result.Errors, te.Field, msg)
		default:
			addValidationError(result.Errors, "_json", err.Error())
		}
		return result, exceptions.ValidateMapException(c, 406, result.Errors)
	}

	// Struct validation
	validate := validator.New()
	helpers.ApplyCustomValidators(validate)

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return fld.Name
		}
		return name
	})

	if err := validate.Struct(newP); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			msg := helpers.GetValidationMessage(locale, e.Tag(), e.Field(), e.Param())
			addValidationError(result.Errors, e.Field(), msg)
		}
		return result, exceptions.ValidateMapException(c, 406, result.Errors)
	}

	result.IsValid = true
	result.Data = newP
	return result, nil
}

func addValidationError(m map[string][]string, field string, msg string) {
	m[field] = append(m[field], msg)
}

// =============================================================================
// AUTO FILL HOOKS
// =============================================================================

// AutoFillCreate fills Created_time, Created_by, Created_from automatically
func AutoFillCreate(model interface{}, tx *gorm.DB) error {
	fiberCtx := FiberCtxFromDB(tx)
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if f := v.FieldByName("Created_time"); f.IsValid() && f.CanSet() {
		setTimeField(f)
	}
	if f := v.FieldByName("Created_from"); f.IsValid() && f.CanSet() {
		if fiberCtx != nil {
			setStringField(f, fiberCtx.IP())
		}
	}
	if f := v.FieldByName("Created_by"); f.IsValid() && f.CanSet() {
		if fiberCtx != nil {
			setIntField(f, GetUserID(fiberCtx))
		}
	}

	return nil
}

// AutoFillUpdate fills Modified_time, Modified_by, Modified_from automatically
func AutoFillUpdate(model interface{}, tx *gorm.DB) error {
	fiberCtx := FiberCtxFromDB(tx)
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if f := v.FieldByName("Modified_time"); f.IsValid() && f.CanSet() {
		setTimeField(f)
	}
	if f := v.FieldByName("Modified_from"); f.IsValid() && f.CanSet() {
		if fiberCtx != nil {
			setStringField(f, fiberCtx.IP())
		}
	}
	if f := v.FieldByName("Modified_by"); f.IsValid() && f.CanSet() {
		if fiberCtx != nil {
			setIntField(f, GetUserID(fiberCtx))
		}
	}

	return nil
}

// AutoFillDelete fills Deleted_time, Deleted_by, Deleted_from automatically
func AutoFillDelete(model interface{}, tx *gorm.DB) error {
	fiberCtx := FiberCtxFromDB(tx)
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if f := v.FieldByName("Deleted_time"); f.IsValid() && f.CanSet() {
		setTimeField(f)
	}
	if f := v.FieldByName("Deleted_from"); f.IsValid() && f.CanSet() {
		if fiberCtx != nil {
			setStringField(f, fiberCtx.IP())
		}
	}
	if f := v.FieldByName("Deleted_by"); f.IsValid() && f.CanSet() {
		if fiberCtx != nil {
			setIntField(f, GetUserID(fiberCtx))
		}
	}

	return nil
}

// GetUserID retrieves user ID from context
func GetUserID(c *fiber.Ctx) int {
	if userID := c.Locals("user_id"); userID != nil {
		return helpers.Conv(userID).Int()
	}
	return 0
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

func setTimeField(f reflect.Value) {
	now := time.Now()
	switch f.Kind() {
	case reflect.Struct:
		if f.Type() == reflect.TypeOf(gorm.DeletedAt{}) {
			f.Set(reflect.ValueOf(gorm.DeletedAt{Time: now, Valid: true}))
		} else {
			f.Set(reflect.ValueOf(now))
		}
	case reflect.String:
		f.SetString(now.Format("2006-01-02 15:04:05"))
	}
}

func setStringField(f reflect.Value, value string) {
	if f.Kind() == reflect.String {
		f.SetString(value)
	}
}

func setIntField(f reflect.Value, value int) {
	switch f.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f.SetInt(int64(value))
	case reflect.String:
		f.SetString(helpers.Conv(value).String())
	}
}
