package router

import (
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestCreateCarrierRequestValidationAcceptsValidRequest(t *testing.T) {
	validate := validator.New()
	request := validCreateCarrierRequest()

	if err := validate.Struct(request); err != nil {
		t.Fatalf("validate.Struct() error = %v, want nil", err)
	}
}

func TestCreateCarrierRequestValidationAcceptsEmptyAircraftList(t *testing.T) {
	validate := validator.New()
	request := validCreateCarrierRequest()
	request.Aircrafts = []CreateAircraftRequest{}

	if err := validate.Struct(request); err != nil {
		t.Fatalf("validate.Struct() error = %v, want nil", err)
	}
}

func TestCreateCarrierRequestValidationAcceptsBoundaryValues(t *testing.T) {
	testCases := []struct {
		name    string
		request CreateCarrierRequest
	}{
		{
			name: "minimum values",
			request: CreateCarrierRequest{
				Name:        "AB",
				Nation:      "DE",
				CarrierType: "HELICOPTER_CARRIER",
				CommandCenter: CreateCommandCenterRequest{
					CodeName:      "CC",
					SecurityLevel: 1,
				},
				Aircrafts: []CreateAircraftRequest{
					{Model: "A", Manufacturer: "AB"},
				},
			},
		},
		{
			name: "maximum values",
			request: CreateCarrierRequest{
				Name:        strings.Repeat("A", 50),
				Nation:      strings.Repeat("B", 50),
				CarrierType: "AIRCRAFT_CARRIER",
				CommandCenter: CreateCommandCenterRequest{
					CodeName:      strings.Repeat("C", 50),
					SecurityLevel: 5,
				},
				Aircrafts: []CreateAircraftRequest{
					{Model: strings.Repeat("D", 50), Manufacturer: strings.Repeat("E", 50)},
				},
			},
		},
	}

	validate := validator.New()
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if err := validate.Struct(testCase.request); err != nil {
				t.Fatalf("validate.Struct() error = %v, want nil", err)
			}
		})
	}
}

func TestCreateCarrierRequestValidationRejectsMissingCommandCenter(t *testing.T) {
	validate := validator.New()
	request := validCreateCarrierRequest()
	request.CommandCenter = CreateCommandCenterRequest{}

	err := validate.Struct(request)
	if err == nil {
		t.Fatal("validate.Struct() error = nil, want validation errors")
	}

	errors := validationErrors(t, err)
	assertValidationError(t, errors, "CodeName", "required")
	assertValidationError(t, errors, "SecurityLevel", "required")
}

func TestCreateCarrierRequestValidationRejectsValuesAboveMaximumLength(t *testing.T) {
	testCases := []struct {
		name   string
		field  string
		change func(*CreateCarrierRequest)
	}{
		{name: "name", field: "Name", change: func(request *CreateCarrierRequest) {
			request.Name = strings.Repeat("A", 51)
		}},
		{name: "nation", field: "Nation", change: func(request *CreateCarrierRequest) {
			request.Nation = strings.Repeat("A", 51)
		}},
		{name: "command center code name", field: "CodeName", change: func(request *CreateCarrierRequest) {
			request.CommandCenter.CodeName = strings.Repeat("A", 51)
		}},
		{name: "aircraft model", field: "Model", change: func(request *CreateCarrierRequest) {
			request.Aircrafts[0].Model = strings.Repeat("A", 51)
		}},
		{name: "aircraft manufacturer", field: "Manufacturer", change: func(request *CreateCarrierRequest) {
			request.Aircrafts[0].Manufacturer = strings.Repeat("A", 51)
		}},
	}

	validate := validator.New()
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request := validCreateCarrierRequest()
			testCase.change(&request)

			err := validate.Struct(request)
			if err == nil {
				t.Fatal("validate.Struct() error = nil, want validation errors")
			}

			assertValidationError(t, validationErrors(t, err), testCase.field, "max")
		})
	}
}

func TestCreateCarrierRequestValidationRejectsMissingRequiredFields(t *testing.T) {
	validate := validator.New()
	request := CreateCarrierRequest{}

	err := validate.Struct(request)
	if err == nil {
		t.Fatal("validate.Struct() error = nil, want validation errors")
	}

	errors := validationErrors(t, err)
	assertValidationError(t, errors, "Name", "required")
	assertValidationError(t, errors, "Nation", "required")
	assertValidationError(t, errors, "CarrierType", "required")
	assertValidationError(t, errors, "Aircrafts", "required")
}

func TestCreateCarrierRequestValidationRejectsInvalidCarrierType(t *testing.T) {
	validate := validator.New()
	request := validCreateCarrierRequest()
	request.CarrierType = "BATTLESHIP"

	err := validate.Struct(request)
	if err == nil {
		t.Fatal("validate.Struct() error = nil, want validation errors")
	}

	assertValidationError(t, validationErrors(t, err), "CarrierType", "oneof")
}

func TestCreateCarrierRequestValidationRejectsInvalidSecurityLevel(t *testing.T) {
	validate := validator.New()
	request := validCreateCarrierRequest()
	request.CommandCenter.SecurityLevel = 6

	err := validate.Struct(request)
	if err == nil {
		t.Fatal("validate.Struct() error = nil, want validation errors")
	}

	assertValidationError(t, validationErrors(t, err), "SecurityLevel", "max")
}

func TestCreateCarrierRequestValidationRejectsSecurityLevelBelowMinimum(t *testing.T) {
	validate := validator.New()
	request := validCreateCarrierRequest()
	request.CommandCenter.SecurityLevel = -1

	err := validate.Struct(request)
	if err == nil {
		t.Fatal("validate.Struct() error = nil, want validation errors")
	}

	assertValidationError(t, validationErrors(t, err), "SecurityLevel", "min")
}

func TestCreateCarrierRequestValidationRejectsInvalidAircraft(t *testing.T) {
	validate := validator.New()
	request := validCreateCarrierRequest()
	request.Aircrafts = []CreateAircraftRequest{
		{
			Model:        "",
			Manufacturer: "A",
		},
	}

	err := validate.Struct(request)
	if err == nil {
		t.Fatal("validate.Struct() error = nil, want validation errors")
	}

	errors := validationErrors(t, err)
	assertValidationError(t, errors, "Model", "required")
	assertValidationError(t, errors, "Manufacturer", "min")
}

func validCreateCarrierRequest() CreateCarrierRequest {
	return CreateCarrierRequest{
		Name:        "Enterprise",
		Nation:      "United States",
		CarrierType: "AIRCRAFT_CARRIER",
		CommandCenter: CreateCommandCenterRequest{
			CodeName:      "Bridge",
			SecurityLevel: 5,
		},
		Aircrafts: []CreateAircraftRequest{
			{
				Model:        "F/A-18",
				Manufacturer: "McDonnell Douglas",
			},
		},
	}
}

func validationErrors(t *testing.T, err error) validator.ValidationErrors {
	t.Helper()

	errors, ok := err.(validator.ValidationErrors)
	if !ok {
		t.Fatalf("error type = %T, want validator.ValidationErrors", err)
	}

	return errors
}

func assertValidationError(t *testing.T, errors validator.ValidationErrors, field string, tag string) {
	t.Helper()

	for _, fieldError := range errors {
		if fieldError.Field() == field && fieldError.Tag() == tag {
			return
		}
	}

	t.Fatalf("validation error for field %q with tag %q not found in %v", field, tag, errors)
}
