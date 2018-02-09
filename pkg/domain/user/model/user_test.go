package model

import (
	"testing"

	"github.com/ld100/goblet/pkg/util/securerandom"
)

func TestHash(t *testing.T) {
	password := "12345zOMFG"
	hash, _ := HashPassword(password)

	if password == hash {
		t.Errorf("Password %s should have hash not equal to the password itself", password)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "12345zOMFG"
	hash, _ := HashPassword(password)
	equals := CheckPasswordHash(password, hash)

	if !equals {
		t.Errorf("Password %s hash %s is incorrect", password, hash)
	}
}

func TestAutoHashPassword(t *testing.T) {
	uuid, _ := securerandom.Uuid()
	password := "12345zOMFG"
	user := &User{
		ID:       1,
		Uuid:     uuid,
		Email:    "test@example.com",
		Password: password,
	}

	err := user.BeforeUpdate()
	if err != nil {
		t.Errorf("Password autohashing does not work")
	} else {
		if user.Password == password {
			t.Errorf("Password autohashing does not work, password was not hashed")
		} else {
			// Password was autohashed successfully
			// running BeforeUpdate for the second time, it should not hash password again
			tempHash := user.Password
			err := user.BeforeUpdate()
			if err != nil {
				t.Errorf("Password autohashing does not work")
			} else {
				if tempHash != user.Password {
					t.Errorf("Password was rehashed while not needed")
				}
			}
		}
	}

}
