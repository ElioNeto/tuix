# tuix — TUI Library for Go

## Project Overview
tuix is a Go library for building Terminal User Interface (TUI) applications by writing HTML and CSS. It parses HTML/CSS and renders them in the terminal using a custom layout and styling engine.

## Tech Stack
- **Language:** Go
- **Module:** `github.com/<user>/tuix` (TBD)
- **Go Version:** 1.22+

## Project Structure (planned)
```
tuix/
├── .teamcode/         # OpenCode/TeamCode configuration
├── cmd/               # CLI entry points
├── internal/
│   ├── parser/        # HTML/CSS parser
│   ├── renderer/      # Terminal renderer
│   ├── layout/        # Box model / flexbox layout engine
│   └── css/           # CSS selector matching & property resolution
├── pkg/
│   └── dom/           # DOM types and manipulation
├── examples/          # Example TUI apps
├── go.mod
└── go.sum
```

## Code Conventions
- **Formatting:** Always run `gofmt` / `go fmt ./...` before committing
- **Imports:** Group standard library, external, and internal imports
- **Errors:** Never ignore errors. Use `fmt.Errorf("...: %w", err)` for wrapping
- **Naming:** Use Go conventions — camelCase for unexported, PascalCase for exported
- **Comments:** Document exported symbols. Use `// Package` doc comments on packages
- **Tests:** Table-driven tests preferred. Test files alongside source files

## Commands
- `dev` — Run the app in development mode
- `build` — Build production binary
- `test` — Run `go test ./...`
- `lint` — Run `go vet ./...` or `golangci-lint`
- `tidy` — Run `go mod tidy`

## LSP
LSP is enabled for Go (`gopls`). The IDE will provide autocompletion, diagnostics, and refactoring support.

## Architecture
- HTML/CSS strings → Parser → DOM tree → Layout engine → Terminal renderer
- The library should support standard CSS layout (block, inline, flexbox)
- Rendering targets ANSI-capable terminals
