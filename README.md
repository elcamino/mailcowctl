# mailcowctl

`mailcowctl` is a command-line client for managing a running mailcow instance through its REST API.
It covers domains, mailboxes, and aliases.

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

## Development

```sh
go test ./...
go vet ./...
go build ./...
```
