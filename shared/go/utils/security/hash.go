package security

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2Params defines the parameters for Argon2id hashing
type Argon2Params struct {
	Memory      uint32 // Memory in KiB (e.g., 64*1024 = 64 MB)
	Iterations  uint32 // Number of iterations (time cost)
	Parallelism uint8  // Number of threads
	SaltLength  uint32 // Length of salt in bytes
	KeyLength   uint32 // Length of derived key in bytes
}

// DefaultArgon2Params returns OWASP recommended parameters for Argon2id
// Memory: 64 MB, Iterations: 3, Parallelism: 4, Salt: 16 bytes, Key: 32 bytes
func DefaultArgon2Params() *Argon2Params {
	return &Argon2Params{
		Memory:      64 * 1024, // 64 MB
		Iterations:  3,
		Parallelism: 4,
		SaltLength:  16,
		KeyLength:   32,
	}
}

// HighSecurityArgon2Params returns parameters for high-security scenarios
// Memory: 256 MB, Iterations: 4, Parallelism: 4
func HighSecurityArgon2Params() *Argon2Params {
	return &Argon2Params{
		Memory:      256 * 1024, // 256 MB
		Iterations:  4,
		Parallelism: 4,
		SaltLength:  16,
		KeyLength:   32,
	}
}

// LowResourceArgon2Params returns parameters for resource-constrained environments
// Memory: 32 MB, Iterations: 2, Parallelism: 2
func LowResourceArgon2Params() *Argon2Params {
	return &Argon2Params{
		Memory:      32 * 1024, // 32 MB
		Iterations:  2,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}

// HashPassword hashes a password using Argon2id with the given parameters
// Returns the hash in PHC string format: $argon2id$v=19$m=65536,t=3,p=4$salt$hash
func HashPassword(password string, params *Argon2Params) (string, error) {
	// Generate random salt
	salt := make([]byte, params.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Generate hash
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	// Encode to PHC string format
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.Memory,
		params.Iterations,
		params.Parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

// HashPasswordDefault hashes a password using default OWASP recommended parameters
func HashPasswordDefault(password string) (string, error) {
	return HashPassword(password, DefaultArgon2Params())
}

// VerifyPassword verifies a password against an Argon2id hash
// Returns true if the password matches, false otherwise
func VerifyPassword(password, encodedHash string) (bool, error) {
	// Parse the encoded hash
	params, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Generate hash from provided password
	testHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	// Use constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare(hash, testHash) == 1 {
		return true, nil
	}

	return false, nil
}

// decodeHash parses a PHC string format hash
func decodeHash(encodedHash string) (*Argon2Params, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, errors.New("invalid hash format")
	}

	if parts[1] != "argon2id" {
		return nil, nil, nil, errors.New("incompatible hash algorithm")
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, fmt.Errorf("invalid version: %w", err)
	}

	if version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible argon2 version")
	}

	params := &Argon2Params{}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d",
		&params.Memory,
		&params.Iterations,
		&params.Parallelism,
	); err != nil {
		return nil, nil, nil, fmt.Errorf("invalid parameters: %w", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("invalid salt: %w", err)
	}
	params.SaltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("invalid hash: %w", err)
	}
	params.KeyLength = uint32(len(hash))

	return params, salt, hash, nil
}

// HashSHA256 creates a SHA-256 hash of input
// WARNING: DO NOT USE FOR PASSWORDS! Only for data integrity checks
func HashSHA256(input string) string {
	// Kept for non-password use cases (checksums, etc.)
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
