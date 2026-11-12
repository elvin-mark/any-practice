package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	AppName          = "linux-practice"
	DefaultImageName = "linux-practice-ubuntu:24.04"
	DefaultUser      = "student"
)

// AppDir returns the root config directory for linux-practice (~/.config/linux-practice).
func AppDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolving user home directory: %w", err)
		}
		configDir = filepath.Join(homeDir, ".config")
	}

	appDir := filepath.Join(configDir, AppName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", fmt.Errorf("creating app config directory at %s: %w", appDir, err)
	}

	return appDir, nil
}

// StateFilePath returns the path to state.yaml.
func StateFilePath() (string, error) {
	appDir, err := AppDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(appDir, "state.yaml"), nil
}

// ScenariosDir returns the path for custom local user scenarios (~/.config/linux-practice/scenarios).
func ScenariosDir() (string, error) {
	appDir, err := AppDir()
	if err != nil {
		return "", err
	}
	scenariosDir := filepath.Join(appDir, "scenarios")
	if err := os.MkdirAll(scenariosDir, 0755); err != nil {
		return "", fmt.Errorf("creating scenarios directory: %w", err)
	}
	return scenariosDir, nil
}
