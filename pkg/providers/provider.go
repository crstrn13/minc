package providers

import (
	"github.com/minc-org/minc/pkg/minc/types"
)

// APIServerContainerPort is the in-container kube-apiserver port, published to
// an ephemeral host port so multiple clusters can coexist.
const APIServerContainerPort = 6443

type Provider interface {
	Name() string
	Info() (*ProviderInfo, error)
	ImageExists(string) bool
	PullImage(image string) error
	Create(cType *types.CreateType) error
	WaitForMicroShiftService() error
	GetKubeConfig() ([]byte, error)
	Delete() error
	List() ([]byte, error)
	// HostPort returns the host port the given container port was published to.
	// Engine-specific: podman vs docker report published ports differently.
	HostPort(containerPort int) (int, error)
}

type ProviderInfo struct {
	Rootless bool
	CGroupV2 bool
}
