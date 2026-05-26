---
name: quail-cli
description: Maintain and extend the Quail CLI Go repository. Use when Codex works on quail-cli commands, API client methods, OAuth or API key authentication, config behavior, reader or author workflows, MCP tools, README command docs, release tags, or any change that needs to stay compatible with Quail API and quaily-server routes.
---

# Quail CLI

## Core Rules

Keep changes small and tied to a real user workflow. Prefer existing command and client patterns over new abstractions.

Do not add a CLI command unless the server API exists or the same change also adds the required server route and authorization. Do not invent unsupported operations.

Never print secrets. Do not write machine-specific paths, usernames, or local environment details into docs, code, examples, commits, or release notes.

## Repo Map

- `cmd/<name>/`: Cobra command implementations. Each package should expose `NewCmd() *cobra.Command`.
- `cmd/root.go`: global flags, config loading, auth selection, command registration, shared context values.
- `cmd/common`: context keys and output format constants.
- `client`: HTTP API methods and response structs.
- `util/login.go`: OAuth and API key login persistence.
- `oauth`: OAuth flow and token refresh.
- `mcp`: MCP server, tools, and resources.
- `docs/feat`: feature notes, named `YYYY-MM-DD_feat_name.md`.

## Command Pattern

Add a command package under `cmd/<domain>` when the domain is new. Register it in `cmd/root.go`.

Inside command handlers, get shared dependencies from context:

```go
cl := cmd.Context().Value(common.CTX_CLIENT{}).(*client.Client)
format := cmd.Context().Value(common.CTX_FORMAT{}).(string)
```

For required flags or args, show `cmd.Help()` and return. For normal API failures, print or log the concrete error and return; do not exit from deep helper functions.

Support `--format json` for new read or mutation commands. Human output should be compact and scan-friendly, usually with `tabwriter` for lists.

## Client Pattern

Put endpoint-specific HTTP calls in `client/*.go`. Keep commands responsible for CLI parsing and printing, not URL building.

Use `Client.sendRequest(method, url, payload)` unless the endpoint needs different status handling. Decode JSON into typed response structs in `client/resp.go` when the response is reused or user-facing.

Accept list and post identifiers as strings when the API supports either id or slug. Use `uint64` only when the API requires numeric ids.

## Auth And Config

Preserve current auth precedence:

1. `QUAIL_API_KEY`
2. `app.api_key` in config
3. OAuth access token and refresh token

Do not force OAuth login for `login` itself or for scriptable API-key usage. Keep `login --api-key` usable with an interactive hidden prompt and with an explicit value.

When adding routes that should work with API keys, update server-side API key authorization as well as OAuth scope routes when needed.

## Reader And Author Scope

Reader commands should model actions a subscriber can perform:

- list subscriptions
- list subscribed posts
- read a post by `--list/--post`
- read a standard `https://quaily.com/{list_slug}/{post_slug}` URL
- list post comments
- create a post comment with `--content`

Author commands should model list owner or cooperator actions:

- list latest comments across owned lists
- list comments for one list by id or slug
- approve, reject, mark spam, or delete a comment

Keep custom domain URL parsing out of `reader read` until the server/API contract makes it reliable.

## Server Contract

If a CLI feature needs quaily-server changes, inspect and patch the server in the same stage. Typical files are:

- route registration
- OAuth scope route map
- API key route allowlist
- handler logic that resolves list id or slug

Run targeted server tests when touching server behavior. If a broader server package fails because of pre-existing unrelated issues, report the exact failing package and reason.

## Validation

For CLI changes, run:

```bash
go test ./...
```

For command wiring, also run focused help or dry-run checks with a disposable config file, for example:

```bash
QUAIL_API_KEY=QK-test go run . --config <temp-config.yaml> reader --help
```

For auth changes, test both config-free API key behavior and the affected command help path. Avoid tests that require a real user token unless the user explicitly provides one.

## Release Work

Before tagging, make sure the worktree is clean, commits are pushed, and the tag version matches the size of the change. Use annotated tags for releases:

```bash
git tag -a vX.Y.Z -m "Release vX.Y.Z"
git push origin vX.Y.Z
```
