package security

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
)

const oidcHTTPTimeout = 10 * time.Second

// Verifier verifies a raw access token and returns application claims.
type Verifier interface {
	Verify(ctx context.Context, rawToken string) (Claims, error)
}

// OIDCVerifier verifies Keycloak tokens through OpenID Connect discovery.
type OIDCVerifier struct {
	verifier *oidc.IDTokenVerifier
	clientID string
}

// NewOIDCVerifier creates a verifier that trusts the configured CA certificate.
func NewOIDCVerifier(ctx context.Context, issuerURL string, clientID string, caCertPath string) (*OIDCVerifier, error) {
	httpClient, err := oidcHTTPClient(caCertPath)
	if err != nil {
		return nil, err
	}

	providerContext := oidc.ClientContext(ctx, httpClient)
	provider, err := oidc.NewProvider(providerContext, issuerURL)
	if err != nil {
		return nil, fmt.Errorf("discover OIDC provider: %w", err)
	}

	return &OIDCVerifier{
		verifier: provider.Verifier(&oidc.Config{SkipClientIDCheck: true}),
		clientID: clientID,
	}, nil
}

// Verify checks the token and decodes only the claims required by the API.
func (verifier *OIDCVerifier) Verify(ctx context.Context, rawToken string) (Claims, error) {
	token, err := verifier.verifier.Verify(ctx, rawToken)
	if err != nil {
		return Claims{}, fmt.Errorf("verify OIDC token: %w", err)
	}

	var claims Claims
	if err := token.Claims(&claims); err != nil {
		return Claims{}, fmt.Errorf("decode OIDC claims: %w", err)
	}
	if !matchesClient(token.Audience, claims.AuthorizedParty, verifier.clientID) {
		return Claims{}, errors.New("OIDC token was not issued for this client")
	}

	return claims, nil
}

func matchesClient(audience []string, authorizedParty string, clientID string) bool {
	for _, value := range audience {
		if value == clientID {
			return true
		}
	}

	return authorizedParty == clientID
}

func oidcHTTPClient(caCertPath string) (*http.Client, error) {
	certificatePEM, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("read OIDC CA certificate: %w", err)
	}

	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("load system CA certificates: %w", err)
	}
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	if ok := rootCAs.AppendCertsFromPEM(certificatePEM); !ok {
		return nil, errors.New("OIDC CA certificate does not contain a valid PEM certificate")
	}

	transport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return nil, errors.New("default HTTP transport has unexpected type")
	}
	transport = transport.Clone()
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    rootCAs,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   oidcHTTPTimeout,
	}, nil
}
