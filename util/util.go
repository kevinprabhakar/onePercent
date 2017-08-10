package util

import (
	"regexp"
	"golang.org/x/crypto/bcrypt"
)

func IsValidEmail(email string) (bool) {
	if (len(email) <= 0) {
		return false
	}
	return regexp.MustCompile(`^([a-zA-Z0-9_\.\-\+])+\@(([a-zA-Z0-9\-])+\.)+([a-zA-Z0-9]{2,4})+$`).MatchString(email)
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if (err != nil){
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}