package setup

import (
	"context"
	"fmt"
	"time"

	openmcpcond "github.com/christophrj/openmcp-testing/pkg/conditions"
	"github.com/christophrj/openmcp-testing/pkg/providers"
	"github.com/christophrj/openmcp-testing/pkg/resources"
	"github.com/vladimirvivien/gexe"
	v1 "k8s.io/api/core/v1"
	apiTypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
	"sigs.k8s.io/e2e-framework/pkg/types"
	"sigs.k8s.io/e2e-framework/support/kind"
)

type OpenMCPSetup struct {
	Namespace       string
	Operator        OpenMCPOperatorSetup
	ClusterProvider providers.CluterProviderSetup
	ServiceProvider providers.ServiceProviderSetup
}

type OpenMCPOperatorSetup struct {
	Name         string
	Namespace    string
	Image        string
	Environment  string
	PlatformName string
}

func (s *OpenMCPSetup) Bootstrap(testenv env.Environment) error {
	if err := PullImage(s.Operator.Image); err != nil {
		return err
	}
	if err := PullImage(s.ClusterProvider.Image); err != nil {
		return err
	}
	if err := PullImage(s.ServiceProvider.Image); err != nil {
		return err
	}
	platformClusterName := envconf.RandomName("platform-cluster", 16)
	s.Operator.Namespace = s.Namespace
	testenv.Setup(envfuncs.CreateClusterWithConfig(kind.NewProvider(), platformClusterName, "../pkg/setup/kind-config.yaml")).
		Setup(envfuncs.CreateNamespace(s.Namespace)).
		Setup(InstallOpenMCPOperator(s.Operator)).
		Setup(providers.InstallClusterProvider(s.ClusterProvider, time.Minute)).
		Setup(s.verifySetup()).
		Setup(providers.InstallServiceProvider(s.ServiceProvider, time.Minute)).
		Finish(providers.DeleteServiceProvider(s.ServiceProvider, time.Minute)).
		Finish(s.cleanup()).
		Finish(envfuncs.DestroyCluster(platformClusterName))
	return nil
}

func (s *OpenMCPSetup) cleanup() types.EnvFunc {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		klog.Info("cleaning up...")
		onboardingCluster := apiTypes.NamespacedName{
			Namespace: s.Namespace,
			Name:      "onboarding",
		}
		err := resources.DeleteObject(ctx, c, onboardingCluster, resources.ClusterGVR)
		if err != nil {
			return ctx, err
		}
		return ctx, wait.For(openmcpcond.New(c, resources.ClusterGVR).
			Deleted(onboardingCluster),
			wait.WithTimeout(time.Minute))
	}
}

func (s *OpenMCPSetup) verifySetup() types.EnvFunc {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		onboardingCluster := "onboarding"
		if err := wait.For(openmcpcond.New(c, resources.ClusterGVR).
			Match(apiTypes.NamespacedName{Namespace: s.Namespace, Name: onboardingCluster}, "Ready", v1.ConditionTrue),
			wait.WithTimeout(time.Minute)); err != nil {
			return ctx, err
		}
		klog.Infof("%s cluster ready", onboardingCluster)
		return ctx, nil
	}
}

func InstallOpenMCPOperator(opts OpenMCPOperatorSetup) types.EnvFunc {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		// apply openmcp operator manifests
		if err := resources.CreateObjectsFromTemplateFile(ctx, c, "../pkg/setup/templates/openmcp-operator.yaml", opts); err != nil {
			return ctx, err
		}
		// wait for deployment to be ready
		if err := wait.For(conditions.New(c.Client().Resources()).
			DeploymentAvailable(opts.Name, opts.Namespace),
			wait.WithTimeout(time.Minute)); err != nil {
			return ctx, err
		}
		klog.Info("openmcp operator ready")
		return ctx, nil
	}
}

func PullImage(image string) error {
	klog.Info("Pulling ", image)
	runner := gexe.New()
	p := runner.RunProc(fmt.Sprintf("docker pull %s", image))
	if p.Err() != nil {
		return fmt.Errorf("docker pull %v failed: %w: %s", image, p.Err(), p.Result())
	}
	return nil
}
