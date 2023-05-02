package config

import "reflect"

type ConfigInterface interface {
	LoadTokenConfig(conf string) (*tokenConfig, error)
}

var _ ConfigInterface = (*FileConfig)(nil)

type FileConfig struct {
	TokenConfig *tokenConfig
}

func (c *FileConfig) LoadTokenConfig(conf string) (*tokenConfig, error) {
	//TODO implement me
	return &tokenConfig{}, nil
}

func (c *FileConfig) parse(confName string) (err error) {
	c.TokenConfig, err = c.LoadTokenConfig(confName)
	if err != nil {
		return err
	}
	if reflect.DeepEqual(c.TokenConfig, &tokenConfig{}) {
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
