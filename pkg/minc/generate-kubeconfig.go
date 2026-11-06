package minc

import (
	"github.com/minc-org/minc/pkg/kubeconfig"
	"github.com/minc-org/minc/pkg/log"
	"github.com/minc-org/minc/pkg/providers/register"
)

func GenerateKubeConfig(provider, name string) error {
	p, err := register.Register(provider, name)
	if err != nil {
		return err
	}
	log.Debug("Provider Info", "Provider", p)
	if _, err := p.List(); err != nil {
		return err
	}
	config, err := resolveKubeConfig(p, name)
	if err != nil {
		return err
	}
	if err := kubeconfig.UpdateKubeConfig(config); err != nil {
		return err
	}
	return nil
}
