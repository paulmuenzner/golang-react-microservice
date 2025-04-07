package crypto

import "crypto/rand"

func GenerateRandomIV() ([]byte, error) {
	iv := make([]byte, 12) // GCM typically uses 12-byte IVs
	_, err := rand.Read(iv)
	if err != nil {
		return nil, err
	}
	return iv, nil
}
