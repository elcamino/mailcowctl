# Security Policy

## Supported Versions

Security fixes are handled on the default branch. Tagged releases should be updated to the latest
available release once one is published.

## Reporting a Vulnerability

Do not open a public GitHub issue for security reports.

Email security reports to:

```text
tobias@scraperwall.com
```

Include:

- A concise description of the issue
- Affected command or code path
- Reproduction steps or a minimal proof of concept
- Whether credentials, passwords, or mail content can be exposed

## Security Expectations

`mailcowctl` handles mail administration credentials and migration data. Changes must preserve
these rules:

- API keys and passwords must not be printed in normal output.
- Structured output must not include extra logs on stdout.
- Auth errors should keep the mailcow IP allow-list and read-only-key hints.
- Tests must not require real production credentials.
