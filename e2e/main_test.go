package e2e

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	ciliumv2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kwait "sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/maruina/kubecon-2024-eu/pkg/config"
)

var testenv *config.Environment

func TestMain(m *testing.M) {
	var err error
	testenv, err = config.NewEnvFromFlags()
	if err != nil {
		panic(err)
	}

	testenv.Setup(
		testenv.Validate(),
		UpdateScheme(),
	).BeforeEachFeature(func(ctx context.Context, cfg *envconf.Config, t *testing.T, f features.Feature) (context.Context, error) {
		// DO NOT enable parallel tests until
		// https://github.com/kubernetes-sigs/e2e-framework/issues/352 is fixed
		// t.Parallel()

		// Create a random name, add it to the context
		// and create a namespace with that name.
		nsPrefix := "e2e-ns"
		if f.Labels().Contains("privileged", "true") {
			nsPrefix += "-privileged"
		}
		ns := envconf.RandomName(nsPrefix, len(nsPrefix)+10)

		return createNamespace(ctx, cfg, ns)
	}).AfterEachFeature(func(ctx context.Context, cfg *envconf.Config, t *testing.T, f features.Feature) (context.Context, error) {
		if testenv.ReportMetrics {
			// Send metrics report here using t.Name() and t.Failed()
			t.Log("send metrics")
		}

		// Delete the current test namespace
		return DeleteTestNamespace(ctx, cfg)
	})

	// Launch package tests
	os.Exit(testenv.Run(m))
}

func UpdateScheme() func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		// This is a central entry point to add custom APIs/schemas
		client, err := cfg.NewClient()
		if err != nil {
			return ctx, err
		}

		// Example of how to add a CRD to the schema
		err = ciliumv2.AddToScheme(client.Resources().GetScheme())
		if err != nil {
			return ctx, err
		}

		return ctx, nil
	}
}

// createNamespace creates a namespace with the required labels and annotations.
// The namespace is added to the context to be reused later.
func createNamespace(ctx context.Context, cfg *envconf.Config, name string) (context.Context, error) {
	client, err := cfg.NewClient()
	if err != nil {
		return ctx, err
	}

	ns := corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name}}
	// Add custom labels and annotations
	labels := map[string]string{
		"app.kubernetes.io/name": "e2e-framework-demo",
	}
	annotations := map[string]string{
		"namespace.validation.example.com/safe-to-delete": "true",
	}
	client.Resources().Label(&ns, labels)
	client.Resources().Annotate(&ns, annotations)

	// Add namespace to context before the namespace is created.
	// Like this, if the namespace is created but there is an error, it will
	// still be cleaned up in the AfterEachFeature hook.
	ctx = config.AddNamespaceToCtx(ctx, name)

	if cErr := client.Resources().Create(ctx, &ns); cErr != nil {
		return ctx, cErr
	}

	// On very busy clusters, it might take some time to get the default SA
	cond := conditions.New(client.Resources())
	err = kwait.For(cond.ResourcesFound(&corev1.ServiceAccountList{
		Items: []corev1.ServiceAccount{{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "default",
				Namespace: name,
			},
		}},
	}), kwait.WithTimeout(2*time.Minute), kwait.WithImmediate(), kwait.WithInterval(1*time.Second))
	if err != nil {
		return ctx, fmt.Errorf("error waiting for default service account: %w", err)
	}

	return ctx, nil
}

// DeleteTestNamespace deletes the namespace associated with the test.
func DeleteTestNamespace(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
	client, err := cfg.NewClient()
	if err != nil {
		return ctx, err
	}

	name, err := config.GetNamespaceFromCtx(ctx)
	if err != nil {
		return ctx, err
	}

	ns := corev1.Namespace{}
	if cErr := client.Resources().Get(ctx, name, name, &ns); cErr != nil {
		return ctx, cErr
	}

	// Delete the namespace object
	if cErr := client.Resources().Delete(ctx, &ns); cErr != nil {
		return ctx, cErr
	}

	return ctx, nil
}
