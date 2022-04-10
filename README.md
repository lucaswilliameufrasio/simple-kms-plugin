# simple-kms-plugin


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
- In another terminal, run `go run cmd/client/unix_client.go`


# (WIP) How i have tested it with Kind


- Attach to kind-control-plane bash
 - I have used lens but you can follow [this tutorial](https://blog.adamgamboa.dev/connecting-shell-to-a-node-in-kubernetes/) to do it in the terminal.
- Navigate to `/etc/kubernetes/pki`
- Create a file called `encryption.yaml` with this content:
```
apiVersion: apiserver.config.k8s.io/v1
kind: EncryptionConfiguration
resources:
  - resources:
      - secrets
    providers:
      - kms:
          name: myKmsPlugin
          endpoint: unix:///tmp/test_kms.sock
          cachesize: 100
          timeout: 3s
      - identity: {}
```
- Add the following flag to `/etc/kubernetes/manifests/kube-apiserver.yaml`: `--encryption-provider-config=/etc/kubernetes/pki/encryption.yaml`
![image](https://user-images.githubusercontent.com/34021576/162605102-220b8021-60cc-462d-b6af-08f3a15e8eaa.png)
- Restart kubelet with `systemctl restart kubelet`


# References
https://blog.logrocket.com/learn-golang-encryption-decryption/
https://www.hairizuan.com/dockerizing-application-that-use-unix-sockets/
https://blog.fearcat.in/a?ID=00700-c7b32931-4077-4fd4-a86d-9f22d7ab9359