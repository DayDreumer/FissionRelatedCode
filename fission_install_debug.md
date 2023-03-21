## DEBUG_RECORD

### K8s上部署pod

#### 0/1 nodes are available: 1 pod has unbound immediate PersistentVolumeClaims.

解：手动创建storageclass pod.

```shell
# step 1
kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/master/deploy/local-path-storage.yaml
# step 2
kubectl patch storageclass local-path -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'
# restart k8s
kubectl reset
kubeadm init --kubernetes-version=v1.23.5 --pod-network-cidr=10.244.0.0/16 --image-repository=registry.aliyuncs.com/google_containers
```

