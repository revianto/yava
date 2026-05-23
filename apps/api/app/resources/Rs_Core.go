package resources

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/models"
)

// ToResource is a generic helper that wraps model data in a standardized JSON response structure.
// It handles paginated models (IndexData), slices of data, and single record maps.
func ToResource[T any](c *fiber.Ctx, data any, singleFn func(*fiber.Ctx, any) T) any {
	switch v := data.(type) {
	case models.IndexData:
		for i, item := range v.Data {
			v.Data[i] = singleFn(c, item)
		}
		return v
	case []any:
		return fiber.Map{"data": Collection(c, v, singleFn)}
	case map[string]any:
		return fiber.Map{"data": singleFn(c, v)}
	case fiber.Map:
		return fiber.Map{"data": singleFn(c, v)}
	default:
		log.Printf("[WARN] ToResource: unhandled type %T, value: %v", data, data)
		return data
	}
}

// Collection transforms a slice of data using the provided single-item transformation function.
func Collection[T any](c *fiber.Ctx, data []any, singleFn func(*fiber.Ctx, any) T) []T {
	result := make([]T, 0, len(data))
	for _, item := range data {
		result = append(result, singleFn(c, item))
	}
	return result
}

// floatToString converts a float64 to a JSON-formatted string.
func floatToString(f float64) string {
	b, _ := json.Marshal(f)
	return string(b)
}

// =============================================================================
// GENERIC RESPONSES (used globally)
// =============================================================================

type ResponseDelete struct {
	Message string `json:"message"`
}

func DeleteResource() fiber.Map {
	return fiber.Map{
		"data": ResponseDelete{
			Message: "Deleted successfully",
		},
	}
}
