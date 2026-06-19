package config

import "testing"

func TestLoadFromEnvironmentUsesDefaultPort(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_URL", "postgres://example")

	config, err := loadFromEnvironment()
	if err != nil {
		t.Fatalf("loadFromEnvironment() error = %v", err)
	}

	if config.Port != defaultPort {
		t.Errorf("Port = %q, want %q", config.Port, defaultPort)
	}
}

func TestLoadFromEnvironmentRequiresDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")

	_, err := loadFromEnvironment()
	if err == nil {
		t.Fatal("loadFromEnvironment() error = nil, want an error")
	}
}

func TestLoadServerFromEnvironmentUsesConfiguredPort(t *testing.T) {
	t.Setenv("PORT", "9090")

	config := loadServerFromEnvironment()

	if config.Port != "9090" {
		t.Errorf("Port = %q, want 9090", config.Port)
	}
}
