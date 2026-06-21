# Installing mailcowctl

## Release Binaries

When releases are published, download the archive for your platform from:

```text
https://github.com/elcamino/mailcowctl/releases
```

Install the binary somewhere on your `PATH`, for example:

```sh
install -m 0755 mailcowctl /usr/local/bin/mailcowctl
```

## Go Install

With Go 1.22 or newer:

```sh
go install github.com/elcamino/mailcowctl@latest
```

This installs `mailcowctl` into `$(go env GOPATH)/bin`.

## Build From Source

```sh
git clone https://github.com/elcamino/mailcowctl.git
cd mailcowctl
go test ./...
go build -o mailcowctl .
```

## Configuration

The CLI reads connection settings from flags, environment variables, or
`$XDG_CONFIG_HOME/mailcowctl/config.yaml`.

```sh
export MAILCOW_HOST=https://mail.example.org
export MAILCOW_API_KEY=...
mailcowctl domain list
```

mailcow API keys are IP allow-listed in the mailcow UI. Write operations require a read-write key.
