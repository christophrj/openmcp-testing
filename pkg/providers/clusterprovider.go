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
	"k8s.io/klog/v2"
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

// InstallClusterProvider creates a cluster provider object on the platform cluster and waits until it is ready
func InstallClusterProvider(opts CluterProviderSetup, timeout time.Duration) env.Func {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		klog.Infof("create cluster provider %s", opts.Name)
		err := resources.CreateObjectsFromTemplate(ctx, c, cpTemplate, opts)
		if err != nil {
			return ctx, err
		}
		return ctx, wait.For(conditions.New(c, resources.ClusterproviderGVR).
			Match(types.NamespacedName{Name: opts.Name}, "Ready", corev1.ConditionTrue),
			wait.WithTimeout(timeout))
	}
}

// CreateMCP creates an MCP object on the onboarding cluster and waits until it is ready
func CreateMCP(name string, timeout time.Duration) features.Func {
	return func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
		klog.Infof("create MCP %s", name)
		onboardingCfg, err := clusterutils.OnboardingConfig()
		if err != nil {
			t.Error(err)
			return ctx
		}
		if err := resources.CreateObjectsFromTemplate(ctx, onboardingCfg, mcpTemplate, struct{ Name string }{Name: name}); err != nil {
			t.Errorf("failed to create MCP: %v", err)
			return ctx
		}
		if err := wait.For(conditions.New(onboardingCfg, resources.ManagedControlPlaneGVR).
			Status(types.NamespacedName{Name: name, Namespace: corev1.NamespaceDefault}, "phase", "Ready"),
			wait.WithTimeout(timeout)); err != nil {
			t.Errorf("MCP failed to get ready: %v", err)
		}
		return ctx
	}
}

// DeleteMCP deletes the MCP object on the onboarding cluster and waits until the object has been deleted
func DeleteMCP(name string, timeout time.Duration) features.Func {
	return func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
		klog.Infof("delete MCP %s", name)
		onboardingCfg, err := clusterutils.OnboardingConfig()
		if err != nil {
			t.Error(err)
			return ctx
		}
		mcp := types.NamespacedName{
			Namespace: corev1.NamespaceDefault,
			Name:      name,
		}
		err = resources.DeleteObject(ctx, onboardingCfg, mcp, resources.ManagedControlPlaneGVR)
		if err != nil {
			t.Errorf("failed to delete MCP %s: %v", name, err)
			return ctx
		}
		if err := wait.For(conditions.New(onboardingCfg, resources.ManagedControlPlaneGVR).Deleted(mcp),
			wait.WithTimeout(timeout)); err != nil {
			t.Errorf("delete MCP %s timed out: %v", name, err)
		}
		return ctx
	}
}
