package helmutil

import (
	"os"
	"path/filepath"
)

const (
	dirSuffix = "k8ssandra"
)

// GetCacheDir returns the caching directory for k8ssandra and creates it if it does not exists
func GetCacheDir(module string) (string, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	targetDir := filepath.Join(userCacheDir, dirSuffix, module)
	return targetDir, nil
}

// GetConfigDir returns the config directory for k8ssandra and creates it if it does not exists
func GetConfigDir(module string) (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	targetDir := filepath.Join(userConfigDir, dirSuffix, module)
	return targetDir, nil
}

func CreateIfNotExistsDir(targetDir string) (string, error) {
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return "", err
		}
	}
	return targetDir, nil
}
