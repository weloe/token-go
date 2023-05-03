package config

import "reflect"

type ConfigInterface interface {
	loadTokenConfig(conf string) (*TokenConfig, error)
}

var _ ConfigInterface = (*FileConfig)(nil)

type FileConfig struct {
	TokenConfig *TokenConfig
}

func (c *FileConfig) loadTokenConfig(conf string) (*TokenConfig, error) {
	//TODO implement me
	return &TokenConfig{}, nil
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
