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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RequestSpec defines the specifications of the backing storage to request or create
type RequestSpec struct {
	// +kubebuilder:validation:MinLength=0

	// Storage size to request
	Storage string `json:"storage,omitempty"`
}

// BackingStorageSpec defines the desired state of the Backing Storage used to
// export the NFS volume
type BackingStorageSpec struct {
	// Defines if the backing storage is created (false) or will use an existing
	// (true) PVC previously created by the user. Default is false
	// +optional
	UseExistingPVC bool `json:"useExistingPVC,omitempty"`

	// +kubebuilder:validation:MinLength=0

	// Name of the PVC to create or use, if exists
	Name string `json:"name,omitempty"`

	// Storage class used to create the PVC. Defaults to ibmc-vpc-block-general-purpose
	// +optional
	StorageClassName string `json:"storageClassName,omitempty"`

	// Specifications to request for the PVC to create
	Request RequestSpec `json:"request,omitempty"`
}

// NfsSpec defines the desired state of Nfs
type NfsSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The storageClass that the provisioner will listen for requests. Defaults to cluster-nfs
	// +optional
	StorageClassName string `json:"storageClassName,omitempty"`

	// Defines the provisioner API. Defaults to cluster.example.com/nfs
	// +optional
	ProvisionerAPI string `json:"provisionerAPI,omitempty"`

	// Specifications for the backing storage used to export the NFS volume
	BackingStorage BackingStorageSpec `json:"backingStorage,omitempty"`
}

// NfsStatus defines the observed state of Nfs
type NfsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Capacity   string `json:"capacity,omitempty"`
	AccessMode string `json:"accessMode,omitempty"`
	Status     string `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Nfs is the Schema for the nfs API
type Nfs struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NfsSpec   `json:"spec,omitempty"`
	Status NfsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NfsList contains a list of Nfs
type NfsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Nfs `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Nfs{}, &NfsList{})
}
