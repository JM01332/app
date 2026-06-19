package config

import "testing"

func TestLoadFromEnvironmentUsesDefaultPort(t *testing.T) {
	setRequiredEnvironment(t)
	t.Setenv("PORT", "")

	config, err := loadFromEnvironment()
	if err != nil {
		t.Fatalf("loadFromEnvironment() error = %v", err)
	}

	if config.Port != defaultPort {
		t.Errorf("Port = %q, want %q", config.Port, defaultPort)
	}
}

func TestLoadFromEnvironmentRequiresDatabaseURL(t *testing.T) {
	setRequiredEnvironment(t)
	t.Setenv("DATABASE_URL", "")

	_, err := loadFromEnvironment()
	if err == nil {
		t.Fatal("loadFromEnvironment() error = nil, want an error")
	}
}

func TestLoadFromEnvironmentRequiresOIDCSettings(t *testing.T) {
	testCases := []string{
		"OIDC_ISSUER_URL",
		"OIDC_CLIENT_ID",
		"OIDC_CA_CERT",
	}

	for _, environmentVariable := range testCases {
		t.Run(environmentVariable, func(t *testing.T) {
			setRequiredEnvironment(t)
			t.Setenv(environmentVariable, "")

			_, err := loadFromEnvironment()
			if err == nil {
				t.Fatalf("loadFromEnvironment() error = nil, want error for %s", environmentVariable)
			}
		})
	}
}

func TestLoadServerFromEnvironmentUsesConfiguredPort(t *testing.T) {
	t.Setenv("PORT", "9090")

	config := loadServerFromEnvironment()

	if config.Port != "9090" {
		t.Errorf("Port = %q, want 9090", config.Port)
	}
}

func setRequiredEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("OIDC_ISSUER_URL", "https://localhost:8843/realms/python")
	t.Setenv("OIDC_CLIENT_ID", "python-client")
	t.Setenv("OIDC_CA_CERT", "certificate.crt")
}
