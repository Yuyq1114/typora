$ErrorActionPreference = "Stop"

kubectl delete trainjob stage2-trainjob stage3-custom-runtime stage4-go-client --ignore-not-found=true --wait=true
kubectl delete job stage1-k8s-job --ignore-not-found=true --wait=true
kubectl delete trainingruntime single-gpu-runtime --ignore-not-found=true --wait=true
kubectl delete configmap kubeflow-workflow-code --ignore-not-found=true
