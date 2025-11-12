package providers

import (
	"context"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func CreateWorkloadCluster() features.Func {
	return func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
		// create workload cluster
		// wait for workload cluster to be ready
		return ctx
	}
}

func CreateMCP() features.Func {
	return func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
		// create mcp
		// wait for mcp to be ready
		return ctx
	}
}
