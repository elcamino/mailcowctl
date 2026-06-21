# Contributing

Thanks for helping improve `mailcowctl`.

## Development Setup

Requirements:

- Go 1.22 or newer
- A shell environment that can run `go test ./...`
- Optional: a disposable mailcow instance for manual API checks

Start with:

```sh
go test ./...
go vet ./...
go build ./...
```

Unit tests must not require a real mailcow instance. Use `httptest.Server` for API behavior.

## Pull Requests

Before opening a pull request:

- Keep changes focused on one feature or fix.
- Add or update tests for changed behavior.
- Run `gofmt`, `go test ./...`, and `go vet ./...`.
- Update `README.md` or `INSTALL.md` when user-facing commands change.
- Do not log API keys, mailbox passwords, relay passwords, app passwords, or exported secrets.

## License

By contributing, you agree that your contribution is licensed under the same license as this
project: GPL-3.0-or-later.

All Go source files should keep the repository GPL header and SPDX line.
