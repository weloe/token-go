package config

import "testing"

func TestNewConfig(t *testing.T) {
	config, err := NewConfig("")
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	t.Log(config)
}

func TestDefaultCookieConfig(t *testing.T) {
	config := DefaultTokenConfig()
	t.Log(config)
}

func TestDefaultTokenConfig(t *testing.T) {
	config := DefaultCookieConfig()
	t.Log(config)
}
