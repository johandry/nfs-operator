# Development Guide of the Operator using the Operator SDK

- [Development Guide of the Operator using the Operator SDK](#development-guide-of-the-operator-using-the-operator-sdk)
  - [QuickStart](#quickstart)
  - [Install the Operator SDK](#install-the-operator-sdk)
  - [Bootstrap the Operator](#bootstrap-the-operator)
  - [Create a Kubernetes cluster for testing](#create-a-kubernetes-cluster-for-testing)
  - [Testing on Kubernetes](#testing-on-kubernetes)
  - [Deployment and Release](#deployment-and-release)
  - [Cleanup](#cleanup)
  - [Reference to Advance Topics](#reference-to-advance-topics)

## QuickStart

1. [Create the test environment](#create-a-kubernetes-cluster-for-testing)
2. [Install the Operator SDK](#install-the-operator-sdk).
3. [Bootstrap the Operator](#bootstrap-the-operator) to create the API and the Controller.
4. Modify on the API file `api/v1alpha1/nfs_types.go` and the Controller file `controllers/nfs_controller.go`
5. Develop the controller modifing the files in the `resources` package
6. (Optional) Update the `VERSION` or the `IMG` variable in the `Makefile`
7. (Optional) Modify `config/samples/nfs.storage.ibmcloud.ibm.com_v1alpha1_nfs.yaml` or `config/samples/kustomization.yaml`.
8. To **build**, **install**, **test**, **release** and **deploy** execute the following commands:

   ```bash
   make
   make install

   # Test Locally: recomended to run this in a different terminal
   make run

   # Test on the IKS Testing Cluster:
   make release
   make deploy

   kubectl get deployments --namespace nfs-operator-system nfs-operator-controller-manager
   kubectl get pods --namespace nfs-operator-system
   kubectl get nfs -A
   ```

   The `make deploy` install the built CRD but the users can install the CRD executing:

   ```bash
   kubectl create -f https://www.johandry.com/nfs-operator/install.yaml
   ```

9. To **use** the NFS Provisioner executing the following commands. To **debug** the NFS Operator make sure `make run` is running in other terminal.

   ```bash
   # Create the NFS CR
   kustomize build config/samples | kubectl -f - apply
   ## Or:
   kubectl apply -f <(echo "
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
   ")
   kubectl get nfs -A

   # Create a PVC and the application (pod or deployment) that uses the NFS CR
   kubectl apply -f <(echo "
    kind: PersistentVolumeClaim
    apiVersion: v1
    metadata:
      name: nfs
    spec:
      storageClassName: cluster-nfs
      accessModes:
        - ReadWriteMany
      resources:
        requests:
          storage: 1Mi
   ")

   kubectl apply -f <(echo "
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
   ")

   kubectl get deployment
   kubectl get pods
   kubectl get nfs
   ```

10. Delete the installed Custom Resource with `make delete`, uninstall the CRD with `make uninstall`, or everything with `make clean`
11. Destroy the test environment with: `cd test; make clean`

## Install the Operator SDK

At this time the latest version of the OperatorSDK is `v0.19.0`. This new version is aligned with **Kubebuilder** which is a big change from previous versions.

```bash
brew install operator-sdk
```

To execute Unit Test the Kubernetes server binaries `etcd`, `kube-apiserver` and `kubectl` are required. Execute `make testbin` to install them into the `./testbin/` directory.

## Bootstrap the Operator

For this project we have the following requirements or parameters:

- Operator name: `nfs-operator`
- Repository: `github.com/johandry/nfs-operator`
- Kind: `Nfs`
- Group: `storage`
- API version: `v1alpha1`
- Domain: `ibmcloud.ibm.com`
- Group: `nfs.storage`

Therefore the following commands are executed to bootstrap the operator:

```bash
mkdir nfs-operator; cd nfs-operator
operator-sdk init --repo="github.com/johandry/nfs-operator" --owner="NFS Operator authors" --plugins="go.kubebuilder.io/v2" --domain="ibmcloud.ibm.com"

git init
git add .
git commit -m "project bootstrap"

operator-sdk create api --group="nfs.storage" --version="v1alpha1" --kind="Nfs" --controller --resource

operator-sdk create webhook --group nfs.storage --version v1alpha1 --kind Nfs --defaulting --programmatic-validation
```

Open the file `api/v1alpha1/nfs_types.go` to edit the struct `NfsSpec` and, optionally, the `NfsStatus`.

The NFS Operator specs are:

- `StorageClassName`: (`string`) the storage class that the provisioner will listen for requests. Default: `example-nfs`.
- `ProvisionerAPI`: (`string`) specify the provisioner API. Default: `example.com/nfs`.
- `BackingStorage`: structure to define the storage used to export the NFS service.
  - `UseExistingPVC`: (`bool`) use an existing claim. Default: `false`.
  - `Name`: (_required_, `string`) Name of the `PersistentVolumeClaim` to use or create.
  - `StorageClassName`: (`string`) cloud provider storage class to use for new PVC
  - `Request`: structure to define the storage to create
    - `Storage`: size of block volume to request from Cloud provider

An example of the NFS CustomResource would be like this:

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

Optionally edit the struct `NfsStatus` to add the parameters to report when you execute `kubectl get nfs`.

Example:

```go
type NfsSpec struct {
  // +optional
  // +kubebuilder:default=example-nfs
  StorageClassName string `json:"storageClassName,omitempty"`

  // +optional
  // +kubebuilder:default=example.com/nfs
  ProvisionerAPI string `json:"provisionerAPI,omitempty"`

  BackingStorage BackingStorageSpec `json:"backingStorage,omitempty"`
}

type NfsStatus struct {
  Capacity   string `json:"capacity,omitempty"`
  AccessMode string `json:"accessMode,omitempty"`
  Status     string `json:"status,omitempty"`
}
```

Execute `make` every time the `*_types.go` file is modified.

Edit the `api/v1alpha1/nfs_webhook.go` to include the validations and defaults values in the functions `Default`, `ValidateCreate` and `ValidateUpdate` functions.

To enable WebHook locally it's required to generate the certificates at `/tmp/k8s-webhook-server/serving-certs/tls.{crt,key}`, that's why WebHook is disabled by default when `make run` is executed. To enable it generate the certificates and execute `make run ENABLE_WEBHOOKS=true`

Open the `controllers/nfs_controller.go` file to edit the content of the `Reconcile` function.

Define the resources the operator creates in the package `resources`. The core of the resources

## Create a Kubernetes cluster for testing

The NFS Operator (at this time) is designed just for IBM Cloud and requires a Kubernetes cluster on IBM Cloud. To create an IKS cluster read the [instructions to setup the environment](./TESTING.md#requirements). Do not forget to:

1. Create the `test/terraform/.target_account` file with the IBM Cloud target account,
2. To export the `IC_API_KEY` environment variable with the API Key, and
3. Create the file `test/terraform/terraform.tfvars` with the variables `project_name` and `owner`.

Then go to the `test/` directory to execute:

```bash
make environment
```

If the tests include to have a pre-created PVC, execute `make pvc` in the `test/` directory.

Get more information about [building the environment](./TESTING.md#build-the-environment).

## Testing on Kubernetes

To test the NFS Custom Resource is required to deploy it. To create a NFS Custom Resource you can use the following example:

```bash
kubectl apply -f <(echo "
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
")
```

Or even better, modify the file `config/samples/nfs.storage.ibmcloud.ibm.com_v1alpha1_nfs.yaml` or the `kustomization.yaml` in the same directory. Then apply it executing `make deploy`.

Verify it's running executing the following `kubectl` command:

```bash
kubectl get deployment nfs-operator
kubectl get pods
kubectl get nfs
```

To view the output of the operator locally execute `make run` in a different terminal/console:

```bash
make run
```

You can use your own application, pod or deployment to test the provisioned PVC, or you can use an existing consumer application that expose a movies API and used a deployed JSON DB file with movies. Check this [testing instructions](./TESTING.md#tests) to know how to use this consumer application.

## Deployment and Release

The NFS Custom Resource requires the Docker image of the operator released to a Docker registry. Modify the `Makefile` to update - if required - the following variables:

```makefile
VERSION       ?= 0.0.1
OPERATOR_NAME ?= nfs-operator
REGISTRY      ?= johandry
```

To release the Docker image with the operator execute:

```bash
make docker-build docker-push
```

The next step is to release the NFS Custom Resource Definition (CRD), this is done with the execution of the following commands:

```bash
TODO
```

The release of the Docker image and the NFS CRD can be done with the execution of `make release`

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

To destroy the testing environment with the IKS cluster and deployed resources:

```bash
cd test
make clean
```

## Reference to Advance Topics

These topics or references may be required for the development of the Operator:

- [Handle Cleanup on Deletion](https://sdk.operatorframework.io/docs/golang/quickstart/#handle-cleanup-on-deletion)
- [Unit Testing](https://sdk.operatorframework.io/docs/golang/unit-testing/)
- [E2E Tests](https://sdk.operatorframework.io/docs/golang/e2e-tests/)
- [Monitoring with Prometheus](https://sdk.operatorframework.io/docs/golang/monitoring/prometheus/)
- [Controller Runtime Client API](https://sdk.operatorframework.io/docs/golang/references/client/)
- [Logging](https://sdk.operatorframework.io/docs/golang/references/logging/)

**Guides and Tutorials**:

- [Kubebuilder Tutorial](https://book.kubebuilder.io/introduction.html)
- GitHub repo [Operator SDK](https://github.com/operator-framework/operator-sdk)
- [Operator SDK](https://sdk.operatorframework.io/docs/golang/quickstart/)
- [Operator SDK Getting Started](https://github.com/operator-framework/getting-started)
- Book [Kubernetes Operators](https://learning.oreilly.com/library/view/kubernetes-operators/9781492048039/ch06.html#adapter_operators)
- Example [Presentation Operator](https://github.com/NautiluX/presentation-example-operator)

**Storage**:

- [External Provisioners](https://github.com/NautiluX/presentation-example-operator)
- [External Storages](https://github.com/kubernetes-incubator/external-storage)
- [NFS Provisioner](https://github.com/kubernetes-incubator/external-storage/tree/master/nfs)
- [External Storage Provisioners](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner)

**Others**:

- CSI implementation for EFS and NFS
- [CSI Driver NFS](https://github.com/kubernetes-csi/csi-driver-nfs)
- Rook NFS operator
- [Rook Operator Kit](https://github.com/rook/operator-kit)
