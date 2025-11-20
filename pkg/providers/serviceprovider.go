package providers

import (
	"context"
	"testing"
	"time"

	"github.com/christophrj/openmcp-testing/pkg/conditions"
	"github.com/christophrj/openmcp-testing/pkg/resources"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
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

// InstallServiceProvider creates a service provider object on the platform cluster and waits until it is ready
func InstallServiceProvider(opts ServiceProviderSetup, timeout time.Duration) env.Func {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		klog.Infof("create service provider %s", opts.Name)
		err := resources.CreateObjectsFromTemplate(ctx, c, spTemplate, opts)
		if err != nil {
			return ctx, err
		}
		return ctx, wait.For(conditions.New(c, resources.ServiceProviderGVR).
			Match(types.NamespacedName{Name: opts.Name}, "Ready", v1.ConditionTrue),
			wait.WithTimeout(timeout))
	}
}

// ImportServiceProviderAPIs iterates over each resource from the passed in directory
// and applies it to the onboarding cluster
func ImportServiceProviderAPIs(directory string, timeout time.Duration) features.Func {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		// import and wait for resources to be synced and ready
		return ctx
	}
}

// ImportDomainAPIs iterates over each resource from the passed in directory
// and applies it to a MCP cluster
func ImportDomainAPIs(directory string, timeout time.Duration) features.Func {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		// import and wait for resources to be synced and ready
		return ctx
	}
}

// DeleteServiceProvider deletes the service provider object on the platform cluster and waits until the object has been deleted
func DeleteServiceProvider(opts ServiceProviderSetup, timeout time.Duration) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		klog.Infof("delete service provider %s", opts.Name)
		serviceProvider := types.NamespacedName{
			Name: opts.Name,
		}
		err := resources.DeleteObject(ctx, cfg, serviceProvider, resources.ServiceProviderGVR)
		if err != nil {
			return ctx, err
		}
		return ctx, wait.For(conditions.New(cfg, resources.ServiceProviderGVR).
			Deleted(serviceProvider),
			wait.WithTimeout(time.Minute))
	}
}
