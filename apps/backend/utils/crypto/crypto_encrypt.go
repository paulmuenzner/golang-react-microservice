package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

func EncryptWithKey(plaintext string, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	// Correct IV size for AES-GCM is 12 bytes
	iv, err := GenerateRandomIV() // GenerateRandomIV() should return a 12-byte IV for AES-GCM
	if err != nil {
		return nil, nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	// Encrypt the plaintext
	ciphertext := aesGCM.Seal(nil, iv, []byte(plaintext), nil)

	return ciphertext, iv, nil
}
