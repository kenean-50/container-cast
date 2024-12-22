package deploy

import (
	"github.com/kenean-50/vm-container-manager/internal/actor/manifest"
)

type DeployConfig struct {
	manifest *manifest.Manifest
}

func NewDeployModule(manifest *manifest.Manifest) *DeployConfig {
	return &DeployConfig{
		manifest: manifest,
	}
}
