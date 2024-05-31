# operator

1. go mod init
```go mod init js.domain/fastapi-operator```
2. operator init
``` operator-sdk init --domain js.domain --project-name js-project --owner "jongsoo"```
3. create api
```operator-sdk create api --group js-group --version v1 --kind JsKind --version v1beta1 ```
4. add spec
5. ...
6. make generate
7. make manifests
8. make docker file ```docker build manager-operator . ```
9. docker tag ```docker tag manager-operator:latest jongsoo/manager-operator:latest ```
10. docker push ```docker push jongsoo/manager-operator:latest ```
11. make deploy 
12. check pod `k get pods -A`
```
NAMESPACE           NAME                                              READY   STATUS      RESTARTS      AGE
js-project-system   js-project-controller-manager-7f747fdd7-rxht2     2/2     Running     0             6m14s
```
13. create crd
```kubectl apply -f config/samples/js_v1_jskind.yaml```
14. validate crd

`kubectl get jskind`
```shell
NAME            AGE
jskind-sample   2m24s
```
`k get deployment`
```shell
NAME            READY   UP-TO-DATE   AVAILABLE   AGE
jskind-sample   1/1     1            1           4m40s
```
`k get pod`
```shell
NAMESPACE           NAME                                              READY   STATUS      RESTARTS      AGE
default             jskind-sample-7f5cf4864b-fvwwc                    1/1     Running     0             54s
```
`k get service`
```shell
NAME            TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
jskind-sample   NodePort       10.100.93.163    <none>        80:32278/TCP                 15h
```
`k get ingress`
```shell
NAME            CLASS   HOSTS   ADDRESS   PORTS   AGE
jskind-sample   nginx   *                 80      15h

```
15. operator log 확인
` k logs js-project-controller-manager-7f747fdd7-rxht2 -n js-project-system`
```shell
2024-05-31T01:26:07Z    INFO    setup   starting manager
2024-05-31T01:26:07Z    INFO    controller-runtime.metrics      Starting metrics server
2024-05-31T01:26:07Z    INFO    starting server {"kind": "health probe", "addr": "[::]:8081"}
2024-05-31T01:26:07Z    INFO    controller-runtime.metrics      Serving metrics server  {"bindAddress": "127.0.0.1:8080", "secure": false}
I0531 01:26:07.190538       1 leaderelection.go:250] attempting to acquire leader lease js-project-system/dfdb12f8.js.domain...
I0531 01:26:07.198344       1 leaderelection.go:260] successfully acquired lease js-project-system/dfdb12f8.js.domain
2024-05-31T01:26:07Z    DEBUG   events  js-project-controller-manager-7f747fdd7-rxht2_490fbf60-4d24-4255-b050-ddc7748ae2e3 became leader        {"type": "Normal", "object": {"kind":"Lease","namespace":"js-project-system","name":"dfdb12f8.js.domain","uid":"114c7e78-1fb9-4ace-bde7-870b22c59f31","apiVersion":"coordination.k8s.io/v1","resourceVersion":"91700"}, "reason": "LeaderElection"}
2024-05-31T01:26:07Z    INFO    Starting EventSource    {"controller": "jskind", "controllerGroup": "js-group.js.domain", "controllerKind": "JsKind", "source": "kind source: *v1beta1.JsKind"}
2024-05-31T01:26:07Z    INFO    Starting Controller     {"controller": "jskind", "controllerGroup": "js-group.js.domain", "controllerKind": "JsKind"}
2024-05-31T01:26:07Z    INFO    Starting workers        {"controller": "jskind", "controllerGroup": "js-group.js.domain", "controllerKind": "JsKind", "worker count": 1}
2024-05-31T01:28:11Z    INFO    Finalizer 추가  {"controller": "jskind", "controllerGroup": "js-group.js.domain", "controllerKind": "JsKind", "JsKind": {"name":"jskind-sample","namespace":"default"}, "namespace": "default", "name": "jskind-sample", "reconcileID": "04026b5e-d12a-4556-9c57-0aabee916c12"}
2024-05-31T01:28:11Z    INFO    JsKind 생성     {"controller": "jskind", "controllerGroup": "js-group.js.domain", "controllerKind": "JsKind", "JsKind": {"name":"jskind-sample","namespace":"default"}, "namespace": "default", "name": "jskind-sample", "reconcileID": "443921c9-0c32-4030-a48b-b78424367938"}
2024-05-31T01:28:11Z    INFO    Deployment 생성 {"controller": "jskind", "controllerGroup": "js-group.js.domain", "controllerKind": "JsKind", "JsKind": {"name":"jskind-sample","namespace":"default"}, "namespace": "default", "name": "jskind-sample", "reconcileID": "443921c9-0c32-4030-a48b-b78424367938"}
2024-05-31T01:28:12Z    INFO    Service 생성    {"controller": "jskind", "controllerGroup": "js-group.js.domain", "controllerKind": "JsKind", "JsKind": {"name":"jskind-sample","namespace":"default"}, "namespace": "default", "name": "jskind-sample", "reconcileID": "443921c9-0c32-4030-a48b-b78424367938"}
2024-05-31T01:28:12Z    INFO    Ingress 생성    {"controller": "jskind", "controllerGroup": "js-group.js.domain", "controllerKind": "JsKind", "JsKind": {"name":"jskind-sample","namespace":"default"}, "namespace": "default", "name": "jskind-sample", "reconcileID": "443921c9-0c32-4030-a48b-b78424367938"}
```