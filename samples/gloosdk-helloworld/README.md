## Creating a VirtualService
### How to create a virtual service for gloo

Prerequisites:
1. Kubernetes cluster with Gloo installed.
2. Create the deployment:

```shell
kubectl -n default apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: echo
  name: echo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: echo
  template:
    metadata:
      labels:
        app: echo
    spec:
      containers:
      - image: ealen/echo-server
        imagePullPolicy: IfNotPresent
        name: echo
        ports:
        - containerPort: 8080
          name: http
          protocol: TCP
        env:
        - name: PORT
          value: '8080'
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: echo
  name: echo
spec:
  ports:
  - name: http
    port: 8080
    protocol: TCP
    targetPort: 8080
  selector:
    app: echo
EOF
```
3. Run the application as follows:

```shell
POD_NAMESPACE=default go run samples/gloosdk-helloworld/main.go

echo-vs virtual service created
```

You can test the result with:

```shell
$ curl $(glooctl proxy url)

{"host":{"hostname":"172.18.0.2","ip":"::ffff:10.244.0.18","ips":[]},"http":{"method":"GET","baseUrl":"","originalUrl":"/","protocol":"http"},"request":{"params":{"0":"/"},"query":{},"cookies":{},"body":{},"headers":{"host":"172.18.0.2:30811","user-agent":"curl/7.74.0","accept":"*/*","x-forwarded-proto":"http","x-request-id":"df6eb31c-308f-40e8-aa22-1a18341911f3","x-envoy-expected-rq-timeout-ms":"15000"}},"environment":{"PATH":"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","HOSTNAME":"echo-568d95b797-456dd","NODE_VERSION":"14.17.1","YARN_VERSION":"1.22.5","PORT":"8080","KUBERNETES_SERVICE_HOST":"10.96.0.1","KUBERNETES_PORT":"tcp://10.96.0.1:443","KUBERNETES_PORT_443_TCP_PROTO":"tcp","KUBERNETES_PORT_443_TCP":"tcp://10.96.0.1:443","ECHO_PORT_8080_TCP_PORT":"8080","ECHO_PORT_8080_TCP_ADDR":"10.96.127.35","ECHO_PORT":"tcp://10.96.127.35:8080","ECHO_PORT_8080_TCP":"tcp://10.96.127.35:8080","KUBERNETES_SERVICE_PORT_HTTPS":"443","ECHO_SERVICE_HOST":"10.96.127.35","ECHO_SERVICE_PORT":"8080","ECHO_SERVICE_PORT_HTTP":"8080","ECHO_PORT_8080_TCP_PROTO":"tcp","KUBERNETES_SERVICE_PORT":"443","KUBERNETES_PORT_443_TCP_PORT":"443","KUBERNETES_PORT_443_TCP_ADDR":"10.96.0.1","HOME":"/root"}}%
```


Cleanup:
```shell
kubectl delete vs echo-vs -n gloo-system
```
