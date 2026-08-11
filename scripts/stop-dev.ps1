<#
 .SYNOPSIS
    停止开发模式进程：go run 后端 + Vite 前端。
 .DESCRIPTION
    不依赖 PID 文件，按进程命令行特征匹配并终止：
      - Go 后端：命令行包含 gm_affair_main.go
      - Vite 前端：命令行包含 vite 且工作目录为本项目 frontend
    不关闭本地 MySQL。
#>

$ErrorActionPreference = "SilentlyContinue"

function Stop-ByCommandLine {
    param(
        [Parameter(Mandatory = $true)]
        [string]$MatchPattern
    )

    $matched = Get-Process | Where-Object {
        $_.CommandLine -and ($_.CommandLine -match $MatchPattern)
    }

    $count = 0
    foreach ($p in $matched) {
        try {
            Write-Host "Stopping PID $($p.Id) : $($p.ProcessName)" -ForegroundColor Gray
            $p.Kill()
            $count++
        } catch {
            # ignore
        }
    }

    return $count
}

$backendCount = Stop-ByCommandLine -MatchPattern "gm_affair_main\.go"
$frontendCount = Stop-ByCommandLine -MatchPattern "vite.*--host|node_modules[\\/]vite[\\/]bin[\\/]vite\.js"

$root = Split-Path -Parent $PSCommandPath
$runtime = Join-Path $root ".runtime"
foreach ($name in @("backend-dev.id", "frontend-dev.id")) {
    $idFile = Join-Path $runtime $name
    if (Test-Path $idFile) {
        Remove-Item -LiteralPath $idFile -Force -ErrorAction SilentlyContinue
    }
}

Write-Host "Dev processes stopped. (backend: $backendCount, frontend: $frontendCount)" -ForegroundColor Yellow
Write-Host "MySQL is left running. To stop MySQL too, use stop-local.ps1." -ForegroundColor Cyan
