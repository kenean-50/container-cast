package deploy

import (
	"github.com/kenean-50/vm-container-manager/internal/actor/config"
)

type DeployConfig struct {
	config *config.Config
}

func NewDeployModule(config *config.Config) *DeployConfig {
	return &DeployConfig{
		config: config,
	}
}
