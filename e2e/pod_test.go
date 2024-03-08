package e2e

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/klient/wait"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	"github.com/maruina/kubecon-2024-eu/pkg/resources"
	"github.com/maruina/kubecon-2024-eu/pkg/resources/pod"
)

func TestPodEnv(t *testing.T) {
	var tc *resources.TestContext
	var err error
	var po corev1.Pod
	projectKey := "project"
	projectValue := "e2e-framework"
	feat := features.New("pod test").
		WithLabel("suite", "provisioning").
		WithLabel("category", "pod").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			// Create the test context
			tc, err = resources.NewTestContext(ctx, config)
			require.NoError(t, err)

			po = pod.NewPod(
				ctx,
				pod.WithAnnotation(projectKey, projectValue),
				pod.WithPodSpec(
					pod.WithEnvFromFieldPath("POD_NAME", "metadata.name"),
					pod.WithEnvFromFieldPath("PROJECT", "metadata.annotations['project']"),
				))
			err = tc.Create(ctx, &po)
			require.NoError(t, err)

			return ctx
		}).
		Assess("pod is ready", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			err = wait.For(tc.PodReady(&po))
			require.NoError(t, err)

			return ctx
		}).
		Assess("environment variables are injected", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			command := []string{"/bin/sh", "-c", "env"}
			stdout, execErr := tc.ExecInPod(ctx, &po, command)
			require.NoError(t, execErr)
			require.Contains(t, stdout, fmt.Sprintf("POD_NAME=%v", po.GetName()))
			require.Contains(t, stdout, fmt.Sprintf("PROJECT=%v", projectValue))

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			err = tc.Delete(ctx, &po)
			require.NoError(t, err)

			return ctx
		}).
		Feature()

	testenv.Test(t, feat)
}

func TestPodDns(t *testing.T) {
	var tc *resources.TestContext
	var err error
	var po corev1.Pod

	feat := features.New("dns test").
		WithLabel("suite", "provisioning").
		WithLabel("category", "dns").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			// Create the test context
			tc, err = resources.NewTestContext(ctx, config)
			require.NoError(t, err)

			po = pod.NewPod(ctx)
			err = tc.Create(ctx, &po)
			require.NoError(t, err)

			return ctx
		}).
		Assess("pod is ready", func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			err = wait.For(tc.PodReady(&po))
			require.NoError(t, err)

			return ctx
		}).
		Assess("dns is resolving", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			command := []string{"/bin/sh", "-c", "getent ahostsv4 datadoghq.com"}
			stdout, execErr := tc.ExecInPod(ctx, &po, command)
			require.NoError(t, execErr)
			require.Contains(t, stdout, "datadoghq.com")

			return ctx
		}).
		Teardown(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			err = tc.Delete(ctx, &po)
			require.NoError(t, err)

			return ctx
		}).
		Feature()

	testenv.Test(t, feat)
}
