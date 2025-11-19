package e2e

import (
	"testing"
	"time"

	"github.com/christophrj/openmcp-testing/pkg/providers"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestServiceProvider(t *testing.T) {
	basicProviderTest := features.New("provider test").
		Setup(providers.CreateWorkloadCluster()).
		Setup(providers.CreateMCP("test-mcp")).
		Assess("verify resources", providers.VerifyServiceProvider([]string{}, time.Minute)).
		Teardown(providers.DeleteMCP("test-mcp")).
		Teardown(providers.DelelteServiceProvider())
	testenv.Test(t, basicProviderTest.Feature())
}
