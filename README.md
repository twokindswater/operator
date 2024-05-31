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
10. 