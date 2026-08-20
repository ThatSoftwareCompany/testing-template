//go:build integration

package errstore

import (
	"context"
	"os"
	"testing"
	"time"

	platformmigrate "github.com/ThatSoftwareCompany/template-go-api/internal/platform/migrate"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresStorePersistsAndListsSafeEvents(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not configured")
	}

	if err := platformmigrate.Up(platformmigrate.Config{
		DatabaseURL:   databaseURL,
		MigrationsDir: "file://../../../migrations",
	}); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		t.Fatalf("open pool: %v", err)
	}
	defer pool.Close()
	store := NewPostgresStore(pool)

	event := ErrorEvent{
		OccurredAt:    time.Now().UTC(),
		CorrelationID: "integration-correlation",
		Method:        "GET",
		Path:          "/api/v1/health",
		Endpoint:      "/api/v1/health",
		StatusCode:    503,
		ErrorCode:     "service_unavailable",
		Message:       "service unavailable",
	}
	if err := store.Persist(context.Background(), event); err != nil {
		t.Fatalf("persist event: %v", err)
	}

	items, err := store.List(context.Background(), Filter{Endpoint: event.Endpoint, Limit: 10})
	if err != nil {
		t.Fatalf("list events: %v", err)
	}
	if len(items) == 0 || items[0].CorrelationID != event.CorrelationID {
		t.Fatalf("unexpected events: %#v", items)
	}
}
