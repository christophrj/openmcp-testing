package setup

import (
	"context"

	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
	"sigs.k8s.io/e2e-framework/pkg/types"
	"sigs.k8s.io/e2e-framework/support/kind"
)

type OpenMCPSetup struct {
}

func (s *OpenMCPSetup) Bootstrap(testenv env.Environment, cluster *kind.Cluster) {
	testenv.Setup(envfuncs.CreateCluster(cluster, "platform-cluster")).
		Setup(InstallOpenMCPOperator()).
		Setup(CreateMCP())
}

func InstallOpenMCPOperator() types.EnvFunc {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		// install openmcp operator
		// wait for onboarding cluster to be ready
		return ctx, nil
	}
}

func CreateMCP() types.EnvFunc {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		// create mcp
		// wait for mcp to be ready
		return ctx, nil
	}
}
