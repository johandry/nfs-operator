package nfs

import (
	"context"
	"fmt"

	nfsstoragev1alpha1 "github.com/johandry/nfs-operator/api/v1alpha1"
	"github.com/johandry/nfs-operator/resources"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ resources.Reconcilable = &StorageClass{}

// StorageClass is the StorageClass resource used by the Nfs controller
type StorageClass struct {
	Owner *nfsstoragev1alpha1.Nfs
}

var yamlStorageClass = []byte(`
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: ibmcloud-nfs
provisioner: ibmcloud/nfs
mountOptions:
	- vers=4.1
`)

// new returns the object as a storage.v1.StorageClass
func (r *StorageClass) new() *storagev1.StorageClass {
	return &storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: storageClassName,
			// Namespace: r.Owner.Namespace,
		},
		Provisioner: provisionerName,
		MountOptions: []string{
			"vers=4.1",
		},
	}
}

// toResource returns the given object as a storage.v1.StorageClass
func (r *StorageClass) toResource(ro runtime.Object) (*storagev1.StorageClass, error) {
	if v, ok := ro.(*storagev1.StorageClass); ok {
		return v, nil
	}
	return nil, fmt.Errorf("the received object is not a storage/v1.StorageClass")
}

// isValid returns true if the given object is valid. If it's valid won't be updated
func (r *StorageClass) isValid(o *storagev1.StorageClass) bool {
	obj := r.new()

	if o.Provisioner != obj.Provisioner {
		return true
	}

	return true
}

// Object implements the Object method of the Reconcilable interface
func (r *StorageClass) Object() runtime.Object {
	return r.new()
}

// SetControllerReference implements the SetControllerReference method of the Reconcilable interface
func (r *StorageClass) SetControllerReference(scheme *runtime.Scheme) error {
	obj := r.new()

	var err error
	if scheme != nil {
		err = ctrl.SetControllerReference(r.Owner, obj, scheme)
	}

	return err
}

// Get implements the Get method of the Reconcilable interface
func (r *StorageClass) Get(ctx context.Context, c client.Client) (runtime.Object, error) {
	found := &storagev1.StorageClass{}
	obj := r.new()

	err := c.Get(ctx, types.NamespacedName{Name: obj.Name /*, Namespace: obj.Namespace */}, found)
	if err == nil {
		return found, nil
	}
	return nil, client.IgnoreNotFound(err)
}

// Validate implements the Validate method of the Reconcilable interface
func (r *StorageClass) Validate(ro runtime.Object) bool {
	current, err := r.toResource(ro)
	if err != nil {
		return false
	}

	return r.isValid(current)
}
