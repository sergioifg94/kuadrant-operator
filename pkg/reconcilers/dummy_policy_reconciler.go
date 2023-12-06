package reconcilers

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DummyPolicy struct {
	client.Object

	Status struct {
		Conditions []metav1.Condition
	}
}

func (p *DummyPolicy) GetStatusConditions() *[]metav1.Condition {
	return &p.Status.Conditions
}

type DummyPolicyReconciler struct {
	*PolicyReconciler[*DummyPolicy]
}

func NewDummyPolicyReconciler() *DummyPolicyReconciler {
	r := &DummyPolicyReconciler{
		PolicyReconciler: &PolicyReconciler[*DummyPolicy]{
			PolicyTemplate: func() *DummyPolicy {
				return &DummyPolicy{}
			},
		},
	}

	r.PolicyReconciler.Validate = r.Validate
	r.PolicyReconciler.ReconcilePolicy = r.ReconcilePolicy

	return r
}

func (r *DummyPolicyReconciler) Validate(ctx context.Context, policy *DummyPolicy) (ValidationResult, error) {
	var relatedResource client.Object
	r.client.Get(ctx, types.NamespacedName{
		Name: "dummy",
	}, relatedResource)

	return ValidationResult{
		Success: true,
	}, nil
}

func (r *DummyPolicyReconciler) ReconcilePolicy(ctx context.Context, policy *DummyPolicy) (ReconciliationResult, error) {
	return ReconciliationResult{}, nil
}

func (r *DummyPolicyReconciler) ServiceProbe(ctx context.Context, policy *DummyPolicy) (ServiceProbeResult, error) {
	return ServiceProbeResult{}, nil
}
