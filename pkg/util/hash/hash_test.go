package hash_test

import (
	"encoding/base64"
	"testing"

	"github.com/ld100/goblet/pkg/util/hash"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestIsBcrypt(t *testing.T) {
	assert := assert.New(t)
	password := "12345zOMFG"
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	hashString := string(bytes)

	assert.True(hash.IsBcrypt(hashString), "IsBcrypt method does not hash as a valid hash")
	assert.False(hash.IsBcrypt(password), "IsBcrypt method treats hash as a valid hash")
}

func TestIsBase64(t *testing.T) {
	assert := assert.New(t)
	data := "12345zOMFG"
	hashString := base64.StdEncoding.EncodeToString([]byte(data))

	assert.True(hash.IsBase64(hashString), "IsBase64 method does not hash as a valid hash")
	assert.False(hash.IsBase64(data), "IsBase64 method treats hash as a valid hash")
}
