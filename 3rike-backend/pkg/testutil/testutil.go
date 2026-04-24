// Package testutil provides an in-memory SQLite database for unit tests.
// Usage:
//
//	db := testutil.NewDB(t)
//	repo := repository.NewDriverRepo(db)
package testutil

import (
	"testing"

	"github.com/3rike12/3rike-backend/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB returns a fresh in-memory SQLite DB with all tables migrated.
// The DB is scoped to the test and cleaned up automatically.
func NewDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("testutil: open sqlite: %v", err)
	}
	if err := repository.Migrate(db); err != nil {
		t.Fatalf("testutil: migrate: %v", err)
	}
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	})
	return db
}
