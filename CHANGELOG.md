# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added
- CONTRIBUTING.md, CODE_OF_CONDUCT.md, SECURITY.md
- GitHub issue and PR templates
- Dependabot configuration for automated dependency updates
- .editorconfig for consistent code style

### Changed
- README.md updated with badges, quick install section, and contributing info

## [1.2.0] - 2025-06-01

### Added
- Cross-platform auto-install and autostart support
- MIT License
- Windows PowerShell one-liner installer (`install.ps1`)
- `install`, `uninstall`, `status` CLI commands

## [1.1.0] - 2025-05-15

### Added
- Improved tray icon quality with dynamic font rendering
- GitHub link in tray menu

### Changed
- Icon generation uses embedded Ubuntu font with OpenType rendering

## [1.0.1] - 2025-05-10

### Fixed
- GitHub Actions: removed deprecated runners, use latest OS versions
- `apt-get update` before installing dependencies

## [1.0.0] - 2025-05-01

### Added
- Initial release
- 5-hour and 7-day rate limit display
- Dynamic color-coded tray icon (green/yellow/red)
- Auto-refresh every 60 seconds
- Disk cache with stale fallback
- OAuth token resolution chain
- Cross-platform support (Linux, macOS, Windows)

[Unreleased]: https://github.com/mrmuminov/claude-usage-tray-go/compare/v1.2.0...HEAD
[1.2.0]: https://github.com/mrmuminov/claude-usage-tray-go/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/mrmuminov/claude-usage-tray-go/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/mrmuminov/claude-usage-tray-go/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/mrmuminov/claude-usage-tray-go/releases/tag/v1.0.0
