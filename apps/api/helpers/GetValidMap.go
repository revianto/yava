package helpers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func GetSimpleMapValue(key string, defaultValue string, data map[string]string) string {
	result, exists := data[key]
	if !exists || result == "" {
		return defaultValue
	}
	return result
}

func GetMap(val any) map[string]any {
	switch v := val.(type) {
	case map[string]any:
		return v
	case fiber.Map:
		return map[string]any(v)
	case []uint8:
		var m map[string]any
		if json.Unmarshal(v, &m) == nil {
			return m
		}
	}
	return nil
}
