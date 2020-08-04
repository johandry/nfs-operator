package nfs

import (
	"context"
	"fmt"

	nfsstoragev1alpha1 "github.com/johandry/nfs-operator/api/v1alpha1"
	"github.com/johandry/nfs-operator/resources"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ resources.Reconcilable = &ServiceAccount{}

// ServiceAccount is the ServiceAccount resource used by the Nfs controller
type ServiceAccount struct {
	Owner *nfsstoragev1alpha1.Nfs
}

var yamlServiceAccount = []byte(`
apiVersion: v1
kind: ServiceAccount
metadata:
  name: nfs-provisioner
`)

// new returns the object as a core.v1.ServiceAccount
func (r *ServiceAccount) new() *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appName,
			Namespace: r.Owner.Namespace,
		},
	}
}

// toResource returns the given object as a core.v1.ServiceAccount
func (r *ServiceAccount) toResource(ro runtime.Object) (*corev1.ServiceAccount, error) {
	if v, ok := ro.(*corev1.ServiceAccount); ok {
		return v, nil
	}
	return nil, fmt.Errorf("the received object is not a core/v1.ServiceAccount")
}

// isValid returns true if the given object is valid. If it's valid won't be updated
func (r *ServiceAccount) isValid(o *corev1.ServiceAccount) bool {
	return true
}

// Object implements the Object method of the Reconcilable interface
func (r *ServiceAccount) Object() runtime.Object {
	return r.new()
}

// SetControllerReference implements the SetControllerReference method of the Reconcilable interface
func (r *ServiceAccount) SetControllerReference(scheme *runtime.Scheme) error {
	obj := r.new()

	var err error
	if scheme != nil {
		err = ctrl.SetControllerReference(r.Owner, obj, scheme)
	}

	return err
}

// Get implements the Get method of the Reconcilable interface
func (r *ServiceAccount) Get(ctx context.Context, c client.Client) (runtime.Object, error) {
	found := &corev1.ServiceAccount{}
	obj := r.new()

	err := c.Get(ctx, types.NamespacedName{Name: obj.Name /*, Namespace: obj.Namespace */}, found)
	if err == nil {
		return found, nil
	}
	return nil, client.IgnoreNotFound(err)
}

// Validate implements the Validate method of the Reconcilable interface
func (r *ServiceAccount) Validate(ro runtime.Object) bool {
	current, err := r.toResource(ro)
	if err != nil {
		return false
	}

	return r.isValid(current)
}
