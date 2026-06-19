package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const defaultPort = "8080"

// Config contains the environment-dependent application settings.
type Config struct {
	Port        string
	DatabaseURL string
}

// ServerConfig contains settings needed before database initialization.
type ServerConfig struct {
	Port string
}

// Load reads an optional .env file and validates the resulting environment.
func Load() (Config, error) {
	if err := loadEnvironmentFile(); err != nil {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	return loadFromEnvironment()
}

// LoadServer reads and validates settings required to start the HTTP server.
func LoadServer() (ServerConfig, error) {
	if err := loadEnvironmentFile(); err != nil {
		return ServerConfig{}, fmt.Errorf("load .env: %w", err)
	}

	return loadServerFromEnvironment(), nil
}

func loadEnvironmentFile() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

func loadFromEnvironment() (Config, error) {
	serverConfig := loadServerFromEnvironment()

	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	return Config{
		Port:        serverConfig.Port,
		DatabaseURL: databaseURL,
	}, nil
}

func loadServerFromEnvironment() ServerConfig {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = defaultPort
	}

	return ServerConfig{Port: port}
}
