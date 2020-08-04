package nfs

import (
	nfsstoragev1alpha1 "github.com/johandry/nfs-operator/api/v1alpha1"
	"github.com/johandry/nfs-operator/resources"
)

const (
	appName                   = "nfs-provisioner"
	imageName                 = "quay.io/kubernetes_incubator/nfs-provisioner:latest"
	storageClassName          = "ibmcloud-nfs"
	provisionerName           = "ibmcloud/nfs"
	persistentVolumeClaimName = "nfs"
)

// Resources return the list of resources to make a NFS Provisioner
func Resources(owner *nfsstoragev1alpha1.Nfs) []resources.Reconcilable {
	return []resources.Reconcilable{
		// Deployment
		&Service{Owner: owner},
		&Deployment{Owner: owner},
		// RBAC
		&ServiceAccount{Owner: owner},
		&Role{Owner: owner},
		&RoleBinding{Owner: owner},
		// StorageClass
		&StorageClass{Owner: owner},
	}
}
