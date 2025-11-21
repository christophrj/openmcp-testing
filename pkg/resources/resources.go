package resources

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/christophrj/openmcp-testing/internal"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/klog/v2"
	"sigs.k8s.io/e2e-framework/klient/decoder"
	"sigs.k8s.io/e2e-framework/klient/k8s"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func DeleteObject(ctx context.Context, c *envconf.Config, obj k8s.Object, options ...wait.Option) error {
	err := c.Client().Resources().Get(ctx, obj.GetName(), obj.GetNamespace(), obj)
	if err != nil {
		return internal.IgnoreNotFound(err)
	}
	if err = c.Client().Resources().Delete(ctx, obj); err != nil {
		return internal.IgnoreNotFound(err)
	}
	if options != nil {
		return wait.For(conditions.New(c.Client().Resources()).ResourceDeleted(obj), options...)
	}
	return nil
}

func CreateObjectsFromTemplateFile(ctx context.Context, cfg *envconf.Config, filePath string, data interface{}) (*unstructured.UnstructuredList, error) {
	manifest, err := internal.ExecTemplateFile(filePath, data)
	if err != nil {
		return nil, err
	}
	return createObjectsFromManifest(ctx, cfg, manifest)
}

func CreateObjectFromTemplate(ctx context.Context, cfg *envconf.Config, template string, data interface{}) (*unstructured.Unstructured, error) {
	manifest, err := internal.ExecTemplate(template, data)
	if err != nil {
		return nil, err
	}
	// TODO single creation
	list, err := createObjectsFromManifest(ctx, cfg, manifest)
	if err != nil {
		return nil, err
	}
	if len(list.Items) < 1 {
		return nil, fmt.Errorf("can't return object from empty list")
	}
	return &list.Items[0], nil
}

func createObjectsFromManifest(ctx context.Context, cfg *envconf.Config, manifest string) (*unstructured.UnstructuredList, error) {
	r := strings.NewReader(manifest)
	list := &unstructured.UnstructuredList{}
	err := decoder.DecodeEach(ctx, r,
		func(ctx context.Context, obj k8s.Object) error {
			return createAndPopulateList(ctx, obj, list, cfg)
		}, decoder.MutateNamespace(cfg.Namespace()))
	return list, err
}

func CreateObjectsFromDir(ctx context.Context, cfg *envconf.Config, dir string) (*unstructured.UnstructuredList, error) {
	list := &unstructured.UnstructuredList{}
	err := decoder.DecodeEachFile(ctx, os.DirFS(dir), "*",
		func(ctx context.Context, obj k8s.Object) error {
			return createAndPopulateList(ctx, obj, list, cfg)
		}, decoder.MutateNamespace(cfg.Namespace()))
	return list, err
}

func createAndPopulateList(ctx context.Context, obj k8s.Object, list *unstructured.UnstructuredList, cfg *envconf.Config) error {
	u, err := internal.ToUnstructured(obj)
	if err != nil {
		return err
	}
	list.Items = append(list.Items, *u)
	klog.Infof("creating object (%s) %s/%s", obj.GetObjectKind().GroupVersionKind(), obj.GetNamespace(), obj.GetName())
	return decoder.CreateIgnoreAlreadyExists(cfg.Client().Resources())(ctx, obj)
}
