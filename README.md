# mailcowctl

[![CI](https://github.com/elcamino/mailcowctl/actions/workflows/ci.yml/badge.svg)](https://github.com/elcamino/mailcowctl/actions/workflows/ci.yml)
[![License: GPL v3 or later](https://img.shields.io/badge/License-GPLv3%2B-blue.svg)](LICENSE.md)
[![Go Reference](https://pkg.go.dev/badge/github.com/elcamino/mailcowctl.svg)](https://pkg.go.dev/github.com/elcamino/mailcowctl)
[![Go Report Card](https://goreportcard.com/badge/github.com/elcamino/mailcowctl)](https://goreportcard.com/report/github.com/elcamino/mailcowctl)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/elcamino/mailcowctl/badge)](https://securityscorecards.dev/viewer/?uri=github.com/elcamino/mailcowctl)
[![Go Version](https://img.shields.io/github/go-mod/go-version/elcamino/mailcowctl)](go.mod)
[![Security Policy](https://img.shields.io/badge/security-policy-brightgreen)](SECURITY.md)

`mailcowctl` is a command-line client for managing a running mailcow instance through its REST API.
It covers domains, mailboxes, aliases, and migration-adjacent mailcow resources from scripts or a terminal.

## Install

See [INSTALL.md](INSTALL.md) for release binaries, `go install`, and source builds.

## Configuration

Connection settings are resolved in this order:

1. Flags: `--host`, `--api-key`, `--profile`
2. Environment: `MAILCOW_HOST`, `MAILCOW_API_KEY`, `MAILCOW_PROFILE`
3. Config file: `$XDG_CONFIG_HOME/mailcowctl/config.yaml` or `~/.config/mailcowctl/config.yaml`

Example config:

```yaml
current_profile: prod
profiles:
  prod:
    host: https://mail.example.org
    api_key: 000000-000000-000000-000000-000000
```

mailcow uses the `X-API-Key` header. API keys are IP allow-listed in mailcow, and write commands
need a read-write key.

### mailcow-dockerized API key

On a `mailcow-dockerized` host, set the API key and allowed source IPs in `mailcow.conf`
(or in `.env` if running mailcow-dockerized):

```sh
API_KEY=$(openssl rand -hex 32)
API_ALLOW_FROM=127.0.0.1,IP1,IP2,...
```

`API_ALLOW_FROM` must include every client or CI runner source IP that will call the API. After
changing `mailcow.conf`, restart or recreate the affected mailcow containers according to the
mailcow-dockerized update/restart procedure. Use the same generated key with `mailcowctl`:

```sh
export MAILCOW_HOST=https://mail.example.org
export MAILCOW_API_KEY="$API_KEY"
```

## Examples

```sh
mailcowctl --host https://mail.example.org --api-key "$MAILCOW_API_KEY" domain list
mailcowctl domain create example.org --description Example --mailboxes 10 --aliases 400 --quota 10240
mailcowctl mailbox create me@example.org --name "Me" --quota 3072 --password-stdin
mailcowctl alias create info@example.org --goto me@example.org
mailcowctl alias delete info@example.org --yes
mailcowctl domain list -o json
```

Quotas are MiB. Passwords should be supplied with `--password-stdin` or `--password-env`; the CLI
does not echo them.

## Commands

```text
mailcowctl domain list
mailcowctl domain get <domain>
mailcowctl domain create <domain>
mailcowctl domain edit <domain>
mailcowctl domain delete <domain> --yes

mailcowctl mailbox list [--domain example.org]
mailcowctl mailbox get <address>
mailcowctl mailbox create <address> --name "Full Name" (--password-stdin | --password-env VAR)
mailcowctl mailbox edit <address>
mailcowctl mailbox delete <address> --yes

mailcowctl alias list [--domain example.org]
mailcowctl alias get <address|id>
mailcowctl alias create <address> (--goto target@example.org | --discard | --to-spam | --to-ham)
mailcowctl alias edit <address|id>
mailcowctl alias delete <address|id> --yes
```

`alias edit` and `alias delete` accept an address and resolve it to the numeric mailcow alias ID
before calling the API.

## Mail Migration

- `dkim get|add|duplicate|delete` — DKIM keys. Private keys are NEVER returned by
  the mailcow API; `dkim get` shows the public record and (with `--dns`) the DNS
  TXT line. Migrating DKIM means regenerating keys on the new server and updating
  each domain's DNS TXT record. `dkim duplicate` only works within one server.
- `syncjob list|get|create|edit|delete` — mailcow's built-in imapsync jobs; the
  mechanism to pull mail content from the old server. mailcow runs them on the
  configured interval (there is no run-now API).
- `apppass list|get|create|edit|delete` — app passwords. Existing secrets are
  hashed and cannot be exported; recreate them with new secrets.
- `filter list|get|create|edit|delete` — sieve filters. `filter get -o json`
  exports the full script for re-creation on the new server.
- `policy list|add|delete` — domain anti-spam allow (wl) / block (bl) lists.

## Routing and Relay

All support list/get/create/delete (no edit -- the mailcow API does not expose
edit for these resources).

- `transport list|get|create|delete` -- outbound transport maps. Passwords are
  masked in table output; use `-o json` to export full credentials.
- `bcc list|get|create|delete` -- BCC maps (`--type sender|rcpt`).
- `recipient list|get|create|delete` -- recipient rewrite maps.
- `tlspolicy list|get|create|delete` -- per-destination TLS policy maps.
- `relayhost list|get|create|delete` -- sender-dependent relay hosts. Passwords
  are masked in table output; `-o json` exports them (the API returns them in
  cleartext).
- `fwdhost list|get|create|delete` -- forwarding hosts. These are keyed by host
  string, not a numeric id: `get` and `delete` take the host.
- `resource list|get|create|delete` -- calendar resources (`--kind location|group|thing`).

## Quarantine Operations

- `quarantine list [--rcpt <addr>]` -- list held mail, optionally filtered by recipient.
- `quarantine get <id>` -- show one quarantined item.
- `quarantine release <id>` -- deliver a held message.
- `quarantine learn-ham <id>` -- deliver a held message and train the spam filter.
- `quarantine delete <id>` -- purge a held item (requires --yes).

Quarantine is operational state, not migration state -- it is not part of a
server-to-server migration. Release and delete act on live mail; the `qitem`
item-id encoding should be confirmed against a single disposable item before
bulk use.

## Development

```sh
go test ./...
go vet ./...
go build ./...
```

## Contributing and Security

Contributions are welcome. Read [CONTRIBUTING.md](CONTRIBUTING.md) before sending changes.

Report vulnerabilities privately using the process in [SECURITY.md](SECURITY.md). Do not open public
issues for security-sensitive findings.

## License

`mailcowctl` is licensed under the GNU General Public License, version 3 or later. See
[LICENSE.md](LICENSE.md).
