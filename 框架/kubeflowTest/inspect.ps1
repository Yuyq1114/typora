$ErrorActionPreference = "Stop"

Write-Host "=== Kubeflow control plane ==="
kubectl get pods -n kubeflow-system

Write-Host "`n=== Learning resources ==="
kubectl get trainingruntime single-gpu-runtime
kubectl get trainjob stage2-trainjob stage3-custom-runtime stage4-go-client
kubectl get jobset stage2-trainjob stage3-custom-runtime stage4-go-client
kubectl get job,pod -o wide | Select-String "stage1-k8s-job|stage2-trainjob|stage3-custom-runtime|stage4-go-client"

Write-Host "`n=== GPU ==="
kubectl get nodes -o custom-columns="NAME:.metadata.name,GPU:.status.allocatable.nvidia\.com/gpu"
