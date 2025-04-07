package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

func DecryptWithKey(ciphertext, key, iv []byte) (string, error) {
	// Create the cipher block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create the GCM cipher mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Decrypt the ciphertext using the IV
	plaintext, err := aesGCM.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
