# Claude Usage Tray - Windows Installer
#
# Install (one-liner):
#   irm https://raw.githubusercontent.com/mrmuminov/claude-usage-tray-go/master/install.ps1 | iex
#
# Install from repo (builds from source if Go is available):
#   .\install.ps1
#
# Uninstall:
#   claude-usage-tray-go.exe uninstall

$ErrorActionPreference = "Stop"

$AppName = "claude-usage-tray-go"
$RepoOwner = "mrmuminov"
$RepoName = "claude-usage-tray-go"
$BinaryName = "claude-usage-tray-go.exe"
$RegistryKey = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Run"
$RegistryValueName = "ClaudeUsageTray"

$localAppData = if ($env:LOCALAPPDATA) { $env:LOCALAPPDATA } else { Join-Path $env:USERPROFILE "AppData\Local" }
$InstallDir = Join-Path $localAppData $AppName
$BinaryPath = Join-Path $InstallDir $BinaryName

Write-Host ""
Write-Host "  Claude Usage Tray - Windows Installer" -ForegroundColor Cyan
Write-Host "  ======================================" -ForegroundColor Cyan
Write-Host ""

# Check if already installed
if (Test-Path $BinaryPath) {
    $answer = Read-Host "  Already installed. Reinstall? [y/N]"
    if ($answer -ne "y" -and $answer -ne "Y" -and $answer -ne "yes") {
        Write-Host "  Cancelled." -ForegroundColor Yellow
        return
    }
    Write-Host ""
}

# Stop running instance if any
$proc = Get-Process -Name ($BinaryName -replace '\.exe$', '') -ErrorAction SilentlyContinue
if ($proc) {
    Write-Host "  [*] Stopping running instance..." -ForegroundColor Yellow
    $proc | Stop-Process -Force
    Start-Sleep -Seconds 2
}

if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

# Determine script directory
$scriptDir = if ($PSScriptRoot) { $PSScriptRoot } else { (Get-Location).Path }
$localBinary = Join-Path $scriptDir $BinaryName
$goModFile = Join-Path $scriptDir "go.mod"

if (Test-Path $localBinary) {
    # 1) Local binary exists — use it directly
    Write-Host "  [*] Using local binary: $localBinary" -ForegroundColor Cyan
    Copy-Item -Path $localBinary -Destination $BinaryPath -Force
    Write-Host "  [OK] Installed to: $BinaryPath" -ForegroundColor Green
}
elseif ((Test-Path $goModFile) -and (Get-Command go -ErrorAction SilentlyContinue)) {
    # 2) Source code + Go available — build from source
    Write-Host "  [*] Building from source..." -ForegroundColor Cyan
    $version = "dev"
    try {
        $version = (git -C $scriptDir describe --tags --always 2>$null)
        if (-not $version) { $version = "dev" }
    } catch { }

    $buildOutput = Join-Path $scriptDir $BinaryName
    & go build -C $scriptDir -ldflags="-X main.Version=$version" -o $buildOutput .
    if ($LASTEXITCODE -ne 0) {
        Write-Host "  [ERROR] Build failed" -ForegroundColor Red
        return
    }
    Write-Host "  [OK] Built successfully ($version)" -ForegroundColor Green

    Copy-Item -Path $buildOutput -Destination $BinaryPath -Force
    Remove-Item -Path $buildOutput -Force
    Write-Host "  [OK] Installed to: $BinaryPath" -ForegroundColor Green
}
else {
    # 3) No local binary, no Go — download from GitHub
    Write-Host "  [*] Downloading from GitHub..." -ForegroundColor Cyan
    $apiUrl = "https://api.github.com/repos/$RepoOwner/$RepoName/releases/latest"
    $release = Invoke-RestMethod -Uri $apiUrl -Headers @{ "User-Agent" = "claude-tray-installer" }
    $asset = $release.assets | Where-Object { $_.name -like "*windows*amd64*" } | Select-Object -First 1

    if (-not $asset) {
        Write-Host "  [ERROR] No Windows binary found in release $($release.tag_name)" -ForegroundColor Red
        return
    }

    Write-Host "  [*] Found: $($release.tag_name)" -ForegroundColor Cyan
    Write-Host "  [*] Downloading..." -ForegroundColor Cyan
    $tempFile = Join-Path $env:TEMP "$BinaryName.tmp"
    Invoke-WebRequest -Uri $asset.browser_download_url -OutFile $tempFile -UseBasicParsing
    Move-Item -Path $tempFile -Destination $BinaryPath -Force
    Write-Host "  [OK] Installed to: $BinaryPath" -ForegroundColor Green
}

# Autostart
Set-ItemProperty -Path $RegistryKey -Name $RegistryValueName -Value $BinaryPath -Type String
Write-Host "  [OK] Autostart enabled" -ForegroundColor Green

# Launch
Start-Process -FilePath $BinaryPath -WindowStyle Hidden
Write-Host "  [OK] Running in system tray" -ForegroundColor Green

Write-Host ""
Write-Host "  Done! Check your system tray." -ForegroundColor Green
Write-Host "  To uninstall: $BinaryPath uninstall" -ForegroundColor Gray
Write-Host ""
