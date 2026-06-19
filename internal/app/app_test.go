package app

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JM01332/app/internal/carrier/model"
	carrierservice "github.com/JM01332/app/internal/carrier/service"
	"github.com/gin-gonic/gin"
)

func TestHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(nil)

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
	router := NewRouter(&fakeCarrierService{})

	request := httptest.NewRequest(http.MethodGet, "/api/carriers", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", response.Code, http.StatusOK)
	}
}

type fakeCarrierService struct{}

func (service *fakeCarrierService) List(ctx context.Context) ([]model.Carrier, error) {
	return []model.Carrier{}, nil
}

func (service *fakeCarrierService) GetByID(ctx context.Context, id int64) (*model.Carrier, error) {
	return nil, nil
}

func (service *fakeCarrierService) Create(ctx context.Context, input carrierservice.CreateCarrierInput) (*model.Carrier, error) {
	return nil, nil
}
