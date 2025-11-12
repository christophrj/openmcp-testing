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
		Setup(InstallClusterProviderKind()).
		Setup(InstallOpenMCPOperator())
}

func InstallClusterProviderKind() types.EnvFunc {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		// install cluster provider kind
		// wait for cluster profiles to be ready
		return ctx, nil
	}
}

func InstallOpenMCPOperator() types.EnvFunc {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		// install openmcp operator
		// wait for onboarding cluster to be ready
		return ctx, nil
	}
}
