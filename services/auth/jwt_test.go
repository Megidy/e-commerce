package auth

import (
	"testing"

	"github.com/google/uuid"
)

func TestCreateJWT(t *testing.T) {
	secret := "superpupersecretkey"
	userid := uuid.NewString()
	token, err := CreateJWT([]byte(secret), userid)

	if err != nil {
		t.Errorf("error while creating token : %v", err)
	}
	if token == "" {
		t.Errorf("expected token to be not empty , error: %v", token)
	}

}
