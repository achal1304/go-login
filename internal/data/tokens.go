package data

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	UserId int
	Email  string
	jwt.StandardClaims
}

func GenerateToken(userID int, ttl time.Duration, email string) (string, error) {
	claims := Token{
		userID,
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte("my-secret-key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DecodeJWT(tokenString string) (*Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Token{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("my-secret-key"), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Token)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
