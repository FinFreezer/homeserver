package auth

import (
	"time"

	"log"

	"github.com/golang-jwt/jwt/v5"
)

func MakeJWT(userName, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "Fin-Homeserver",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(tokenSecret))
	return tokenStr, err
}

func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		log.Printf("An error occurred: %s", err)
		return "Unauthorized", err
	} else if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		validUser := claims.Subject
		return validUser, nil
	} else {
		log.Printf("An error occurred: %s", err)
	}
	return "Unauthorized", err
}
