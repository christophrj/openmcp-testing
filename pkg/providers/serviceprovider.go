package providers

import (
	"context"
	"testing"
	"time"

	"github.com/christophrj/openmcp-testing/pkg/clusterutils"
	"github.com/christophrj/openmcp-testing/pkg/conditions"
	"github.com/christophrj/openmcp-testing/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"
	"sigs.k8s.io/e2e-framework/klient/wait"
	e2efwconditions "sigs.k8s.io/e2e-framework/klient/wait/conditions"
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

func serviceProviderRef(name string) *unstructured.Unstructured {
	obj := &unstructured.Unstructured{}
	obj.SetName(name)
	obj.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "openmcp.cloud",
		Version: "v1alpha1",
		Kind:    "ServiceProvider",
	})
	return obj
}

// InstallServiceProvider creates a service provider object on the platform cluster and waits until it is ready
func InstallServiceProvider(opts ServiceProviderSetup, timeout time.Duration) env.Func {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		klog.Infof("create service provider: %s", opts.Name)
		obj, err := resources.CreateObjectFromTemplate(ctx, c, spTemplate, opts)
		if err != nil {
			return ctx, err
		}
		return ctx, wait.For(conditions.Match(obj, c, "Ready", corev1.ConditionTrue), wait.WithTimeout(timeout))
	}
}

// ImportServiceProviderAPIs iterates over each resource from the passed in directory
// and applies it to the onboarding cluster
func ImportServiceProviderAPIs(directory string, timeout time.Duration) features.Func {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		klog.Infof("apply service provider resources to onboarding cluster from %s ...", directory)
		c, err := clusterutils.OnboardingConfig()
		if err != nil {
			t.Errorf("failed to retrieve onboarding cluster config: %v", err)
			return ctx
		}
		objList, err := resources.CreateObjectsFromDir(ctx, c, directory)
		if err != nil {
			t.Errorf("failed to create objects from %s: %v", directory, err)
		}
		if err := wait.For(e2efwconditions.New(c.Client().Resources()).ResourcesFound(objList)); err != nil {
			t.Error(err)
		}
		return ctx
	}
}

// ImportDomainAPIs iterates over each resource from the passed in directory
// and applies it to a MCP cluster
func ImportDomainAPIs(directory string, timeout time.Duration) features.Func {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		klog.Infof("apply service provider resources to MCP cluster from %s ...", directory)
		c, err := clusterutils.McpConfig()
		if err != nil {
			t.Errorf("failed to retrieve MCP cluster config: %v", err)
			return ctx
		}
		objList, err := resources.CreateObjectsFromDir(ctx, c, directory)
		if err != nil {
			t.Errorf("failed to create objects from %s: %v", directory, err)
		}
		if err := wait.For(e2efwconditions.New(c.Client().Resources()).ResourcesFound(objList)); err != nil {
			t.Error(err)
		}
		return ctx
	}
}

// DeleteServiceProvider deletes the service provider object on the platform cluster and waits until the object has been deleted
func DeleteServiceProvider(opts ServiceProviderSetup, timeout time.Duration) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		klog.Infof("delete service provider: %s", opts.Name)
		return ctx, resources.DeleteObject(ctx, cfg, serviceProviderRef(opts.Name), wait.WithTimeout(time.Minute))
	}
}
