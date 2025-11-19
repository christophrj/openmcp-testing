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

// prefix = kind cluster name prefix, namespace = config namespace
func ConfigByPrefix(prefix string, namespace string) (*envconf.Config, error) {
	kind := cluster.NewProvider()
	clusterName, err := retrieveKindClusterNameByPrefix(prefix)
	if err != nil {
		return nil, err
	}
	kubeConfig, err := kind.KubeConfig(clusterName, false)
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
	return onboardingCfg.WithNamespace(namespace), nil
}

// set up onboarding cluster client
func OnboardingConfig() (*envconf.Config, error) {
	return ConfigByPrefix("onboarding", corev1.NamespaceDefault)
}

func retrieveKindClusterNameByPrefix(prefix string) (string, error) {
	kind := cluster.NewProvider()
	clusters, err := kind.List()
	if err != nil {
		return "", err
	}
	for _, clusterName := range clusters {
		if strings.HasPrefix(clusterName, prefix) {
			return clusterName, nil
		}
	}
	return "", fmt.Errorf("no cluster found with prefix %s", prefix)
}
