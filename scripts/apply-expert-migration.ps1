$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
$runtime = Join-Path $root ".runtime"
$mysql = Join-Path $runtime "mysql\bin\mysql.exe"
$mysqlIni = Join-Path $runtime "mysql.ini"
$migration = Join-Path $root "backend\scripts\update_260727_expert_management.sql"
$env:MYSQL_PWD = "123456"

if (-not (Test-Path $mysql)) {
    throw "Local MySQL client was not found."
}
if (-not (Test-Path $migration)) {
    throw "Expert management migration file was not found."
}

$command = "set MYSQL_PWD=123456&& `"$mysql`" --defaults-file=`"$mysqlIni`" -uroot yx_exam_affair_gm < `"$migration`""
& cmd.exe /d /c $command
if ($LASTEXITCODE -ne 0) {
    throw "Failed to apply the expert management migration."
}

Write-Host "Expert management migration completed." -ForegroundColor Green
