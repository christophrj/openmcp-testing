package e2e

import (
	"testing"
	"time"

	"github.com/christophrj/openmcp-testing/pkg/providers"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestServiceProvider(t *testing.T) {
	basicProviderTest := features.New("provider test").
		Setup(providers.CreateMCP("test-mcp", time.Minute)).
		Assess("verify resources", providers.ImportServiceProviderAPIs("", time.Minute)).
		Teardown(providers.DeleteMCP("test-mcp", time.Minute))
	testenv.Test(t, basicProviderTest.Feature())
}
