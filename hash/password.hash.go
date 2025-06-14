package hash

import (
	"crypto/sha512"
	"encoding/hex"
)

func HashPassword(password string) string {
	hasher := sha512.New()
	hasher.Write([]byte(password))

	hashedPassword := hasher.Sum(nil)
	return hex.EncodeToString(hashedPassword)
}

func ComparePasswords(hashedPassword, password string) bool {
	return HashPassword(password) == hashedPassword
}
