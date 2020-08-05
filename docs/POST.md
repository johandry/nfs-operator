# NFS Provisioner Operator with OperatorSDK & Kubebuilder

`defer Conclusion()`

- [NFS Provisioner Operator with OperatorSDK & Kubebuilder](#nfs-provisioner-operator-with-operatorsdk--kubebuilder)
  - [Kubernetes Operators](#kubernetes-operators)
  - [Use Case](#use-case)
  - [Design](#design)
    - [Deployment](#deployment)
    - [StorageClass](#storageclass)
    - [NFS Provisioner PersistentVolumeClaim](#nfs-provisioner-persistentvolumeclaim)
    - [Back Storage PersistentVolumeClaim](#back-storage-persistentvolumeclaim)
  - [Use](#use)
    - [Container and Deployment](#container-and-deployment)
    - [Custom Resource Definition](#custom-resource-definition)
  - [Development](#development)
  - [Conclusion](#conclusion)
  - [Reference](#reference)

## Kubernetes Operators

An **Operator** is a way to package, run, and maintain a Kubernetes application, acting as a _controller_ for a [Custom Resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) to manage this application and their components.

A **Resource** is an endpoint of the [Kubernetes API](https://kubernetes.io/docs/reference/using-api/api-overview/) that store a collection of [API Objects](https://kubernetes.io/docs/concepts/overview/working-with-objects/kubernetes-objects/). For example, the [Pod](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#pod-v1-core) resource endpoint contain a collection of Pod Objects, each one with it's own specifications and status.

A **Custom Resource** (CR) is a resource that is not in the default installation of Kubernetes. It can be used to install an application, a service or Kubernetes functions. For example, an application to [deploy and host presentations](https://opensource.com/article/20/3/kubernetes-operator-sdk) or to [deploy and manage Redis clusters](https://medium.com/manikkothu/build-kubernetes-operator-using-kubebuilder-4bfef299757d) and many [more](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/#example).

To represent a Custom Resource is similar to any other resources: using a YAML file. This YAML file can be deployed to the cluster to execute its function. For example, the YAML file `nfs-presentation.yaml` may be a presentation custom resource:

```yaml
apiVersion: presentation.johandry.com/v1alpha1
kind: Presentation
metadata:
  name: nfs-operator-presentation
spec:
  markdown: |
    # NFS Provisioner Operator with OperatorSDK & Kubebuilder
    ---
    ## Kubernetes Operators

    * Test
    ---
```

In order to Kubernetes understand and validate this Custom Resource it needs a **Custom Resource Definition** (CRD) to the CR into the Kubernetes API. Also, in order to execute the logic of the Custom Resource it needs a **controller** which is a containerized application running in one Pod controlled by a **Deployment** (or Job or DaemonSet) and exposed by a **Service**. Finally, **RBAC** rules are required to authorize access to the controller in the cluster.

_TODO: Image with the Operator Architecture_

Every resource (custom or not) has **specifications** and **status**. When the Kubernetes user/developer execute `kubectl create -f nfs-presentation.yaml`, Kubernetes through the operator controller make sure the CR objects specifications are set as defined, in this example by making sure a deployment exists running Nginx, a collection of ConfigMap exists with the web page formed by the given markdown content and a Service is there to expose this web pages. If I change the presentation Custom Resource the operator controller will update the previous objects to make the changes to the specifications happen. Also, when the Kubernetes user/developer execute `kubectl get presentation nfs-operator-presentation`, the output is the status defined by the operator controller of this presentation Custom Resource.

## Use Case

The `PersistentVolumeClaim` available on IBM Cloud Gen 2 - at this time - only allows access mode `ReadWriteOnce` because it's a BlockStorage type volume. There is also a limit of block storages you can have in a VPC, as well as there is a limit of the total volume size requested. You can read more about [IBMCloud VPC Block Storage](https://cloud.ibm.com/docs/containers?topic=containers-vpc-block) in the IBMCloud documentation site.

In a regular application you may have a deployment with multiple replicas of your Pod or containerized application. If your application needs an external storage you can provide it with a PVC using a external provisioner. The limitation, as mentioned above, is that this PVC is RWO and cannot be used by multiple containers restricting the number of replicas or containers to just one.

![Application using a PVC](./images/NFS_Provisioner-PVC.png)

This is a simple example of the PVC and the deployment.

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: movies
spec:
  storageClassName: ibmc-vpc-block-general-purpose
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: movies
  name: movies
spec:
  replicas: 1
  selector:
    matchLabels:
      app: movies
  template:
    metadata:
      labels:
        app: movies
    spec:
      volumes:
        - name: movies-volume
          persistentVolumeClaim:
            claimName: movies
      containers:
        - image: us.icr.io/iac-registry/movies:1.1
          name: movies
          volumeMounts:
            - name: movies-volume
              mountPath: "/data"
```

The `PersistentVolume` is [dynamically provisioned](https://kubernetes.io/blog/2017/03/dynamic-provisioning-and-storage-classes-kubernetes/), this means it's created by demand. The developer only needs to create the PVC using a predefined StorageClass that is used by a Provisioner to create the volume. The `StorageClass` `ibmc-vpc-block-general-purpose` and others can be found in the [IBMCloud VPC Block Storage](https://cloud.ibm.com/docs/containers?topic=containers-vpc-block) documentation.

As a side note, even if the application needs and request 1Gb of storage the minimum provided size is 10Gb. A more complex and detailed example can be found in the [IBM Cloud Containers Patterns](https://ibm.github.io/cloud-enterprise-examples/iac-resources/container#persistent-volumes)

There are multiple solutions to this problem. For example, to use [StatefulSet](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) using `volumeClaimTemplates` which will create a `PersistentVolume` for every replica provisioned by a **PersistentVolume Provisioner**.

## Design

Other solution is to use a [External Storage Provisioners](https://github.com/kubernetes-incubator/external-storage). There are many external provisioners such as NFS, EFS, Local Volume, and many others maintained by the community. There is also a [library](https://github.com/kubernetes-sigs/sig-storage-lib-external-provisioner) to create your own [external provisioner](https://kubernetes.io/docs/concepts/storage/storage-classes/#provisioner).

The following diagram shows how to use the NFS Provisioner to allow multiple containers to read from a `ReadWriteMany` PVC. This PVC is managed by the NFS Provisioner exposing a NFS Service and this NFS uses a backend storage that you should previously create. The backend storage is accessible through a `PersistentVolumeClaim`, in this example, it's the previously created IBMCloud VPC Block Storage.

![Application using the NFS Provisioner](./images/NFS_Provisioner-NFS%20Provisioner.png)

The NFS Provisioner consist of a **Deployment**, **Service** to allow external resources to access the container NFS service, a **PVC** to access backend storage (if any), **ServiceClass**, **Service Account**, **Roles** and **RoleBinding** (cluster and namespaced) to manage RBAC security. All these files are in the `test/kubernetes/nfs-provisioner` directory of the [GitHub repo](https://github.com/johandry/nfs-operator/tree/master/test/kubernetes/nfs-provisioner) but here are segments of code of some of them.

The NFS External Provisioner, like any other external provisioner, is a dynamic `PersistentVolume` provisioner. A `StorageClass` is defined to be its `provisioner`, the created instance then watch for `PersistentVolumeClaims` asking for this `StorageClass` to automatically create the `PersistentVolume` for the user resource.

### Deployment

This deployment has one replica of a container running the image `quay.io/kubernetes_incubator/nfs-provisioner` which serve the NFS service.

```yaml
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
          # Container that serving the NFS service
          image: quay.io/kubernetes_incubator/nfs-provisioner:latest
          ports:
            # List of ports related to NFS service
            ...
          securityContext:
            capabilities:
              add:
                - DAC_READ_SEARCH
                - SYS_RESOURCE
          args:
            # Name of the provisioner defined in the StorageClass
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
            # Name of the PersistentVolumeClaim previously created
            claimName: nfs-block-custom
```

### StorageClass

The NFS Provisioner will monitor for any PVC request using this StorageClass. The StorageClasses are required for [Dynamic Provisioning](https://kubernetes.io/blog/2017/03/dynamic-provisioning-and-storage-classes-kubernetes/).

```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: ibmcloud-nfs
provisioner: ibmcloud/nfs
mountOptions:
  - vers=4.1
```

### NFS Provisioner PersistentVolumeClaim

Do not confuse this PVC withe the PVC used for Backend Storage. This is the PVC to use the NFS volume expose by the containers. This example PVC only request 1Mb of the NFS storage.

```yaml
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: nfs
spec:
  storageClassName: ibmcloud-nfs
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 1Mi
```

### Back Storage PersistentVolumeClaim

This is the PVC used by the container to use as volume for the NFS service. This should be the first resource to create and last to delete. If you want to reuse the files and data stored by the containers make sure to create a backup, snapshot or do not delete it. This PVC request 10Gb from the volume, there will be plenty of storage if the containers use just 1Mb.

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: movies
spec:
  storageClassName: ibmc-vpc-block-general-purpose
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
```

## Use

To use the NFS Provisioner you need to create several resources, so wouldn't be nice to have an operator to manage everything for you?

This operator deploy all the resources previously mentioned to have a backend storage and the NFS Provisioner. So, you only need to deploy the operator when the cluster is created, then use the `nfs` custom resource and a PVC to request storage for your containers, as many as you want.

![Application using the NFS Provisioner Operator](./images/NFS_Provisioner-NFS%20Provisioner%20Operator.png)

The operator, like most of the operators, is formed by a **Container** and a **Deployment** to have this containerized operator running in the cluster, also a set of resources to setup RBAC for the operator and the **Custom Resource Definition** to define the Nfs resource.

To use the operator you need the **NFS Custom Resource** and a **PVC** to request storage from the NFS Provisioner.

### Container and Deployment

The operator application is containerized and available on any reachable Container Registry, in this case we use DockerHub and it's available on `docker.io/johandry/nfs-operator`. The operator application is the one that deploy or creates all the required resources mentioned before. The operator, among other things, is watching all the resources regularly and keep them with the defined specifications.

The deployment only has one replica and it's used to have the operator container running on the cluster. The deployment can be view at `deploy/operator.yaml` in the [GitHub repo](https://github.com/johandry/nfs-operator/tree/master/deploy).

### Custom Resource Definition

The Custom Resource Definition (CRD) allow us to define the **Nfs** object to be used like any other Kubernetes object. The CRD can be found in the `deploy/crds/*_nfs_crd.yaml` in the [GitHub repo](https://github.com/johandry/nfs-operator/tree/master/deploy/crds).

**Custom Resource** and **PersistentVolumeClaim**

Once the CRD is created the Nfs object kind can be used by the kubernetes developer or admin using a **Custom Resource** file. The CR is defined in the `pkg/apis/ibmcloud/v1alpha1/nfs_types.go` file in the `NfsSpec` struct.

A demo for a NFS CR would be like the following.

```yaml
apiVersion: ibmcloud.ibm.com/v1alpha1
kind: Nfs
metadata:
  name: nfs
  namespace: nfs-test
spec:
  storageClass: ibmcloud-nfs
  provisionerAPI: ibmcloud/nfs
  backingStorage:
    pvcName: nfs-block-custom
```

Creating this CR allow us to claim a volume using the `storageClassName` **ibmcloud-nfs**

## Development

The Operators are usually developed using a development tool, there are multiple of them such as [KUDO](https://kudo.dev/), [Metacontroller](https://metacontroller.app/), [Kopf](https://github.com/zalando-incubator/kopf) for Python and [Java Operator SDK](https://github.com/ContainerSolutions/java-operator-sdk), but most of the operators are managed by the [controller runtime](https://github.com/kubernetes-sigs/controller-runtime) which is used and supported by [Kubebuilder](https://book.kubebuilder.io) (from Kubernetes community) and [Operator SDK](https://github.com/operator-framework/getting-started) (from CoreOS, now RedHat).

With **Operator SDK** you can develop operators using Go, Ansible or Helm, and since version [0.19.0](https://github.com/operator-framework/operator-sdk/releases/tag/v0.19.0) it is integrated with **Kubebuilder**, so basically you'll be working with both and the developers can refer to their documentation.

The first step working with Operator SDK is to install it. If you are on macOS and have `brew` you can easily install it with the following command, otherwise check the [installation guide](https://sdk.operatorframework.io/docs/install-operator-sdk/).

```bash
brew install operator-sdk
```

## Conclusion

Operators are software SRE, everyone is building operators now and most of the times it's good to have one to package, run, and maintain an application on Kubernetes. Once the CRD with the operator or controller is deployed the users or developers only needs to create/deploy, get status and (if needed) modify a **Custom Resource**.

Using **Operator SDK** is an excellent tool to develop your Operators, you can use Helm or Ansible but for better control and understanding of operators, I'd suggest to build them on Go. Also, with the recent integration with **Kubebuilder** the tool is even better and stronger.

The development of a Operator begins with (1) knowing why it's needed, (1) develop or design the resources required to accomplish your task, then (3) use Operator SDK/Kubebuilder to define the resource(s) API Specs and Status, and (4) develop the controller(s).

You may test your controller locally using the Operator SDK/Kubebuilder tools or on an external Kubernetes cluster, it could be local with KinD or MiniKube.

Finally, you can release the operator on your own web site, github repo or on [OperatorHub](https://operatorhub.io).

## Reference

- [Kubernetes Operators](https://learning.oreilly.com/library/view/kubernetes-operators/9781492048039/) book at O'Reilly.
- [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) at kubernetes.io.
- [Build a Kubernetes Operator in 10 minutes with Operator SDK](https://opensource.com/article/20/3/kubernetes-operator-sdk)
- [Build Kubernetes Operator using Kubebuilder](https://medium.com/manikkothu/build-kubernetes-operator-using-kubebuilder-4bfef299757d)
- [Controller Runtime](https://github.com/kubernetes-sigs/controller-runtime) GitHub repository
- [Kubebuilder](https://book.kubebuilder.io/) online book.
- [Operator SDK](https://sdk.operatorframework.io) website.
- [Go Based Operators](https://sdk.operatorframework.io/docs/golang/quickstart/) using Operator SDK (_and Kubebuilder_)
