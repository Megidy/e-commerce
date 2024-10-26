package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func ComparePassword(hashedPassword, defaultPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(defaultPassword))
	return err == nil
}
