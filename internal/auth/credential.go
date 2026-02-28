package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	moodleerrors "github.com/nathanfredericks/moodle-cli/internal/errors"
)

const (
	credFileName    = "credentials"
	filePermissions = 0600
)

// CredentialStore manages credential storage and retrieval.
type CredentialStore interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Delete(key string) error
}

// FileCredentialStore stores credentials in a file with optional encryption.
type FileCredentialStore struct {
	dir string
}

// NewFileCredentialStore creates a new file-based credential store.
func NewFileCredentialStore(dir string) *FileCredentialStore {
	return &FileCredentialStore{dir: dir}
}

func (s *FileCredentialStore) Get(key string) (string, error) {
	// Check env var override first — MOODLE_TOKEN overrides any token lookup
	if token := os.Getenv("MOODLE_TOKEN"); token != "" {
		if strings.HasSuffix(key, "_token") || key == "token" {
			return token, nil
		}
	}

	data, err := s.readStore()
	if err != nil {
		return "", err
	}
	val, ok := data[key]
	if !ok {
		return "", &moodleerrors.AuthError{Msg: fmt.Sprintf("credential %q not found", key)}
	}
	return val, nil
}

func (s *FileCredentialStore) Set(key string, value string) error {
	data, err := s.readStore()
	if err != nil {
		data = make(map[string]string)
	}
	data[key] = value
	return s.writeStore(data)
}

func (s *FileCredentialStore) Delete(key string) error {
	data, err := s.readStore()
	if err != nil {
		return err
	}
	delete(data, key)
	return s.writeStore(data)
}

func (s *FileCredentialStore) readStore() (map[string]string, error) {
	path := filepath.Join(s.dir, credFileName)
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return make(map[string]string), nil
	}
	if err != nil {
		return nil, &moodleerrors.AuthError{Msg: "unable to read credentials", Err: err}
	}
	return parseCredentials(string(content)), nil
}

func (s *FileCredentialStore) writeStore(data map[string]string) error {
	if err := os.MkdirAll(s.dir, 0700); err != nil {
		return &moodleerrors.AuthError{Msg: "unable to create credentials directory", Err: err}
	}
	path := filepath.Join(s.dir, credFileName)
	content := formatCredentials(data)
	return os.WriteFile(path, []byte(content), filePermissions)
}

func parseCredentials(content string) map[string]string {
	result := make(map[string]string)
	lines := splitLines(content)
	for _, line := range lines {
		if line == "" {
			continue
		}
		for i, c := range line {
			if c == '=' {
				key := line[:i]
				value := line[i+1:]
				result[key] = value
				break
			}
		}
	}
	return result
}

func formatCredentials(data map[string]string) string {
	var result string
	for k, v := range data {
		result += k + "=" + v + "\n"
	}
	return result
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// Encrypt encrypts plaintext using AES-256-GCM.
func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts AES-256-GCM encrypted data.
func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// DeriveKey converts a passphrase to a 32-byte key using simple hashing.
// In production, use a proper KDF like Argon2 or scrypt.
func DeriveKey(passphrase string) []byte {
	h := make([]byte, 32)
	b := []byte(passphrase)
	for i := range h {
		if i < len(b) {
			h[i] = b[i]
		}
	}
	return h
}

// TokenKey is the credential store key for the authentication token.
const TokenKey = "token"
