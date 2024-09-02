package jwtFunc

import (
	"encoding/base64"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Build-D-An-Ki-n-Truc/auth/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

var CFG = config.CFG

// generateToken generates a JWT for the given username and role.
func GenerateToken(username string, role string) (string, error) {

	// Define the claim
	claim := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// Sign the token
	hmacSecret, err := base64.StdEncoding.DecodeString(CFG.Secret)
	if err != nil {
		log.Panic("Error decoding Base64 string:", err)
	}
	tokenString, err := token.SignedString(hmacSecret)

	if err != nil {
		log.Println("Error creating token: ", err)
		return "", err
	}

	return tokenString, nil
}

// ExtractToken extract claims from the given JWT token string.
func ExtractToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(CFG.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
