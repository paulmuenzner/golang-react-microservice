package crypto

import "crypto/rand"

type AllowedByteSize int

const (
	ByteSize18  AllowedByteSize = 18
	ByteSize32  AllowedByteSize = 32
	ByteSize64  AllowedByteSize = 64
	ByteSize128 AllowedByteSize = 128
	ByteSize256 AllowedByteSize = 256
)

// GenerateKey generates a key with the specified byte size
func (bs AllowedByteSize) GenerateKey() ([]byte, error) {
	// Cast bs to int (size of byte array)
	byteKey := make([]byte, int(bs))

	// Fill the byte slice with random data
	_, err := rand.Read(byteKey)
	if err != nil {
		return nil, err
	}
	return byteKey, nil
}
