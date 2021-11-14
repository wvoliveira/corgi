package server

import "testing"

func TestNewConfig(t *testing.T) {
	logger := NewLogger()
	config := NewConfig(logger, "../")

	/*
		Test default values.
	*/

	if config.DBType != "persistent" {
		t.Errorf("DBType was incorrect, got: %s, want: %s.", config.DBType, "persistent")
	}

	if config.DBDriver != "sqlite" {
		t.Errorf("DBDriver was incorrect, got: %s, want: %s.", config.DBDriver, "sqlite")
	}

	if config.DBSource != "redir.db" {
		t.Errorf("DBSource was incorrect, got: %s, want: %s.", config.DBSource, "redir.db")
	}

	if config.ServerAddress != "0.0.0.0:8080" {
		t.Errorf("ServerAddress was incorrect, got: %s, want: %s.", config.ServerAddress, "0.0.0.0:8080")
	}

	if config.SecretKey != "changeme" {
		t.Errorf("SecretKey was incorrect, got: %s, want: %s.", config.SecretKey, "changeme")
	}
}
