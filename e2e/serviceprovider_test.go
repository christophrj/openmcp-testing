package e2e

import (
	"testing"
	"time"

	"github.com/christophrj/openmcp-testing/pkg/providers"
	"github.com/christophrj/openmcp-testing/pkg/resources"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestServiceProvider(t *testing.T) {
	basicProviderTest := features.New("provider test").
		Setup(providers.InstallServiceProvider()).
		Setup(resources.ImportResources([]string{})).
		Assess("verify resources", providers.VerifyServiceProvider([]string{}, time.Minute)).
		Teardown(providers.DelelteServiceProvider())
	testenv.Test(t, basicProviderTest.Feature())
}
