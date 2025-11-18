package e2e

import (
	"fmt"
	"os"
	"testing"

	"github.com/christophrj/openmcp-testing/pkg/providers"
	"github.com/christophrj/openmcp-testing/pkg/setup"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

var testenv env.Environment

func TestMain(m *testing.M) {
	openmcp := setup.OpenMCPSetup{
		Namespace: "openmcp-system",
		Operator: setup.OpenMCPOperatorSetup{
			Name:         "openmcp-operator",
			Image:        "ghcr.io/openmcp-project/images/openmcp-operator:v0.13.0",
			Environment:  "debug",
			PlatformName: "platform",
		},
		ClusterProvider: providers.CluterProviderSetup{
			Name:  "kind",
			Image: "ghcr.io/openmcp-project/images/cluster-provider-kind:v0.0.15",
		},
		ServiceProvider: providers.ServiceProviderSetup{
			Name:  "crossplane",
			Image: "ghcr.io/openmcp-project/images/service-provider-crossplane:v0.0.4",
		},
	}
	testenv = env.NewWithConfig(envconf.New().WithNamespace(openmcp.Namespace))
	if err := openmcp.Bootstrap(testenv); err != nil {
		panic(fmt.Errorf("openmcp bootstrap failed: %v", err))
	}
	os.Exit(testenv.Run(m))
}
