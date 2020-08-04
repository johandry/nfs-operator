package nfs

import (
	"context"
	"fmt"

	nfsstoragev1alpha1 "github.com/johandry/nfs-operator/api/v1alpha1"
	"github.com/johandry/nfs-operator/resources"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ resources.Reconcilable = &Deployment{}

// Deployment is the Deployment resource used by the Nfs controller
type Deployment struct {
	Owner *nfsstoragev1alpha1.Nfs
}

var yamlDeployment = []byte(`
kind: Deployment
apiVersion: apps/v1
metadata:
  name: nfs-provisioner
spec:
  selector:
    matchLabels:
      app: nfs-provisioner
  replicas: 1
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: nfs-provisioner
    spec:
      serviceAccount: nfs-provisioner
      containers:
        - name: nfs-provisioner
          image: quay.io/kubernetes_incubator/nfs-provisioner:latest
          ports:
            - name: nfs
              containerPort: 2049
            - name: nfs-udp
              containerPort: 2049
              protocol: UDP
            - name: nlockmgr
              containerPort: 32803
            - name: nlockmgr-udp
              containerPort: 32803
              protocol: UDP
            - name: mountd
              containerPort: 20048
            - name: mountd-udp
              containerPort: 20048
              protocol: UDP
            - name: rquotad
              containerPort: 875
            - name: rquotad-udp
              containerPort: 875
              protocol: UDP
            - name: rpcbind
              containerPort: 111
            - name: rpcbind-udp
              containerPort: 111
              protocol: UDP
            - name: statd
              containerPort: 662
            - name: statd-udp
              containerPort: 662
              protocol: UDP
          securityContext:
            capabilities:
              add:
                - DAC_READ_SEARCH
                - SYS_RESOURCE
          args:
            - "-provisioner=ibmcloud/nfs"
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: SERVICE_NAME
              value: nfs-provisioner
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
          imagePullPolicy: "IfNotPresent"
          volumeMounts:
            - name: export-volume
              mountPath: /export
      volumes:
        - name: export-volume
          persistentVolumeClaim:
            claimName: nfs-block-custom
`)

// new returns the object as a apps.v1.Deployment
func (r *Deployment) new() *appsv1.Deployment {
	replicas := int32(1)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appName,
			Namespace: r.Owner.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": appName,
				},
			},
			Replicas: &replicas,
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": appName,
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: appName,
					Containers: []corev1.Container{
						{
							Name:  appName,
							Image: imageName,
							Ports: []corev1.ContainerPort{
								{
									Name:          "nfs",
									ContainerPort: 2049,
								},
								{
									Name:          "nfs-udp",
									ContainerPort: 2049,
									Protocol:      corev1.ProtocolUDP,
								},
								{
									Name:          "nlockmgr",
									ContainerPort: 32803,
								},
								{
									Name:          "nlockmgr-udp",
									ContainerPort: 32803,
									Protocol:      corev1.ProtocolUDP,
								},
								{
									Name:          "mountd",
									ContainerPort: 20048,
								},
								{
									Name:          "mountd-udp",
									ContainerPort: 20048,
									Protocol:      corev1.ProtocolUDP,
								},
								{
									Name:          "rquotad",
									ContainerPort: 875,
								},
								{
									Name:          "rquotad-udp",
									ContainerPort: 875,
									Protocol:      corev1.ProtocolUDP,
								},
								{
									Name:          "rpcbind",
									ContainerPort: 111,
								},
								{
									Name:          "rpcbind-udp",
									ContainerPort: 111,
									Protocol:      corev1.ProtocolUDP,
								},
								{
									Name:          "statd",
									ContainerPort: 662,
								},
								{
									Name:          "statd-udp",
									ContainerPort: 662,
									Protocol:      corev1.ProtocolUDP,
								},
							},
							SecurityContext: &corev1.SecurityContext{
								Capabilities: &corev1.Capabilities{
									Add: []corev1.Capability{
										"DAC_READ_SEARCH",
										"SYS_RESOURCE",
									},
								},
							},
							Args: []string{
								fmt.Sprintf("-provisioner=%s", provisionerName),
							},
							Env: []corev1.EnvVar{
								{
									Name: "POD_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
								{
									Name:  "SERVICE_NAME",
									Value: appName,
								},
								{
									Name: "POD_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
							ImagePullPolicy: corev1.PullIfNotPresent,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "export-volume",
									MountPath: "/export",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "export-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: r.Owner.Spec.BackingStorage.Name,
								},
							},
						},
					},
				},
			},
		},
	}
}

// toResource returns the given object as a apps.v1.Deployment
func (r *Deployment) toResource(ro runtime.Object) (*appsv1.Deployment, error) {
	if v, ok := ro.(*appsv1.Deployment); ok {
		return v, nil
	}
	return nil, fmt.Errorf("the received object is not a apps/v1.Deployment")
}

// isValid returns true if the given object is valid. If it's valid won't be updated
func (r *Deployment) isValid(o *appsv1.Deployment) bool {
	return true
}

// Object implements the Object method of the Reconcilable interface
func (r *Deployment) Object() runtime.Object {
	return r.new()
}

// SetControllerReference implements the SetControllerReference method of the Reconcilable interface
func (r *Deployment) SetControllerReference(scheme *runtime.Scheme) error {
	obj := r.new()

	var err error
	if scheme != nil {
		err = ctrl.SetControllerReference(r.Owner, obj, scheme)
	}

	return err
}

// Get implements the Get method of the Reconcilable interface
func (r *Deployment) Get(ctx context.Context, c client.Client) (runtime.Object, error) {
	found := &appsv1.Deployment{}
	obj := r.new()

	err := c.Get(ctx, types.NamespacedName{Name: obj.Name /*, Namespace: obj.Namespace */}, found)
	if err == nil {
		return found, nil
	}
	return nil, client.IgnoreNotFound(err)
}

// Validate implements the Validate method of the Reconcilable interface
func (r *Deployment) Validate(ro runtime.Object) bool {
	current, err := r.toResource(ro)
	if err != nil {
		return false
	}

	return r.isValid(current)
}
