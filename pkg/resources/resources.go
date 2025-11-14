package resources

import (
	"context"
	"io"
	"os"

	"github.com/christophrj/openmcp-testing/internal"
	"sigs.k8s.io/e2e-framework/klient/decoder"
	"sigs.k8s.io/e2e-framework/klient/k8s"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

func ResClient(cfg *envconf.Config) (*resources.Resources, error) {
	r, err := resources.New(cfg.Client().RESTConfig())
	if err != nil {
		return nil, err
	}
	r.WithNamespace(cfg.Namespace())
	return r, nil
}

func CreateObjectsFromFile(ctx context.Context, cfg *envconf.Config, filePath string) error {
	substFile, err := substitute(filePath)
	if err != nil {
		return err
	}
	cl, err := ResClient(cfg)
	if err != nil {
		return err
	}
	return decoder.DecodeEach(ctx, substFile, decoder.CreateIgnoreAlreadyExists(cl), decoder.MutateNamespace(cfg.Namespace()))
}

func GetObjectsFromFile(ctx context.Context, cfg *envconf.Config, filePath string) ([]k8s.Object, error) {
	substFile, err := substitute(filePath)
	if err != nil {
		return nil, err
	}
	cl, err := ResClient(cfg)
	if err != nil {
		return nil, err
	}
	objects := make([]k8s.Object, 0)
	err = decoder.DecodeEach(ctx, substFile, decoder.ReadHandler(cl, func(ctx context.Context, obj k8s.Object) error {
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
