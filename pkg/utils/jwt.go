package utils

import (
	"errors"
	"health_backend/config"
	"health_backend/internal/domain/entity"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(user *entity.User, config *config.Config) (string, error) {
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

func VerifyJWTToken(tokenString string, config *config.Config) (*Claims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(config.Server.JwtSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
