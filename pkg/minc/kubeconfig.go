package minc

import (
	"fmt"

	"github.com/minc-org/minc/pkg/constants"
	"github.com/minc-org/minc/pkg/kubeconfig"
	"github.com/minc-org/minc/pkg/providers"
)

// resolveKubeConfig reads the cluster's kubeconfig, finds the ephemeral host
// port its apiserver was published to, and rewrites it for host access.
func resolveKubeConfig(p providers.Provider, name string) ([]byte, error) {
	raw, err := p.GetKubeConfig()
	if err != nil {
		return nil, err
	}
	port, err := p.HostPort(providers.APIServerContainerPort)
	if err != nil {
		return nil, err
	}
	server := fmt.Sprintf("https://%s:%d", constants.HostName, port)
	return kubeconfig.Rewrite(raw, name, server)
}
