package pod

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/maruina/kubecon-2024-eu/pkg/config"
	"github.com/maruina/kubecon-2024-eu/pkg/resources/container"
)

const Name = "pod"

// NewPod creates a new instance of Pod.
func NewPod(ctx context.Context, opts ...OptsPod) corev1.Pod {
	p := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        Name,
			Namespace:   config.MustGetNamespaceFromCtx(ctx),
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
		},
		Spec: NewPodSpec(),
	}

	for _, opt := range opts {
		opt(&p)
	}

	return p
}

type OptsPod func(pod *corev1.Pod)

func WithAnnotation(key, value string) OptsPod {
	return func(pod *corev1.Pod) {
		pod.ObjectMeta.Annotations[key] = value
	}
}

func WithPodSpec(opts ...OptsPodSpec) OptsPod {
	return func(pod *corev1.Pod) {
		pod.Spec = NewPodSpec(opts...)
	}
}

type OptsPodSpec func(pod *corev1.PodSpec)

// NewPodSpec creates a new instance of PodSpec with default values.
func NewPodSpec(opts ...OptsPodSpec) corev1.PodSpec {
	podSpec := corev1.PodSpec{
		Containers: []corev1.Container{container.NewContainer()},
		// Exit pod quickly since using bash -c does not catch SIGTERM correctly
		TerminationGracePeriodSeconds: ptr.To(int64(1)),
	}

	for _, opt := range opts {
		opt(&podSpec)
	}

	return podSpec
}

// WithEnvFromFieldPath adds an environment variable to the first container based on a field path.
func WithEnvFromFieldPath(name, fieldPath string) OptsPodSpec {
	return func(podSpec *corev1.PodSpec) {
		container.WithEnvFromFieldPath(name, fieldPath)(&podSpec.Containers[0])
	}
}
