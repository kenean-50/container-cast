package config

import "github.com/rs/zerolog"

type ConfigLoader interface {
	LoadFromYAML(filePath string) (*Config, error)
	// LoadFromGRPC(address string) (*Config, error)
}

type Config struct {
	Values Values
	logger zerolog.Logger
}

type Values struct {
	Services map[string]Service `mapstructure:"services"`
	Servers  map[string]Server  `mapstructure:"servers"`
}

type Service struct {
	Image string   `mapstructure:"image"`
	Ports []string `mapstructure:"ports"`
}

type Server struct {
	Host           string `mapstructure:"host"`
	User           string `mapstructure:"user"`
	SSHPort        int    `mapstructure:"ssh_port"`
	PrivateKeyPath string `mapstructure:"private_key_path"`
}
