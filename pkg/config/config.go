package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/caarlos0/env/v9"
	"github.com/rk-the-dev/golib-core/pkg/logger"
	"gopkg.in/yaml.v2"
)

const (
	FileJSON = "json"
	FileYAML = "yaml"
)

// LoadConfig loads configuration from environment variables into the given struct
func LoadConfig(config interface{}) error {
	if config == nil {
		return fmt.Errorf("config object cannot be nil")
	}
	err := env.Parse(config)
	if err != nil {
		logger.Error("Failed to load environment configuration", map[string]interface{}{"error": err})
		return err
	}
	logger.Info("Environment configuration loaded successfully", nil)
	return nil
}

// LoadConfig loads configuration from ENV, JSON, and YAML
func LoadConfigFromFile(config interface{}, filePath string, fileType string) error {
	if config == nil {
		return fmt.Errorf("config object cannot be nil")
	}
	// Load from YAML file (if provided)
	if fileType == FileYAML {
		err := loadFromYAML(config, filePath)
		if err != nil {
			logger.Warn("Could not load YAML config, proceeding with JSON/ENV", map[string]interface{}{"error": err})
		}
	}
	// Load from JSON file (if provided)
	if fileType == FileJSON {
		err := loadFromJSON(config, filePath)
		if err != nil {
			logger.Warn("Could not load JSON config, proceeding with ENV", map[string]interface{}{"error": err})
		}
	}
	// Load from environment variables (overrides JSON and YAML)
	err := env.Parse(config)
	if err != nil {
		logger.Error("Failed to load environment configuration", map[string]interface{}{"error": err})
		return err
	}
	logger.Info("Configuration loaded successfully", nil)
	return nil
}

// loadFromJSON reads and unmarshals JSON config into the given struct
func loadFromJSON(config interface{}, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open JSON config file: %w", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return fmt.Errorf("failed to decode JSON config file: %w", err)
	}
	logger.Info("Loaded configuration from JSON file", map[string]interface{}{"file": filePath})
	return nil
}

// loadFromYAML reads and unmarshals YAML config into the given struct
func loadFromYAML(config interface{}, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open YAML config file: %w", err)
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return fmt.Errorf("failed to decode YAML config file: %w", err)
	}
	logger.Info("Loaded configuration from YAML file", map[string]interface{}{"file": filePath})
	return nil
}
