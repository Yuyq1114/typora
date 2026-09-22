param(
    [Parameter(Mandatory = $true)]
    [ValidateSet(1, 2, 3, 4)]
    [int]$Stage
)

$ErrorActionPreference = "Stop"
$projectDir = Split-Path -Parent $MyInvocation.MyCommand.Path

function Invoke-Kubectl {
    param([Parameter(ValueFromRemainingArguments = $true)][string[]]$Arguments)
    & kubectl @Arguments
    if ($LASTEXITCODE -ne 0) {
        throw "kubectl failed: $($Arguments -join ' ')"
    }
}

Push-Location $projectDir
try {
    Write-Host "Updating training code ConfigMap..."
    kubectl create configmap kubeflow-workflow-code `
        --from-file=train.py=training/train.py `
        --dry-run=client -o yaml | kubectl apply -f -
    if ($LASTEXITCODE -ne 0) { throw "Failed to update training code" }

    switch ($Stage) {
        1 {
            Invoke-Kubectl delete job stage1-k8s-job --ignore-not-found=true --wait=true
            Invoke-Kubectl apply -f manifests/01-k8s-job.yaml
            Invoke-Kubectl wait job/stage1-k8s-job --for=condition=Complete --timeout=10m
            Invoke-Kubectl logs job/stage1-k8s-job
        }
        2 {
            Invoke-Kubectl delete trainjob stage2-trainjob --ignore-not-found=true --wait=true
            Invoke-Kubectl apply -f manifests/02-trainjob.yaml
            Invoke-Kubectl wait trainjob/stage2-trainjob --for=condition=Complete --timeout=10m
            Invoke-Kubectl logs -l jobset.sigs.k8s.io/jobset-name=stage2-trainjob --all-containers=true
        }
        3 {
            Invoke-Kubectl apply -f manifests/03-runtime.yaml
            Invoke-Kubectl delete trainjob stage3-custom-runtime --ignore-not-found=true --wait=true
            Invoke-Kubectl apply -f manifests/03-trainjob.yaml
            Invoke-Kubectl wait trainjob/stage3-custom-runtime --for=condition=Complete --timeout=10m
            Invoke-Kubectl logs -l jobset.sigs.k8s.io/jobset-name=stage3-custom-runtime --all-containers=true
        }
        4 {
            Invoke-Kubectl apply -f manifests/03-runtime.yaml
            Push-Location go-client
            try {
                go run .
                if ($LASTEXITCODE -ne 0) { throw "Go client failed" }
            }
            finally {
                Pop-Location
            }
        }
    }
}
finally {
    Pop-Location
}
