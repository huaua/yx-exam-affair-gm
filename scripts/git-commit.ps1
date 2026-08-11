# git-commit.ps1 - one-click commit script
# Usage:
#   .\scripts\git-commit.ps1                 # commit with default message (date)
#   .\scripts\git-commit.ps1 "fix bug"       # commit with custom message
#   .\scripts\git-commit.ps1 "msg" -Push     # commit and push to remote (needs remote)
param(
    [string]$Message = "",
    [switch]$Push
)

$ErrorActionPreference = "Stop"
$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
Set-Location $repoRoot

try { $gitVer = & git --version 2>&1 } catch { Write-Error "git not found in PATH"; exit 1 }
Write-Host "git: $gitVer"

& git add -A
$status = & git status --porcelain
if (-not $status) {
    Write-Host "No changes to commit." -ForegroundColor Yellow
    exit 0
}

if ([string]::IsNullOrWhiteSpace($Message)) {
    $date = Get-Date -Format "yyyy-MM-dd"
    $Message = "chore: auto commit $date"
}

& git commit -m $Message
if ($LASTEXITCODE -ne 0) { Write-Error "commit failed"; exit $LASTEXITCODE }
Write-Host "committed: $Message" -ForegroundColor Green

if ($Push) {
    $remote = & git remote get-url origin 2>$null
    if (-not $remote) {
        Write-Warning "No remote 'origin' configured, skip push. Run: git remote add origin <repo-url>"
        exit 0
    }
    $branch = & git branch --show-current
    & git push origin $branch
    if ($LASTEXITCODE -eq 0) { Write-Host "pushed to $remote ($branch)" -ForegroundColor Green }
    else { Write-Error "push failed"; exit $LASTEXITCODE }
} else {
    Write-Host "Local commit done. Add -Push to push to remote." -ForegroundColor Cyan
}
