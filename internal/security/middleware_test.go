package security

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthenticateAcceptsValidBearerToken(t *testing.T) {
	verifier := &fakeVerifier{claims: adminClaims()}
	router := securityTestRouter(Authenticate(verifier))

	response := performSecurityRequest(router, "Bearer valid-token")

	if response.Code != http.StatusNoContent {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusNoContent)
	}
	if verifier.rawToken != "valid-token" {
		t.Errorf("verified token = %q, want valid-token", verifier.rawToken)
	}
}

func TestAuthenticateRejectsMissingToken(t *testing.T) {
	router := securityTestRouter(Authenticate(&fakeVerifier{}))

	response := performSecurityRequest(router, "")

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusUnauthorized)
	}
	if response.Header().Get("WWW-Authenticate") != "Bearer" {
		t.Error("WWW-Authenticate header is missing")
	}
}

func TestAuthenticateRejectsInvalidToken(t *testing.T) {
	router := securityTestRouter(Authenticate(&fakeVerifier{err: errors.New("invalid signature")}))

	response := performSecurityRequest(router, "Bearer invalid-token")

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}

func TestRequireAnyRoleAcceptsNormalizedAdminRole(t *testing.T) {
	verifier := &fakeVerifier{claims: adminClaims()}
	router := securityTestRouter(
		Authenticate(verifier),
		RequireAnyRole("python-client", RoleUser, RoleAdmin),
	)

	response := performSecurityRequest(router, "Bearer valid-token")

	if response.Code != http.StatusNoContent {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusNoContent)
	}
}

func TestRequireAnyRoleRejectsMissingRole(t *testing.T) {
	verifier := &fakeVerifier{claims: Claims{
		ResourceAccess: map[string]ClientAccess{
			"python-client": {Roles: []string{"user"}},
		},
	}}
	router := securityTestRouter(
		Authenticate(verifier),
		RequireAnyRole("python-client", RoleAdmin),
	)

	response := performSecurityRequest(router, "Bearer valid-token")

	if response.Code != http.StatusForbidden {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusForbidden)
	}
}

type fakeVerifier struct {
	claims   Claims
	err      error
	rawToken string
}

func (verifier *fakeVerifier) Verify(ctx context.Context, rawToken string) (Claims, error) {
	verifier.rawToken = rawToken
	return verifier.claims, verifier.err
}

func adminClaims() Claims {
	return Claims{
		Subject:           "subject-1",
		PreferredUsername: "admin",
		ResourceAccess: map[string]ClientAccess{
			"python-client": {Roles: []string{"admin"}},
		},
	}
}

func securityTestRouter(middleware ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", append(middleware, func(context *gin.Context) {
		context.Status(http.StatusNoContent)
	})...)
	return router
}

func performSecurityRequest(router *gin.Engine, authorizationHeader string) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	if authorizationHeader != "" {
		request.Header.Set("Authorization", authorizationHeader)
	}
	router.ServeHTTP(response, request)
	return response
}
