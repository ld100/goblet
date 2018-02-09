package hash

import (
	"encoding/base64"
	"strings"
)

func IsBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

func IsBcrypt(s string) bool {
	HASH_LENGTH := 60
	HASH_SIGNATURE := "$2"

	if len(s) == HASH_LENGTH && strings.HasPrefix(s, HASH_SIGNATURE) {
		return true
	}
	return false
}
