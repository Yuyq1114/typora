$ErrorActionPreference = "Stop"

kubectl delete trainjob kubeflow-train-demo --ignore-not-found=true
kubectl delete configmap kubeflow-train-demo-code --ignore-not-found=true
