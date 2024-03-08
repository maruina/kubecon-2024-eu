package config

import (
	"context"
	"flag"
	"fmt"

	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

var environment = Environment{}

// Environment represents the global test configuration.
// Wraps the e2e-framework Environment.
type Environment struct {
	env.Environment
	DefaultContainerRegistry string
	DefaultContainerImage    string
	ReportMetrics            bool
}

// NewEnvFromFlags configure the global Environment from flags.
// This is a wrapper around e2e-framework's NewFromFlags to allow
// for additional flags.
// The function returns a pointer so that the global Environment can be
// tuned as needed after creation.
func NewEnvFromFlags() (*Environment, error) {
	flag.StringVar(&environment.DefaultContainerRegistry, "default-container-registry", "registry.k8s.io/e2e-test-images", "The registry to pull the default container image from. Default to registry.k8s.io/e2e-test-images.")
	flag.StringVar(&environment.DefaultContainerImage, "default-container-image", "agnhost:2.47", "The image for the default container. Default to agnhost:2.47.")
	flag.BoolVar(&environment.ReportMetrics, "report-metrics", false, "Enable/disable reporting tests metrics. Default to false.")

	e, err := env.NewFromFlags()
	environment.Environment = e

	return &environment, err
}

// GetEnv returns a copy of the global Environment.
func GetEnv() Environment {
	return environment
}

func (e *Environment) Validate() func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		if e.DefaultContainerRegistry == "" {
			return ctx, fmt.Errorf("--default-container-registry flag provided but set to empty string")
		}
		if e.DefaultContainerImage == "" {
			return ctx, fmt.Errorf("--default-container-image flag provided but set to empty string")
		}

		return ctx, nil
	}
}
