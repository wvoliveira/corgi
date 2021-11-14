package server

import (
	"testing"
)

// TestInitDatabase test for startup database.
func TestInitDatabase(t *testing.T) {
	logger := NewLogger()
	config := NewConfig(logger, "../")

	// Set database in memory.
	config.DBType = "memory"

	_ = NewDatabase(logger, config)
}

// TestLoadSqlite test loadSqlite to run sqlite database.
func TestLoadSqlite(t *testing.T) {
	logger := NewLogger()
	config := NewConfig(logger, "../")

	// Set database in memory.
	config.DBType = "memory"

	_, err := loadSqlite(config)

	if err != nil {
		t.Fatalf("error to execute loadSqlite(config): %s", err)
	}
}
