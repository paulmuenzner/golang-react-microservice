package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// EncryptAES encrypts plaintext using AES-GCM (Galois/Counter Mode).
// It generates a random Nonce and prepends it to the ciphertext before Base64 encoding.
func EncryptAES(key []byte, plaintext string) (string, error) {
	// Create a new cipher block from the key. This initializes the AES algorithm.
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create the GCM cipher mode instance. GCM provides authenticated encryption.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a byte slice for the Nonce (Initialization Vector) of the required size.
	nonce := make([]byte, gcm.NonceSize())

	// Fill the Nonce with cryptographically secure random bytes.
	// This Nonce MUST be unique for every encryption using the same key.
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the plaintext. The nonce is prepended to the ciphertext
	// because the receiver (DecryptAES) needs it for decryption and authentication.
	// The last argument (additionalData) is nil because it's not used here.
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Encode the final result (Nonce + Ciphertext) to a Base64 string for safe transmission.
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES decrypts ciphertext using AES-GCM (Galois/Counter Mode).
// The input ciphertext is expected to be a Base64-encoded string containing
// the nonce followed immediately by the actual encrypted data.
func DecryptAES(key []byte, ciphertext string) (string, error) {
	// Decode the Base64-encoded ciphertext string into a byte slice (data).
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// Create a new cipher block from the key. This establishes the AES algorithm.
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create the GCM cipher mode instance. GCM provides authenticated encryption.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Determine the size of the Nonce (a fixed size for GCM).
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// Split the 'data' byte slice into the Nonce (at the start) and the actual encrypted data.
	nonce, encryptedData := data[:nonceSize], data[nonceSize:]

	// Decrypt the data and authenticate the integrity using the Nonce.
	// The last argument (additionalData) is nil because it's not used here.
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		// This error often indicates a failed authentication tag (tampered data).
		return "", err
	}

	// Convert the decrypted byte slice back to a string and return.
	return string(plaintext), nil
}
