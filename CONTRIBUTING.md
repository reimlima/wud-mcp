# Contributing

Contributions are welcome. Please read the guidelines below before opening a PR.

## Local development

```bash
git clone https://github.com/reimlima/wud-mcp.git
cd wud-mcp
go mod tidy
```

Copy the example env file and point it at your WUD instance:

```bash
cp .env.example .env
# edit .env
```

Run the server:

```bash
go run .
```

## Testing interactively with MCP Inspector

Build the binary and launch the inspector (requires Node.js):

```bash
go build -o wud-mcp .

npx @modelcontextprotocol/inspector \
  -e WUD_BASE_URL=http://your-wud-host:3000 \
  -e WUD_API_USER=admin \
  -e WUD_API_PASSWORD=secret \
  ./wud-mcp
```

Opens a browser UI at `http://localhost:6274` where you can browse all tools, supply parameters, and inspect raw JSON responses.

## Running tests

```bash
go test ./...
```

## Running the linter

Go linting is handled by [MegaLinter](https://megalinter.io/) in CI via `golangci-lint`. To run it locally:

```bash
npx mega-linter-runner
```

## Vulnerability scanning

CI runs two vulnerability scanners on every push and PR:

- **[Trivy](https://github.com/aquasecurity/trivy)** — scans the filesystem and Go dependencies for known CVEs.
- **[govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)** — checks Go source against the Go vulnerability database.

To run `govulncheck` locally:

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

## Pull request guidelines

- Keep PRs small and focused on a single concern.
- **PR titles must follow [Conventional Commits](https://conventionalcommits.org/)** — enforced automatically by CI.
  Allowed types: `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `ci`, `perf`, `build`.
  Example: `feat: add get_container_logs tool`
- Link the related issue in the PR description when applicable.
- Add or update tests for any changed behaviour.
- Ensure `go test ./...` passes before requesting review.

## Dependency updates

[Dependabot](https://docs.github.com/en/code-security/dependabot) automatically opens weekly PRs for outdated Go modules and GitHub Actions. No manual action is required from contributors — maintainers review and merge these.

## Code style

- Idiomatic Go — follow the patterns already in the codebase.
- No unnecessary external dependencies; prefer the standard library.
- No comments that merely restate what the code already says.
