package auth

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
)

// TokenVerifier checks whether a bearer token is valid for this API.
type TokenVerifier interface {
	Verify(ctx context.Context, token string) error
}

// OIDCVerifier verifies bearer tokens against an OpenID Connect provider.
type OIDCVerifier struct {
	verifier *oidc.IDTokenVerifier
	clientID string
}

// NewOIDCVerifier creates a verifier from the provider discovery document.
func NewOIDCVerifier(ctx context.Context, issuerURL string, clientID string) (*OIDCVerifier, error) {
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, fmt.Errorf("create OIDC provider: %w", err)
	}

	return &OIDCVerifier{
		verifier: provider.Verifier(&oidc.Config{SkipClientIDCheck: true}),
		clientID: clientID,
	}, nil
}

// Verify checks the token signature, issuer, expiry and Keycloak client.
func (verifier *OIDCVerifier) Verify(ctx context.Context, token string) error {
	idToken, err := verifier.verifier.Verify(ctx, token)
	if err != nil {
		return fmt.Errorf("verify OIDC token: %w", err)
	}

	var claims struct {
		AuthorizedParty string `json:"azp"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return fmt.Errorf("read OIDC token claims: %w", err)
	}

	if hasAudience(idToken.Audience, verifier.clientID) || claims.AuthorizedParty == verifier.clientID {
		return nil
	}

	return fmt.Errorf("token is not issued for client %q", verifier.clientID)
}

func hasAudience(audiences []string, expected string) bool {
	for _, audience := range audiences {
		if audience == expected {
			return true
		}
	}

	return false
}
