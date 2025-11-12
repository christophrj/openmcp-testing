package resources

import (
	"context"
	"os"
	"testing"

	"sigs.k8s.io/e2e-framework/klient/decoder"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func ResClient(cfg *envconf.Config) (*resources.Resources, error) {
	r, err := resources.New(cfg.Client().RESTConfig())
	if err != nil {
		return nil, err
	}
	return r, nil
}

func ImportResources(directories []string) features.Func {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		for _, dir := range directories {
			if err := importResources(ctx, cfg, dir); err != nil {
				t.Fatalf("import failed: %v", err)
			}
		}
		return ctx
	}
}

func importResources(ctx context.Context, cfg *envconf.Config, directory string) error {
	cl, err := ResClient(cfg)
	if err != nil {
		return err
	}
	if err := decoder.DecodeEachFile(ctx, os.DirFS(directory), "*", decoder.CreateIgnoreAlreadyExists(cl)); err != nil {
		return err
	}
	return nil
}
