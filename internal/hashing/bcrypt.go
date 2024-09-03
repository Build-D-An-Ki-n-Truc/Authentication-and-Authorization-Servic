package hashing

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// generate HashedPassword from password ([]byte)
func GenerateHash(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	if err != nil {
		logrus.Println(err)
		return nil, err
	}

	return hashedPassword, nil
}

// Comparing Password
func ComparePassword(hashedPassword []byte, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	return err == nil
}
