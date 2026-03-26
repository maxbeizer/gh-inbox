# Copilot Instructions for gh-inbox

## Project Overview
gh-inbox is a rich TUI (Terminal User Interface) for managing GitHub notifications, built as a `gh` CLI extension. Inspired by gh-dash's visual style but focused entirely on the notifications workflow.

## Technology Stack
- **Language**: Go 1.24+
- **TUI Framework**: Bubble Tea v2 (`charm.land/bubbletea/v2`)
- **Styling**: Lipgloss v2 (`charm.land/lipgloss/v2`)
- **Markdown**: Glamour v2 (`charm.land/glamour/v2`)
- **Mouse Support**: bubblezone v2 (`github.com/lrstanley/bubblezone/v2`)
- **GitHub API**: go-gh v2 (`github.com/cli/go-gh/v2`)
- **CLI Framework**: Cobra (`github.com/spf13/cobra`)

## Architecture
- `main.go` — Cobra root command that launches the Bubble Tea program
- `internal/api/` — GitHub REST/GraphQL client wrappers via go-gh
- `internal/model/` — Domain types (Notification, Filter, Sort)
- `internal/tui/` — Bubble Tea application (root model, keys, styles, messages)
- `internal/tui/components/` — Composable UI components (table, preview, header, statusbar, filter, help)

## Go Development Guidelines
- Follow idiomatic Go: explicit error handling, composition over inheritance, short receiver names
- Use table-driven tests with `t.Run()` subtests
- Keep functions focused and under 50 lines where possible
- Use interfaces for abstraction and testability
- Packages organized by functionality, not by layer

## Bubble Tea Patterns
- Root model (`*App`) uses pointer receivers for all tea.Model methods
- Root composes child component models (header, table, preview, statusbar, filter, help)
- Custom `tea.Msg` types in `tui/messages.go` for inter-component communication
- Keybindings defined centrally in `tui/keys.go` (mirrors GitHub web shortcuts)
- Theme/styles defined in `tui/theme/theme.go` (separate package to avoid import cycles)

## API Patterns
- REST for all GitHub API access (notifications CRUD + subject detail fetches)
- Pagination capped at 10 pages (500 notifications) to avoid rate limiting
- Always fetches with `all=true`; filtering is done client-side
- Preview detail cache bounded at 100 entries, cleared on refresh
- All API access through `go-gh` which handles auth via `gh auth`

## Testing
- Table-driven tests for model parsing and API response mapping
- Test files alongside source: `*_test.go`
- Run: `make test` or `make ci`
