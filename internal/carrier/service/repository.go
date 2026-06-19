package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/JM01332/app/internal/carrier/model"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const carrierNameConstraint = "carrier_name_key"

// CarrierRepository provides persistence operations for carriers.
type CarrierRepository struct {
	db *gorm.DB
}

// NewCarrierRepository creates a repository backed by GORM.
func NewCarrierRepository(db *gorm.DB) *CarrierRepository {
	return &CarrierRepository{db: db}
}

// List returns all carriers with their relationships ordered by ID.
func (repository *CarrierRepository) List(ctx context.Context) ([]model.Carrier, error) {
	var carriers []model.Carrier

	result := repository.db.WithContext(ctx).
		Preload("CommandCenter").
		Preload("Aircrafts", func(database *gorm.DB) *gorm.DB {
			return database.Order("aircraft.id ASC")
		}).
		Order("carrier.id ASC").
		Find(&carriers)
	if result.Error != nil {
		return nil, fmt.Errorf("list carriers: %w", result.Error)
	}

	return carriers, nil
}

// Create stores a carrier and its relationships in one transaction.
func (repository *CarrierRepository) Create(ctx context.Context, carrier *model.Carrier) error {
	err := repository.db.WithContext(ctx).Transaction(func(transaction *gorm.DB) error {
		return transaction.Create(carrier).Error
	})
	if err == nil {
		return nil
	}

	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) && postgresError.ConstraintName == carrierNameConstraint {
		return fmt.Errorf("%w: %s", ErrCarrierNameExists, carrier.Name)
	}

	return fmt.Errorf("create carrier: %w", err)
}
