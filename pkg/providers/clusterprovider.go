package providers

import (
	"context"
	"testing"
	"time"

	"github.com/christophrj/openmcp-testing/pkg/conditions"
	"github.com/christophrj/openmcp-testing/pkg/resources"
	v1 "k8s.io/api/core/v1"
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
			Match(types.NamespacedName{Name: opts.Name}, "Ready", v1.ConditionTrue),
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

func CreateMCP() features.Func {
	return func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
		// create mcp
		// wait for mcp to be ready
		return ctx
	}
}
