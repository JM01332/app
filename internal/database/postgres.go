package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Postgres owns the GORM database handle and its underlying connection pool.
type Postgres struct {
	DB   *gorm.DB
	pool *sql.DB
}

// OpenPostgres opens and verifies a PostgreSQL connection.
func OpenPostgres(ctx context.Context, databaseURL string) (*Postgres, error) {
	if strings.TrimSpace(databaseURL) == "" {
		return nil, errors.New("database URL is required")
	}

	db, err := gorm.Open(postgres.Open(databaseURL), newGormConfig())
	if err != nil {
		return nil, fmt.Errorf("open PostgreSQL connection: %w", err)
	}

	pool, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("access PostgreSQL connection pool: %w", err)
	}

	if err := pool.PingContext(ctx); err != nil {
		_ = pool.Close()
		return nil, fmt.Errorf("ping PostgreSQL: %w", err)
	}

	return &Postgres{DB: db, pool: pool}, nil
}

func newGormConfig() *gorm.Config {
	return &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}
}

// Close releases all connections in the PostgreSQL pool.
func (postgres *Postgres) Close() error {
	if postgres == nil || postgres.pool == nil {
		return nil
	}

	return postgres.pool.Close()
}
