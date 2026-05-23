package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    time.Time
	RtExpires    time.Time
}

func CreateTokenJWT() (*TokenDetails, error) {

	// create ssid
	suuid := uuid.Must(uuid.NewV4()).String()
	rsuuid := uuid.Must(uuid.NewV4()).String()

	// get aes key
	key := GetEnv("AES_KEY")

	// create token
	// encrypt token aes
	aes, err := AesEncrypt(suuid, key)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt access uuid: %w", err)
	}

	expiary_duration, _ := time.ParseDuration(GetEnv("JWT_TOKEN_EXPIRY"))
	token_expiary := time.Now().Add(expiary_duration)
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["dtid"] = aes
	// atClaims["exp"] = time.Now().Add(token_expiary).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(GetEnv("JWT_KEY")))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// create refresh token
	// encrypt token aes
	raes, err := AesEncrypt(rsuuid, key)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt refresh uuid: %w", err)
	}

	// get refresh duration
	r_expiary_duration, _ := time.ParseDuration(GetEnv("JWT_REFRESH_TOKEN_EXPIRY"))
	refresh_token_expiary := time.Now().Add(r_expiary_duration)
	ratClaims := jwt.MapClaims{}
	ratClaims["authorized"] = true
	ratClaims["dtid"] = raes

	// ratClaims["exp"] = time.Now().Add(refresh_token_expiary).Unix()
	rat := jwt.NewWithClaims(jwt.SigningMethodHS256, ratClaims)
	refresh_token, err2 := rat.SignedString([]byte(GetEnv("JWT_KEY")))
	if err2 != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err2)
	}
	tokenDetails := &TokenDetails{
		AccessToken:  token,
		RefreshToken: refresh_token,
		AccessUuid:   suuid,
		RefreshUuid:  rsuuid,
		AtExpires:    token_expiary,
		RtExpires:    refresh_token_expiary,
	}

	return tokenDetails, nil
}

func RequestTokenJwt(authorizationHeader string) (interface{}, error) {
	// validate token
	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(GetEnv("JWT_KEY")), nil
	})

	if token == nil {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		// decrypt token aes
		aes := fmt.Sprintf("%v", claims["dtid"])
		key := GetEnv("AES_KEY")
		suuid, errDecrypt := AesDecrypt(aes, key)
		if errDecrypt != nil {
			return nil, fmt.Errorf("failed to decrypt uuid: %w", errDecrypt)
		}

		return suuid, nil
	} else {
		return nil, err
	}
}
