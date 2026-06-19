package router

import (
	"github.com/JM01332/app/internal/carrier/model"
	carrierservice "github.com/JM01332/app/internal/carrier/service"
)

func mapCreateCarrierInput(request CreateCarrierRequest) carrierservice.CreateCarrierInput {
	aircrafts := make([]carrierservice.CreateAircraftInput, len(request.Aircrafts))
	for index, aircraft := range request.Aircrafts {
		aircrafts[index] = carrierservice.CreateAircraftInput{
			Model:        aircraft.Model,
			Manufacturer: aircraft.Manufacturer,
		}
	}

	return carrierservice.CreateCarrierInput{
		Name:        request.Name,
		Nation:      request.Nation,
		CarrierType: model.CarrierType(request.CarrierType),
		CommandCenter: carrierservice.CreateCommandCenterInput{
			CodeName:      request.CommandCenter.CodeName,
			SecurityLevel: request.CommandCenter.SecurityLevel,
		},
		Aircrafts: aircrafts,
	}
}

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
