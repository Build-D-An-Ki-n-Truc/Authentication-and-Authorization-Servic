package auth

import (
	"github.com/Build-D-An-Ki-n-Truc/auth/internal/hashing"
)

// find hashedPassword with username in Database
func Login(username string, password string) (string, bool) {
	// sample password until have database
	hashedPassword := "sample"
	check := hashing.ComparePasswrod([]byte(hashedPassword), []byte(password))

	// Correct password then return user role and true
	if check {
		role := "sampleRole" // sample
		return role, check
	}
	return "", check
}
