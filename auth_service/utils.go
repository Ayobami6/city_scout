package authservice

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateAPIKey() string {
	return uuid.New().String()

}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
