---
name: quail-cli
description: Help users operate Quaily from the command line with quail-cli. Use when the user wants to log in with OAuth or API key, configure quail-cli, read subscriptions and posts, read or write comments, manage author comments, publish or manage posts, or output JSON for scripts.
---

# Quail CLI

## What To Do

Help the user use `quail-cli` to work with Quaily. Give commands they can run, explain required flags, and keep defaults simple.

Default output is human-readable. Use `--json` only when the user asks for JSON, scripting, piping, or `jq`.

Never ask the user to paste an API key or token into chat. Tell them to use `quail-cli login --api-key` or `QUAIL_API_KEY` in their shell.

## Install

For macOS and Linux, install the latest release with:

```bash
curl -fsSL https://raw.githubusercontent.com/quailyquaily/quail-cli/master/scripts/install.sh | bash
```

For Windows PowerShell, install the latest release with:

```powershell
irm https://raw.githubusercontent.com/quailyquaily/quail-cli/master/scripts/install.ps1 | iex
```

If the user does not want to run a remote install script, point them to the official GitHub releases page:

```text
https://github.com/quailyquaily/quail-cli/releases
```

Tell them to download the binary for their operating system, extract it, and put the `quail-cli` binary in `PATH`.

If the user has Go installed, they can install from source:

```bash
go install github.com/quailyquaily/quail-cli@latest
```

Verify installation:

```bash
quail-cli version
```

If `quail-cli` is not found, the binary is not in `PATH`. Explain how to fix `PATH` for the user's shell only when they ask for OS-specific help.

## Authentication

For normal interactive use:

```bash
quail-cli login
```

For API key login:

```bash
quail-cli login --api-key
```

For scripts and temporary sessions:

```bash
QUAIL_API_KEY=QK-... quail-cli me
```

Auth priority is:

1. `QUAIL_API_KEY`
2. saved `app.api_key`
3. saved OAuth token

## Global Options

Use these options before or after the command:

```bash
quail-cli --json me
quail-cli --api-base https://api.quail.ink me
quail-cli --auth-base https://quaily.com login
quail-cli --config ./config.yaml me
```

Do not use `--format`; `--json` is the supported JSON switch.

## Reader Tasks

List subscriptions:

```bash
quail-cli reader subscriptions
```

List posts from subscribed lists:

```bash
quail-cli reader posts --limit 20 --offset 0
```

Read a post by Quaily URL:

```bash
quail-cli reader read https://quaily.com/list-slug/post-slug
```

Read a post by list and post id or slug:

```bash
quail-cli reader read --list list-slug --post post-slug
```

Read comments for a post:

```bash
quail-cli reader comments --post 123
```

Write a comment:

```bash
quail-cli reader comment --post 123 --content "Thanks for the post."
```

`reader read <URL>` supports standard `https://quaily.com/{list_slug}/{post_slug}` URLs. Do not assume custom domains are supported.

## Author Tasks

List latest comments across the user's lists:

```bash
quail-cli comments latest --limit 50
```

List comments for one list:

```bash
quail-cli comments list --list list-slug --limit 20 --offset 0
```

Moderate comments:

```bash
quail-cli comments approve 123
quail-cli comments reject 123
quail-cli comments spam 123
quail-cli comments delete 123
```

Use list id or list slug when a command accepts `--list`.

## Post Tasks

Create or update a post from Markdown frontmatter:

```bash
quail-cli post upsert post.md --list list-slug
```

Publish while upserting:

```bash
quail-cli post upsert post.md --list list-slug --publish
```

Operate on an existing post:

```bash
quail-cli post publish --list list-slug --post post-slug
quail-cli post unpublish --list list-slug --post post-slug
quail-cli post deliver --list list-slug --post post-slug
quail-cli post delete --list list-slug --post post-slug
```

## JSON Output

Use `--json` for automation:

```bash
quail-cli --json reader subscriptions
quail-cli --json reader read https://quaily.com/list-slug/post-slug
quail-cli --json comments latest --limit 10
```

When explaining JSON usage, prefer commands that can be piped into `jq`.

## Troubleshooting

For `unauthorized user`, check login state first:

```bash
quail-cli me
```

If using an API key, confirm the shell variable is present without printing the key:

```bash
test -n "$QUAIL_API_KEY"
```

For scripts, prefer `QUAIL_API_KEY` over writing a config file.

If a command needs access to private or paid content, the authenticated user must have permission to read it.
