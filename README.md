# Quail CLI

`quail-cli` is a command-line interface for interacting with [Quail](https://quaily.com), designed to simplify and automate operations such as user authentication, managing posts, and fetching user details.

Quail CLI interacts with the Quail API at `https://api.quail.ink`.

## Install the latest release

Please Check the latest release [here](https://github.com/quailyquaily/quail-cli/releases), and download the binary for your platform.

Please extract the binary and put it in your `PATH`.

## Install manually

To manually install `quail-cli`, you can use the following command:

```bash
$ go install github.com/quailyquaily/quail-cli@latest
```

> you need to have `go` installed to install `quail-cli` manually. Please refer to the [official Go installation guide](https://go.dev/doc/install) for more information.

## Usage (command line)

After installation, you can start using `quail-cli` by calling the following command:

```bash
$ quail-cli [command]
```

### Available Commands

- **help**: Get help about any command.
- **login**: Authenticate with Quail using OAuth.
- **me**: Retrieve current user information.
- **post**: Create, update, delete, or retrieve posts.

### Global Flags

- `--api-base string`: Quail API base URL (default: `https://api.quail.ink`).
- `--auth-base string`: Quail Auth base URL (default: `https://quaily.com`).
- `--config string`: Path to the configuration file (default: `$HOME/.config/quail-cli/config.yaml`).
- `--format string`: Specify output format, either `human` (human-readable) or `json` (default: `human`).
- `-h, --help`: Display help information for the `quail-cli`.

### Authenticate with Quail

```bash
$ quail-cli login
```

This will initiate OAuth login to authenticate with Quail. Please follow the instructions to complete the authentication process.

1. visit the URL provided in the terminal.
2. Authorize the application.

### Retrieve Current User Information

```bash
$ quail-cli me
```

Get the details of the currently authenticated user.

### Post Operations

#### Upsert a Post

```bash
$ quail-cli post upsert your_markdown_file.md -l your_list_slug
```

in which,

- `your_markdown_file.md` is the path to the markdown file.
- `your_list_slug` is the slug of the list where the post will be created or updated. For example, if the list URL is `https://quaily.com/my_list`, then `your_list_slug` is `my_list`.

quail-cli will read the frontmatter from the markdown file to create or update a post. If the post does not exist, it will be created. If it exists, it will be updated.

Here is an example of a markdown file:

```markdown
---
title: "Here is the title"
slug: your-post-slug
datetime: 2024-09-30 18:42
summary: "This is a summary of the post."
tags: tag1, tag2, tag3
cover_image_url: "your-image-url.jpg"
---

> Any sufficiently advanced technology is indistinguishable from magic.
>
> -- Arthur C. Clarke

This is the body of the post.

    int main() {
        printf("Hello, World!");
        return 0;
    }

## Section Title

This is the last section of the post.
```

#### Publish/Unpublish/Deliver/Delete a Post

```bash
$ quail-cli post publish -l your_list_slug -p your_post_slug
```

```bash
$ quail-cli post unpublish -l your_list_slug -p your_post_slug
```

```bash
$ quail-cli post deliver -l your_list_slug -p your_post_slug
```

```bash
$ quail-cli post delete -l your_list_slug -p your_post_slug
```

## Usage (MCP server)

> [!WARNING]
> This feature is experimental.

MCP(Model Context Protocol) is a protocol for interacting with AI models (like Claude, GPT, etc.). quail-cli can be used as a MCP server.

Before that, you need to login to quail-cli by running the following command:

```bash
$ quail-cli login
```

**SSE mode**

Run the following command to start the MCP server in SSE mode.

```bash
$ quail-cli mcp --sse
```

The server will start a HTTP server and listen on port 8083. You can use any MCP client to connect to this server by providing the URL `http://localhost:8083/sse`.

**Stdio mode**

Before configuring the quail-cli as a MCP server in stdio mode, you need to install quail-cli. And then, find the path of the `quail-cli` binary and pass `mcp` as the argument to the quail-cli.

```bash
/your_install_path/quail-cli mcp
```

### Interact with quail-cli via MCP client

Here is an example prompt:

> what's the url of the 1th post of my first channel?

The LLM will call a series of tools from quail-cli to get the answer.

### Supported tools

- `login`: login to quail.
- `search_post`: search a post.
- `get_my_channels`: get all channels of the current user.
- `get_my_channel_posts`: get posts of a specific channel.
- `get_post`: get a post.
- `publish_post`: publish a post.
- `get_quaily_url`: get the url of a post or a channel.

## Configuration

By default, `quail-cli` reads from `$HOME/.config/quail-cli/config.yaml`. if the file does not exist, it will be created after the first login.

You can specify a different configuration file by using the `--config` flag.

```bash
$ quail-cli --config /path/to/config.yaml
```

### Configuration File Example

```yaml
# DO NOT modify `app` section, quail-cli will manage it.
app:
  access_token: ""
  expiry: ""
  refresh_token: ""
  token_type: ""
  user:
    id: 1
    name: "your_name"
    bio: "your_bio"
    
post:
  # frontmatter_mapping is used to map the frontmatter keys
  # for this example:
  # you can use`featureImage` in the frontmatter and it will be mapped to `cover_image_url`
  frontmatter_mapping:
    cover_image_url: featureImage
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

## License

This project is licensed under the GNU Affero General Public License v3.0.
