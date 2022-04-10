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


# How to use with K8S
https://kubernetes.io/docs/tasks/administer-cluster/kms-provider/


# How to test the plugin listening to UNIX network with client example

- In one terminal, run `PROTOCOL="unix" ENCRYPTION_SECRET="@pEg<P+lRi<G>?e,dZpWJxyj" make run`
- In another terminal, run `go run cmd/client/unix_client.go`


# References
https://blog.logrocket.com/learn-golang-encryption-decryption/
https://www.hairizuan.com/dockerizing-application-that-use-unix-sockets/
https://blog.fearcat.in/a?ID=00700-c7b32931-4077-4fd4-a86d-9f22d7ab9359