<#
 .SYNOPSIS
    启动国美考务系统（局域网共享模式），让同一网络下的其他电脑可以访问。
 .DESCRIPTION
    依次启动本地 MySQL、Go 后端、Vite 前端。
    前端绑定 0.0.0.0 并自动在 Windows 防火墙放行 8848 端口，
    其他同事通过 http://<你的局域网IP>:8848 即可访问。
    注意：添加防火墙规则需要以管理员身份运行本脚本。
#>

# ---------- 自动以管理员身份重新运行（防火墙规则需要管理员权限） ----------
if (-not ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Host "重新以管理员身份启动以配置防火墙..." -ForegroundColor Yellow
    Start-Process -FilePath "powershell.exe" -Verb RunAs -ArgumentList "-NoProfile -ExecutionPolicy Bypass -File `"$PSCommandPath`""
    exit
}

$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
$runtime = Join-Path $root ".runtime"
$mysqlIni = Join-Path $runtime "mysql.ini"
$mysqlServer = Join-Path $runtime "mysql\bin\mysqld.exe"
$mysqlAdmin = Join-Path $runtime "mysql\bin\mysqladmin.exe"
$backendExe = Join-Path $runtime "bin\gm-affair.exe"
$nodeExe = (Get-Command node.exe -ErrorAction Stop).Source
$env:MYSQL_PWD = "123456"

$frontendPort = 8848

function Test-Http([string]$Url) {
    try {
        $response = Invoke-WebRequest -UseBasicParsing -Uri $Url -TimeoutSec 2
        return $response.StatusCode -eq 200
    } catch {
        return $false
    }
}

function Test-MySql {
    $previousPreference = $ErrorActionPreference
    $ErrorActionPreference = "SilentlyContinue"
    & $mysqlAdmin --defaults-file=$mysqlIni -uroot ping 2>$null | Out-Null
    $exitCode = $LASTEXITCODE
    $ErrorActionPreference = $previousPreference
    return $exitCode -eq 0
}

function Wait-Until([scriptblock]$Check, [string]$Name, [int]$Seconds = 60) {
    foreach ($index in 1..($Seconds * 2)) {
        if (& $Check) {
            return
        }
        Start-Sleep -Milliseconds 500
    }
    throw "$Name startup timed out."
}

# ---------- 1. MySQL ----------
if (-not (Test-Path $mysqlServer)) {
    throw "Local MySQL runtime was not found under .runtime\mysql."
}
if (-not (Test-MySql)) {
    $mysqlProcess = Start-Process -FilePath $mysqlServer -ArgumentList "--defaults-file=$mysqlIni" -WindowStyle Hidden -PassThru
    Set-Content -LiteralPath (Join-Path $runtime "mysql-process.id") -Value $mysqlProcess.Id -Encoding ascii
    Wait-Until { Test-MySql } "MySQL"
}

# ---------- 2. Go 后端 ----------
if (-not (Test-Path $backendExe)) {
    throw "Backend executable .runtime\bin\gm-affair.exe was not found."
}
if (-not (Test-Http "http://127.0.0.1:30000/captcha")) {
    $backendProcess = Start-Process -FilePath $backendExe `
        -WorkingDirectory (Join-Path $root "backend") `
        -WindowStyle Hidden `
        -PassThru
    Set-Content -LiteralPath (Join-Path $runtime "backend-process.id") -Value $backendProcess.Id -Encoding ascii
    Wait-Until { Test-Http "http://127.0.0.1:30000/captcha" } "Go backend"
}

# ---------- 3. Vite 前端（绑定 0.0.0.0，允许局域网访问） ----------
if (-not (Test-Path (Join-Path $root "frontend\node_modules\vite\bin\vite.js"))) {
    throw "Frontend dependencies are missing. Run pnpm install in frontend."
}
if (-not (Test-Http "http://127.0.0.1:$frontendPort")) {
    $frontendProcess = Start-Process -FilePath $nodeExe `
        -WorkingDirectory (Join-Path $root "frontend") `
        -ArgumentList "node_modules/vite/bin/vite.js", "--host", "0.0.0.0", "--port", "$frontendPort" `
        -WindowStyle Hidden `
        -PassThru
    Set-Content -LiteralPath (Join-Path $runtime "frontend-process.id") -Value $frontendProcess.Id -Encoding ascii
    Wait-Until { Test-Http "http://127.0.0.1:$frontendPort" } "Vue frontend"
}

# ---------- 4. 防火墙放行前端端口 ----------
$fwRuleName = "GM-Exam-Frontend-$frontendPort"
$existing = netsh advfirewall firewall show rule name="$fwRuleName" 2>$null
if ($existing -notmatch "规则名称") {
    netsh advfirewall firewall add rule name="$fwRuleName" dir=in action=allow protocol=TCP localport=$frontendPort | Out-Null
    Write-Host "已添加防火墙入站规则: $fwRuleName" -ForegroundColor Cyan
} else {
    Write-Host "防火墙入站规则已存在: $fwRuleName" -ForegroundColor Cyan
}

# ---------- 5. 打印局域网访问地址 ----------
$lanIp = (Get-NetIPAddress -AddressFamily IPv4 |
    Where-Object { $_.InterfaceAlias -notmatch "Loopback" -and $_.IPAddress -notmatch "^169\.254" -and $_.PrefixOrigin -ne "WellKnown" } |
    Select-Object -First 1).IPAddress
if (-not $lanIp) {
    $lanIp = (Test-Connection -ComputerName (Hostname) -Count 1).IPv4Address.IPAddressToString
}

Write-Host ""
Write-Host "GM exam system is running (共享模式)." -ForegroundColor Green
Write-Host "本机访问:    http://127.0.0.1:$frontendPort"
if ($lanIp) {
    Write-Host "同事访问:    http://$lanIp`:$frontendPort" -ForegroundColor Green
}
Write-Host "Backend:     http://127.0.0.1:30000"
Write-Host "Database:    127.0.0.1:13306 / yx_exam_affair_gm"
Write-Host "Username:    admin"
Write-Host "Password:    Admin123"
