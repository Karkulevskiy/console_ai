package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

var encryptionKey = []byte("your-32-byte-long-key-goes-here!") // 32 байта для AES-256

func Encrypt(text string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}
	plaintext := []byte(text)
	ctr := cipher.NewCTR(block, encryptionKey[:block.BlockSize()])
	ciphertext := make([]byte, len(plaintext))
	ctr.XORKeyStream(ciphertext, plaintext)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(cryptoText string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}
	ctr := cipher.NewCTR(block, encryptionKey[:block.BlockSize()])
	plaintext := make([]byte, len(ciphertext))
	ctr.XORKeyStream(plaintext, ciphertext)
	return string(plaintext), nil
}
