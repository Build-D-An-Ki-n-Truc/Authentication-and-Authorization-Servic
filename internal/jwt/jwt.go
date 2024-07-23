package jwt

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Build-D-An-Ki-n-Truc/auth/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// Global config variable

// generateToken generates a JWT for the given username and role.
func generateToken(username string, role string) (string, error) {
	var cfg = config.LoadConfig()
	// Define the claim
	claim := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claim)

	// Sign the token
	tokenString, err := token.SignedString(cfg.Secret)

	if err != nil {
		log.Println("Error creating token: ", err)
		return "", err
	}

	return tokenString, nil
}

// verifyToken verifies the given JWT token string.
func verifyToken(tokenString string) (jwt.MapClaims, error) {
	var cfg = config.LoadConfig()
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.Secret), nil
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

// func getToken() (string, error){

// }