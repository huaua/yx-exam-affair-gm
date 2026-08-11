$ErrorActionPreference = "Stop"

# 项目根目录（脚本位于 <root>/scripts/ 下）
$root = Split-Path -Parent $PSScriptRoot
$bin  = Join-Path $root ".runtime/bin/gm-affair.exe"
$new  = Join-Path $root "backend/gm-affair-new.exe"

if (-not (Test-Path $new)) {
    Write-Host "未找到新版本 backend/gm-affair-new.exe，请先构建。" -ForegroundColor Red
    exit 1
}

# 1. 停服
Write-Host "正在停止旧后端进程..." -ForegroundColor Cyan
taskkill /F /IM gm-affair.exe 2>&1 | Out-Null

# 等待进程真正退出（最多 10 秒）
$waited = 0
while (Get-Process -Name "gm-affair" -ErrorAction SilentlyContinue) {
    Start-Sleep -Seconds 1
    $waited++
    if ($waited -ge 10) { break }
}

if (Get-Process -Name "gm-affair" -ErrorAction SilentlyContinue) {
    Write-Host "" -ForegroundColor Red
    Write-Host "无法终止 gm-affair.exe（权限不足）。请按以下步骤手动处理：" -ForegroundColor Red
    Write-Host "  1) 按 Ctrl+Shift+Esc 打开任务管理器" -ForegroundColor Yellow
    Write-Host "  2) 切到“详细信息”页，找到 gm-affair.exe，右键“结束任务”" -ForegroundColor Yellow
    Write-Host "  3) 重新运行本脚本： powershell -ExecutionPolicy Bypass -File `"$PSCommandPath`"" -ForegroundColor Yellow
    Write-Host "" -ForegroundColor Red
    exit 1
}

# 2. 备份旧 exe
$ts = Get-Date -Format "yyyyMMddHHmmss"
$backup = Join-Path $root ".runtime/bin/gm-affair-prev-$ts.exe"
Copy-Item $bin $backup -Force
Write-Host "已备份旧版本 -> $backup" -ForegroundColor Green

# 3. 替换为新 exe
Copy-Item $new $bin -Force
Write-Host "已替换为新版本。" -ForegroundColor Green

# 4. 启动新进程（必须在 backend/ 工作目录下运行，exe 读取 ./config/config.yaml）
Start-Process -FilePath $bin -WorkingDirectory (Join-Path $root "backend") -WindowStyle Hidden
Write-Host "已启动新后端进程。" -ForegroundColor Green

# 5. 健康检查（最多等待 20 秒）
$ok = $false
for ($i = 0; $i -lt 20; $i++) {
    try {
        $r = Invoke-WebRequest -Uri "http://127.0.0.1:30000/captcha" -UseBasicParsing -TimeoutSec 2
        if ($r.StatusCode -eq 200) { $ok = $true; break }
    } catch { }
    Start-Sleep -Seconds 1
}

if ($ok) {
    Write-Host "" -ForegroundColor Green
    Write-Host "部署成功！后端已用新版本运行，刷新浏览器即可看到所属学院名称。" -ForegroundColor Green
} else {
    Write-Host "后端已启动但健康检查未完成，请稍候在浏览器访问确认。" -ForegroundColor Yellow
}
