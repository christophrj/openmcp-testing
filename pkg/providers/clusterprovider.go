package providers

import (
	"context"
	"testing"
	"time"

	"github.com/christophrj/openmcp-testing/pkg/clusterutils"
	"github.com/christophrj/openmcp-testing/pkg/conditions"
	"github.com/christophrj/openmcp-testing/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

const cpTemplate = `
apiVersion: openmcp.cloud/v1alpha1
kind: ClusterProvider
metadata:
  name: {{.Name}}
spec:
  image: {{.Image}}
  extraVolumeMounts:
    - mountPath: /var/run/docker.sock
      name: docker
  extraVolumes:
    - name: docker
      hostPath:
        path: /var/run/host-docker.sock
        type: Socket
`

const mcpTemplate = `
apiVersion: core.openmcp.cloud/v2alpha1
kind: ManagedControlPlaneV2
metadata:
  name: {{.Name}}
spec:
  iam: {}
`

type CluterProviderSetup struct {
	Name  string
	Image string
}

func InstallClusterProvider(opts CluterProviderSetup) env.Func {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		// install cluster provider kind
		err := resources.CreateObjectsFromTemplate(ctx, c, cpTemplate, opts)
		if err != nil {
			return ctx, err
		}
		// wait for cluster provider to be ready
		return ctx, wait.For(conditions.New(c, resources.ClusterproviderGVR).
			Match(types.NamespacedName{Name: opts.Name}, "Ready", corev1.ConditionTrue),
			wait.WithTimeout(time.Minute))
	}
}

func CreateWorkloadCluster() features.Func {
	return func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
		// create workload cluster
		// wait for workload cluster to be ready
		return ctx
	}
}

func CreateMCP(name string) features.Func {
	return func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
		onboardingCfg, err := clusterutils.OnboardingConfig()
		if err != nil {
			t.Error(err)
			return ctx
		}
		// create MCP
		if err := resources.CreateObjectsFromTemplate(ctx, onboardingCfg, mcpTemplate, struct{ Name string }{Name: name}); err != nil {
			t.Errorf("failed to create MCP: %v", err)
		}
		// wait for MCP to get ready
		if err := wait.For(conditions.New(onboardingCfg, resources.ManagedControlPlaneGVR).
			Match(types.NamespacedName{Name: name, Namespace: corev1.NamespaceDefault}, "Ready", corev1.ConditionTrue),
			wait.WithTimeout(time.Minute)); err != nil {
			t.Errorf("MCP failed to get ready: %v", err)
		}
		return ctx
	}
}
