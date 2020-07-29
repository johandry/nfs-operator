# NFS Provisioner

This directory contain all the resources required to deploy the NFS Provisioner. The NFS operator deploy all these resources and the PVC if it does not exists.

To **deploy** manually the NFS Provisioner execute:

```bash
cd nfs-provisioner
kubectl apply -f deployment.yaml
kubectl apply -f rbac.yaml
kubectl apply -f class.yaml
kubectl apply -f claim.yaml
```

To **verify** it's already deployed, either manually or by the operator, execute:

```bash
kubectl get sa nfs-provisioner
kubectl get svc nfs-provisioner
kubectl get deploy nfs-provisioner
kubectl get clusterrole nfs-provisioner-runner
kubectl get clusterrolebinding run-nfs-provisioner
kubectl get role leader-locking-nfs-provisioner
kubectl get rolebinding leader-locking-nfs-provisioner
kubectl get storageclass ibmcloud-nfs

kubectl get pvc nfs
kubectl get pvc nfs | awk '{print $3}' | grep -v VOLUME | while read pv; do kubectl get pv $$pv; done

if kubectl get pvc nfs | grep -q 'Bound'; then echo "NFS PVC is Bound"; else echo "NFS PVC is still Pending"; fi
```

If the operator deploy the resources with other name, replace `nfs-provisioner` with the defined name.

You can **wait** for the provisioned PVC to be ready executing the following `while` command or any of the following `watch` commands waiting for the PV to be `Bound`:

```bash
while kubectl get pvc nfs | grep -q 'Pending'; do printf .; sleep 3; done
# Or
watch kubectl get pvc nfs
# Or
kubectl get pvc nfs -w
```

To **delete** everything deployed, execute:

```bash
cd nfs-provisioner
kubectl delete -f claim.yaml
kubectl delete -f class.yaml
kubectl delete -f rbac.yaml
kubectl delete -f deployment.yaml
```
