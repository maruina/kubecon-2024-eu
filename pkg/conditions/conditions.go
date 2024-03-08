package conditions

import (
	appsv1 "k8s.io/api/apps/v1"
	apimachinerywait "k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/e2e-framework/klient/k8s"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
	"sigs.k8s.io/e2e-framework/klient/wait/conditions"
)

type Condition struct {
	*conditions.Condition
	resources *resources.Resources
}

// New is used to create a new Condition that can be used to perform
// a series of wait checks against a resource in question.
// Wraps the e2e-framework Condition.
func New(r *resources.Resources) *Condition {
	return &Condition{
		Condition: conditions.New(r),
		resources: r,
	}
}

func (c *Condition) StsReady(sts *appsv1.StatefulSet) apimachinerywait.ConditionWithContextFunc {
	return c.ResourceMatch(sts, func(object k8s.Object) bool {
		obj, ok := object.(*appsv1.StatefulSet)
		if !ok {
			return false
		}

		return (obj.Status.UpdatedReplicas == *obj.Spec.Replicas) && (obj.Status.ReadyReplicas == *obj.Spec.Replicas)
	})
}
