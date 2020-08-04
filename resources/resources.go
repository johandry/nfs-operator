package resources

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Reconcilable is a resource that can be reconciled by the controller
type Reconcilable interface {
	// Object return the resource as it should be in the cluster
	Object() runtime.Object
	// SetControllerReference set the controller reference to the owner
	SetControllerReference(*runtime.Scheme) error
	// Get retreive the resource from the cluster, if not found return (nil, nil)
	// if fail to retreive return (nil, err)
	Get(context.Context, client.Client) (runtime.Object, error)
	// Validate return true if the given object is different to the object as it should
	// be in the cluster. Return 'false' to always leave the resource as it is in
	// the cluster
	Validate(runtime.Object) bool
}
