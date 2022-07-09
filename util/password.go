package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashPwd), nil
}

func CheckPassword(password string, hashPwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(password))
}
