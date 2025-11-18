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
		Namespace:                "openmcp-system",
		ClusterProviderManifests: "crs/setup/cluster-provider-kind.yaml",
		OperatorManifests:        "crs/setup/openmcp-operator.yaml",
		OperatorName:             "openmcp-operator",
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
