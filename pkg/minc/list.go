package minc

import (
	"github.com/minc-org/minc/pkg/log"
	"github.com/minc-org/minc/pkg/providers/register"
)

func List(provider, name string) ([]byte, error) {
	p, err := register.Register(provider, name)
	if err != nil {
		return nil, err
	}
	log.Debug("Provider Info", "Provider", p)
	return p.List()
}
