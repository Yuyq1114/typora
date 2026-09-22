$ErrorActionPreference = "Stop"

$projectDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$jobName = "kubeflow-train-demo"
$configMapName = "kubeflow-train-demo-code"

Push-Location $projectDir
try {
    Write-Host "[1/4] Updating training code in Kubernetes..."
    kubectl create configmap $configMapName `
        --from-file=train.py `
        --dry-run=client `
        -o yaml | kubectl apply -f -
    if ($LASTEXITCODE -ne 0) { throw "Failed to update ConfigMap" }

    Write-Host "[2/4] Submitting Kubeflow TrainJob..."
    kubectl delete trainjob $jobName --ignore-not-found=true --wait=true
    kubectl apply -f trainjob.yaml
    if ($LASTEXITCODE -ne 0) { throw "Failed to create TrainJob" }

    Write-Host "[3/4] Waiting for training to finish..."
    kubectl wait trainjob/$jobName `
        --for=condition=Complete `
        --timeout=10m
    if ($LASTEXITCODE -ne 0) {
        kubectl get trainjob,jobset,pods
        throw "Training did not complete successfully"
    }

    Write-Host "[4/4] Training logs:"
    kubectl logs `
        -l "jobset.sigs.k8s.io/jobset-name=$jobName" `
        --all-containers=true
}
finally {
    Pop-Location
}
