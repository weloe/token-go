package config

import (
	"github.com/spf13/viper"
	"log"
	"reflect"
)

type ConfigInterface interface {
	loadTokenConfig(conf string) error
}

var _ ConfigInterface = (*FileConfig)(nil)

type FileConfig struct {
	TokenConfig *TokenConfig
}

// Use viper to load config with file.
func (c *FileConfig) loadTokenConfig(conf string) error {
	viper.SetConfigFile(conf)
	var err error
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&c)
	if err != nil {
		return err
	}

	return nil
}

func (c *FileConfig) parse(confName string) (err error) {
	err = c.loadTokenConfig(confName)
	if err != nil {
		return err
	}
	if c.TokenConfig == nil || reflect.DeepEqual(c.TokenConfig, &TokenConfig{}) {
		c.TokenConfig = DefaultTokenConfig()
		log.Print("Token-go read empty or error file config, use default config")
	}

	return err
}

// ReadConfig create from file.
func ReadConfig(confName string) (*FileConfig, error) {
	c := &FileConfig{}
	err := c.parse(confName)
	return c, err
}
