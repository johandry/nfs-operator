/*
Copyright 2020 NFS Operator authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	nfsstoragev1alpha1 "github.com/johandry/nfs-operator/api/v1alpha1"
	"github.com/johandry/nfs-operator/resources"
	vpcblockbackend "github.com/johandry/nfs-operator/resources/backends/vpc-block"
	nfsprovisioner "github.com/johandry/nfs-operator/resources/provisioners/nfs"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// NfsReconciler reconciles a Nfs object
type NfsReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=nfs.storage.ibmcloud.ibm.com,resources=nfs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nfs.storage.ibmcloud.ibm.com,resources=nfs/status,verbs=get;update;patch

// Reconcile ...
func (r *NfsReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("nfs", req.NamespacedName)

	log.Info("reconciling NFS")

	nfs := &nfsstoragev1alpha1.Nfs{}
	if err := r.Get(ctx, req.NamespacedName, nfs); err != nil {
		log.Info("NFS not found")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info(fmt.Sprintf("Found NFS %s/%s", nfs.Namespace, nfs.Name))

	resourcesToReconcile := []resources.Reconcilable{}
	resourcesToReconcile = append(resourcesToReconcile, vpcblockbackend.Resources(nfs)...)
	resourcesToReconcile = append(resourcesToReconcile, nfsprovisioner.Resources(nfs)...)

	for _, res := range resourcesToReconcile {
		if err := r.apply(ctx, res); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager ...
func (r *NfsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nfsstoragev1alpha1.Nfs{}).
		// Backend Storage: VPC Block
		// Owns(&corev1.PersistentVolumeClaim{}).
		// NFS Provisioner: Deployment
		Owns(&corev1.Service{}).
		Owns(&appsv1.Deployment{}).
		// NFS Provisioner: RBAC
		// Owns(&corev1.ServiceAccount{}).
		// Owns(&rbacv1.Role{}).
		// Owns(&rbacv1.RoleBinding{}).
		// NFS Provisioner: StorageClass (non namespaced)
		// Owns(&storagev1.StorageClass{}).
		Complete(r)
}

func (r *NfsReconciler) apply(ctx context.Context, res resources.Reconcilable) error {
	current, err := res.Get(ctx, r.Client)
	if err != nil {
		r.Log.Error(err, "Fail to get the resource")
		return err
	}

	obj := res.Object()

	if err := res.SetControllerReference(r.Scheme); err != nil {
		return err
	}

	if current != nil {
		if ok := res.Validate(current); !ok {
			r.Log.Info("Resource already exists but invalid, updating the resource to original state")
			return r.Update(ctx, obj)
		}
		r.Log.Info("Skip reconcile: Resource already exists")
		return nil
	}

	// create it if not found and no error
	r.Log.Info("Resource created")
	return r.Create(ctx, obj)
}
