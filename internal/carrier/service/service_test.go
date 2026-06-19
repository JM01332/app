package service

import (
	"context"
	"errors"
	"testing"

	"github.com/JM01332/app/internal/carrier/model"
)

type fakeCarrierStore struct {
	carriers       []model.Carrier
	listError      error
	listContext    context.Context
	createdCarrier *model.Carrier
	createError    error
	createContext  context.Context
}

func (store *fakeCarrierStore) List(ctx context.Context) ([]model.Carrier, error) {
	store.listContext = ctx
	return store.carriers, store.listError
}

func (store *fakeCarrierStore) Create(ctx context.Context, carrier *model.Carrier) error {
	store.createContext = ctx
	store.createdCarrier = carrier
	return store.createError
}

func TestCarrierServiceListReturnsRepositoryResult(t *testing.T) {
	store := &fakeCarrierStore{
		carriers: []model.Carrier{{ID: 1000, Name: "Enterprise"}},
	}
	service := NewCarrierService(store)
	ctx := context.Background()

	carriers, err := service.List(ctx)
	if err != nil {
		t.Fatalf("List() error = %v, want nil", err)
	}
	if store.listContext != ctx {
		t.Error("List() did not forward the context")
	}
	if len(carriers) != 1 || carriers[0].Name != "Enterprise" {
		t.Errorf("List() carriers = %v, want Enterprise", carriers)
	}
}

func TestCarrierServiceListReturnsRepositoryError(t *testing.T) {
	repositoryError := errors.New("database unavailable")
	store := &fakeCarrierStore{listError: repositoryError}
	service := NewCarrierService(store)

	carriers, err := service.List(context.Background())
	if carriers != nil {
		t.Errorf("List() carriers = %v, want nil", carriers)
	}
	if !errors.Is(err, repositoryError) {
		t.Errorf("List() error = %v, want %v", err, repositoryError)
	}
}

func TestCarrierServiceCreateMapsInputAndForwardsContext(t *testing.T) {
	store := &fakeCarrierStore{}
	service := NewCarrierService(store)
	ctx := context.Background()
	input := CreateCarrierInput{
		Name:        "Enterprise",
		Nation:      "United States",
		CarrierType: model.CarrierTypeAircraft,
		CommandCenter: CreateCommandCenterInput{
			CodeName:      "Bridge",
			SecurityLevel: 5,
		},
		Aircrafts: []CreateAircraftInput{
			{Model: "F/A-18", Manufacturer: "McDonnell Douglas"},
		},
	}

	carrier, err := service.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create() error = %v, want nil", err)
	}
	if store.createContext != ctx {
		t.Error("Create() did not forward the context")
	}
	if carrier != store.createdCarrier {
		t.Error("Create() did not return the carrier passed to the repository")
	}
	if carrier.Name != input.Name || carrier.Nation != input.Nation {
		t.Errorf("Create() carrier = %+v, want name and nation from input", carrier)
	}
	if carrier.CarrierType != input.CarrierType {
		t.Errorf("Create() carrier type = %q, want %q", carrier.CarrierType, input.CarrierType)
	}
	if carrier.CommandCenter.CodeName != input.CommandCenter.CodeName ||
		carrier.CommandCenter.SecurityLevel != input.CommandCenter.SecurityLevel {
		t.Errorf("Create() command center = %+v, want input values", carrier.CommandCenter)
	}
	if len(carrier.Aircrafts) != 1 ||
		carrier.Aircrafts[0].Model != input.Aircrafts[0].Model ||
		carrier.Aircrafts[0].Manufacturer != input.Aircrafts[0].Manufacturer {
		t.Errorf("Create() aircrafts = %+v, want input values", carrier.Aircrafts)
	}
}

func TestCarrierServiceCreateReturnsRepositoryError(t *testing.T) {
	store := &fakeCarrierStore{createError: ErrCarrierNameExists}
	service := NewCarrierService(store)

	carrier, err := service.Create(context.Background(), CreateCarrierInput{Name: "Enterprise"})
	if carrier != nil {
		t.Errorf("Create() carrier = %+v, want nil", carrier)
	}
	if !errors.Is(err, ErrCarrierNameExists) {
		t.Errorf("Create() error = %v, want ErrCarrierNameExists", err)
	}
}
