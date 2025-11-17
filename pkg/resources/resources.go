package resources

import (
	"context"
	"io"
	"os"

	"github.com/christophrj/openmcp-testing/internal"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/e2e-framework/klient/decoder"
	"sigs.k8s.io/e2e-framework/klient/k8s"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

var (
	ClusterproviderGVR = schema.GroupVersionResource{Group: "openmcp.cloud", Version: "v1alpha1", Resource: "clusterproviders"}
	ClusterGVR         = schema.GroupVersionResource{Group: "clusters.openmcp.cloud", Version: "v1alpha1", Resource: "clusters"}
)

type Client struct {
	Resources *resources.Resources
}

func New(cfg *rest.Config) (*Client, error) {
	r, err := resources.New(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{
		Resources: r,
	}, nil
}

func (c *Client) WithNamespace(namespace string) *Client {
	c.Resources = c.Resources.WithNamespace(namespace)
	return c
}

func NewFromEnvConfig(cfg *envconf.Config) (*Client, error) {
	c, err := New(cfg.Client().RESTConfig())
	if err != nil {
		return nil, err
	}
	return c.WithNamespace(cfg.Namespace()), nil
}

func (c *Client) GetObject(ctx context.Context, ref types.NamespacedName, gvr schema.GroupVersionResource) (*unstructured.Unstructured, error) {
	cl, err := dynamic.NewForConfig(c.Resources.GetConfig())
	if err != nil {
		return nil, err
	}
	res := cl.Resource(gvr)
	return res.Namespace(ref.Namespace).Get(ctx, ref.Name, metav1.GetOptions{})
}

func (c *Client) DeleteObject(ctx context.Context, ref types.NamespacedName, gvr schema.GroupVersionResource) error {
	cl, err := dynamic.NewForConfig(c.Resources.GetConfig())
	if err != nil {
		return err
	}
	res := cl.Resource(gvr)
	return res.Namespace(ref.Namespace).Delete(ctx, ref.Name, metav1.DeleteOptions{})
}

func CreateObjectsFromFile(ctx context.Context, cfg *envconf.Config, filePath string) error {
	substFile, err := substitute(filePath)
	if err != nil {
		return err
	}
	cl, err := NewFromEnvConfig(cfg)
	if err != nil {
		return err
	}
	return decoder.DecodeEach(ctx, substFile, decoder.CreateIgnoreAlreadyExists(cl.Resources), decoder.MutateNamespace(cfg.Namespace()))
}

func GetObjectsFromFile(ctx context.Context, cfg *envconf.Config, filePath string) ([]k8s.Object, error) {
	substFile, err := substitute(filePath)
	if err != nil {
		return nil, err
	}
	cl, err := NewFromEnvConfig(cfg)
	if err != nil {
		return nil, err
	}
	objects := make([]k8s.Object, 0)
	err = decoder.DecodeEach(ctx, substFile, decoder.ReadHandler(cl.Resources, func(ctx context.Context, obj k8s.Object) error {
		objects = append(objects, obj)
		return nil
	}), decoder.MutateNamespace(cfg.Namespace()))
	return objects, err
}

func substitute(filePath string) (io.Reader, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return internal.SubstitutePlaceholders(f)
}
