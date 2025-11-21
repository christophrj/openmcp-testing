package setup

import (
	"context"
	"fmt"
	"time"

	"github.com/christophrj/openmcp-testing/pkg/providers"
	"github.com/christophrj/openmcp-testing/pkg/resources"
	"github.com/vladimirvivien/gexe"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
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
	testenv.Setup(createPlatformCluster(platformClusterName)).
		Setup(envfuncs.CreateNamespace(s.Namespace)).
		Setup(InstallOpenMCPOperator(s.Operator)).
		Setup(providers.InstallClusterProvider(s.ClusterProvider, time.Minute)).
		Setup(s.verifyEnvironment()).
		Setup(providers.InstallServiceProvider(s.ServiceProvider, time.Minute)).
		Finish(providers.DeleteServiceProvider(s.ServiceProvider, time.Minute)).
		Finish(s.cleanup()).
		Finish(envfuncs.DestroyCluster(platformClusterName))
	return nil
}

func createPlatformCluster(name string) types.EnvFunc {
	klog.Info("create platform cluster...")
	return envfuncs.CreateClusterWithConfig(kind.NewProvider(), name, "../pkg/setup/kind-config.yaml")
}

func (s *OpenMCPSetup) cleanup() types.EnvFunc {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		klog.Info("cleaning up environment...")
		onboardingCluster := apimachinerytypes.NamespacedName{
			Namespace: s.Namespace,
			Name:      "onboarding",
		}
		return ctx, providers.DeleteCluster(ctx, c, onboardingCluster, wait.WithTimeout(time.Second*10))
	}
}

func (s *OpenMCPSetup) verifyEnvironment() types.EnvFunc {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		klog.Info("verify environment...")
		obj := apimachinerytypes.NamespacedName{
			Namespace: s.Namespace,
			Name:      "onboarding",
		}
		return ctx, providers.ClusterReady(ctx, c, obj, wait.WithTimeout(time.Minute))
	}
}

func InstallOpenMCPOperator(opts OpenMCPOperatorSetup) types.EnvFunc {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		// apply openmcp operator manifests
		if _, err := resources.CreateObjectsFromTemplateFile(ctx, c, "../pkg/setup/templates/openmcp-operator.yaml", opts); err != nil {
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
