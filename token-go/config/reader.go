package config

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
)

type ConfigInterface interface {
	loadTokenConfig(conf string) (*TokenConfig, error)
}

var _ ConfigInterface = (*FileConfig)(nil)

type FileConfig struct {
	TokenConfig *TokenConfig
}

func (c *FileConfig) loadTokenConfig(conf string) (*TokenConfig, error) {
	var config *FileConfig
	viper.SetConfigFile(conf)
	var err error
	err = viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error config file: %s \n", err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("error viper unmarshal config: %s \n", err)
	}

	return config.TokenConfig, nil
}

func (c *FileConfig) parse(confName string) (err error) {
	c.TokenConfig, err = c.loadTokenConfig(confName)
	if err != nil {
		return err
	}
	if reflect.DeepEqual(c.TokenConfig, &TokenConfig{}) {
		c.TokenConfig = DefaultTokenConfig()
	}

	return err
}

// NewConfig create from file.
func NewConfig(confName string) (ConfigInterface, error) {
	c := &FileConfig{}
	err := c.parse(confName)
	return c, err
}
