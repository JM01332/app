package app

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JM01332/app/internal/carrier/model"
	carrierservice "github.com/JM01332/app/internal/carrier/service"
	"github.com/JM01332/app/internal/security"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(&fakeCarrierService{}, zap.NewNop(), &fakeTokenVerifier{}, "python-client")

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", response.Code, http.StatusOK)
	}

	contentType := response.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q, want application/json", contentType)
	}

	var body healthResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	if body.Status != "ok" {
		t.Errorf("status = %q, want ok", body.Status)
	}
}

func TestCarrierRoutesAreRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(&fakeCarrierService{}, zap.NewNop(), &fakeTokenVerifier{}, "python-client")

	request := httptest.NewRequest(http.MethodGet, "/api/carriers", nil)
	request.Header.Set("Authorization", "Bearer valid-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", response.Code, http.StatusOK)
	}
}

func TestCarrierRoutesRequireToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(&fakeCarrierService{}, zap.NewNop(), &fakeTokenVerifier{}, "python-client")

	request := httptest.NewRequest(http.MethodGet, "/api/carriers", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Errorf("status code = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}

func TestUserRoleCanReadCarriers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(
		&fakeCarrierService{},
		zap.NewNop(),
		&fakeTokenVerifier{roles: []string{"user"}},
		"python-client",
	)

	request := httptest.NewRequest(http.MethodGet, "/api/carriers", nil)
	request.Header.Set("Authorization", "Bearer user-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", response.Code, http.StatusOK)
	}
}

func TestUserRoleCannotCreateCarrier(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(
		&fakeCarrierService{},
		zap.NewNop(),
		&fakeTokenVerifier{roles: []string{"user"}},
		"python-client",
	)

	request := httptest.NewRequest(http.MethodPost, "/api/carriers", nil)
	request.Header.Set("Authorization", "Bearer user-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Errorf("status code = %d, want %d", response.Code, http.StatusForbidden)
	}
}

type fakeCarrierService struct{}

type fakeTokenVerifier struct {
	roles []string
}

func (verifier *fakeTokenVerifier) Verify(ctx context.Context, rawToken string) (security.Claims, error) {
	roles := verifier.roles
	if roles == nil {
		roles = []string{"admin"}
	}

	return security.Claims{
		ResourceAccess: map[string]security.ClientAccess{
			"python-client": {Roles: roles},
		},
	}, nil
}

func (service *fakeCarrierService) List(ctx context.Context) ([]model.Carrier, error) {
	return []model.Carrier{}, nil
}

func (service *fakeCarrierService) GetByID(ctx context.Context, id int64) (*model.Carrier, error) {
	return nil, nil
}

func (service *fakeCarrierService) Create(ctx context.Context, input carrierservice.CreateCarrierInput) (*model.Carrier, error) {
	return nil, nil
}
