package resources

import (
	"context"
	"strings"

	"github.com/christophrj/openmcp-testing/internal"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/e2e-framework/klient/decoder"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

var (
	ClusterproviderGVR     = schema.GroupVersionResource{Group: "openmcp.cloud", Version: "v1alpha1", Resource: "clusterproviders"}
	ServiceProviderGVR     = schema.GroupVersionResource{Group: "openmcp.cloud", Version: "v1alpha1", Resource: "serviceproviders"}
	ClusterGVR             = schema.GroupVersionResource{Group: "clusters.openmcp.cloud", Version: "v1alpha1", Resource: "clusters"}
	ManagedControlPlaneGVR = schema.GroupVersionResource{Group: "core.openmcp.cloud", Version: "v2alpha1", Resource: "managedcontrolplanev2s"}
)

func GetObject(ctx context.Context, c *envconf.Config, ref types.NamespacedName, gvr schema.GroupVersionResource) (*unstructured.Unstructured, error) {
	cl, err := dynamic.NewForConfig(c.Client().RESTConfig())
	if err != nil {
		return nil, err
	}
	res := cl.Resource(gvr)
	return res.Namespace(ref.Namespace).Get(ctx, ref.Name, metav1.GetOptions{})
}

func DeleteObject(ctx context.Context, c *envconf.Config, ref types.NamespacedName, gvr schema.GroupVersionResource) error {
	cl, err := dynamic.NewForConfig(c.Client().RESTConfig())
	if err != nil {
		return err
	}
	res := cl.Resource(gvr)
	return res.Namespace(ref.Namespace).Delete(ctx, ref.Name, metav1.DeleteOptions{})
}

func CreateObjectsFromTemplateFile(ctx context.Context, cfg *envconf.Config, filePath string, data interface{}) error {
	manifest, err := internal.ExecTemplateFile(filePath, data)
	if err != nil {
		return err
	}
	return CreateObjectsFromManifest(ctx, cfg, manifest)
}

func CreateObjectsFromTemplate(ctx context.Context, cfg *envconf.Config, template string, data interface{}) error {
	manifest, err := internal.ExecTemplate(template, data)
	if err != nil {
		return err
	}
	return CreateObjectsFromManifest(ctx, cfg, manifest)
}

func CreateObjectsFromManifest(ctx context.Context, cfg *envconf.Config, manifest string) error {
	r := strings.NewReader(manifest)
	return decoder.DecodeEach(ctx, r, decoder.CreateIgnoreAlreadyExists(cfg.Client().Resources()), decoder.MutateNamespace(cfg.Namespace()))
}
