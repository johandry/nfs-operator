package vpcblock

import (
	"context"
	"fmt"

	nfsstoragev1alpha1 "github.com/johandry/nfs-operator/api/v1alpha1"
	"github.com/johandry/nfs-operator/resources"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ resources.Reconcilable = &PersistentVolumeClaim{}

// PersistentVolumeClaim is the PersistentVolumeClaim resource used by the Nfs controller
type PersistentVolumeClaim struct {
	Owner *nfsstoragev1alpha1.Nfs
}

var yamlPersistentVolumeClaim = []byte(`
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: nfs-block-custom
spec:
  storageClassName: ibmc-vpc-block-general-purpose
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
`)

// new returns the object as a core.v1.PersistentVolumeClaim
func (r *PersistentVolumeClaim) new() *corev1.PersistentVolumeClaim {
	storageClassNameStr := r.Owner.Spec.BackingStorage.StorageClassName
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      persistentVolumeClaimName,
			Namespace: r.Owner.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassNameStr,
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(r.Owner.Spec.BackingStorage.Request.Storage),
				},
			},
		},
	}
}

// toResource returns the given object as a core.v1.PersistentVolumeClaim
func (r *PersistentVolumeClaim) toResource(ro runtime.Object) (*corev1.PersistentVolumeClaim, error) {
	if v, ok := ro.(*corev1.PersistentVolumeClaim); ok {
		return v, nil
	}
	return nil, fmt.Errorf("the received object is not a core/v1.PersistentVolumeClaim")
}

// isValid returns true if the given object is valid. If it's valid won't be updated
func (r *PersistentVolumeClaim) isValid(o *corev1.PersistentVolumeClaim) bool {
	return true
}

// Object implements the Object method of the Reconcilable interface
func (r *PersistentVolumeClaim) Object() runtime.Object {
	return r.new()
}

// SetControllerReference implements the SetControllerReference method of the Reconcilable interface
func (r *PersistentVolumeClaim) SetControllerReference(scheme *runtime.Scheme) error {
	obj := r.new()

	var err error
	if scheme != nil {
		err = ctrl.SetControllerReference(r.Owner, obj, scheme)
	}

	return err
}

// Get implements the Get method of the Reconcilable interface
func (r *PersistentVolumeClaim) Get(ctx context.Context, c client.Client) (runtime.Object, error) {
	found := &corev1.PersistentVolumeClaim{}
	obj := r.new()

	err := c.Get(ctx, types.NamespacedName{Name: obj.Name /*, Namespace: obj.Namespace */}, found)
	if err == nil {
		return found, nil
	}
	return nil, client.IgnoreNotFound(err)
}

// Validate implements the Validate method of the Reconcilable interface
func (r *PersistentVolumeClaim) Validate(ro runtime.Object) bool {
	current, err := r.toResource(ro)
	if err != nil {
		return false
	}

	return r.isValid(current)
}
