package utils

import (
	"medods_auth/authService/internal/entities"
	"time"

	"github.com/golang-jwt/jwt"
)

// TODO
// 1. Create JWT

type JWTService struct {
	secretKey string
}

func NewJWTService(Key string) *JWTService {
	return &JWTService{
		secretKey: Key,
	}
}

func (j *JWTService) GenerateToken(user entities.User, ip string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user_id": user.ID,
		"ip":      ip,
		"exp":     time.Now().Add(time.Hour * 3).Unix(),
	})
	return token.SignedString([]byte(j.secretKey))
}
