package config

import "testing"

func TestLoadFromEnvironmentUsesDefaultPort(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("OIDC_ENABLED", "")

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
	t.Setenv("OIDC_ENABLED", "")

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

func TestLoadFromEnvironmentReadsOIDCConfig(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("OIDC_ENABLED", "true")
	t.Setenv("OIDC_ISSUER_URL", "http://localhost:8880/realms/python")
	t.Setenv("OIDC_CLIENT_ID", "python-client")

	config, err := loadFromEnvironment()
	if err != nil {
		t.Fatalf("loadFromEnvironment() error = %v", err)
	}

	if !config.OIDC.Enabled {
		t.Fatal("OIDC.Enabled = false, want true")
	}
	if config.OIDC.IssuerURL != "http://localhost:8880/realms/python" {
		t.Errorf("OIDC.IssuerURL = %q, want Keycloak issuer", config.OIDC.IssuerURL)
	}
	if config.OIDC.ClientID != "python-client" {
		t.Errorf("OIDC.ClientID = %q, want python-client", config.OIDC.ClientID)
	}
}

func TestLoadFromEnvironmentRequiresOIDCIssuerWhenEnabled(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("OIDC_ENABLED", "true")
	t.Setenv("OIDC_CLIENT_ID", "python-client")

	_, err := loadFromEnvironment()
	if err == nil {
		t.Fatal("loadFromEnvironment() error = nil, want an error")
	}
}

func TestLoadFromEnvironmentRejectsInvalidOIDCEnabled(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("OIDC_ENABLED", "sometimes")

	_, err := loadFromEnvironment()
	if err == nil {
		t.Fatal("loadFromEnvironment() error = nil, want an error")
	}
}
