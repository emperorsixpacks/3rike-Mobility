// Package testutil provides a Postgres test database bootstrap.
// Requires DATABASE_URL env var pointing to a test Postgres instance.
// Usage:
//
//	db := testutil.NewDB(t)
//	repo := repository.NewDriverRepo(db)
package testutil

import (
	"os"
	"testing"

	"github.com/3rike12/3rike-backend/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB returns a Postgres DB with all tables migrated.
// Skips the test if DATABASE_URL is not set.
func NewDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("skipping: DATABASE_URL not set")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("testutil: open postgres: %v", err)
	}
	if err := repository.Migrate(db); err != nil {
		t.Fatalf("testutil: migrate: %v", err)
	}
	return db
}
