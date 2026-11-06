package kubeconfig

import (
	"testing"

	"k8s.io/client-go/tools/clientcmd"
)

// rawKubeConfig mimics the kubeconfig MicroShift writes inside the container:
// everything is named "microshift" and the server points at the in-container port.
const rawKubeConfig = `
apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://127.0.0.1.nip.io:6443
    insecure-skip-tls-verify: true
  name: microshift
users:
- name: microshift
  user:
    token: abc
contexts:
- context:
    cluster: microshift
    user: microshift
  name: microshift
current-context: microshift
`

func TestRewriteRenamesAndRepoints(t *testing.T) {
	out, err := Rewrite([]byte(rawKubeConfig), "dev2", "https://127.0.0.1.nip.io:49200")
	if err != nil {
		t.Fatalf("Rewrite: %v", err)
	}
	cfg, err := clientcmd.Load(out)
	if err != nil {
		t.Fatalf("load rewritten: %v", err)
	}

	if _, ok := cfg.Clusters["dev2"]; !ok {
		t.Fatalf("cluster not renamed to dev2: %+v", cfg.Clusters)
	}
	if _, ok := cfg.Clusters["microshift"]; ok {
		t.Fatalf("old cluster name should be gone")
	}
	if got := cfg.Clusters["dev2"].Server; got != "https://127.0.0.1.nip.io:49200" {
		t.Fatalf("server not repointed: %q", got)
	}
	if _, ok := cfg.Contexts["dev2"]; !ok {
		t.Fatalf("context not renamed to dev2")
	}
	if cfg.CurrentContext != "dev2" {
		t.Fatalf("current-context=%q want dev2", cfg.CurrentContext)
	}
	if cfg.Contexts["dev2"].Cluster != "dev2" || cfg.Contexts["dev2"].AuthInfo != "dev2" {
		t.Fatalf("context refs not renamed: %+v", cfg.Contexts["dev2"])
	}
}

// Two rewritten configs with different names must not collide when merged.
func TestRewriteTwoClustersMergeCleanly(t *testing.T) {
	a, _ := Rewrite([]byte(rawKubeConfig), "a", "https://127.0.0.1.nip.io:1111")
	b, _ := Rewrite([]byte(rawKubeConfig), "b", "https://127.0.0.1.nip.io:2222")

	ca, _ := clientcmd.Load(a)
	cb, _ := clientcmd.Load(b)
	merged := mergeConfigs(ca, cb)

	if len(merged.Clusters) != 2 {
		t.Fatalf("expected 2 clusters after merge, got %d", len(merged.Clusters))
	}
	if merged.Clusters["a"].Server != "https://127.0.0.1.nip.io:1111" {
		t.Fatalf("cluster a server wrong: %q", merged.Clusters["a"].Server)
	}
	if merged.Clusters["b"].Server != "https://127.0.0.1.nip.io:2222" {
		t.Fatalf("cluster b server wrong: %q", merged.Clusters["b"].Server)
	}
}
