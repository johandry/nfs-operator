package nfs

import (
	"context"
	"fmt"

	nfsstoragev1alpha1 "github.com/johandry/nfs-operator/api/v1alpha1"
	"github.com/johandry/nfs-operator/resources"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	rbacv1 "k8s.io/api/rbac/v1"
)

var _ resources.Reconcilable = &RoleBinding{}

// RoleBinding is the RoleBinding resource used by the Nfs controller
type RoleBinding struct {
	Owner *nfsstoragev1alpha1.Nfs
}

var yamlRoleBinding = []byte(`
kind: RoleBindingBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: leader-locking-nfs-provisioner
subjects:
  - kind: ServiceAccount
    name: nfs-provisioner
    # replace with namespace where provisioner is deployed
    namespace: default
roleRef:
  kind: RoleBinding
  name: leader-locking-nfs-provisioner
  apiGroup: rbac.authorization.k8s.io
`)

// new returns the object as a rbac.v1.RoleBinding
func (r *RoleBinding) new() *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "leader-locking-" + appName,
			Namespace: r.Owner.Namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      appName,
				Namespace: r.Owner.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     "leader-locking-" + appName,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
}

// toResource returns the given object as a rbac.v1.RoleBinding
func (r *RoleBinding) toResource(ro runtime.Object) (*rbacv1.RoleBinding, error) {
	if v, ok := ro.(*rbacv1.RoleBinding); ok {
		return v, nil
	}
	return nil, fmt.Errorf("the received object is not a rbac/v1.RoleBinding")
}

// isValid returns true if the given object is valid. If it's valid won't be updated
func (r *RoleBinding) isValid(o *rbacv1.RoleBinding) bool {
	return true
}

// Object implements the Object method of the Reconcilable interface
func (r *RoleBinding) Object() runtime.Object {
	return r.new()
}

// SetControllerReference implements the SetControllerReference method of the Reconcilable interface
func (r *RoleBinding) SetControllerReference(scheme *runtime.Scheme) error {
	obj := r.new()

	var err error
	if scheme != nil {
		err = ctrl.SetControllerReference(r.Owner, obj, scheme)
	}

	return err
}

// Get implements the Get method of the Reconcilable interface
func (r *RoleBinding) Get(ctx context.Context, c client.Client) (runtime.Object, error) {
	found := &rbacv1.RoleBinding{}
	obj := r.new()

	err := c.Get(ctx, types.NamespacedName{Name: obj.Name, Namespace: obj.Namespace}, found)
	if err == nil {
		return found, nil
	}
	return nil, client.IgnoreNotFound(err)
}

// Validate implements the Validate method of the Reconcilable interface
func (r *RoleBinding) Validate(ro runtime.Object) bool {
	current, err := r.toResource(ro)
	if err != nil {
		return false
	}

	return r.isValid(current)
}
