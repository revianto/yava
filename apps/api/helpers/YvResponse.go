package helpers

import "github.com/gofiber/fiber/v2"

func YvSuccess(data interface{}) fiber.Map {
	return fiber.Map{"success": true, "data": data}
}

func YvList(data interface{}, page, limit, total int64) fiber.Map {
	return fiber.Map{
		"success": true,
		"data":    data,
		"meta": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}
}

func YvError(code string, message string) fiber.Map {
	return fiber.Map{
		"success": false,
		"error":   fiber.Map{"code": code, "message": message},
	}
}
