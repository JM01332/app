package router

import "time"

type CarrierResponse struct {
	ID            int64                 `json:"id"`
	Name          string                `json:"name"`
	Nation        string                `json:"nation"`
	CarrierType   string                `json:"carrierType"`
	CommandCenter CommandCenterResponse `json:"commandCenter"`
	Aircrafts     []AircraftResponse    `json:"aircrafts"`
	CreatedAt     time.Time             `json:"createdAt"`
	UpdatedAt     time.Time             `json:"updatedAt"`
}

type CommandCenterResponse struct {
	ID            int64  `json:"id"`
	CodeName      string `json:"codeName"`
	SecurityLevel int    `json:"securityLevel"`
}

type AircraftResponse struct {
	ID           int64  `json:"id"`
	Model        string `json:"model"`
	Manufacturer string `json:"manufacturer"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

type APIError struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Fields  []FieldError `json:"fields"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
