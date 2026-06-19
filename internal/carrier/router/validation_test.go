package router

import (
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
