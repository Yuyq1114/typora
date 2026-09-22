# 扩展三：Admission Webhook / API Server

状态：已实现，服务端代码在 `webhook.go`。

目标：理解请求写入 etcd 前的认证、授权、准入和校验流程。

```text
Request
   ↓
Authentication
   ↓
Authorization
   ↓
Mutating Admission
   ↓
Schema Validation
   ↓
Validating Admission
   ↓
Conversion / Storage Version
   ↓
Storage
   ↓
etcd
```

最小实现分两步：

1. Validating Webhook：拒绝没有资源限制的 Pod。
2. Mutating Webhook：为没有标签的 Pod 自动添加标签。

当前服务会自动生成本机 HTTPS 证书；WebhookConfiguration 使用 `host.docker.internal` 访问 GoLand 中运行的服务，只作用于带 `mini-webhook=enabled` 标签的 namespace。

1. 将 `../main.go` 的 `controllerChoice` 改成 `13` 并启动。
2. 在另一个终端执行：

```powershell
.\13AdmissionWebhook\install.ps1
kubectl apply -f .\13AdmissionWebhook\pod-valid.yaml
kubectl apply -f .\13AdmissionWebhook\pod-invalid.yaml
```

合法 Pod 会得到 `mini-webhook=mutated` 标签；缺少 CPU/Memory limits 的 Pod 会被拒绝。

清理：

```powershell
kubectl delete mutatingwebhookconfiguration mini-mutating-webhook
kubectl delete validatingwebhookconfiguration mini-validating-webhook
kubectl label namespace test1 mini-webhook-
```
