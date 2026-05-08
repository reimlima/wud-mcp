# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build binary
go build -o wud-mcp .

# Verify all packages compile (no binary output)
go build ./...

# Run tests
go test ./...

# Run a single test
go test ./client/ -run TestGetApp_Success -v

# Run tests with coverage
go test -coverprofile=coverage.out -covermode=atomic ./...

# Vet
go vet ./...

# Lint (via MegaLinter locally)
npx mega-linter-runner

# Vulnerability scan
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Run the server (reads config from env)
WUD_BASE_URL=http://localhost:3000 go run .

# Tidy dependencies
go mod tidy
```

## CI jobs (ci.yml)

| Job | Trigger | What it does |
|---|---|---|
| `ci` | push/PR | test, coverage, go vet, build |
| `trivy` | push/PR | filesystem + dep CVE scan (SARIF → Security tab) |
| `golangci-lint` | push/PR | golangci-lint with correct Go version |
| `megalinter` | push/PR | lint all file types including Go revive (golangci-lint disabled; runs as separate job) |
| `lint-pr-title` | PR only | enforces Conventional Commits on PR title |
| `govulncheck` | push/PR (after ci) | Go vulnerability database scan |
| `provenance` | tag push | SLSA Level 3 attestation, uploaded to GitHub Release |

## Architecture

Thin MCP server over the WUD REST API. Three layers:

**`client/client.go`** — all HTTP calls to WUD. `New()` reads config from env vars (`WUD_BASE_URL`, `WUD_API_USER`, `WUD_API_PASSWORD`, `WUD_TIMEOUT`). The `do()` method handles auth, error decoding, and empty-body normalization (204 → `{"status":"ok"}`). Path segments are `url.PathEscape`d before concatenation.

**`tools/*.go`** — one `Register*` function per file, each calling `s.AddTool(...)` for every MCP tool in that domain. Tool parameter descriptions use `mcp.Description(...)` (a `PropertyOption`), not `mcp.WithDescription(...)` (which is a `ToolOption` for the tool itself). Arguments are accessed via `req.GetArguments()["key"].(string)`, not `req.Params.Arguments`.

**`main.go`** — wires client → server, calls all `Register*` functions, then `server.ServeStdio(s)`.

## Key API facts

- `mcp.WithDescription("...")` — sets description on the **tool**
- `mcp.Description("...")` — sets description on a **property** inside `mcp.WithString(...)`
- `req.GetArguments()` returns `map[string]any`; direct indexing of `req.Params.Arguments` fails (typed as `any`)
- Tests use `t.Setenv("WUD_BASE_URL", srv.URL)` + `client.New()` to point the client at an `httptest.Server`
- Tool handlers are testable directly via `s.GetTool("name").Handler(ctx, req)` — no stdio transport needed
