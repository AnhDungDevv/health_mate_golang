package utils

import (
	"health_backend/config"
	"health_backend/internal/models"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(user *models.User, config *config.Config) (string, error) {
	claims := &Claims{
		Email: user.Email,
		ID:    strconv.Itoa(int(user.ID)),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "health_backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Server.JwtSecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
