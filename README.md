# wud-mcp

[![CI](https://github.com/reimlima/wud-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/reimlima/wud-mcp/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/reimlima/wud-mcp)](https://github.com/reimlima/wud-mcp/releases)
[![Go 1.22+](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/reimlima/wud-mcp)](https://goreportcard.com/report/github.com/reimlima/wud-mcp)
[![SLSA 3](https://slsa.dev/images/gh-badge-level3.svg)](https://slsa.dev)
[![MegaLinter](https://img.shields.io/badge/MegaLinter-enabled-brightgreen)](https://megalinter.io)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
[![Dependabot](https://img.shields.io/badge/dependabot-enabled-025E8C?logo=dependabot)](https://github.com/reimlima/wud-mcp/blob/main/.github/dependabot.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

An MCP (Model Context Protocol) server that wraps the [WUD (What's Up Docker)](https://getwud.github.io/wud/) REST API, enabling AI assistants like Claude to interact with your self-hosted WUD instance.

## Why

WUD watches your Docker containers and notifies you when new image versions are available — a staple tool for homelabbers and self-hosters. This MCP server bridges WUD and AI assistants, so you can ask Claude things like:

- "Which of my containers have updates available?"
- "Trigger a watch on all containers."
- "Run the Slack trigger for my nginx container."

## Prerequisites

- Go 1.22+ (for building from source)
- A running [WUD](https://getwud.github.io/wud/) instance

## Installation

### Download pre-built binary

Grab the latest binary for your platform from the [GitHub Releases](https://github.com/reimlima/wud-mcp/releases) page.

```bash
# Linux/macOS example
tar -xzf wud-mcp_*.tar.gz
mv wud-mcp /usr/local/bin/
```

### go install

```bash
go install github.com/reimlima/wud-mcp@latest
```

### Docker

Docker support is planned for a future release.

## Configuration

All configuration is done via environment variables:

| Variable | Default | Description |
|---|---|---|
| `WUD_BASE_URL` | `http://localhost:3000` | Base URL of your WUD instance |
| `WUD_API_USER` | _(empty)_ | Basic auth username (optional) |
| `WUD_API_PASSWORD` | _(empty)_ | Basic auth password (optional) |
| `WUD_TIMEOUT` | `10` | HTTP client timeout in seconds |

Copy `.env.example` to `.env` and fill in your values.

## Connecting to Claude Desktop

Add the following to your `claude_desktop_config.json` (usually at `~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "wud": {
      "command": "/usr/local/bin/wud-mcp",
      "env": {
        "WUD_BASE_URL": "http://your-wud-host:3000",
        "WUD_API_USER": "admin",
        "WUD_API_PASSWORD": "secret"
      }
    }
  }
}
```

## Connecting to Claude Code

```bash
claude mcp add wud /usr/local/bin/wud-mcp \
  --env WUD_BASE_URL=http://your-wud-host:3000
```

Or add it to your project's `.mcp.json`:

```json
{
  "mcpServers": {
    "wud": {
      "command": "/usr/local/bin/wud-mcp",
      "env": {
        "WUD_BASE_URL": "http://your-wud-host:3000"
      }
    }
  }
}
```

## Available Tools

| Tool | Description |
|---|---|
| `get_app_info` | Get WUD application name and version |
| `list_containers` | List all containers currently watched by WUD |
| `watch_all_containers` | Trigger a manual watch on all containers |
| `get_container` | Get details of a specific container by ID |
| `get_container_triggers` | List all triggers associated with a container |
| `watch_container` | Trigger a manual watch on a specific container |
| `run_container_trigger` | Manually run a trigger on a container |
| `delete_container` | Delete a container from WUD by ID |
| `list_registries` | List all configured container registries |
| `get_registry` | Get a specific registry by type and name |
| `list_triggers` | List all configured triggers |
| `get_trigger` | Get a specific trigger by type and name |
| `run_trigger` | Run a trigger with optional simulated container data |
| `list_watchers` | List all configured watchers |
| `get_watcher` | Get a specific watcher by type and name |
| `get_store_config` | Get WUD store configuration |
| `get_log_config` | Get WUD logger configuration |

## CI / Quality

Every push and pull request runs:

| Check | Tool |
|---|---|
| Build & test | `go test`, `go vet`, `go build` |
| Linting | [MegaLinter](https://megalinter.io) (includes `golangci-lint`) |
| Vulnerability scan (deps) | [Trivy](https://github.com/aquasecurity/trivy) |
| Vulnerability scan (Go) | [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) |
| PR title format | [Conventional Commits](https://conventionalcommits.org) |
| Supply-chain provenance | [SLSA Level 3](https://slsa.dev) (on release tags) |
| Dependency updates | [Dependabot](https://docs.github.com/en/code-security/dependabot) (weekly) |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT — see [LICENSE](LICENSE).
