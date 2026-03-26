# gh-inbox

📬 A rich TUI for managing GitHub notifications — keyboard-native, visually appealing, inspired by [gh-dash](https://github.com/dlvhdr/gh-dash).

<img src="demo.gif" alt="gh-inbox demo" width="100%">

## Features

- **Browse notifications** with a beautiful table view showing reason, type, repo, title, and relative time
- **Preview panel** with markdown-rendered PR/issue details
- **Full lifecycle management**: mark read, mark done, unsubscribe
- **Filter & sort**: by read/unread, participating, reason, repo, and fuzzy text search
- **Keyboard-native** with vim-style navigation (j/k, g/G) plus mouse support
- **Open in browser** or copy URL to clipboard
- **Catppuccin-inspired** color palette

## Install

```bash
gh extension install maxbeizer/gh-inbox
```

## Usage

```bash
gh inbox
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
| `?` | Help |
| `q` | Quit |

> Action keys match [GitHub's web notification shortcuts](https://docs.github.com/en/get-started/accessibility/keyboard-shortcuts#notifications).

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
git tag v0.1.0
git push origin v0.1.0
```

## License

MIT
