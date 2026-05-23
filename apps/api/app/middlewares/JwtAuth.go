package middlewares

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/repositories"
	"github.com/revianto/yava/api/exceptions"
	"github.com/revianto/yava/api/helpers"
	"gorm.io/gorm"
)

// JwtAuth middleware validates the JWT token in the Authorization header, httpOnly cookies,
// or token query parameter (for <img> tags).
func CheckClientToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenString string

		// 1. Authorization header
		authHeader := c.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 {
				tokenString = "Bearer " + parts[1]
			}
		}

		// 2. httpOnly cookies (user_access_token takes priority over client_access_token)
		if tokenString == "" {
			for _, name := range []string{"user_access_token", "client_access_token"} {
				if cookie := c.Cookies(name); cookie != "" {
					tokenString = "Bearer " + cookie
					break
				}
			}
		}

		// 3. Query parameter fallback (useful for <img> tags rendering static images)
		if tokenString == "" {
			if queryToken := c.Query("token"); queryToken != "" {
				tokenString = "Bearer " + queryToken
			}
		}

		locale := getLocaleFromCtx(c)

		if tokenString == "" {
			return exceptions.ResponseErrorException(c, exceptions.AuthException(c, 401, helpers.TransError(locale, "token_missing", nil)))
		}

		// Verify token
		uuid, err := helpers.RequestTokenJwt(tokenString)
		if err != nil {
			return exceptions.ResponseErrorException(c, exceptions.AuthException(c, 401, helpers.TransError(locale, "token_invalid", nil)))
		}
		db, _ := c.Locals("DB").(*gorm.DB)
		userToken, _ := repositories.AuthClientTokenSingle(db, fiber.Map{}, c, locale, func(db *gorm.DB) *gorm.DB {
			return db.Where(`
			tr_oauth_tokens.id = ?
			AND tr_oauth_tokens.revoked != 1
			AND tr_oauth_tokens.expiary_time > ?`,
				uuid,
				time.Now(),
			)
		})

		if userToken == nil {
			return exceptions.ResponseErrorException(c, exceptions.AuthException(c, 401, helpers.TransError(locale, "token_invalid", nil)))
		}

		c.Locals("tokenId", uuid)
		c.Locals("clientId", userToken["client_id"])

		return c.Next()
	}
}

// CheckUserToken validates the JWT and verifies that the token record has a user_id (i.e. user has logged in).
// Stores tokenId and authedUserId in locals for downstream handlers.
func CheckUserToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		locale := getLocaleFromCtx(c)

		var tokenString string

		// 1. Authorization header
		authHeader := c.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 {
				tokenString = "Bearer " + parts[1]
			}
		}

		// 2. httpOnly user_access_token cookie
		if tokenString == "" {
			if cookie := c.Cookies("user_access_token"); cookie != "" {
				tokenString = "Bearer " + cookie
			}
		}

		if tokenString == "" {
			return exceptions.ResponseErrorException(c, exceptions.AuthException(c, 401, helpers.TransError(locale, "token_missing", nil)))
		}

		uuid, err := helpers.RequestTokenJwt(tokenString)
		if err != nil {
			return exceptions.ResponseErrorException(c, exceptions.AuthException(c, 401, helpers.TransError(locale, "token_invalid", nil)))
		}

		db, _ := c.Locals("DB").(*gorm.DB)
		if db == nil {
			return exceptions.ResponseErrorException(c, exceptions.AuthException(c, 500, helpers.TransError(locale, "db_unavailable", nil)))
		}

		userToken, _ := repositories.AuthClientTokenSingle(db, fiber.Map{}, c, locale, func(db *gorm.DB) *gorm.DB {
			return db.Where(`
			tr_oauth_tokens.id = ?
			AND tr_oauth_tokens.revoked != 1
			AND tr_oauth_tokens.expiary_time > ?
			AND tr_oauth_tokens.user_id > 0`,
				uuid,
				time.Now(),
			)
		})

		if len(userToken) == 0 {
			return exceptions.ResponseErrorException(c, exceptions.AuthException(c, 401, helpers.TransError(locale, "unauthorized", nil)))
		}

		userId := helpers.Conv(userToken["user_id"]).Int64()

		userData, _ := repositories.UserSingle(db, fiber.Map{}, c, locale, func(db *gorm.DB) *gorm.DB {
			return db.Where("tr_users.id = ?", userId)
		})

		if len(userData) == 0 {
			return exceptions.ResponseErrorException(c, exceptions.AuthException(c, 401, helpers.TransError(locale, "unauthorized", nil)))
		}

		c.Locals("tokenId", uuid)
		c.Locals("userId", userId)
		c.Locals("userData", fiber.Map(userData))

		return c.Next()
	}
}

func getLocaleFromCtx(c *fiber.Ctx) string {
	if locale, ok := c.Locals("locale").(string); ok && locale != "" {
		return locale
	}
	return "id"
}
