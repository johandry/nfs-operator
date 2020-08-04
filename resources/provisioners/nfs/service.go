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

var _ resources.Reconcilable = &Service{}

// Service is the Service resource used by the Nfs controller
type Service struct {
	Owner *nfsstoragev1alpha1.Nfs
}

var yamlService = []byte(`
kind: Service
apiVersion: v1
metadata:
  name: nfs-provisioner
  labels:
    app: nfs-provisioner
spec:
  ports:
    - name: nfs
      port: 2049
    - name: nfs-udp
      port: 2049
      protocol: UDP
    - name: nlockmgr
      port: 32803
    - name: nlockmgr-udp
      port: 32803
      protocol: UDP
    - name: mountd
      port: 20048
    - name: mountd-udp
      port: 20048
      protocol: UDP
    - name: rquotad
      port: 875
    - name: rquotad-udp
      port: 875
      protocol: UDP
    - name: rpcbind
      port: 111
    - name: rpcbind-udp
      port: 111
      protocol: UDP
    - name: statd
      port: 662
    - name: statd-udp
      port: 662
      protocol: UDP
  selector:
    app: nfs-provisioner
`)

// new returns the object as a core.v1.Service
func (r *Service) new() *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appName,
			Namespace: r.Owner.Namespace,
			Labels: map[string]string{
				"app": appName,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: "nfs",
					Port: int32(2049),
				},
				{
					Name:     "nfs-udp",
					Port:     int32(2049),
					Protocol: corev1.ProtocolUDP,
				},
				{
					Name: "nlockmgr",
					Port: int32(32803),
				},
				{
					Name:     "nlockmgr-udp",
					Port:     int32(32803),
					Protocol: corev1.ProtocolUDP,
				},
				{
					Name: "mountd",
					Port: int32(20048),
				},
				{
					Name:     "mountd-udp",
					Port:     int32(20048),
					Protocol: corev1.ProtocolUDP,
				},
				{
					Name: "rquotad",
					Port: int32(875),
				},
				{
					Name:     "rquotad-udp",
					Port:     int32(875),
					Protocol: corev1.ProtocolUDP,
				},
				{
					Name: "rpcbind",
					Port: int32(111),
				},
				{
					Name:     "rpcbind-udp",
					Port:     int32(111),
					Protocol: corev1.ProtocolUDP,
				},
				{
					Name: "statd",
					Port: int32(662),
				},
				{
					Name:     "statd-udp",
					Port:     int32(662),
					Protocol: corev1.ProtocolUDP,
				},
			},
			Selector: map[string]string{
				"app": appName,
			},
		},
	}
}

// toResource returns the given object as a core.v1.Service
func (r *Service) toResource(ro runtime.Object) (*corev1.Service, error) {
	if v, ok := ro.(*corev1.Service); ok {
		return v, nil
	}
	return nil, fmt.Errorf("the received object is not a core/v1.Service")
}

// isValid returns true if the given object is valid. If it's valid won't be updated
func (r *Service) isValid(o *corev1.Service) bool {
	return true
}

// Object implements the Object method of the Reconcilable interface
func (r *Service) Object() runtime.Object {
	return r.new()
}

// SetControllerReference implements the SetControllerReference method of the Reconcilable interface
func (r *Service) SetControllerReference(scheme *runtime.Scheme) error {
	obj := r.new()

	var err error
	if scheme != nil {
		err = ctrl.SetControllerReference(r.Owner, obj, scheme)
	}

	return err
}

// Get implements the Get method of the Reconcilable interface
func (r *Service) Get(ctx context.Context, c client.Client) (runtime.Object, error) {
	found := &corev1.Service{}
	obj := r.new()

	err := c.Get(ctx, types.NamespacedName{Name: obj.Name /*, Namespace: obj.Namespace */}, found)
	if err == nil {
		return found, nil
	}
	return nil, client.IgnoreNotFound(err)
}

// Validate implements the Validate method of the Reconcilable interface
func (r *Service) Validate(ro runtime.Object) bool {
	current, err := r.toResource(ro)
	if err != nil {
		return false
	}

	return r.isValid(current)
}
