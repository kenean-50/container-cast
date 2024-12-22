package manifest

import (
	"github.com/kenean-50/vm-container-manager/internal/util/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Manifest struct {
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

func NewManifest(path, file string) (*Manifest, error) {
	var m Manifest
	m.logger = log.
		With().
		Str("actor", "manifest").
		Logger()

	y := config.NewConfig()
	err := y.LoadConfigFile(path, "yaml", file)

	if err != nil {
		m.logger.
			Fatal().
			Str("status", "failed to load config file").
			Str("reason", ""+err.Error()).
			Send()
	}

	if err := y.Unmarshal(&m.Values); err != nil {
		m.logger.
			Fatal().
			Str("status", "failed to nmarshal").
			Str("reason", ""+err.Error()).
			Send()
		return nil, err
	}

	return &m, nil
}

// type PortMapping struct {
// 	HostPort      string
// 	ContainerPort string
// }

// func (p *PortMapping) UnmarshalText(text []byte) error {
// 	parts := strings.Split(string(text), ":")
// 	if len(parts) == 2 {
// 		p.HostPort = parts[0]
// 		p.ContainerPort = parts[1]
// 	}
// 	return nil
// }
