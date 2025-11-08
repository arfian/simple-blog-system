package middleware

import (
	"errors"
	"simple-blog-system/config"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	ID       string `json:"id"`
	Username string `json:"username"`
}

func ParseJWTToken(tokenString string) (*JWTClaims, error) {
	configData := config.GetConfig()
	secretKey := configData.JWT.SigningKey

	tokenString = strings.Split(tokenString, "Bearer ")[1]
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid jwt token")
}
