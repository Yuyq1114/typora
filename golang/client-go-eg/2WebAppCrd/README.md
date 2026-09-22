# 第二阶段：WebApp CRD

目标：让 Kubernetes API Server 认识 `WebApp`，并由 controller 维护同名 Deployment 和 Service。

现有代码：

- `crd.yaml`：定义 `demo.example.com/v1` 的 WebApp CRD。
- `webapp.yaml`：`test1/webapp-demo` 示例。
- `controller.go`：WebApp controller。
- `../main.go`：将 `controllerChoice` 改成 `2` 后启动本阶段。

第一次使用需要安装 CRD 和示例：

```powershell
kubectl apply -f .\2WebAppCrd\crd.yaml
kubectl apply -f .\2WebAppCrd\webapp.yaml
```
