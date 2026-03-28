package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	// Default parameters for Argon2id
	defaultMemory      = 64 * 1024 // 64 MB
	defaultIterations  = 3
	defaultParallelism = 2
	defaultSaltLength  = 16
	defaultKeyLength   = 32
)

var (
	ErrInvalidHash         = errors.New("invalid hash format")
	ErrIncompatibleVersion = errors.New("incompatible hash version")
)

// HashConfig holds the configuration for password hashing
type HashConfig struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultHashConfig returns the default hash configuration
func DefaultHashConfig() *HashConfig {
	return &HashConfig{
		Memory:      defaultMemory,
		Iterations:  defaultIterations,
		Parallelism: defaultParallelism,
		SaltLength:  defaultSaltLength,
		KeyLength:   defaultKeyLength,
	}
}

// HashPassword hashes a password using Argon2id
func HashPassword(password string) (string, error) {
	return HashPasswordWithConfig(password, DefaultHashConfig())
}

// HashPasswordWithConfig hashes a password using Argon2id with custom config
func HashPasswordWithConfig(password string, config *HashConfig) (string, error) {
	// Generate a random salt
	salt := make([]byte, config.SaltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Generate the hash
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		config.Iterations,
		config.Memory,
		config.Parallelism,
		config.KeyLength,
	)

	// Encode the hash in the standard format
	// $argon2id$v=19$m=65536,t=3,p=2$salt$hash
	encodedHash := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		config.Memory,
		config.Iterations,
		config.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return encodedHash, nil
}

// VerifyPassword verifies a password against a hash
func VerifyPassword(password, encodedHash string) (bool, error) {
	// Parse the encoded hash
	config, salt, hash, err := parseEncodedHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Generate the hash from the provided password
	computedHash := argon2.IDKey(
		[]byte(password),
		salt,
		config.Iterations,
		config.Memory,
		config.Parallelism,
		config.KeyLength,
	)

	// Use constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare(hash, computedHash) == 1 {
		return true, nil
	}
	return false, nil
}

// parseEncodedHash parses an encoded hash string
func parseEncodedHash(encodedHash string) (*HashConfig, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse version: %w", err)
	}
	if version != 19 {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	config := &HashConfig{}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d",
		&config.Memory,
		&config.Iterations,
		&config.Parallelism,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse parameters: %w", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to decode salt: %w", err)
	}
	config.SaltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to decode hash: %w", err)
	}
	config.KeyLength = uint32(len(hash))

	return config, salt, hash, nil
}

// GenerateRandomString generates a random string of the specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// GenerateRandomBytes generates random bytes of the specified length
func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return bytes, nil
}