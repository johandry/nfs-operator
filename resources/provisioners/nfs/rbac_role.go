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

var _ resources.Reconcilable = &Role{}

// Role is the Role resource used by the Nfs controller
type Role struct {
	Owner *nfsstoragev1alpha1.Nfs
}

var yamlRole = []byte(`
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: leader-locking-nfs-provisioner
rules:
  - apiGroups: [""]
    resources: ["endpoints"]
    verbs: ["get", "list", "watch", "create", "update", "patch"]
`)

// new returns the object as a rbac.v1.Role
func (r *Role) new() *rbacv1.Role {
	return &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "leader-locking-" + appName,
			Namespace: r.Owner.Namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"endpoints"},
				Verbs:     []string{"get", "list", "watch", "create", "update", "patch"},
			},
		},
	}
}

// toResource returns the given object as a rbac.v1.Role
func (r *Role) toResource(ro runtime.Object) (*rbacv1.Role, error) {
	if v, ok := ro.(*rbacv1.Role); ok {
		return v, nil
	}
	return nil, fmt.Errorf("the received object is not a rbac/v1.Role")
}

// isValid returns true if the given object is valid. If it's valid won't be updated
func (r *Role) isValid(o *rbacv1.Role) bool {
	return true
}

// Object implements the Object method of the Reconcilable interface
func (r *Role) Object() runtime.Object {
	return r.new()
}

// SetControllerReference implements the SetControllerReference method of the Reconcilable interface
func (r *Role) SetControllerReference(scheme *runtime.Scheme) error {
	obj := r.new()

	var err error
	if scheme != nil {
		err = ctrl.SetControllerReference(r.Owner, obj, scheme)
	}

	return err
}

// Get implements the Get method of the Reconcilable interface
func (r *Role) Get(ctx context.Context, c client.Client) (runtime.Object, error) {
	found := &rbacv1.Role{}
	obj := r.new()

	err := c.Get(ctx, types.NamespacedName{Name: obj.Name, Namespace: obj.Namespace}, found)
	if err == nil {
		return found, nil
	}
	return nil, client.IgnoreNotFound(err)
}

// Validate implements the Validate method of the Reconcilable interface
func (r *Role) Validate(ro runtime.Object) bool {
	current, err := r.toResource(ro)
	if err != nil {
		return false
	}

	return r.isValid(current)
}
