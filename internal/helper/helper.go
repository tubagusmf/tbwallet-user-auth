package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tubagusmf/tbwallet-user-auth/internal/config"
	"github.com/tubagusmf/tbwallet-user-auth/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func HashRequestPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(userID int64) (strToken string, err error) {
	expiredAt := time.Now().UTC().Add(config.JWTExp())
	strToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":     expiredAt.Unix(),
		"user_id": userID,
	}).SignedString([]byte(config.JWTSigningKey()))
	return
}

func DecodeToken(token string, claim *model.CustomClaims) (err error) {
	jwt.ParseWithClaims(token, claim, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSigningKey()), nil
	})
	return
}
