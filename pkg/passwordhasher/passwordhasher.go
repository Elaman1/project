package passwordhasher

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	// кодируем salt + hash вместе
	saltEncoded := base64.RawStdEncoding.EncodeToString(salt)
	hashEncoded := base64.RawStdEncoding.EncodeToString(hash)
	return fmt.Sprintf("%s:%s", saltEncoded, hashEncoded), nil
}

func CheckPassword(password, hashWithSalt string) (bool, error) {
	// Разделяем строку на соль и хеш
	parts := strings.Split(hashWithSalt, ":")
	if len(parts) != 2 {
		return false, errors.New("invalid hash format")
	}

	saltEncoded := parts[0]
	hashEncoded := parts[1]

	// Декодируем base64 обратно в байты
	salt, err := base64.RawStdEncoding.DecodeString(saltEncoded)
	if err != nil {
		return false, err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(hashEncoded)
	if err != nil {
		return false, err
	}

	// Генерируем новый хеш с теми же параметрами
	computedHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, uint32(len(expectedHash)))

	// Сравниваем хеши
	if len(computedHash) != len(expectedHash) {
		return false, nil
	}

	// Побайтовое сравнение (безопасное)
	for i := range computedHash {
		if computedHash[i] != expectedHash[i] {
			return false, nil
		}
	}

	return true, nil
}
