package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hashed := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return fmt.Sprintf("%s$%s", base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hashed)), nil
}

func VerifyPassword(password, hash string) bool {
	parts := strings.Split(hash, "$")
	if len(parts) != 2 {
		return false
	}
	salt, _ := base64.RawStdEncoding.DecodeString(parts[0])
	storedHash, _ := base64.RawStdEncoding.DecodeString(parts[1])
	hashed := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return subtle.ConstantTimeCompare(hashed, storedHash) == 1
}
