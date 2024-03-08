package container

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/maruina/kubecon-2024-eu/pkg/config"
)

const Name = "basic"

func NewContainer(opts ...Opts) corev1.Container {
	c := corev1.Container{
		Name:    Name,
		Image:   fmt.Sprintf("%s/%s", config.GetEnv().DefaultContainerRegistry, config.GetEnv().DefaultContainerImage),
		Command: []string{"sleep", "infinity"},
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				"cpu":    resource.MustParse("100m"),
				"memory": resource.MustParse("256Mi"),
			},
		},
	}
	for _, opt := range opts {
		opt(&c)
	}

	return c
}

type Opts func(container *corev1.Container)

// WithEnvFromFieldPath adds an environment variable to the container based on a field path.
func WithEnvFromFieldPath(name, fieldPath string) Opts {
	return func(container *corev1.Container) {
		container.Env = append(container.Env, corev1.EnvVar{
			Name: name,
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: fieldPath,
				},
			},
		})
	}
}
