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
	"github.com/revianto/yava/api/exceptions"
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
		return exceptions.ResponseErrorException(c, exceptions.ErrorException(c, fiber.StatusBadRequest, "Authorization code diperlukan"))
	}

	cfg := googleOAuthConfig()
	token, err := cfg.Exchange(context.Background(), code)
	if err != nil {
		return exceptions.ResponseErrorException(c, exceptions.ErrorException(c, fiber.StatusUnauthorized, fmt.Sprintf("Gagal exchange token: %s", err.Error())))
	}

	client := cfg.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil || resp.StatusCode != http.StatusOK {
		return exceptions.ResponseErrorException(c, exceptions.ErrorException(c, fiber.StatusInternalServerError, "Gagal mengambil data pengguna dari Google"))
	}
	defer resp.Body.Close()

	var gUser googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&gUser); err != nil {
		return exceptions.ResponseErrorException(c, exceptions.ErrorException(c, fiber.StatusInternalServerError, "Gagal memproses data pengguna"))
	}
	if !gUser.VerifiedEmail {
		return exceptions.ResponseErrorException(c, exceptions.ErrorException(c, fiber.StatusForbidden, "Email Google belum diverifikasi"))
	}

	data := fiber.Map{
		"google_id":  gUser.Sub,
		"email":      gUser.Email,
		"name":       gUser.Name,
		"avatar_url": gUser.Picture,
	}

	result, svcErr := services.AuthProcessGoogleUser(getDB(c), data, c, getLocale(c))
	if svcErr != nil {
		return exceptions.ResponseErrorException(c, svcErr.(exceptions.AppError))
	}

	resultMap, _ := result.(map[string]any)
	jwtToken, _ := resultMap["token"].(string)

	dur, _ := time.ParseDuration(helpers.GetEnv("JWT_TOKEN_EXPIRY"))
	if dur == 0 {
		dur = 24 * time.Hour
	}
	c.Cookie(&fiber.Cookie{
		Name:     "yv_token",
		Value:    jwtToken,
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
	return c.JSON(fiber.Map{"data": fiber.Map{"message": "Berhasil logout"}})
}

// GET /v1/auth/me
func AuthMe(c *fiber.Ctx) error {
	result, err := services.AuthMe(getDB(c), getBodyData(c), c, getLocale(c))
	if err != nil {
		return exceptions.ResponseErrorException(c, err.(exceptions.AppError))
	}
	return c.JSON(resources.UserResource(c, result))
}
