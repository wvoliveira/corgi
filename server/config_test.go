package server

import "testing"

func TestNewConfig(t *testing.T) {
	logger := NewLogger()
	config := NewConfig(*logger, "../")

	if config.ServerAddress != "0.0.0.0:8000" {
		t.Errorf("ServerAddress was incorrect, got: %s, want: %s.", config.ServerAddress, "0.0.0.0:8000")
	}

	if config.SecretKey != "changeme" {
		t.Errorf("SecretKey was incorrect, got: %s, want: %s.", config.SecretKey, "changeme")
	}
}
