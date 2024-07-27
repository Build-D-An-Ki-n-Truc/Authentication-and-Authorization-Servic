package jwtFunc

import (
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"

	"github.com/Build-D-An-Ki-n-Truc/auth/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// generateToken generates a JWT for the given username and role.
func GenerateToken(username string, role string) (string, error) {
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
func VerifyToken(tokenString string) (jwt.MapClaims, error) {
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

// getTokenFromRequest extracts the JWT from the Authorization header of the NATS request.
func GetTokenFromRequest(m *nats.Msg) (string, error) {

	// Fix later : Token will in header with the bearer name
	// Get the Authorization header
	authHeader := m.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header missing")
	}

	// Split the header to get the token part
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	// Return the token part
	return parts[1], nil
}
