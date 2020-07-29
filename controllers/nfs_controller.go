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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	nfsstoragev1alpha1 "github.com/johandry/nfs-operator/api/v1alpha1"
)

// NfsReconciler reconciles a Nfs object
type NfsReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=nfs.storage.ibmcloud.ibm.com,resources=nfs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nfs.storage.ibmcloud.ibm.com,resources=nfs/status,verbs=get;update;patch

func (r *NfsReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("nfs", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *NfsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nfsstoragev1alpha1.Nfs{}).
		Complete(r)
}
