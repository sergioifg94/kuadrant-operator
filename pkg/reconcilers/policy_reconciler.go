package reconcilers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type PolicyReconciler[T PolicyStatus] struct {
	*TargetRefReconciler

	PolicyTemplate func() T

	Validate        func(ctx context.Context, policy T) (ValidationResult, error)
	ReconcilePolicy func(ctx context.Context, policy T) (ReconciliationResult, error)
	ServiceProbe    func(ctx context.Context, policy T) (ServiceProbeResult, error)
}

var _ reconcile.Reconciler = &PolicyReconciler[PolicyStatus]{}

type PolicyStageFn[T PolicyStatus] func(context.Context, client.Client, T) error

type ValidationResult struct {
	Success bool
	Error   *ValidationError
}

type ValidationError struct {
	Type    ValidationErrorType
	Message string
}

type ValidationErrorType string

type ReconciliationResult struct {
}

type ServiceProbeResult struct {
}

func (r *PolicyReconciler[T]) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	obj := r.PolicyTemplate()

	if err := r.client.Get(ctx, req.NamespacedName, obj); err != nil {
		return ctrl.Result{}, err
	}

	conditions := obj.GetStatusConditions()

	validationResult, err := r.Validate(ctx, obj)
	if err != nil {
		return ctrl.Result{}, err
	}

	// If the validation fails, set the Accepted status condition to invalid
	// and finish reconciliation
	if !validationResult.Success {
		condition := metav1.Condition{
			Type:    "Accepted",
			Status:  metav1.ConditionFalse,
			Reason:  "Invalid",
			Message: validationResult.Error.Message,
		}

		meta.SetStatusCondition(conditions, condition)

		return ctrl.Result{}, r.client.Status().Update(ctx, obj)
	}

	_, err = r.ReconcilePolicy(ctx, obj)
	if err != nil {
		return ctrl.Result{}, err
	}

	meta.SetStatusCondition(conditions, metav1.Condition{
		Type:    "Accepted",
		Status:  metav1.ConditionTrue,
		Reason:  "Accepted",
		Message: "Policy has been accepted",
	})

	_, err = r.ServiceProbe(ctx, obj)

	return reconcile.Result{}, nil
}
