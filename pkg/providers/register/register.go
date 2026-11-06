package register

import (
	"github.com/minc-org/minc/pkg/providers"
	"github.com/minc-org/minc/pkg/providers/moby"
	"github.com/minc-org/minc/pkg/providers/podman"
	"github.com/minc-org/minc/pkg/rootlessmarker"
	"github.com/spf13/viper"
)

// Register builds the container provider for the given engine, scoped to the
// named cluster.
func Register(provider, name string) (providers.Provider, error) {
	allowRootless := viper.GetBool("allow-rootless") || rootlessmarker.Present(name)
	switch provider {
	case "podman":
		return podman.New(allowRootless, name)
	case "docker":
		return moby.New(name)
	default:
		return podman.New(allowRootless, name)
	}
}
