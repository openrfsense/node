package nats

import (
	"os"
	"path/filepath"
)

var DefaultTokenFile = ""

// Writes token to token file.
func SaveTokenFile(path string, token string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(token), 0o755)
}

// Retrieves token from a token file.
func GetToken(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// TODO: ensure token is not kept in memory
	data, err := os.ReadFile(path)
	return string(data), err
}
