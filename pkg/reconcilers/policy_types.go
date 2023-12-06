package reconcilers

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PolicyStatus interface {
	client.Object

	GetStatusConditions() *[]metav1.Condition
}
