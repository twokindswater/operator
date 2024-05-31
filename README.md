# operator

1. go mod init
```go mod init js.domain/fastapi-operator```
2. operator init
``` operator-sdk init --domain js.domain --project-name js-project --owner "jongsoo"```
3. create api
```operator-sdk create api --group js-group --version v1 --kind JsKind --version v1beta1 ```
4. 