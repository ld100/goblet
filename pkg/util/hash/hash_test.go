package hash

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestIsBcrypt(t *testing.T) {
	assert := assert.New(t)
	password := "12345zOMFG"
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	hash := string(bytes)

	assert.True(IsBcrypt(hash), "IsBcrypt method does not hash as a valid hash")
	assert.False(IsBcrypt(password), "IsBcrypt method treats hash as a valid hash")
}

func TestIsBase64(t *testing.T) {
	assert := assert.New(t)
	data := "12345zOMFG"
	hash := base64.StdEncoding.EncodeToString([]byte(data))

	assert.True(IsBase64(hash), "IsBase64 method does not hash as a valid hash")
	assert.False(IsBase64(data), "IsBase64 method treats hash as a valid hash")
}
