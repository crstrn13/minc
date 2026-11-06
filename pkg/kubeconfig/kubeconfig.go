package kubeconfig

import (
	"fmt"
	"os"

	"github.com/minc-org/minc/pkg/log"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// Rewrite renames the cluster/user/context entries in the raw (in-container)
// kubeconfig to name and repoints the API server at server, so named clusters
// merge into the user's kubeconfig without colliding. Only the port changes, so
// the cert SAN (the hostname) still validates.
func Rewrite(raw []byte, name, server string) ([]byte, error) {
	cfg, err := clientcmd.Load(raw)
	if err != nil {
		return nil, fmt.Errorf("loading kubeconfig: %w", err)
	}

	renamed := api.NewConfig()
	for _, cluster := range cfg.Clusters {
		c := cluster.DeepCopy()
		c.Server = server
		renamed.Clusters[name] = c
		break
	}
	for _, auth := range cfg.AuthInfos {
		renamed.AuthInfos[name] = auth.DeepCopy()
		break
	}
	renamed.Contexts[name] = &api.Context{
		Cluster:   name,
		AuthInfo:  name,
		Namespace: "default",
	}
	renamed.CurrentContext = name

	return clientcmd.Write(*renamed)
}

func UpdateKubeConfig(config []byte) error {
	kubeConfigPath := getKubeConfigPath()
	log.Debug(fmt.Sprintf("Updating kubeconfig at %s", kubeConfigPath))
	defaultConfig, err := clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		defaultConfig = api.NewConfig()
	}
	uShiftConfig, err := clientcmd.Load(config)
	if err != nil {
		return err
	}
	mergedConfig := mergeConfigs(defaultConfig, uShiftConfig)

	err = clientcmd.WriteToFile(*mergedConfig, kubeConfigPath)
	if err != nil {
		return err
	}
	return nil
}

func getKubeConfigPath() string {
	kubeConfigPath := clientcmd.RecommendedHomeFile
	if os.Getenv("KUBECONFIG") != "" {
		kubeConfigPath = os.Getenv("KUBECONFIG")
	}
	return kubeConfigPath
}

// mergeConfigs merges two kubeconfig files and returns a single merged configuration.
func mergeConfigs(config1, config2 *api.Config) *api.Config {
	mergedConfig := config1.DeepCopy() // Copy config1 as base

	// Merge clusters
	for name, cluster := range config2.Clusters {
		mergedConfig.Clusters[name] = cluster
	}

	// Merge contexts
	for name, context := range config2.Contexts {
		mergedConfig.Contexts[name] = context

	}

	// Merge AuthInfos (users)
	for name, auth := range config2.AuthInfos {
		mergedConfig.AuthInfos[name] = auth
	}

	// update the current context
	mergedConfig.CurrentContext = config2.CurrentContext

	return mergedConfig
}

// RemoveClusterFromConfig removes the named cluster's cluster/context/user
// entries from the user's kubeconfig.
func RemoveClusterFromConfig(name string) error {
	kubeConfigPath := getKubeConfigPath()
	log.Debug(fmt.Sprintf("Updating kubeconfig at %s", kubeConfigPath))
	config, err := clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		return err
	}
	// Check if the cluster exists
	if _, exists := config.Clusters[name]; !exists {
		log.Debug(fmt.Sprintf("cluster %s not found in kubeconfig", name))
		return nil
	}

	// Remove cluster, context, and user
	delete(config.Clusters, name)

	// Find and remove the associated context
	for ctxName, ctx := range config.Contexts {
		if ctx.Cluster == name {
			delete(config.Contexts, ctxName)
			// Remove associated user
			delete(config.AuthInfos, ctx.AuthInfo)
			// If it's the current context, reset it
			if config.CurrentContext == ctxName {
				config.CurrentContext = ""
			}
		}
	}

	// Save updated kubeconfig
	err = clientcmd.WriteToFile(*config, kubeConfigPath)
	if err != nil {
		return fmt.Errorf("failed to save kubeconfig: %v", err)
	}

	log.Debug(fmt.Sprintf("Cluster %s removed successfully from kubeconfig\n", name))
	return nil
}
