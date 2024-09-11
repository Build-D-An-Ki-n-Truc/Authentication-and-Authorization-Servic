package auth

import (
	"github.com/Build-D-An-Ki-n-Truc/auth/internal/db/mongodb"
	"github.com/Build-D-An-Ki-n-Truc/auth/internal/hashing"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Check for username and password in corresponding database. Return Role and True if password correct, False otherwise
func Login(username string, password string) (string, bool) {

	// Get hashed password from database
	user, err := mongodb.ReadUser(username)
	if err != nil {
		logrus.Println("Error reading user from database: ", err)
		return "", false
	}

	check := hashing.ComparePassword([]byte(user.Password), []byte(password))

	// Correct password then return user role and true
	if check {
		logrus.Println(check)
		role := user.Role
		return role, check
	}

	return "", check
}

// Check for username and password in corresponding database. Return Role and True if password correct, False otherwise
func LoginBrand(username string, password string) (string, bool, string) {

	// Get hashed password from database
	user, err := mongodb.ReadUser(username)
	if err != nil {
		logrus.Println("Error reading user from database: ", err)
		return "", false, ""
	}

	check := hashing.ComparePassword([]byte(user.Password), []byte(password))

	// Correct password then return user role and true
	if check {
		logrus.Println(check)
		role := user.Role
		if user.BrandID != primitive.NilObjectID {
			return role, check, user.BrandID.Hex()
		}
		return role, check, ""
	}

	return "", check, ""
}
