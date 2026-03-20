# hizal-tui

**A beautiful terminal UI for [Hizal](https://github.com/XferOps/hizal) — structured memory for AI agents.**

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

> ⚠️ Early development. Core search and browse are working; more views coming fast.

---

## Install

```bash
go install github.com/XferOps/hizal-tui/cmd/hizal@latest
```

Or build from source:

```bash
git clone git@github.com:XferOps/hizal-tui.git
cd hizal-tui
go build -o hizal ./cmd/hizal
```

## Usage

```bash
hizal
```

On first run, you'll be prompted for your Hizal API URL and API key. Config is saved to `~/.config/hizal/config.toml`.

## Keybindings

| Key | Action |
|-----|--------|
| `1-4` | Switch views (Search / Sessions / Projects / Agents) |
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `Enter` | Open / select |
| `Esc` | Back |
| `Tab` | Cycle scope in search (All / Project / Agent / Org) |
| `?` | Toggle keybind help |
| `q` | Quit |

## Contributing

**This is a great place to contribute.** The TUI is self-contained — you don't need deep knowledge of the Hizal API to work on it. If you know Go and have opinions about what a beautiful terminal app looks like, we want your help.

Ideas for contributions:
- Session view (browse active and past agent sessions)
- Chunk detail view with markdown rendering
- inject_audience viewer
- Chunk create / edit flow
- Theme customization

Open an issue or jump straight to a PR.

## Stack

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) — pre-built components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — styling and layout
- [Glamour](https://github.com/charmbracelet/glamour) — markdown rendering

## License

Apache 2.0
