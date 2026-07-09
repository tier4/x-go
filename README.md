# x-go

[![GoDoc](https://godoc.org/github.com/tier4/x-go?status.svg)](https://godoc.org/github.com/tier4/x-go)
[![Test](https://github.com/tier4/x-go/actions/workflows/test.yml/badge.svg)](https://github.com/tier4/x-go/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tier4/x-go)](https://goreportcard.com/report/github.com/tier4/x-go)

Shared libraries used in the TIER IV. Use at your own risk. Breaking changes should be anticipated.

## Development

### Setup

This repository uses [mise](https://mise.jdx.dev/) to manage development tooling. The standard setup is:

```sh
mise install
```

`mise install` installs the pinned tools (currently `pre-commit`) and runs `pre-commit install` through the `postinstall` hook, so the local secret-scan hook is enabled in one step.

Pin your personal tool versions (Go, etc.) in `mise.local.toml`; the shared `mise.toml` intentionally manages only `pre-commit`.

> If your local git has `core.hooksPath` set, `pre-commit install` is refused and the hook is not installed automatically. In that case run `pre-commit run --all-files` manually (or unset `core.hooksPath`).

If you don't use mise, run the equivalent:

```sh
make setup
```

### Secret scan (gitleaks)

Layer 1 of the secret-scanning guardrail: a local [gitleaks](https://github.com/gitleaks/gitleaks) hook runs on `git commit` via [pre-commit](https://pre-commit.com/) to block secrets before they enter the history.

- Config: `.pre-commit-config.yaml`, `.gitleaks.toml` (extends the default ruleset with `useDefault = true`).
- Enabled automatically by `mise install` (or `make setup`).
- Run manually: `pre-commit run gitleaks --all-files`.

This local hook is best-effort and can be bypassed with `git commit --no-verify`. Effective enforcement is delegated to Layer 2 (CI required check / org ruleset) and GitHub Push Protection.
