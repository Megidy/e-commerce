package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "Password"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("error while hashing password : %v", err)
	}
	if hashedPassword == nil || string(hashedPassword) == (password) {
		t.Errorf("error while hashing password: expected to not empty or not to be as unhashed: %v", err)
	}
}

func TestComparePassword(t *testing.T) {
	password := "password"
	hashedPassword, _ := HashPassword(password)

	ok := ComparePassword(string(hashedPassword), password)

	if !ok {
		t.Errorf("failed password to be as planned ")
	}
}
