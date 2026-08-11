$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
$runtime = Join-Path $root ".runtime"
$mysqlIni = Join-Path $runtime "mysql.ini"
$mysqlServer = Join-Path $runtime "mysql\bin\mysqld.exe"
$mysqlAdmin = Join-Path $runtime "mysql\bin\mysqladmin.exe"
$backendExe = Join-Path $runtime "bin\gm-affair.exe"
$nodeExe = (Get-Command node.exe -ErrorAction Stop).Source
$env:MYSQL_PWD = "123456"

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

if (-not (Test-Path $mysqlServer)) {
    throw "Local MySQL runtime was not found under .runtime\mysql."
}

if (-not (Test-MySql)) {
    $mysqlProcess = Start-Process -FilePath $mysqlServer -ArgumentList "--defaults-file=$mysqlIni" -WindowStyle Hidden -PassThru
    Set-Content -LiteralPath (Join-Path $runtime "mysql-process.id") -Value $mysqlProcess.Id -Encoding ascii
    Wait-Until { Test-MySql } "MySQL"
}

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

if (-not (Test-Path (Join-Path $root "frontend\node_modules\vite\bin\vite.js"))) {
    throw "Frontend dependencies are missing. Run pnpm install in frontend."
}

if (-not (Test-Http "http://127.0.0.1:8848")) {
    $frontendProcess = Start-Process -FilePath $nodeExe `
        -WorkingDirectory (Join-Path $root "frontend") `
        -ArgumentList "node_modules/vite/bin/vite.js", "--host", "127.0.0.1" `
        -WindowStyle Hidden `
        -PassThru
    Set-Content -LiteralPath (Join-Path $runtime "frontend-process.id") -Value $frontendProcess.Id -Encoding ascii
    Wait-Until { Test-Http "http://127.0.0.1:8848" } "Vue frontend"
}

Write-Host ""
Write-Host "GM exam system is running." -ForegroundColor Green
Write-Host "Frontend: http://127.0.0.1:8848"
Write-Host "Backend:  http://127.0.0.1:30000"
Write-Host "Database: 127.0.0.1:13306 / yx_exam_affair_gm"
Write-Host "Username: admin"
Write-Host "Password: Admin123"
