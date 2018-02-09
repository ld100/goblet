package hash

import (
	"encoding/base64"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestIsBcrypt(t *testing.T) {
	password := "12345zOMFG"
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	hash := string(bytes)

	if !IsBcrypt(hash) {
		t.Errorf("IsBcrypt method does not treat %s as a valid hash", hash)
	}

	if IsBcrypt(password) {
		t.Errorf("IsBcrypt method treats %s as a valid hash", password)
	}
}

func TestIsBase64(t *testing.T) {
	data := "12345zOMFG"
	hash := base64.StdEncoding.EncodeToString([]byte(data))

	if !IsBase64(hash) {
		t.Errorf("IsBase64 method does not treat %s as a valid hash", hash)
	}

	if IsBase64(data) {
		t.Errorf("IsBase64 method treats %s as a valid hash", data)
	}
}
