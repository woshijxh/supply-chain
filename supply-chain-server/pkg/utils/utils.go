package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"supply-chain-server/internal/config"
)

func GenerateToken(userID uint, username, role string) (string, error) {
	claims := jwt.MapClaims{
		"userId":   userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Duration(config.AppConfig.JWT.ExpireTime) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
