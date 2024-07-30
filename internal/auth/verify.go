package auth

import (
	"fmt"

	"github.com/Build-D-An-Ki-n-Truc/auth/internal/jwtFunc"
)

// Claim structure
/*
	claim := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}
*/

// Check the token in request Header. Return true if token is valid and user have enough authority
func VerifyRequest(tokenString string, usernameAuth string, roleAuth string, roleRequired []string) (bool, error) {
	//extract the token
	claims, err := jwtFunc.ExtractToken(tokenString)
	if err != nil {
		return false, err
	}

	if claims["username"] != usernameAuth {
		return false, fmt.Errorf("username didn't match")
	}

	if claims["role"] != roleAuth {
		return false, fmt.Errorf("role didn't match")
	}
	checkRoleFunc := func(a string, list []string) bool {
		for _, b := range list {
			if b == a {
				return true
			}
		}
		return false
	}
	if !checkRoleFunc(roleAuth, roleRequired) {
		return false, fmt.Errorf("unathorized request")
	}

	return true, nil
}
