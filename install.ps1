# Claude Usage Tray - Windows Installer
#
# Install:
#   irm https://raw.githubusercontent.com/mrmuminov/claude-usage-tray-go/master/install.ps1 | iex
#
# Or run from the repo directory (uses local binary if available):
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

# Check for local binary in current directory
$localBinary = Join-Path $PSScriptRoot $BinaryName
if (-not $PSScriptRoot) {
    $localBinary = Join-Path (Get-Location) $BinaryName
}

if (Test-Path $localBinary) {
    # Use local binary from repo
    Write-Host "  [*] Using local binary: $localBinary" -ForegroundColor Cyan

    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    Copy-Item -Path $localBinary -Destination $BinaryPath -Force
    Write-Host "  [OK] Installed to: $BinaryPath" -ForegroundColor Green
}
else {
    # Download from GitHub releases
    Write-Host "  [*] No local binary found, downloading from GitHub..." -ForegroundColor Cyan
    $apiUrl = "https://api.github.com/repos/$RepoOwner/$RepoName/releases/latest"
    $release = Invoke-RestMethod -Uri $apiUrl -Headers @{ "User-Agent" = "claude-tray-installer" }
    $asset = $release.assets | Where-Object { $_.name -like "*windows*amd64*" } | Select-Object -First 1

    if (-not $asset) {
        Write-Host "  [ERROR] No Windows binary found in release $($release.tag_name)" -ForegroundColor Red
        return
    }

    Write-Host "  [*] Found: $($release.tag_name)" -ForegroundColor Cyan

    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

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
