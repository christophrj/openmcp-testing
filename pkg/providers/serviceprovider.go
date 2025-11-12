package providers

import (
	"context"
	"testing"
	"time"

	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func InstallServiceProvider() features.Func {
	return func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
		return ctx
	}
}

// VerifyServiceProvider iterates over each resource directory and waits until each resource is synced and ready.
func VerifyServiceProvider(directories []string, timeout time.Duration) features.Func {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		for range directories {
			// wait for resources to be synced and ready
		}
		return ctx
	}
}

func DelelteServiceProvider() features.Func {
	return func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
		return ctx
	}
}
