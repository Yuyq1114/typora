# 第七阶段：Reconcile

本阶段代码已经包含在第二阶段。

位置：`../2WebAppCrd/controller.go`

- `reconcile()`：读取 WebApp 的 `spec.image`、`spec.replicas`、`spec.port`。
- `ensureDeployment()`：创建 Deployment，或把副本数、镜像和端口恢复到期望状态。
- `ensureService()`：创建或更新 Service，同时保留 API Server 分配的 ClusterIP。

核心关系：

```text
WebApp spec（期望状态） → Reconcile → Deployment / Service（实际状态）
```
