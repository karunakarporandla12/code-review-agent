package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

func deriveKey(secret string) []byte {
	normalized := strings.TrimSpace(secret)
	if normalized == "" {
		normalized = "code-review-agent-default-encryption-key"
	}
	sum := sha256.Sum256([]byte(normalized))
	return sum[:]
}

func EncryptString(secret, plaintext string) (string, error) {
	if strings.TrimSpace(plaintext) == "" {
		return "", nil
	}

	block, err := aes.NewCipher(deriveKey(secret))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptString(secret, encryptedText string) (string, error) {
	trimmed := strings.TrimSpace(encryptedText)
	if trimmed == "" {
		return "", nil
	}

	raw, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil {
		return encryptedText, nil
	}

	block, err := aes.NewCipher(deriveKey(secret))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(raw) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := raw[:nonceSize], raw[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return encryptedText, nil
	}
	return string(plaintext), nil
}
