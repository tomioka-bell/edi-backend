package utils

import (
	"backend/internal/core/domains"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWTFromPrincipal(principal *domains.EDI_Principal) (string, error) {
	secret := os.Getenv("TOKEN_SECRET_KEY")
	if secret == "" {
		return "", errors.New("missing TOKEN_SECRET_KEY")
	}
	jwtSecretKey := []byte(secret)

	claims := jwt.MapClaims{
		"user_id":      principal.ExternalID,
		"username":     principal.Username,
		"display_name": principal.DisplayName,
		"profile":      principal.Profile,
		"group":        principal.Group,
		"status":       principal.Status,
		"iat":          time.Now().Unix(),
		"exp":          time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", errors.New("เกิดข้อผิดพลาดในการเซ็นชื่อ JWT")
	}
	return signedToken, nil
}
