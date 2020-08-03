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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var nfslog = logf.Log.WithName("nfs-resource")

// SetupWebhookWithManager setup the WebHook with the Manager
func (r *Nfs) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-nfs-storage-ibmcloud-ibm-com-v1alpha1-nfs,mutating=true,failurePolicy=fail,groups=nfs.storage.ibmcloud.ibm.com,resources=nfs,verbs=create;update,versions=v1alpha1,name=mnfs.kb.io

var _ webhook.Defaulter = &Nfs{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Nfs) Default() {
	nfslog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
	if r.Spec.StorageClassName == "" {
		r.Spec.StorageClassName = "cluster-nfs"
	}
	if r.Spec.ProvisionerAPI == "" {
		r.Spec.ProvisionerAPI = "cluster.example.com/nfs"
	}

	// if r.Spec.BackingStorage == nil {
	// 	// TODO: If this is nil, add a default BackingStorageSpec and return
	// }
	if r.Spec.BackingStorage.StorageClassName == "" {
		r.Spec.BackingStorage.StorageClassName = "ibmc-vpc-block-general-purpose"
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-nfs-storage-ibmcloud-ibm-com-v1alpha1-nfs,mutating=false,failurePolicy=fail,groups=nfs.storage.ibmcloud.ibm.com,resources=nfs,versions=v1alpha1,name=vnfs.kb.io

var _ webhook.Validator = &Nfs{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Nfs) ValidateCreate() error {
	nfslog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Nfs) ValidateUpdate(old runtime.Object) error {
	nfslog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Nfs) ValidateDelete() error {
	nfslog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
