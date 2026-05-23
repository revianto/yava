package helpers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// NormalizePhone menormalisasi nomor HP ke format E.164 (+62xxx)
// Menerima: 081xxx, 6281xxx, 81xxx, +6281xxx → +6281xxx
func NormalizePhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if strings.HasPrefix(phone, "+62") {
		return phone
	}
	if strings.HasPrefix(phone, "62") {
		return "+" + phone
	}
	if strings.HasPrefix(phone, "0") {
		return "+62" + phone[1:]
	}
	return "+62" + phone
}

// GetLocale mengambil locale dari context, fallback ke "id"
func GetLocale(c *fiber.Ctx) string {
	if locale, ok := c.Locals("locale").(string); ok && locale != "" {
		return locale
	}
	return "id"
}

// GetTokenId mengambil token ID dari context (di-set oleh CheckClientToken / CheckUserToken)
func GetTokenId(c *fiber.Ctx) string {
	return Conv(c.Locals("tokenId")).String()
}

// GetUserId mengambil user ID dari context (di-set oleh CheckUserToken)
func GetUserId(c *fiber.Ctx) int64 {
	return Conv(c.Locals("userId")).Int64()
}

// GetClientId mengambil client ID dari context (di-set oleh CheckClientToken)
func GetClientId(c *fiber.Ctx) string {
	return Conv(c.Locals("clientId")).String()
}

// GetUserData mengambil user data dari context (di-set oleh CheckUserToken)
func GetUserData(c *fiber.Ctx) fiber.Map {
	if userData, ok := c.Locals("userData").(fiber.Map); ok {
		return userData
	}
	return nil
}
