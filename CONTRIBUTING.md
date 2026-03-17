# Contributing to Claude Usage Tray

Thank you for your interest in contributing! Here's how to get started.

## Development Setup

1. **Prerequisites:** Go 1.25+, Git
2. **Linux only:** `sudo apt-get install libayatana-appindicator3-dev`

```bash
git clone https://github.com/mrmuminov/claude-usage-tray-go.git
cd claude-usage-tray-go
make build
make run
```

## How to Contribute

### Reporting Bugs

- Open an [issue](https://github.com/mrmuminov/claude-usage-tray-go/issues/new?template=bug_report.md)
- Include your OS, Go version, and steps to reproduce

### Suggesting Features

- Open an [issue](https://github.com/mrmuminov/claude-usage-tray-go/issues/new?template=feature_request.md)
- Describe the use case and expected behavior

### Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make your changes
4. Test on your platform: `make run`
5. Commit with a clear message: `git commit -m "feat: add your feature"`
6. Push and open a Pull Request

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` — new feature
- `fix:` — bug fix
- `refactor:` — code change that neither fixes a bug nor adds a feature
- `docs:` — documentation only
- `ci:` — CI/CD changes
- `chore:` — maintenance tasks

## Project Structure

```
├── main.go              # Entry point, systray UI
├── data.go              # API communication, caching
├── auth.go              # OAuth token resolution
├── auth_{platform}.go   # Platform-specific auth
├── format.go            # UI text formatting
├── icon_gen.go          # Dynamic icon generation
├── install.go           # Install/uninstall logic
├── install_{platform}.go # Platform-specific install
├── browser_{platform}.go # Platform-specific browser open
├── version.go           # Version string
├── install.ps1          # Windows PowerShell installer
├── fonts/               # Embedded font assets
└── .github/workflows/   # CI/CD pipelines
```

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Keep platform-specific code in `_darwin.go`, `_linux.go`, `_windows.go`, or `_other.go` files
- Use build tags (`//go:build`) for platform targeting

## Releasing

Releases are automated via GitHub Actions. To create a new release:

1. Update the changelog
2. Tag: `git tag v1.x.x`
3. Push: `git push origin v1.x.x`

The CI pipeline builds for Linux, macOS, and Windows and publishes to GitHub Releases.
