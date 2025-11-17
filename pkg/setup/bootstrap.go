package setup

import (
	"context"
	"fmt"
	"os"
	"time"

	openmcpcond "github.com/christophrj/openmcp-testing/pkg/conditions"
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
	Namespace                string
	ClusterProviderManifests string
	OperatorManifests        string
	OperatorName             string
}

func (s *OpenMCPSetup) Bootstrap(testenv env.Environment) error {
	if err := PullImage(os.Getenv("OPENMCP_OPERATOR_IMAGE")); err != nil {
		return err
	}
	if err := PullImage(os.Getenv("OPENMCP_CP_KIND_IMAGE")); err != nil {
		return err
	}
	platformClusterName := envconf.RandomName("platform-cluster", 16)
	testenv.Setup(envfuncs.CreateClusterWithConfig(kind.NewProvider(), platformClusterName, "kind-config.yaml")).
		Setup(envfuncs.CreateNamespace(s.Namespace)).
		Setup(InstallOpenMCPOperator(s.OperatorName, s.Namespace, s.OperatorManifests)).
		Setup(InstallClusterProvider("kind", s.ClusterProviderManifests)).
		Setup(s.verifySetup()).
		Finish(s.cleanup()).
		Finish(envfuncs.DestroyCluster(platformClusterName))
	return nil
}

func (s *OpenMCPSetup) cleanup() types.EnvFunc {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		klog.Info("cleaning up...")
		cl, err := resources.NewFromEnvConfig(c)
		if err != nil {
			return ctx, err
		}
		onboardingCluster := apiTypes.NamespacedName{
			Namespace: s.Namespace,
			Name:      "onboarding",
		}
		err = cl.DeleteObject(ctx, onboardingCluster, resources.ClusterGVR)
		if err != nil {
			return ctx, err
		}
		return ctx, wait.For(openmcpcond.New(cl.Resources).ClusterDelete(onboardingCluster), wait.WithTimeout(time.Minute))
	}
}

func (s *OpenMCPSetup) verifySetup() types.EnvFunc {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		cl, err := resources.NewFromEnvConfig(c)
		if err != nil {
			return ctx, err
		}
		onboardingCluster := "onboarding"
		if err := wait.For(openmcpcond.New(cl.Resources).ClusterConditionMatch(apiTypes.NamespacedName{
			Namespace: s.Namespace,
			Name:      onboardingCluster,
		}, "Ready", v1.ConditionTrue), wait.WithTimeout(time.Minute)); err != nil {
			return ctx, err
		}
		klog.Infof("%s cluster ready", onboardingCluster)
		return ctx, nil
	}
}

func InstallClusterProvider(name string, clusterProviderManifestPath string) types.EnvFunc {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		// install cluster provider kind
		err := resources.CreateObjectsFromFile(ctx, c, clusterProviderManifestPath)
		if err != nil {
			return ctx, err
		}
		// wait for cluster provider to be ready
		cl, err := resources.NewFromEnvConfig(c)
		if err != nil {
			return ctx, err
		}
		if err := wait.For(openmcpcond.New(cl.Resources).ClusterProviderConditionMatch(name, "Ready", v1.ConditionTrue), wait.WithTimeout(time.Minute)); err != nil {
			return ctx, err
		}
		klog.Infof("cluster provider %s ready", name)
		return ctx, nil
	}
}

func InstallOpenMCPOperator(name string, nameSpace string, manifests string) types.EnvFunc {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		// apply openmcp operator manifests
		if err := resources.CreateObjectsFromFile(ctx, c, manifests); err != nil {
			return ctx, err
		}
		cl, err := resources.NewFromEnvConfig(c)
		if err != nil {
			return ctx, err
		}
		// wait for deployment to be ready
		if err := wait.For(conditions.New(cl.Resources).DeploymentAvailable(name, nameSpace), wait.WithTimeout(time.Minute)); err != nil {
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
