package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const salt = "kL9jasA2ksd1aAkd6ak3s"

func EncodeJwtToken(data map[string]interface{}, exp time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  float64(time.Now().Add(exp).Unix()),
		"data": data,
	})
	return token.SignedString([]byte(salt))
}

func DecodeJwtToken(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(salt), nil
	})

	if err != nil {
		if jve, ok := err.(*jwt.ValidationError); ok {
			if jve.Errors&jwt.ValidationErrorExpired == 1 {
				return nil, errors.New("token is expired")
			}
		}
		return nil, err
	}

	return token.Claims.(jwt.MapClaims)["data"].(map[string]interface{}), nil
}
