package database

import (
	"context"
	"testing"
)

func TestOpenPostgresRequiresDatabaseURL(t *testing.T) {
	_, err := OpenPostgres(context.Background(), "")
	if err == nil {
		t.Fatal("OpenPostgres() error = nil, want an error")
	}
}
