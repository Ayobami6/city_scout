package authservice

import "github.com/google/uuid"

func GenerateAPIKey() string {
	return uuid.New().String()

}

func HashPassword(password string) string {
	return password
}
