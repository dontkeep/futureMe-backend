package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
)

func Encrypt(plaintext string) (string, error) {
	keyB64 := os.Getenv("ENCRYPTION_KEY")
	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return "", errors.New("invalid encryption key encoding")
	}
	if len(key) != 32 {
		return "", errors.New("encryption key must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	plaintextBytes := []byte(plaintext)
	ciphertext := make([]byte, aes.BlockSize+len(plaintextBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintextBytes)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(ciphertextB64 string) (string, error) {
	keyB64 := os.Getenv("ENCRYPTION_KEY")
	key, err := base64.StdEncoding.DecodeString(keyB64)
	if err != nil {
		return "", errors.New("invalid encryption key encoding")
	}
	if len(key) != 32 {
		return "", errors.New("encryption key must be 32 bytes for AES-256")
	}
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}
