<#
.SYNOPSIS
    Claude Usage Tray - Windows Installer

.DESCRIPTION
    Downloads and installs claude-usage-tray-go on Windows.
    Sets up autostart via Windows Registry.

.PARAMETER Uninstall
    Remove the application and autostart entry.

.EXAMPLE
    # Install (run in PowerShell as current user):
    irm https://raw.githubusercontent.com/mrmuminov/claude-usage-tray-go/master/install.ps1 | iex

    # Or run locally:
    .\install.ps1

    # Uninstall:
    .\install.ps1 -Uninstall
#>

param(
    [switch]$Uninstall
)

$ErrorActionPreference = "Stop"

$AppName = "claude-usage-tray-go"
$RepoOwner = "mrmuminov"
$RepoName = "claude-usage-tray-go"
$BinaryName = "claude-usage-tray-go.exe"
$RegistryKey = "HKCU:\Software\Microsoft\Windows\CurrentVersion\Run"
$RegistryValueName = "ClaudeUsageTray"

function Get-InstallDir {
    $localAppData = $env:LOCALAPPDATA
    if (-not $localAppData) {
        $localAppData = Join-Path $env:USERPROFILE "AppData\Local"
    }
    return Join-Path $localAppData $AppName
}

function Write-Status {
    param([string]$Message, [string]$Type = "Info")
    switch ($Type) {
        "Success" { Write-Host "[OK] $Message" -ForegroundColor Green }
        "Error"   { Write-Host "[ERROR] $Message" -ForegroundColor Red }
        "Warning" { Write-Host "[WARN] $Message" -ForegroundColor Yellow }
        "Info"    { Write-Host "[*] $Message" -ForegroundColor Cyan }
    }
}

function Get-LatestReleaseUrl {
    Write-Status "Fetching latest release info from GitHub..."
    $apiUrl = "https://api.github.com/repos/$RepoOwner/$RepoName/releases/latest"

    try {
        $release = Invoke-RestMethod -Uri $apiUrl -Headers @{ "User-Agent" = "claude-usage-tray-installer" }
    }
    catch {
        throw "Failed to fetch release info: $_"
    }

    $asset = $release.assets | Where-Object { $_.name -like "*windows*amd64*" } | Select-Object -First 1

    if (-not $asset) {
        throw "No Windows amd64 binary found in the latest release ($($release.tag_name))."
    }

    Write-Status "Found release: $($release.tag_name) - $($asset.name)"
    return $asset.browser_download_url
}

function Stop-RunningInstance {
    $process = Get-Process -Name ($BinaryName -replace '\.exe$', '') -ErrorAction SilentlyContinue
    if ($process) {
        Write-Status "Stopping running instance..." "Warning"
        $process | Stop-Process -Force
        Start-Sleep -Seconds 1
    }
}

function Install-App {
    $installDir = Get-InstallDir
    $binaryPath = Join-Path $installDir $BinaryName

    Write-Host ""
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "  Claude Usage Tray - Windows Installer" -ForegroundColor Cyan
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host ""

    # Stop any running instance
    Stop-RunningInstance

    # Get download URL
    $downloadUrl = Get-LatestReleaseUrl

    # Create install directory
    if (-not (Test-Path $installDir)) {
        New-Item -ItemType Directory -Path $installDir -Force | Out-Null
        Write-Status "Created directory: $installDir"
    }

    # Download binary
    Write-Status "Downloading binary..."
    $tempFile = Join-Path $env:TEMP "$BinaryName.tmp"
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile -UseBasicParsing
    }
    catch {
        throw "Failed to download binary: $_"
    }

    # Move binary to install directory
    Move-Item -Path $tempFile -Destination $binaryPath -Force
    Write-Status "Installed binary to: $binaryPath" "Success"

    # Set up autostart via registry
    try {
        Set-ItemProperty -Path $RegistryKey -Name $RegistryValueName -Value $binaryPath -Type String
        Write-Status "Autostart configured (runs on login)" "Success"
    }
    catch {
        Write-Status "Failed to configure autostart: $_" "Error"
    }

    # Launch the application
    Write-Host ""
    Write-Status "Launching $AppName..."
    Start-Process -FilePath $binaryPath -WindowStyle Hidden

    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "  Installation complete!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "  Install path : $binaryPath" -ForegroundColor Gray
    Write-Host "  Autostart    : Enabled" -ForegroundColor Gray
    Write-Host "  Status       : Running" -ForegroundColor Gray
    Write-Host ""
    Write-Host "  The app is now in your system tray." -ForegroundColor White
    Write-Host "  It will start automatically on login." -ForegroundColor White
    Write-Host ""
    Write-Host "  To uninstall:" -ForegroundColor Gray
    Write-Host "    .\install.ps1 -Uninstall" -ForegroundColor Yellow
    Write-Host "    # or: $binaryPath uninstall" -ForegroundColor Yellow
    Write-Host ""
}

function Uninstall-App {
    $installDir = Get-InstallDir
    $binaryPath = Join-Path $installDir $BinaryName

    Write-Host ""
    Write-Host "========================================" -ForegroundColor Yellow
    Write-Host "  Claude Usage Tray - Uninstaller" -ForegroundColor Yellow
    Write-Host "========================================" -ForegroundColor Yellow
    Write-Host ""

    # Stop any running instance
    Stop-RunningInstance

    # Remove autostart registry entry
    try {
        $regValue = Get-ItemProperty -Path $RegistryKey -Name $RegistryValueName -ErrorAction SilentlyContinue
        if ($regValue) {
            Remove-ItemProperty -Path $RegistryKey -Name $RegistryValueName -Force
            Write-Status "Autostart entry removed" "Success"
        }
        else {
            Write-Status "No autostart entry found (already clean)" "Info"
        }
    }
    catch {
        Write-Status "Failed to remove autostart entry: $_" "Warning"
    }

    # Remove binary
    if (Test-Path $binaryPath) {
        try {
            Remove-Item -Path $binaryPath -Force
            Write-Status "Binary removed: $binaryPath" "Success"
        }
        catch {
            Write-Status "Failed to remove binary (may be in use): $_" "Warning"
            Write-Status "You can delete it manually: $binaryPath" "Info"
        }
    }
    else {
        Write-Status "Binary not found at: $binaryPath" "Info"
    }

    # Remove install directory if empty
    if ((Test-Path $installDir) -and -not (Get-ChildItem $installDir)) {
        Remove-Item -Path $installDir -Force
        Write-Status "Removed empty directory: $installDir" "Success"
    }

    # Clean up cache
    $cacheDir = Join-Path $env:TEMP "claude"
    $cacheFile = Join-Path $cacheDir "tray-cache.json"
    if (Test-Path $cacheFile) {
        Remove-Item -Path $cacheFile -Force
        Write-Status "Cache file removed" "Success"
    }

    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "  Uninstall complete!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host ""
}

# Main
if ($Uninstall) {
    Uninstall-App
}
else {
    Install-App
}
