package service

import "github.com/JM01332/app/internal/carrier/model"

// CreateCarrierInput contains the validated data required to create a carrier.
type CreateCarrierInput struct {
	Name          string
	Nation        string
	CarrierType   model.CarrierType
	CommandCenter CreateCommandCenterInput
	Aircrafts     []CreateAircraftInput
}

// CreateCommandCenterInput contains data for the carrier's command center.
type CreateCommandCenterInput struct {
	CodeName      string
	SecurityLevel int
}

// CreateAircraftInput contains data for an aircraft assigned to the carrier.
type CreateAircraftInput struct {
	Model        string
	Manufacturer string
}
