package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/christophrj/openmcp-testing/pkg/providers"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestServiceProvider(t *testing.T) {
	basicProviderTest := features.New("provider test").
		Setup(providers.CreateMCP("test-mcp", time.Minute)).
		Setup(providers.ImportServiceProviderAPIs("serviceproviderobjects", time.Minute)).
		Setup(providers.ImportDomainAPIs("domainobjects", time.Minute)).
		Assess("verify API status conditions", func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			return ctx
		}).
		Teardown(providers.DeleteMCP("test-mcp", time.Minute))
	testenv.Test(t, basicProviderTest.Feature())
}
