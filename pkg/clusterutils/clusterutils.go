package clusterutils

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/e2e-framework/klient"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/kind/pkg/cluster"
)

// set up onboarding cluster client
func OnboardingConfig() (*envconf.Config, error) {
	kubeConfig, err := retrieveKindClusterNameByPrefix("onboarding")
	if err != nil {
		return nil, err
	}
	restConfig, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to create rest config based on kubeconfig: %v", err)
	}
	onboardingClient, err := klient.New(restConfig)
	if err != nil {
		return nil, err
	}
	onboardingCfg := envconf.New().WithClient(onboardingClient)
	return onboardingCfg.WithNamespace(corev1.NamespaceDefault), nil
}

func retrieveKindClusterNameByPrefix(prefix string) (string, error) {
	kind := cluster.NewProvider()
	clusters, err := kind.List()
	if err != nil {
		return "", err
	}
	for _, clusterName := range clusters {
		if strings.HasPrefix(clusterName, prefix) {
			return kind.KubeConfig(clusterName, false)
		}
	}
	return "", fmt.Errorf("no cluster found with prefix %s", prefix)
}
