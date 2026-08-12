<#
 .SYNOPSIS
    停止开发模式整套：go run 后端 + Vite 前端 + 本地 MySQL。
 .DESCRIPTION
    停止策略（多保险，确保真正停干净）：
      1. 读取 .runtime 下的 backend-dev.id / frontend-dev.id，杀掉记录的 PID；
      2. 按真实进程名匹配：gm_affair_main（go run 编译出的子进程）、gm-affair.exe（编译版）；
         node 启动 vite.js 的前端进程；
      3. 按端口兜底：强制释放 30000（后端）与 8848（前端）的占用进程；
      4. 关闭本地 MySQL（mysqladmin shutdown），并清理 mysql-process.id。
    对应启动脚本为 start-dev.ps1，二者配套使用。
#>

$ErrorActionPreference = "SilentlyContinue"

$root = Split-Path -Parent $PSScriptRoot
$runtime = Join-Path $root ".runtime"
$mysqlIni    = Join-Path $runtime "mysql.ini"
$mysqlAdmin  = Join-Path $runtime "mysql\bin\mysqladmin.exe"
$env:MYSQL_PWD = "123456"

$frontendPort = 8848
$backendPort  = 30000

function Stop-Pid {
    param([int]$Id, [string]$Label = "")
    if (-not $Id -or $Id -le 0) { return }
    try {
        $proc = Get-Process -Id $Id -ErrorAction SilentlyContinue
        if ($proc) {
            Write-Host ("Stopping PID {0} ({1}) : {2}" -f $Id, $Label, $proc.ProcessName) -ForegroundColor Gray
            $proc.Kill()
        }
    } catch { }
}

function Kill-ByExactName {
    param([string]$Name)
    $procs = @(Get-Process -Name $Name -ErrorAction SilentlyContinue)
    foreach ($p in $procs) {
        Stop-Pid -Id $p.Id -Label $Name
    }
    return $procs.Count
}

function Kill-ViteNode {
    $count = 0
    $nodes = @(Get-Process -Name node -ErrorAction SilentlyContinue)
    foreach ($n in $nodes) {
        $cl = ""
        try { $cl = $n.CommandLine } catch { $cl = "" }
        if ($cl -and ($cl -match 'vite[\\/]bin[\\/]vite\.js' -or $cl -match '--host')) {
            Stop-Pid -Id $n.Id -Label "vite(node)"
            $count++
        }
    }
    return $count
}

function Kill-ByPort {
    param([int]$Port)
    $count = 0
    try {
        $conns = @(Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue)
    } catch {
        $conns = @()
    }
    foreach ($c in $conns) {
        Stop-Pid -Id $c.OwningProcess -Label ("port {0}" -f $Port)
        $count++
    }
    return $count
}

# ---------- 1. 读取并杀掉记录在 id 文件里的 PID ----------
$root = Split-Path -Parent $PSScriptRoot
$runtime = Join-Path $root ".runtime"
$idFiles = @("backend-dev.id", "frontend-dev.id")
foreach ($name in $idFiles) {
    $idFile = Join-Path $runtime $name
    if (Test-Path -LiteralPath $idFile) {
        $pidVal = 0
        try { $pidVal = [int](Get-Content -LiteralPath $idFile -ErrorAction SilentlyContinue | Select-Object -First 1) } catch { $pidVal = 0 }
        Stop-Pid -Id $pidVal -Label $name
        Remove-Item -LiteralPath $idFile -Force -ErrorAction SilentlyContinue
    }
}

# ---------- 2. 按真实进程名 ----------
$byName = 0
$byName += Kill-ByExactName -Name "gm_affair_main"
$byName += Kill-ByExactName -Name "gm-affair"
$byName += Kill-ViteNode

# ---------- 3. 按端口兜底释放 ----------
$byPort = 0
$byPort += Kill-ByPort -Port $backendPort
$byPort += Kill-ByPort -Port $frontendPort

# 缓冲一下，确保子进程（go run 拉起的 gm_affair_main）也退出
Start-Sleep -Milliseconds 800

# 二次兜底：端口若仍被占，再杀一次
$byPort2 = 0
$byPort2 += Kill-ByPort -Port $backendPort
$byPort2 += Kill-ByPort -Port $frontendPort

# ---------- 4. 关闭本地 MySQL ----------
$mysqlStopped = $false
if (Test-Path -LiteralPath $mysqlAdmin) {
    try {
        & $mysqlAdmin --defaults-extra-file=$mysqlIni -uroot shutdown 2>$null
        Write-Host "MySQL shutdown requested." -ForegroundColor Gray
        $mysqlStopped = $true
    } catch { }
    # 兜底：按端口 13306 再杀 mysqld 进程
    $mysqlConns = @(Get-NetTCPConnection -LocalPort 13306 -State Listen -ErrorAction SilentlyContinue)
    foreach ($mc in $mysqlConns) {
        Stop-Pid -Id $mc.OwningProcess -Label "mysql(13306)"
    }
}
$mysqlIdFile = Join-Path $runtime "mysql-process.id"
if (Test-Path -LiteralPath $mysqlIdFile) {
    Remove-Item -LiteralPath $mysqlIdFile -Force -ErrorAction SilentlyContinue
}

Write-Host ("Dev processes stopped. (name-matched: {0}, port-released: {1}, mysql-stopped: {2})" -f $byName, ($byPort + $byPort2), $mysqlStopped) -ForegroundColor Yellow