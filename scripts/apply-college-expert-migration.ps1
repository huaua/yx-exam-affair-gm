$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
$mysql = Join-Path $root ".runtime\mysql\bin\mysql.exe"
$mysqlIni = Join-Path $root ".runtime\mysql.ini"
$migration = Join-Path $root "backend\scripts\update_260727_college_expert_management.sql"
$env:MYSQL_PWD = "123456"

if (-not (Test-Path $mysql)) {
    throw "Local MySQL client was not found."
}
if (-not (Test-Path $migration)) {
    throw "College expert management migration file was not found."
}

$command = "set MYSQL_PWD=123456&& `"$mysql`" --defaults-file=`"$mysqlIni`" -uroot yx_exam_affair_gm < `"$migration`""
& cmd.exe /d /c $command
if ($LASTEXITCODE -ne 0) {
    throw "Failed to apply the college expert management migration."
}

Write-Host "College expert management migration completed." -ForegroundColor Green
