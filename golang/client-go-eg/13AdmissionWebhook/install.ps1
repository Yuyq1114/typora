$ErrorActionPreference = "Stop"
$caPath = Join-Path $PSScriptRoot "tls\ca.crt"
if (-not (Test-Path -LiteralPath $caPath)) {
    throw "先把 main.go 的 controllerChoice 设为 13 并启动一次，以生成 tls/ca.crt"
}
$caBundle = [Convert]::ToBase64String([IO.File]::ReadAllBytes($caPath))
$wslAddress = Get-NetIPAddress -AddressFamily IPv4 -ErrorAction SilentlyContinue |
    Where-Object { $_.InterfaceAlias -like "*WSL*" } |
    Select-Object -First 1 -ExpandProperty IPAddress
$webhookHost = if ($wslAddress) { $wslAddress } else { "host.docker.internal" }
$template = Get-Content -LiteralPath (Join-Path $PSScriptRoot "webhook.yaml") -Raw
$rendered = $template.Replace("CA_BUNDLE", $caBundle).Replace("WEBHOOK_HOST", $webhookHost)
$rendered | kubectl apply -f -
kubectl label namespace test1 mini-webhook=enabled --overwrite
Write-Output "Webhook URL host: $webhookHost"
