package config

import "testing"

func TestNewConfig(t *testing.T) {
	config, err := ReadConfig("../examples/token_conf.yaml")
	if err != nil {
		t.Errorf("read error: %v", err)
	}
	tokenConfig := config.TokenConfig

	t.Log(tokenConfig)
}

func TestNewIniConfig(t *testing.T) {
	config, err := ReadConfig("../examples/token_conf.ini")
	if err != nil {
		t.Errorf("read error: %v", err)
	}
	tokenConfig := config.TokenConfig

	t.Log(tokenConfig)
}

func TestDefaultCookieConfig(t *testing.T) {
	config := DefaultTokenConfig()
	t.Log(config)
}

func TestDefaultTokenConfig(t *testing.T) {
	config := DefaultCookieConfig()
	t.Log(config)
}
