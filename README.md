# gh-inbox

📬 A rich TUI for managing GitHub notifications.

<img src="demo.gif" alt="gh-inbox demo" width="100%">

## Features

- **Browse notifications** in a table view with reason icons, type, repo, title, and relative time
- **Preview panel** — toggle a side panel to read PR/issue details without leaving the terminal
- **Full lifecycle management** — mark read, unread, done, and unsubscribe
- **Filter & sort** — cycle through unread/all/participating, sort by updated/repo/reason, fuzzy search
- **Keyboard-native** — vim-style navigation (j/k, g/G) with action keys matching [GitHub's web shortcuts](https://docs.github.com/en/get-started/accessibility/keyboard-shortcuts#notifications)
- **Mouse support** — click to select rows
- **Powerline footer** — always-visible status bar with filter state, counts, and key hints
- **Catppuccin Mocha** color palette

## Install

```bash
gh extension install maxbeizer/gh-inbox
```

## Usage

```bash
gh inbox
```

Try it without GitHub auth using demo mode:

```bash
gh inbox --demo
```

## Key Bindings

| Key | Action |
|-----|--------|
| `j/k` or `↑/↓` | Navigate |
| `g/G` | Top/bottom |
| `p` or `Enter` | Toggle preview |
| `e` | Mark as done |
| `⇧I` | Mark as read |
| `⇧U` | Mark as unread |
| `⇧M` | Unsubscribe |
| `o` | Open in browser |
| `y` | Copy URL |
| `r/R` | Refresh / refresh all |
| `/` | Search |
| `f` | Cycle filter (unread → all → participating) |
| `s` | Cycle sort (updated → repo → reason) |
| `?` | Help + icon legend |
| `q` | Quit |

> Action keys (e, ⇧I, ⇧U, ⇧M) match [GitHub's web notification shortcuts](https://docs.github.com/en/get-started/accessibility/keyboard-shortcuts#notifications).

## Development

```bash
make help          # see all targets
make build         # build binary
make test          # run tests
make ci            # build + vet + test-race
make install-local # install extension from checkout
make relink-local  # reinstall after changes
```

## Releasing

```bash
git tag v0.2.0
git push origin v0.2.0
```

## Inspiration

gh-inbox wouldn't exist without these projects:

- **[gh-dash](https://github.com/dlvhdr/gh-dash)** — the gold standard for GitHub TUIs. gh-inbox borrows its visual language: Bubble Tea composition, lipgloss styling, table + preview layout, and the general "make the terminal feel good" philosophy.
- **[gh-not](https://github.com/nobe4/gh-not)** — a powerful notifications manager with rule-based filtering and local caching. gh-not proved that the GitHub notifications API is workable and that there's a real need for better notification tooling. gh-inbox takes a different approach (interactive TUI vs. rule engine) but shares the same frustration with GitHub's default notification experience.

## License

MIT
