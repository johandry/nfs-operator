# NFS Provisioner Operator

- [NFS Provisioner Operator](#nfs-provisioner-operator)
  - [Usage](#usage)
  - [Build](#build)
  - [Testing](#testing)
  - [Cleanup](#cleanup)

The NFS Provisioner Operator creates a **NFS External Provisioner** which creates a `ReadWriteMany` `PersistentVolumeClaim` to be consumed by any Pod/Container in the cluster. The backend block storage to be exposed by NFS, could be previously created and specified, or could be created by the operator. At this time, this backend block storage is an IBM Cloud VPC Block.

The goal of this NFS Provisioner Operator is to make it easier to Kubernetes developers to have a PVC that can be used by many pods (`ReadWriteMany`) using the same volume, saving resources and money.

Refer to the [documentation](./docs/index.md) for more information about the design and architecture of the NFS Provisioner Operator.

## Usage

Before use it you need to deploy the NFS Provisioner Operator, this is usually done, but not necessarily, when the cluster is created. The deployment can be done with the following `kubectl` command:

```bash
kubectl create -f https://www.johandry.com/nfs-operator/nfs_provisioner_operator_install.yaml
```

The first step after the NFS Provisioner Operator is deployed is to create a **NFS CustomResource** defining the `storageClassName` and the backing block storage. The backend block storage is created by the operator or you can provide an existing storage accessible through a PVC.

An example of a regular NFS CustomResource could be like this.

```yaml
apiVersion: nfs.storage.ibmcloud.ibm.com/v1alpha1
kind: Nfs
metadata:
  name: cluster-nfs
spec:
  storageClassName: cluster-nfs
  provisionerAPI: example.com/nfs
  backingStorage:
    name: export-nfs-block
    storageClassName: ibmc-vpc-block-general-purpose
    request:
      storage: 10Gi
```

Notice the value of `spec.storageClassName` and the values of the `backingStorage` specification. The backend block storage will be of `backingStorage.storageClassName` name `ibmc-vpc-block-general-purpose` with **10Gb**.

If you have your own block storage to be used by the NFS Provisioner, read the [documentation](./docs/index.md).

To use the storage, create a PVC using the given storage class, in this example it is `example-nfs`. The VPC for this example would be like this.

```yaml
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: nfs
spec:
  storageClassName: example-nfs
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 1Mi
```

This VPC request 1Mb and the name is `nfs`, as its access mode is `ReadWriteMany` many containers or Pods can use it.

The following Pod example uses the NFS Provider creating a volume from the PVC using the claim name `nfs`. Then mount the volume in any directory of the container with `volumeMounts`.

This is a simple example of a container that is creating a file in the NFS volume:

```yaml
kind: Pod
apiVersion: v1
metadata:
  name: consumer
spec:
  containers:
    - name: consumer
      image: busybox
      command:
        - "/bin/sh"
      args:
        - "-c"
        - "touch /mnt/SUCCESS && exit 0 || exit 1"
      volumeMounts:
        - name: nfs-pvc
          mountPath: "/mnt"
  restartPolicy: "Never"
  volumes:
    - name: nfs-pvc
      persistentVolumeClaim:
        claimName: nfs
```

A demo application to use the NFS service can be found in the `kubernetes/consumer` folder and it's a simple API for movies. The database - a single JSON file - is stored in the shared volume. The deployment uses a initContainer to move the JSON database/file to the shared volume.

More information can be found in the [documentation](./docs/index.md).

## Build

The development of the operator focus on basically the files:

- `api/v1alpha1/nfs_types.go`: defines the operator specs and status, modifying this file requires to execute `make`
- `controllers/nfs_controller.go`: contains the `Reconcile` function to create or delete all the required resources.
- `resources/`: directory and packages with all the logic to create NFS Provisioner and the Backing Storage (PVC)

After modify any of the files it's recommended to execute `make` to generate the CR and CRD's, and to build the Docker container with the NFS Operator and finally push it to the Docker Registry.

To quick test the operator (build, deploy and test locally), execute:

```bash
make
make run
make deploy
```

To test using the consumer application on the cluster, execute:

```bash
cd test
make consumer
make test
```

Refer to the [DEVELOPMENT](./docs/DEVELOPMENT.md) and [TESTING](./docs/TESTING.md) documents for more information. Optionally, read the `Makefile`'s to be familiar with all the tasks you can execute for testing.

## Testing

The tests require a Kubernetes cluster on IBM Cloud, to get one follow the instructions from the testing [TESTING](./docs/TESTING.md) document or follow the following quick start instructions.

```bash
make environment

# Optionally, create the PVC
make pvc

# Edit the Custom Resource or, at least, confirm the specifications
vim config/samples/nfs.storage.ibmcloud.ibm.com_v1alpha1_nfs.yaml
vim kustomization.yaml

make deploy

# Test the Operator locally
make run

# Or test it with the consumer application
cd test
make consumer
make test
```

To know the status of the resources created by the NFS Operator or the NFS Provisioner, execute `make list` and to know all the resources in the cluster, either created by the code or external, execute `make list-all`.

## Cleanup

To delete the CR and CRD execute:

```bash
make clean
```

To delete the resources for testing execute:

```bash
cd test
make delete
```

To wipe out everything, including the IKS cluster execute, form the `test` directory, `make clean` or to get the cluster to the original state (not recommended) execute `make purge`.

Refer to the `test/Makefile` file to be familiar with all the tasks you can execute for testing and the [TESTING](./docs/TESTING.md) document for more information.
