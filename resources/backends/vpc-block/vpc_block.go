package vpcblock

import (
	nfsstoragev1alpha1 "github.com/johandry/nfs-operator/api/v1alpha1"
	"github.com/johandry/nfs-operator/resources"
)

const (
	persistentVolumeClaimName = "nfs-block-custom"
)

// Resources return the list of resources to make a NFS Provisioner
func Resources(owner *nfsstoragev1alpha1.Nfs) []resources.Reconcilable {
	return []resources.Reconcilable{
		&PersistentVolumeClaim{Owner: owner},
	}
}
