$ErrorActionPreference = "SilentlyContinue"

$root = Split-Path -Parent $PSScriptRoot
$runtime = Join-Path $root ".runtime"
$env:MYSQL_PWD = "123456"

foreach ($name in @("frontend-process.id", "backend-process.id")) {
    $pidFile = Join-Path $runtime $name
    if (Test-Path $pidFile) {
        $processId = Get-Content -LiteralPath $pidFile
        Stop-Process -Id $processId -Force -ErrorAction SilentlyContinue
        Remove-Item -LiteralPath $pidFile -Force -ErrorAction SilentlyContinue
    }
}

$mysqlAdmin = Join-Path $runtime "mysql\bin\mysqladmin.exe"
$mysqlIni = Join-Path $runtime "mysql.ini"
if (Test-Path $mysqlAdmin) {
    & $mysqlAdmin --defaults-file=$mysqlIni -uroot shutdown 2>$null
}

Remove-Item -LiteralPath (Join-Path $runtime "mysql-process.id") -Force -ErrorAction SilentlyContinue
Write-Host "GM exam system local processes have stopped." -ForegroundColor Yellow
