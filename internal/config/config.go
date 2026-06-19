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

// Load reads an optional .env file and validates the resulting environment.
func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	return loadFromEnvironment()
}

func loadFromEnvironment() (Config, error) {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = defaultPort
	}

	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	return Config{
		Port:        port,
		DatabaseURL: databaseURL,
	}, nil
}
