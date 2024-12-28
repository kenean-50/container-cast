package config

import (
	"github.com/kenean-50/vm-container-manager/internal/util/config"
	"github.com/rs/zerolog/log"
)

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) LoadFromYAML(filePath, fileName string) (*Config, error) {
	c.logger = log.
		With().
		Str("actor", "config").
		Logger()
	y := config.NewConfig()
	err := y.LoadConfigFile(filePath, "yaml", fileName)

	if err != nil {
		c.logger.
			Fatal().
			Str("status", "failed to load config file").
			Str("reason", ""+err.Error()).
			Send()
	}

	if err := y.Unmarshal(&c.Values); err != nil {
		c.logger.
			Fatal().
			Str("status", "failed to nmarshal").
			Str("reason", ""+err.Error()).
			Send()
		return nil, err
	}
	return c, nil
}
