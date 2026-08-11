$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
$mysql = Join-Path $root ".runtime\mysql\bin\mysql.exe"
$mysqlIni = Join-Path $root ".runtime\mysql.ini"
$migration = Join-Path $root "backend\scripts\update_260728_expert_ban_management.sql"
$env:MYSQL_PWD = "123456"

$command = "set MYSQL_PWD=123456&& `"$mysql`" --defaults-file=`"$mysqlIni`" -uroot yx_exam_affair_gm < `"$migration`""
& cmd.exe /d /c $command
if ($LASTEXITCODE -ne 0) {
    throw "Failed to apply the expert ban management migration."
}

Write-Host "Expert ban management migration completed." -ForegroundColor Green
