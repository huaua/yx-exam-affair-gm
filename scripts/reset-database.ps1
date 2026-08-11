$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
$runtime = Join-Path $root ".runtime"
$mysql = Join-Path $runtime "mysql\bin\mysql.exe"
$mysqlIni = Join-Path $runtime "mysql.ini"
$sourceSql = Join-Path $root "国美考务-源码+演示库\yx_exam_affair_gm_demo-20260721.sql"
$runtimeSql = Join-Path $runtime "demo.sql"
$expertMigration = Join-Path $root "backend\scripts\update_260727_expert_management.sql"
$expertBanMigration = Join-Path $root "backend\scripts\update_260728_expert_ban_management.sql"
$env:MYSQL_PWD = "123456"

if (-not (Test-Path $mysql)) {
    throw "Local MySQL client was not found."
}

Copy-Item -LiteralPath $sourceSql -Destination $runtimeSql -Force

& $mysql --defaults-file=$mysqlIni -uroot -e "DROP DATABASE IF EXISTS yx_exam_affair_gm; CREATE DATABASE yx_exam_affair_gm CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;"
if ($LASTEXITCODE -ne 0) {
    throw "Failed to create the database."
}

$command = "set MYSQL_PWD=123456&& `"$mysql`" --defaults-file=`"$mysqlIni`" -uroot yx_exam_affair_gm < `"$runtimeSql`""
& cmd.exe /d /c $command
if ($LASTEXITCODE -ne 0) {
    throw "Failed to import the demo database."
}

$migrationCommand = "set MYSQL_PWD=123456&& `"$mysql`" --defaults-file=`"$mysqlIni`" -uroot yx_exam_affair_gm < `"$expertMigration`""
& cmd.exe /d /c $migrationCommand
if ($LASTEXITCODE -ne 0) {
    throw "Failed to apply the expert management migration."
}

$expertBanMigrationCommand = "set MYSQL_PWD=123456&& `"$mysql`" --defaults-file=`"$mysqlIni`" -uroot yx_exam_affair_gm < `"$expertBanMigration`""
& cmd.exe /d /c $expertBanMigrationCommand
if ($LASTEXITCODE -ne 0) {
    throw "Failed to apply the expert ban management migration."
}

$adminHash = '$2a$14$r/gT5wOy/KlkYic7LT5Qreos5c1lEJcJ9k8d68PD/J7waZmRuUxnC'
& $mysql --defaults-file=$mysqlIni -uroot -e "UPDATE yx_exam_affair_gm.sys_user SET password='$adminHash', state='1', deleted=1 WHERE user_name='admin';"
if ($LASTEXITCODE -ne 0) {
    throw "Failed to restore the local admin password."
}

Write-Host "Demo database reset completed. Login with admin / Admin123." -ForegroundColor Green
