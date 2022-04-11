# Simple KMS Plugin

A project developed for `learning purposes` and to be used on a k3s cluster.

# How to run

## TCP

``` bash
ENCRYPTION_SECRET="@pEg<P+lRi<G>?e,dZpWJxyj" make run
```

## Unix

``` bash
PROTOCOL="unix" SOCKET_SERVER_FILE="/tmp/production_socket_file.sock" ENCRYPTION_SECRET="@pEg<P+lRi<G>?e,dZpWJxyj" make run
```


### With Docker


```
docker build -t lucaswilliameufrasio/simple-kms-plugin:0.0.1 .
```

```
mkdir ./tmp
```

```
docker run -it -v $(pwd)/tmp:/tmp -e PROTOCOL="unix" -e ENCRYPTION_SECRET="@pEg<P+lRi<G>?e,dZpWJxyj" lucaswilliameufrasio/simple-kms-plugin:0.0.1
```

```
SOCKET_SERVER_FILE=./tmp/test_kms.sock go run cmd/client/unix_client.go
```


# How to use with K8S
https://kubernetes.io/docs/tasks/administer-cluster/kms-provider/


# How to test the plugin listening to UNIX network with client example

- In one terminal, run `PROTOCOL="unix" ENCRYPTION_SECRET="@pEg<P+lRi<G>?e,dZpWJxyj" make run`
- In another terminal, run `go run cmd/unix_client/client.go`

# How to test the plugin listening to TCP network with client example

- In one terminal, run `ENCRYPTION_SECRET="@pEg<P+lRi<G>?e,dZpWJxyj" make run`
- In another terminal, run `go run cmd/tcp_client/client.go`


# How i have tested it with Kind

- Create a DaemonSet resource with this command: `kubectl apply -f k8s/daemon.yaml`
- Override a label of master node that on my system had an empty value. I fixed it with `kubectl label nodes kind-control-plane node-role.kubernetes.io/master=true --overwrite`
- Attach to kind-control-plane bash
 - I have used `docker exec -it ad619f1862a9 bash` but you can follow [this tutorial](https://blog.adamgamboa.dev/connecting-shell-to-a-node-in-kubernetes/) to do it in the terminal or just use [Lens](https://k8slens.dev/).
- Navigate to `/etc/kubernetes/pki`
- Create a file called `encryption-provider-config.yaml` with this content:
```
apiVersion: apiserver.config.k8s.io/v1
kind: EncryptionConfiguration
resources:
  - resources:
      - secrets
    providers:
      - kms:
          name: simple-kms-plugin
          endpoint: unix:///var/run/simple-kms-plugin/server.sock
          cachesize: 100
          timeout: 3s
      - identity: {}
```
- Add the following flag to `/etc/kubernetes/manifests/kube-apiserver.yaml`: `- --encryption-provider-config=/etc/kubernetes/pki/encryption-provider-config.yaml`
![image](https://user-images.githubusercontent.com/34021576/162663659-7b6b491d-d282-44af-af7a-289781cdc267.png)
- Add the following volume directives at `kube-apiserver.yaml`:
``` yaml
...
volumeMounts:
  - mountPath: /var/run/simple-kms-plugin
    name: simple-kms-plugin-dir
...

volumes:
  - hostPath:
      path: /var/run/simple-kms-plugin
      type: DirectoryOrCreate
    name: simple-kms-plugin-dir
```
![image](https://user-images.githubusercontent.com/34021576/162662626-acaecbc7-8fef-4d4b-9981-d84dc173e2b2.png)

- Save `/etc/kubernetes/manifests/kube-apiserver.yaml` and wait for kube-apiserver restart
- (Optional) Restart kubelet with `systemctl restart kubelet`


# Verifying the whole thing

- Create a secret
```
kubectl create secret generic simple-kms-plugin-secret -n default --from-literal=simple=data
```

- See if the secret is correctly decrypted
```
kubectl get secret simple-kms-plugin-secret -o=jsonpath='{.data.simple}' | base64 -d

```


# References
- https://blog.logrocket.com/learn-golang-encryption-decryption/
- https://www.hairizuan.com/dockerizing-application-that-use-unix-sockets/
- https://blog.fearcat.in/a?ID=00700-c7b32931-4077-4fd4-a86d-9f22d7ab9359
- https://github.com/GoogleCloudPlatform/k8s-cloudkms-plugin
- https://github.com/Tencent/tke-kms-plugin
- https://kubernetes.io/docs/tasks/administer-cluster/kms-provider/
- https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/