package router

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/JM01332/app/internal/carrier/model"
	carrierservice "github.com/JM01332/app/internal/carrier/service"
	"github.com/gin-gonic/gin"
)

func TestListReturnsCarriers(t *testing.T) {
	router := newTestRouter(&fakeCarrierService{
		listResult: []model.Carrier{testCarrier()},
	})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/carriers", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusOK)
	}

	var body []CarrierResponse
	decodeResponse(t, response, &body)
	if len(body) != 1 {
		t.Fatalf("response length = %d, want 1", len(body))
	}
	if body[0].ID != 1000 {
		t.Fatalf("id = %d, want 1000", body[0].ID)
	}
}

func TestListReturnsInternalError(t *testing.T) {
	router := newTestRouter(&fakeCarrierService{
		listError: errors.New("database unavailable"),
	})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/carriers", nil)
	router.ServeHTTP(response, request)

	assertErrorResponse(t, response, http.StatusInternalServerError, errorCodeInternal)
}

func TestCreateReturnsCreatedCarrier(t *testing.T) {
	service := &fakeCarrierService{
		createResult: ptr(testCarrier()),
	}
	router := newTestRouter(service)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/carriers", bytes.NewBufferString(validCreateCarrierJSON()))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusCreated)
	}
	if location := response.Header().Get("Location"); location != "/api/carriers/1000" {
		t.Fatalf("Location = %q, want /api/carriers/1000", location)
	}
	if service.createInput.Name != "Enterprise" {
		t.Fatalf("create input name = %q, want Enterprise", service.createInput.Name)
	}
	if service.createInput.CarrierType != model.CarrierTypeAircraft {
		t.Fatalf("create input carrier type = %q, want %q", service.createInput.CarrierType, model.CarrierTypeAircraft)
	}
	if service.createInput.CommandCenter.CodeName != "Bridge" {
		t.Fatalf("create input command center = %+v, want Bridge", service.createInput.CommandCenter)
	}
	if len(service.createInput.Aircrafts) != 1 || service.createInput.Aircrafts[0].Model != "F/A-18" {
		t.Fatalf("create input aircrafts = %+v, want F/A-18", service.createInput.Aircrafts)
	}
}

func TestCreateRejectsInvalidJSON(t *testing.T) {
	router := newTestRouter(&fakeCarrierService{})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/carriers", bytes.NewBufferString(`{"name":`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)

	assertErrorResponse(t, response, http.StatusBadRequest, errorCodeInvalidJSON)
}

func TestCreateRejectsUnknownJSONField(t *testing.T) {
	router := newTestRouter(&fakeCarrierService{})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/carriers", bytes.NewBufferString(`{"unknown":"field"}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)

	assertErrorResponse(t, response, http.StatusBadRequest, errorCodeInvalidJSON)
}

func TestCreateRejectsValidationError(t *testing.T) {
	router := newTestRouter(&fakeCarrierService{})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/carriers", bytes.NewBufferString(`{"name":"","nation":"U","carrierType":"BATTLESHIP","commandCenter":{"codeName":"","securityLevel":6},"aircrafts":[]}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)

	assertErrorResponse(t, response, http.StatusUnprocessableEntity, errorCodeValidationFailed)
}

func TestCreateReturnsConflictForDuplicateName(t *testing.T) {
	router := newTestRouter(&fakeCarrierService{
		createError: carrierservice.ErrCarrierNameExists,
	})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/carriers", bytes.NewBufferString(validCreateCarrierJSON()))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)

	assertErrorResponse(t, response, http.StatusConflict, errorCodeCarrierNameExists)
}

func TestCreateReturnsInternalError(t *testing.T) {
	router := newTestRouter(&fakeCarrierService{
		createError: errors.New("database unavailable"),
	})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/carriers", bytes.NewBufferString(validCreateCarrierJSON()))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)

	assertErrorResponse(t, response, http.StatusInternalServerError, errorCodeInternal)
}

type fakeCarrierService struct {
	listResult   []model.Carrier
	listError    error
	createResult *model.Carrier
	createError  error
	createInput  carrierservice.CreateCarrierInput
}

func (service *fakeCarrierService) List(ctx context.Context) ([]model.Carrier, error) {
	return service.listResult, service.listError
}

func (service *fakeCarrierService) Create(ctx context.Context, input carrierservice.CreateCarrierInput) (*model.Carrier, error) {
	service.createInput = input
	return service.createResult, service.createError
}

func newTestRouter(service CarrierService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api")
	RegisterRoutes(api, service)

	return router
}

func testCarrier() model.Carrier {
	now := time.Date(2026, 6, 19, 12, 0, 0, 0, time.UTC)

	return model.Carrier{
		ID:          1000,
		Name:        "Enterprise",
		Nation:      "United States",
		CarrierType: model.CarrierTypeAircraft,
		CommandCenter: model.CommandCenter{
			ID:            2000,
			CodeName:      "Bridge",
			SecurityLevel: 5,
			CarrierID:     1000,
		},
		Aircrafts: []model.Aircraft{
			{
				ID:           3000,
				Model:        "F/A-18",
				Manufacturer: "McDonnell Douglas",
				CarrierID:    1000,
			},
		},
		Erzeugt:      now,
		Aktualisiert: now,
	}
}

func validCreateCarrierJSON() string {
	return `{"name":"Enterprise","nation":"United States","carrierType":"AIRCRAFT_CARRIER","commandCenter":{"codeName":"Bridge","securityLevel":5},"aircrafts":[{"model":"F/A-18","manufacturer":"McDonnell Douglas"}]}`
}

func assertErrorResponse(t *testing.T, response *httptest.ResponseRecorder, status int, code string) {
	t.Helper()

	if response.Code != status {
		t.Fatalf("status code = %d, want %d", response.Code, status)
	}

	var body ErrorResponse
	decodeResponse(t, response, &body)
	if body.Error.Code != code {
		t.Fatalf("error code = %q, want %q", body.Error.Code, code)
	}
	if body.Error.Fields == nil {
		t.Fatal("error fields = nil, want empty slice or field errors")
	}
}

func decodeResponse(t *testing.T, response *httptest.ResponseRecorder, target any) {
	t.Helper()

	if err := json.Unmarshal(response.Body.Bytes(), target); err != nil {
		t.Fatalf("decode response body: %v", err)
	}
}

func ptr(carrier model.Carrier) *model.Carrier {
	return &carrier
}
