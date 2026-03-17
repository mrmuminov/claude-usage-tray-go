# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| Latest release | Yes |
| Older versions | No |

## Reporting a Vulnerability

If you discover a security vulnerability, please report it responsibly:

1. **Do not** open a public GitHub issue
2. Email **darkshadeuz@gmail.com** with details
3. Include steps to reproduce if possible

You should receive a response within 48 hours. We will work with you to understand and address the issue before any public disclosure.

## Security Considerations

- OAuth tokens are read from the system keychain or local config files — they are never stored by this application
- API responses are cached locally in the system temp directory
- No data is sent to third parties; the app only communicates with `api.anthropic.com`
