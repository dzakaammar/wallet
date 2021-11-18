package internal

import (
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenClaims struct {
	UserID string
	jwt.StandardClaims
}

func generateToken(userID string) *jwt.Token {
	claims := TokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(6 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token
}

type AuthToken struct {
	secret string
}

func NewAuthToken(secret string) *AuthToken {
	return &AuthToken{secret: secret}
}

func (a AuthToken) GetTokenFor(userID string) (string, error) {
	token := generateToken(userID)
	str, err := token.SignedString([]byte(a.secret))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return str, nil
}

func (a AuthToken) ParseToken(token string) (*TokenClaims, error) {
	var claims TokenClaims
	tok, err := jwt.ParseWithClaims(token, &claims, func(*jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !tok.Valid {
		return nil, fmt.Errorf("unauthenticated")
	}

	return &claims, nil
}
