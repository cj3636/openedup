# OpenedUp v2

OpenedUp is a keyboard-first TUI filesystem control center for people who live in terminals, SSH sessions, and tmux panes. It stays intentionally small: it helps you navigate files, directories, favorites, history, commands, and optional system services, then launches the existing tools you already use.

## Highlights

- Bubble Tea/Bubbles/Lipgloss TUI with header, breadcrumbs, global search, list, optional preview, and status footer.
- Unified `Entry` model for groups, files, folder browsing/viewing, history, favorites, commands, and search results.
- Lazy directory loading; no recursive disk scans at startup.
- Editor detection using `$EDITOR`, VS Code, Nano, Vim, then Vi.
- Skate-path-compatible storage rooted at `~/.local/share/charm/kv/openedup` for config-adjacent app data, favorites, and history.
- Optional zoxide and systemd providers when those tools are installed.
- Config file at `~/.config/openedup/config.json` with sensible defaults.

## Install

```bash
go install .
```

## Run

```bash
openedup [start-directory]
```

## Keyboard shortcuts

| Shortcut | Action |
| --- | --- |
| `j` / `k`, arrows | Move selection |
| `enter` / `l` | Open selected entry |
| `backspace` / `h` | Go back |
| `ctrl+f` | Global search |
| `ctrl+h` | History |
| `ctrl+d` | Favorites |
| `ctrl+g` | Home/groups |
| `ctrl+s` | Settings hint |
| `?` | Help |
| `q` | Quit |

## Non-goals

OpenedUp is not a shell replacement, text editor, full file manager, terminal emulator, Git client, network browser, or plugin host. It is a fast launcher and navigation hub for existing Unix tools.
