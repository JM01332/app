package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JM01332/app/internal/app"
	"github.com/JM01332/app/internal/carrier/model"
	carrierservice "github.com/JM01332/app/internal/carrier/service"
	"github.com/JM01332/app/internal/security"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const oidcClientID = "python-client"

func TestCarrierListRequiresAuthentication(t *testing.T) {
	router := newIntegrationRouter([]string{"user"})

	request := httptest.NewRequest(http.MethodGet, "/api/carriers", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}

func TestCarrierListReturnsDataForUserRole(t *testing.T) {
	router := newIntegrationRouter([]string{"user"})

	request := httptest.NewRequest(http.MethodGet, "/api/carriers", nil)
	request.Header.Set("Authorization", "Bearer valid-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusOK)
	}

	var body []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
	if len(body) != 1 || body[0].ID != 1000 || body[0].Name != "USS Gerald R. Ford" {
		t.Fatalf("response body = %+v, want seeded carrier", body)
	}
}

func TestCarrierCreateRequiresAdminRole(t *testing.T) {
	router := newIntegrationRouter([]string{"user"})

	request := httptest.NewRequest(http.MethodPost, "/api/carriers", nil)
	request.Header.Set("Authorization", "Bearer valid-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusForbidden)
	}
}

type fakeCarrierService struct{}

func (service *fakeCarrierService) List(ctx context.Context) ([]model.Carrier, error) {
	return []model.Carrier{
		{
			ID:          1000,
			Name:        "USS Gerald R. Ford",
			Nation:      "USA",
			CarrierType: model.CarrierTypeAircraft,
		},
	}, nil
}

func (service *fakeCarrierService) GetByID(ctx context.Context, id int64) (*model.Carrier, error) {
	return &model.Carrier{ID: id, Name: "USS Gerald R. Ford"}, nil
}

func (service *fakeCarrierService) Create(ctx context.Context, input carrierservice.CreateCarrierInput) (*model.Carrier, error) {
	return &model.Carrier{ID: 1001, Name: input.Name}, nil
}

type fakeTokenVerifier struct {
	roles []string
}

func (verifier *fakeTokenVerifier) Verify(ctx context.Context, rawToken string) (security.Claims, error) {
	return security.Claims{
		ResourceAccess: map[string]security.ClientAccess{
			oidcClientID: {Roles: verifier.roles},
		},
	}, nil
}

func newIntegrationRouter(roles []string) *gin.Engine {
	gin.SetMode(gin.TestMode)

	return app.NewRouter(
		&fakeCarrierService{},
		zap.NewNop(),
		&fakeTokenVerifier{roles: roles},
		oidcClientID,
	)
}
