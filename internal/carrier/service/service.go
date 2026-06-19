package service

import (
	"context"

	"github.com/JM01332/app/internal/carrier/model"
)

// CarrierStore defines the persistence operations required by CarrierService.
type CarrierStore interface {
	List(ctx context.Context) ([]model.Carrier, error)
	Create(ctx context.Context, carrier *model.Carrier) error
}

// CarrierService provides the application logic for carriers.
type CarrierService struct {
	repository CarrierStore
}

// NewCarrierService creates a service with its persistence dependency.
func NewCarrierService(repository CarrierStore) *CarrierService {
	return &CarrierService{repository: repository}
}

// List returns all carriers from the repository.
func (service *CarrierService) List(ctx context.Context) ([]model.Carrier, error) {
	return service.repository.List(ctx)
}

// Create maps validated input and stores the resulting carrier.
func (service *CarrierService) Create(ctx context.Context, input CreateCarrierInput) (*model.Carrier, error) {
	carrier := carrierFromCreateInput(input)
	if err := service.repository.Create(ctx, carrier); err != nil {
		return nil, err
	}

	return carrier, nil
}

func carrierFromCreateInput(input CreateCarrierInput) *model.Carrier {
	aircrafts := make([]model.Aircraft, len(input.Aircrafts))
	for index, aircraft := range input.Aircrafts {
		aircrafts[index] = model.Aircraft{
			Model:        aircraft.Model,
			Manufacturer: aircraft.Manufacturer,
		}
	}

	return &model.Carrier{
		Name:        input.Name,
		Nation:      input.Nation,
		CarrierType: input.CarrierType,
		CommandCenter: model.CommandCenter{
			CodeName:      input.CommandCenter.CodeName,
			SecurityLevel: input.CommandCenter.SecurityLevel,
		},
		Aircrafts: aircrafts,
	}
}
