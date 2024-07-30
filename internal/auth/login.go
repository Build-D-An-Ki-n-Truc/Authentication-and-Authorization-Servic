package auth

import (
	"fmt"

	"github.com/Build-D-An-Ki-n-Truc/auth/internal/hashing"
)

// Check for username and password in corresponding database. Return Role and True if password correct, False otherwise
func Login(username string, password string, collection string) (string, bool) {
	// sample password until have database
	hashedPassword := "$2a$10$1En6mrfnzK6PqAlRch5MzuP1k3e3gBcEvIYG4t8Zyayalx14Xs.Lu"
	check := hashing.ComparePassword([]byte(hashedPassword), []byte(password))

	// Correct password then return user role and true
	if check {
		fmt.Println(check)
		role := "sampleRole" // sample
		return role, check
	}

	return "", check
}
