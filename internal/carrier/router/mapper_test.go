package router

import (
	"testing"
	"time"

	"github.com/JM01332/app/internal/carrier/model"
)

func TestMapCreateCarrierInputMapsNestedRequest(t *testing.T) {
	request := CreateCarrierRequest{
		Name:        "Enterprise",
		Nation:      "United States",
		CarrierType: "AIRCRAFT_CARRIER",
		CommandCenter: CreateCommandCenterRequest{
			CodeName:      "Bridge",
			SecurityLevel: 5,
		},
		Aircrafts: []CreateAircraftRequest{
			{Model: "F/A-18", Manufacturer: "McDonnell Douglas"},
		},
	}

	input := mapCreateCarrierInput(request)

	if input.Name != request.Name || input.Nation != request.Nation {
		t.Errorf("carrier input = %+v, want name and nation from request", input)
	}
	if input.CarrierType != model.CarrierTypeAircraft {
		t.Errorf("carrier type = %q, want %q", input.CarrierType, model.CarrierTypeAircraft)
	}
	if input.CommandCenter.CodeName != request.CommandCenter.CodeName ||
		input.CommandCenter.SecurityLevel != request.CommandCenter.SecurityLevel {
		t.Errorf("command center input = %+v, want request values", input.CommandCenter)
	}
	if len(input.Aircrafts) != 1 ||
		input.Aircrafts[0].Model != request.Aircrafts[0].Model ||
		input.Aircrafts[0].Manufacturer != request.Aircrafts[0].Manufacturer {
		t.Errorf("aircraft input = %+v, want request values", input.Aircrafts)
	}
}

func TestMapCreateCarrierInputKeepsEmptyAircraftList(t *testing.T) {
	request := CreateCarrierRequest{
		Aircrafts: []CreateAircraftRequest{},
	}

	input := mapCreateCarrierInput(request)

	if input.Aircrafts == nil {
		t.Fatal("aircraft input = nil, want empty slice")
	}
	if len(input.Aircrafts) != 0 {
		t.Errorf("aircraft input length = %d, want 0", len(input.Aircrafts))
	}
}

func TestMapCarrierResponseMapsRelationshipsAndTimestamps(t *testing.T) {
	createdAt := time.Date(2026, 6, 19, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	carrier := model.Carrier{
		ID:           1000,
		Name:         "Enterprise",
		Nation:       "United States",
		CarrierType:  model.CarrierTypeAircraft,
		Erzeugt:      createdAt,
		Aktualisiert: updatedAt,
		CommandCenter: model.CommandCenter{
			ID:            2000,
			CodeName:      "Bridge",
			SecurityLevel: 5,
		},
		Aircrafts: []model.Aircraft{
			{ID: 3000, Model: "F/A-18", Manufacturer: "McDonnell Douglas"},
		},
	}

	response := mapCarrierResponse(carrier)

	if response.ID != carrier.ID || response.Name != carrier.Name || response.Nation != carrier.Nation {
		t.Errorf("carrier response = %+v, want carrier values", response)
	}
	if response.CarrierType != string(carrier.CarrierType) {
		t.Errorf("carrier type = %q, want %q", response.CarrierType, carrier.CarrierType)
	}
	if response.CommandCenter.ID != carrier.CommandCenter.ID ||
		response.CommandCenter.CodeName != carrier.CommandCenter.CodeName ||
		response.CommandCenter.SecurityLevel != carrier.CommandCenter.SecurityLevel {
		t.Errorf("command center response = %+v, want model values", response.CommandCenter)
	}
	if len(response.Aircrafts) != 1 ||
		response.Aircrafts[0].ID != carrier.Aircrafts[0].ID ||
		response.Aircrafts[0].Model != carrier.Aircrafts[0].Model ||
		response.Aircrafts[0].Manufacturer != carrier.Aircrafts[0].Manufacturer {
		t.Errorf("aircraft response = %+v, want model values", response.Aircrafts)
	}
	if !response.CreatedAt.Equal(createdAt) || !response.UpdatedAt.Equal(updatedAt) {
		t.Errorf("response timestamps = %v/%v, want %v/%v", response.CreatedAt, response.UpdatedAt, createdAt, updatedAt)
	}
}

func TestMapCarrierResponsesReturnsEmptySlice(t *testing.T) {
	responses := mapCarrierResponses([]model.Carrier{})

	if responses == nil {
		t.Fatal("carrier responses = nil, want empty slice")
	}
	if len(responses) != 0 {
		t.Errorf("carrier response length = %d, want 0", len(responses))
	}
}
