package resources

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/klient/k8s"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
	"sigs.k8s.io/e2e-framework/pkg/envconf"

	"github.com/maruina/kubecon-2024-eu/pkg/conditions"
	"github.com/maruina/kubecon-2024-eu/pkg/config"
)

// TestContext holds test-specific context. It can be extended on a
// per-test basis to hold additional information related to the test.
// Wraps resources.Resources and conditions.Condition.
type TestContext struct {
	*resources.Resources
	*conditions.Condition

	namespace string
}

// NewTestContext creates a new test context that can be used to create and
// interact with kubernetes resources from inside a test.
func NewTestContext(ctx context.Context, cfg *envconf.Config) (*TestContext, error) {
	ns, err := config.GetNamespaceFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	client, err := cfg.NewClient()
	if err != nil {
		return nil, err
	}

	res := client.Resources(ns)
	tc := &TestContext{
		Resources: res,
		Condition: conditions.New(res),
		namespace: ns,
	}

	return tc, nil
}

// Get returns the object with the given name in the test context's namespace.
// Wraps resources.Resources.Get. To get resources in a different namespace,
// call tc.Resources.Get directly.
func (tc *TestContext) Get(ctx context.Context, name string, obj k8s.Object) error {
	return tc.Resources.Get(ctx, name, tc.namespace, obj)
}

// ExecInPod returns the result of an exec in the pod.
// Wraps resources.Resources.ExecInPod to not manually manage stdout and stderr.
func (tc *TestContext) ExecInPod(ctx context.Context, pod *corev1.Pod, cmd []string) (output string, err error) {
	stdOut := bytes.Buffer{}
	stdErr := bytes.Buffer{}
	err = tc.Resources.ExecInPod(ctx, pod.Namespace, pod.Name, pod.Spec.Containers[0].Name, cmd, &stdOut, &stdErr)
	output = strings.TrimSuffix(stdOut.String(), "\n")
	if err != nil {
		return
	}
	if stdErr.String() != "" {
		err = fmt.Errorf("error running command: %s", stdErr.String())

		return
	}

	return
}
