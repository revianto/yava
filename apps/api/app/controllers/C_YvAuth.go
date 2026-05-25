package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/revianto/yava/api/app/resources"
	"github.com/revianto/yava/api/app/services"
	"github.com/revianto/yava/api/helpers"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func googleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     helpers.GetEnv("GOOGLE_CLIENT_ID"),
		ClientSecret: helpers.GetEnv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  helpers.GetEnv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

type googleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"email_verified"`
}

// GET /v1/auth/google
func AuthGoogle(c *fiber.Ctx) error {
	url := googleOAuthConfig().AuthCodeURL("yava-state", oauth2.AccessTypeOnline)
	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

// GET /v1/auth/google/callback
func AuthGoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.YvError("MISSING_CODE", "Authorization code diperlukan"))
	}

	cfg := googleOAuthConfig()
	token, err := cfg.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(helpers.YvError("OAUTH_EXCHANGE_FAILED", fmt.Sprintf("Gagal exchange token: %s", err.Error())))
	}

	client := cfg.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil || resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.YvError("USERINFO_FAILED", "Gagal mengambil data pengguna dari Google"))
	}
	defer resp.Body.Close()

	var gUser googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&gUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.YvError("USERINFO_PARSE_FAILED", "Gagal memproses data pengguna"))
	}
	if !gUser.VerifiedEmail {
		return c.Status(fiber.StatusForbidden).JSON(helpers.YvError("EMAIL_NOT_VERIFIED", "Email Google belum diverifikasi"))
	}

	var avatarURL *string
	if gUser.Picture != "" {
		avatarURL = &gUser.Picture
	}

	result, svcErr := services.AuthProcessGoogleUser(getDB(c), gUser.Sub, gUser.Email, gUser.Name, avatarURL)
	if svcErr != nil {
		se := svcErr.(*services.ServiceError)
		return c.Status(se.Code).JSON(helpers.YvError(se.ErrCode, se.Message))
	}

	dur, _ := time.ParseDuration(helpers.GetEnv("JWT_TOKEN_EXPIRY"))
	if dur == 0 {
		dur = 24 * time.Hour
	}
	c.Cookie(&fiber.Cookie{
		Name:     "yv_token",
		Value:    result.Token,
		Expires:  time.Now().Add(dur),
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   helpers.GetEnv("APP_ENV") == "production",
	})

	frontendURL := helpers.GetEnv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	return c.Redirect(frontendURL, fiber.StatusTemporaryRedirect)
}

// POST /v1/auth/logout
func AuthLogout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "yv_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HTTPOnly: true,
		SameSite: "Lax",
	})
	return c.JSON(helpers.YvSuccess(fiber.Map{"message": "Berhasil logout"}))
}

// GET /v1/auth/me
func AuthMe(c *fiber.Ctx) error {
	userID, _ := c.Locals("yv_user_id").(int64)
	user, svcErr := services.AuthGetMe(getDB(c), userID)
	if svcErr != nil {
		se := svcErr.(*services.ServiceError)
		return c.Status(se.Code).JSON(helpers.YvError(se.ErrCode, se.Message))
	}
	return c.JSON(helpers.YvSuccess(resources.UserResource(*user)))
}
