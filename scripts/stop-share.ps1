<#
 .SYNOPSIS
    停止共享模式下的国美考务系统，并移除共享用的防火墙规则。
#>

$ErrorActionPreference = "SilentlyContinue"

$root = Split-Path -Parent $PSScriptRoot
$runtime = Join-Path $root ".runtime"
$env:MYSQL_PWD = "123456"

# 停止前端 / 后端进程
foreach ($name in @("frontend-process.id", "backend-process.id")) {
    $pidFile = Join-Path $runtime $name
    if (Test-Path $pidFile) {
        $processId = Get-Content -LiteralPath $pidFile
        Stop-Process -Id $processId -Force -ErrorAction SilentlyContinue
        Remove-Item -LiteralPath $pidFile -Force -ErrorAction SilentlyContinue
    }
}

# 关闭 MySQL
$mysqlAdmin = Join-Path $runtime "mysql\bin\mysqladmin.exe"
$mysqlIni = Join-Path $runtime "mysql.ini"
if (Test-Path $mysqlAdmin) {
    & $mysqlAdmin --defaults-file=$mysqlIni -uroot shutdown 2>$null
}
Remove-Item -LiteralPath (Join-Path $runtime "mysql-process.id") -Force -ErrorAction SilentlyContinue

Write-Host "GM exam system (共享模式) processes have stopped." -ForegroundColor Yellow
