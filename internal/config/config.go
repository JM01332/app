package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const defaultPort = "8080"

// Config contains the environment-dependent application settings.
type Config struct {
	Port        string
	DatabaseURL string
	OIDC        OIDCConfig
}

// OIDCConfig contains optional OpenID Connect settings.
type OIDCConfig struct {
	Enabled   bool
	IssuerURL string
	ClientID  string
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

	oidcConfig, err := loadOIDCFromEnvironment()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Port:        serverConfig.Port,
		DatabaseURL: databaseURL,
		OIDC:        oidcConfig,
	}, nil
}

func loadServerFromEnvironment() ServerConfig {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = defaultPort
	}

	return ServerConfig{Port: port}
}

func loadOIDCFromEnvironment() (OIDCConfig, error) {
	enabledText := strings.TrimSpace(os.Getenv("OIDC_ENABLED"))
	if enabledText == "" {
		return OIDCConfig{}, nil
	}

	enabled, err := strconv.ParseBool(enabledText)
	if err != nil {
		return OIDCConfig{}, fmt.Errorf("OIDC_ENABLED must be true or false: %w", err)
	}
	if !enabled {
		return OIDCConfig{Enabled: false}, nil
	}

	issuerURL := strings.TrimSpace(os.Getenv("OIDC_ISSUER_URL"))
	if issuerURL == "" {
		return OIDCConfig{}, errors.New("OIDC_ISSUER_URL is required when OIDC is enabled")
	}

	clientID := strings.TrimSpace(os.Getenv("OIDC_CLIENT_ID"))
	if clientID == "" {
		return OIDCConfig{}, errors.New("OIDC_CLIENT_ID is required when OIDC is enabled")
	}

	return OIDCConfig{
		Enabled:   true,
		IssuerURL: issuerURL,
		ClientID:  clientID,
	}, nil
}
