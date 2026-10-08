<#
 .SYNOPSIS
    本地开发启动（不编译 exe）：用 go run 起后端 + Vite 起前端。
 .DESCRIPTION
    1. 确保本地 MySQL 已运行（若未运行则拉起，不关闭）。
    2. 停止可能残留的 go run 后端与 vite 前端（保证每次都跑最新代码）。
    3. 用 `go run gm_affair_main.go` 直接启动 Go 后端（内存编译，不生成 exe）。
    4. 用 Vite dev 启动前端（热更新，不打包）。
    适合开发调试，无需编译部署产物。
#>

$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
$runtime = Join-Path $root ".runtime"
$mysqlIni = Join-Path $runtime "mysql.ini"
$mysqlServer = Join-Path $runtime "mysql\bin\mysqld.exe"
$mysqlAdmin = Join-Path $runtime "mysql\bin\mysqladmin.exe"
$backendDir = Join-Path $root "backend"
$goExe = (Get-Command go.exe -ErrorAction Stop).Source
$nodeExe = (Get-Command node.exe -ErrorAction Stop).Source
$env:MYSQL_PWD = "123456"

$frontendPort = 8848
$backendPort = 30000

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

function Wait-Until([scriptblock]$Check, [string]$Name, [int]$Seconds = 90) {
    foreach ($index in 1..($Seconds * 2)) {
        if (& $Check) {
            return
        }
        Start-Sleep -Milliseconds 500
    }
    throw "$Name startup timed out."
}

function Stop-ByCommandLine {
    param([Parameter(Mandatory = $true)][string]$MatchPattern)
    $count = 0
    Get-Process | Where-Object {
        $_.CommandLine -and ($_.CommandLine -match $MatchPattern)
    } | ForEach-Object {
        try {
            Write-Host "Stopping existing PID $($_.Id) : $($_.ProcessName)" -ForegroundColor Gray
            $_.Kill()
            $count++
        } catch {
            # ignore
        }
    }
    return $count
}

# ---------- 1. MySQL（确保 13306 上跑的是本项目实例；否则拉起，但不在此脚本中关闭） ----------
if (-not (Test-Path $mysqlServer)) {
    throw "Local MySQL runtime was not found under .runtime\mysql."
}

# 端口互串防护：若 13306 被其他项目的 MySQL 实例占用，先停掉再拉起本项目实例。
$projectDataDir = (Join-Path $runtime "mysql-data").Replace('/', '\').TrimEnd('\').ToLower()
$mysqlListeners = $null
try { $mysqlListeners = Get-NetTCPConnection -LocalPort 13306 -State Listen -ErrorAction SilentlyContinue } catch { }
if ($mysqlListeners) {
    foreach ($l in $mysqlListeners) {
        $cmd = $null
        try { $cmd = (Get-CimInstance Win32_Process -Filter "ProcessId=$($l.OwningProcess)" -ErrorAction SilentlyContinue).CommandLine } catch { }
        if ($cmd) {
            $m = [regex]::Match($cmd, '--datadir[=\s]+"?([^"\s]+)"?')
            if ($m.Success) {
                $dd = $m.Groups[1].Value.Replace('/', '\').TrimEnd('\').ToLower()
                if ($dd -ne $projectDataDir) {
                    Write-Host "Stopping foreign MySQL (PID $($l.OwningProcess), datadir=$dd) holding port 13306..." -ForegroundColor Yellow
                    Stop-Process -Id $l.OwningProcess -Force -ErrorAction SilentlyContinue
                    Start-Sleep -Seconds 2
                }
            }
        }
    }
}

if (-not (Test-MySql)) {
    $mysqlProcess = Start-Process -FilePath $mysqlServer -ArgumentList "--defaults-file=$mysqlIni" -WindowStyle Hidden -PassThru
    Set-Content -LiteralPath (Join-Path $runtime "mysql-process.id") -Value $mysqlProcess.Id -Encoding ascii
    Wait-Until { Test-MySql } "MySQL"
}

# ---------- 2. 先停止已有的 go run / 编译版 exe / vite 进程，确保端口释放且跑的是最新代码 ----------
Write-Host "Cleaning up existing dev processes..." -ForegroundColor Cyan
Stop-ByCommandLine -MatchPattern "gm_affair_main\.go" | Out-Null
Stop-ByCommandLine -MatchPattern "gm-affair\.exe" | Out-Null
Stop-ByCommandLine -MatchPattern "vite.*--host|node_modules[\\/]vite[\\/]bin[\\/]vite\.js" | Out-Null
# 同时按端口 30000 占用者强制释放（兜底：杀掉占用 30000 的任意进程）
$holder = (Get-NetTCPConnection -LocalPort $backendPort -ErrorAction SilentlyContinue).OwningProcess
if ($holder) {
    try {
        Write-Host "Releasing port $backendPort (PID $holder)..." -ForegroundColor Gray
        Stop-Process -Id $holder -Force -ErrorAction SilentlyContinue
    } catch { }
}

# 给端口释放一点缓冲时间
Start-Sleep -Milliseconds 500

# 等待确认旧进程已退出
$wait = 0
while ((Test-Http "http://127.0.0.1:$backendPort/captcha") -and $wait -lt 20) {
    Start-Sleep -Milliseconds 500
    $wait++
}

# ---------- 3. Go 后端（go run，不编译 exe） ----------
Write-Host "Starting Go backend with go run..." -ForegroundColor Cyan
$backendProcess = Start-Process -FilePath $goExe `
    -ArgumentList "run", "gm_affair_main.go" `
    -WorkingDirectory $backendDir `
    -WindowStyle Hidden `
    -PassThru
Set-Content -LiteralPath (Join-Path $runtime "backend-dev.id") -Value $backendProcess.Id -Encoding ascii
Wait-Until { Test-Http "http://127.0.0.1:$backendPort/captcha" } "Go backend (go run)"

# ---------- 4. Vite 前端 ----------
$viteJs = Join-Path $root "frontend\node_modules\vite\bin\vite.js"
if (-not (Test-Path $viteJs)) {
    throw "Frontend dependencies are missing. Run 'pnpm install' (or 'npm install') in frontend first."
}
Write-Host "Starting Vue frontend with vite..." -ForegroundColor Cyan
$frontendProcess = Start-Process -FilePath $nodeExe `
    -WorkingDirectory (Join-Path $root "frontend") `
    -ArgumentList $viteJs, "--host", "127.0.0.1", "--port", "$frontendPort" `
    -WindowStyle Hidden `
    -PassThru
Set-Content -LiteralPath (Join-Path $runtime "frontend-dev.id") -Value $frontendProcess.Id -Encoding ascii
Wait-Until { Test-Http "http://127.0.0.1:$frontendPort" } "Vue frontend (vite)"

Write-Host ""
Write-Host "GM exam system is running (dev mode, no build)." -ForegroundColor Green
Write-Host "Frontend: http://127.0.0.1:$frontendPort"
Write-Host "Backend:  http://127.0.0.1:$backendPort"
Write-Host "Database: 127.0.0.1:13306 / yx_exam_affair_gm"
Write-Host "Username: admin"
Write-Host "Password: Admin123"
