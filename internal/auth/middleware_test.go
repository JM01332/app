package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMiddlewareRejectsMissingBearerToken(t *testing.T) {
	router := testRouter(&fakeVerifier{})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}

func TestMiddlewareRejectsInvalidBearerToken(t *testing.T) {
	router := testRouter(&fakeVerifier{err: errors.New("invalid token")})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer invalid")
	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}

func TestMiddlewareAcceptsValidBearerToken(t *testing.T) {
	verifier := &fakeVerifier{}
	router := testRouter(verifier)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer valid-token")
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusOK)
	}
	if verifier.token != "valid-token" {
		t.Fatalf("verified token = %q, want valid-token", verifier.token)
	}
}

type fakeVerifier struct {
	token string
	err   error
}

func (verifier *fakeVerifier) Verify(ctx context.Context, token string) error {
	verifier.token = token
	return verifier.err
}

func testRouter(verifier TokenVerifier) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Middleware(verifier))
	router.GET("/protected", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})

	return router
}
