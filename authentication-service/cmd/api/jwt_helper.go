package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var secretKey = []byte("secret-key")

type tokenData struct {
	Username string
	Email    string
	UserId   int
	Type     string
}

func createToken(username, email string, userID int, userType string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"email":    email,
			"type":     userType,
			"id":       userID,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func parseTokenClaims(tokenString string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func extractFieldsFromToken(tokenString string) (tokenData, error) {

	tkData := tokenData{}

	// Verify the token
	claims, err := parseTokenClaims(tokenString)
	if err != nil {
		return tkData, err
	}

	// Access fields from claims
	tkData.Username = claims["username"].(string)
	tkData.Email = claims["email"].(string)
	tkData.UserId = int(claims["id"].(float64))
	tkData.Type = claims["type"].(string)

	return tkData, nil
}

func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
