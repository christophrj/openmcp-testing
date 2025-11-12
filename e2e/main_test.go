package e2e

import (
	"os"
	"testing"

	"github.com/christophrj/openmcp-testing/pkg/setup"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/third_party/kind"
)

var testenv env.Environment

func TestMain(m *testing.M) {
	testenv = env.New()
	openmcp := setup.OpenMCPSetup{}
	openmcp.Bootstrap(testenv, &kind.Cluster{})
	os.Exit(testenv.Run(m))
}
