package router

import "github.com/JM01332/app/internal/carrier/model"

func mapCarrierResponses(carriers []model.Carrier) []CarrierResponse {
	responses := make([]CarrierResponse, 0, len(carriers))
	for _, carrier := range carriers {
		responses = append(responses, mapCarrierResponse(carrier))
	}

	return responses
}

func mapCarrierResponse(carrier model.Carrier) CarrierResponse {
	return CarrierResponse{
		ID:          carrier.ID,
		Name:        carrier.Name,
		Nation:      carrier.Nation,
		CarrierType: string(carrier.CarrierType),
		CommandCenter: CommandCenterResponse{
			ID:            carrier.CommandCenter.ID,
			CodeName:      carrier.CommandCenter.CodeName,
			SecurityLevel: carrier.CommandCenter.SecurityLevel,
		},
		Aircrafts: mapAircraftResponses(carrier.Aircrafts),
		CreatedAt: carrier.Erzeugt,
		UpdatedAt: carrier.Aktualisiert,
	}
}

func mapAircraftResponses(aircrafts []model.Aircraft) []AircraftResponse {
	responses := make([]AircraftResponse, 0, len(aircrafts))
	for _, aircraft := range aircrafts {
		responses = append(responses, AircraftResponse{
			ID:           aircraft.ID,
			Model:        aircraft.Model,
			Manufacturer: aircraft.Manufacturer,
		})
	}

	return responses
}
