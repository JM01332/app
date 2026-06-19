package database

import (
	"context"
	"testing"
	"time"
)

func TestOpenPostgresRequiresDatabaseURL(t *testing.T) {
	_, err := OpenPostgres(context.Background(), "")
	if err == nil {
		t.Fatal("OpenPostgres() error = nil, want an error")
	}
}

func TestGormConfigUsesUTCNowFunc(t *testing.T) {
	now := newGormConfig().NowFunc()

	if now.Location() != time.UTC {
		t.Fatalf("NowFunc location = %v, want UTC", now.Location())
	}
}
