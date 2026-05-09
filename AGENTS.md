# Repository Guidelines

## Project Structure & Module Organization

This is a single Go module for the `mann` CLI/TUI. `main.go` contains the command parser, embedded catalog loading, `fzf` integration, and clipboard backend selection. Cheat-sheet content lives in `data/`; every `data/*.md` file except `data/README.md` is embedded into the binary via `//go:embed data/*.md`. Add tests beside the code as `*_test.go` files; there is no separate test directory today.

## Build, Test, and Development Commands

- `go mod tidy`: normalize module metadata after dependency changes.
- `go build -o mann .`: build the local binary and embed the current `data/` catalog.
- `go install .`: install the tool on your `GOBIN`/`GOPATH` path.
- `go test ./...`: run all Go tests. In restricted environments with a read-only home cache, use `GOCACHE=/tmp/go-build go test ./...`.
- `./mann`: run the picker locally. Interactive mode requires `fzf` on `PATH`; without it, the app prints the catalog as a fallback.

Rebuild after editing `data/` files, because the catalog is compiled into the executable.

## Coding Style & Naming Conventions

Use standard Go formatting: run `gofmt` on changed `.go` files before committing. Keep package-level names unexported unless they are intentionally part of a public API; this module currently uses simple lowerCamelCase helpers and short structs such as `Item`. Prefer small functions that match existing responsibilities: parsing catalog files, choosing clipboard backends, and rendering `fzf` input. Keep Markdown catalog filenames lowercase with hyphens, for example `network-scan.md`.

## Testing Guidelines

Use Go's built-in `testing` package. Focus unit tests on pure behavior such as `parseFile`, `readCommand`, sorting/loading rules, and clipboard candidate selection where practical. Name tests by behavior, for example `TestParseFileAcceptsFencedCommand`. Avoid tests that require real `fzf` or clipboard utilities unless they are explicitly integration tests and can be skipped when tools are missing.

## Commit & Pull Request Guidelines

The current history is minimal (`Initial commit`, `2`), so use clear, imperative commit subjects going forward, such as `Add catalog parser tests` or `Update nmap cheat entries`. Pull requests should describe the change, list validation performed (`go test ./...`, manual `./mann` run), and call out any changes to embedded catalog format or external tool expectations.

## Security & Configuration Tips

The catalog contains defensive and authorized-testing commands. Do not add destructive, credential-stealing, or unauthorized-use examples. Keep commands explicit, documented, and scoped to legitimate administration or testing workflows.
