package helpers

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type YvClaims struct {
	UserID    int64  `json:"user_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	jwt.RegisteredClaims
}

func YvCreateToken(userID int64, email, name, avatarURL string) (string, error) {
	dur, err := time.ParseDuration(GetEnv("JWT_TOKEN_EXPIRY"))
	if err != nil {
		dur = 24 * time.Hour
	}
	claims := YvClaims{
		UserID:    userID,
		Email:     email,
		Name:      name,
		AvatarURL: avatarURL,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(dur)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(GetEnv("JWT_KEY")))
}

func YvParseToken(tokenStr string) (*YvClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &YvClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(GetEnv("JWT_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*YvClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
