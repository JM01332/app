package router

type CreateCarrierRequest struct {
	Name          string                       `json:"name" validate:"required,min=2,max=50"`
	Nation        string                       `json:"nation" validate:"required,min=2,max=50"`
	CarrierType   string                       `json:"carrierType" validate:"required,oneof=AIRCRAFT_CARRIER HELICOPTER_CARRIER"`
	CommandCenter CreateCommandCenterRequest   `json:"commandCenter" validate:"required"`
	Aircrafts     []CreateAircraftRequest      `json:"aircrafts" validate:"required,dive"`
}

type CreateCommandCenterRequest struct {
	CodeName      string `json:"codeName" validate:"required,min=2,max=50"`
	SecurityLevel int    `json:"securityLevel" validate:"required,min=1,max=5"`
}

type CreateAircraftRequest struct {
	Model        string `json:"model" validate:"required,min=1,max=50"`
	Manufacturer string `json:"manufacturer" validate:"required,min=2,max=50"`
}
