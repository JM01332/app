package model

// CarrierType represents a carrier_type value from PostgreSQL.
type CarrierType string

const (
	CarrierTypeAircraft   CarrierType = "AIRCRAFT_CARRIER"
	CarrierTypeHelicopter CarrierType = "HELICOPTER_CARRIER"
)
