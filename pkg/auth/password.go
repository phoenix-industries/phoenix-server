package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/scrypt"
)

const (
	scryptMemCost = 1 << 15
	scryptRounds  = 8
	scryptSaltLen = 16
	scryptKeyLen  = 32
)

func (a *Auth) Hash(password string) (string, string, error) {
	salt, err := a.GenerateSalt()
	if err != nil {
		return "", "", err
	}
	return a.HashWithSalt([]byte(password), salt)
}

func (a *Auth) HashWithSalt(password, salt []byte) (string, string, error) {
	hash, err := scrypt.Key(
		password,
		append(salt, a.passwordSaltSeparator...),
		scryptMemCost,
		scryptRounds,
		1,
		scryptKeyLen,
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash password: %w", err)
	}

	block, err := aes.NewCipher(hash)
	if err != nil {
		return "", "", fmt.Errorf("failed to create cipher: %w", err)
	}

	cipherBytes := make([]byte, aes.BlockSize+len(a.passwordSigningKey))
	blk, key := cipherBytes[:aes.BlockSize], cipherBytes[aes.BlockSize:]

	stream := cipher.NewCTR(block, blk)
	stream.XORKeyStream(key, a.passwordSigningKey)

	encodedHash := base64.StdEncoding.EncodeToString(key)
	encodedSalt := base64.StdEncoding.EncodeToString(salt)

	return encodedHash, encodedSalt, nil
}

func (a *Auth) GenerateSalt() ([]byte, error) {
	salt := make([]byte, scryptSaltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

func (a *Auth) Compare(password, hash, salt string) (bool, error) {
	ph, _, err := a.HashWithSalt([]byte(password), []byte(salt))
	if err != nil {
		return false, err
	}
	return ph == hash, nil
}
