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

const spTemplate = `
apiVersion: openmcp.cloud/v1alpha1
kind: ServiceProvider
metadata:
  name: {{.Name}}
spec:
  image: {{.Image}}
`

type ServiceProviderSetup struct {
	Name  string
	Image string
}

func InstallServiceProvider(opts ServiceProviderSetup) env.Func {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		err := resources.CreateObjectsFromTemplate(ctx, c, spTemplate, opts)
		if err != nil {
			return ctx, err
		}
		return ctx, wait.For(conditions.New(c, resources.ServiceProviderGVR).
			Match(types.NamespacedName{Name: opts.Name}, "Ready", v1.ConditionTrue))
	}
}

// VerifyServiceProvider iterates over each resource directory and waits until each resource is synced and ready.
func VerifyServiceProvider(directories []string, timeout time.Duration) features.Func {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		for range directories {
			// wait for resources to be synced and ready
		}
		return ctx
	}
}

func DelelteServiceProvider() features.Func {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		return ctx
	}
}
