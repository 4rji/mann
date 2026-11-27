# mann

`mann` is a terminal command reference for system administration, networking, and authorized security testing. Search the catalog, preview a command, and press Enter to print and copy it. Each session starts in English. Press **Shift+S** for Spanish or **Shift+E** for English.

The Go port of `mann.sh` uses `fzf` for the interactive picker and embeds the Markdown manuals in the binary.

## Requirements

- Go 1.22+
- `fzf` 0.35.0 or later on `$PATH` for the interactive picker
- A clipboard helper for your platform (see install table below). On a headless box without one, the command still prints to stdout — the binary just reports the missing backend.

If `fzf` is missing, `mann` prints the full cheat list to stdout and exits — so it’s still usable as a static reference.

## Install

Install directly from GitHub:

```sh
go install github.com/4rji/mann@latest
```

Ensure your Go binary directory (`GOBIN`, or `$(go env GOPATH)/bin` by default) is on `$PATH`, then run `mann`.

## Build

This repo follows the project convention: each tool is its own Go module. Build from inside `mann/`:

```sh
cd mann
go mod tidy
go build -o mann .
```

To install on `$PATH`:

```sh
go install .
```

The catalog under `data/` is embedded at compile time, so you must re-build after editing or adding markdown files there.

## Usage

```sh
./mann
```

| Key | Action |
| --- | --- |
| Type | Filter entries by section or description. Lowercase `e` and `s` are ordinary search characters. |
| `↑` / `↓` | Move through matches. |
| `Shift+S` | Show the picker and catalog headings in Spanish. |
| `Shift+E` | Show them in English. |
| `Tab` | Show or hide the command preview. |
| `Enter` | Print the selected command and copy it to the clipboard. |
| `Esc` / `Ctrl+C` | Quit the picker. |

The language shortcuts appear on the bottom border of the picker. Switching languages keeps your search text. Matches can change when the section names and descriptions change language. For example, search for `tmux`, then press `Shift+S` to browse the same topic in Spanish.

The preview shows the command and a matching example when available. After selection, the terminal shows which clipboard tool handled the copy. Paste the command into your notes or terminal and fill in any placeholders before use.

### Updating the catalog (`-u`)

The binary ships with an embedded catalog, but you can pull the latest manuals
from [`github.com/4rji/mann`](https://github.com/4rji/mann) without rebuilding:

```sh
mann -u            # download data/*.md into ~/.local/share/mann
```

No `sudo` needed — the catalog lives in a per-user directory following the
XDG spec: `$XDG_DATA_HOME/mann`, or `~/.local/share/mann` by default.
Downloaded files are read at runtime and **override the built-in ones by
filename**, so new or updated manuals take effect immediately — no rebuild.
For older downloads, the picker uses the built-in translations where the command and heading still match.

Override the location with `MANN_DATA_DIR` (used both to download and to read):

```sh
MANN_DATA_DIR=/some/dir mann -u    # download here
MANN_DATA_DIR=/some/dir mann       # and read from here
```

`mann -h` prints the available flags.

## Clipboard backends

The binary tries the following in order, picking the first one available on `$PATH`:

| OS / session     | Tool                       | Install (Arch)              | Install (Debian / Kali)       | Install (Fedora)              | Install (macOS) |
|------------------|----------------------------|-----------------------------|-------------------------------|-------------------------------|-----------------|
| macOS            | `pbcopy`                   | n/a                         | n/a                           | n/a                           | built-in        |
| Linux / Wayland  | `wl-copy`                  | `sudo pacman -S wl-clipboard` | `sudo apt install wl-clipboard` | `sudo dnf install wl-clipboard` | n/a             |
| Linux / X11      | `xclip`                    | `sudo pacman -S xclip`      | `sudo apt install xclip`      | `sudo dnf install xclip`      | n/a             |
| Linux / X11 alt  | `xsel`                     | `sudo pacman -S xsel`       | `sudo apt install xsel`       | `sudo dnf install xsel`       | n/a             |
| Windows          | `clip`                     | n/a                         | n/a                           | n/a                           | n/a             |
| Android (Termux) | `termux-clipboard-set`     | n/a                         | n/a                           | n/a                           | `pkg install termux-api` |

If none are found, `mann` prints an install hint for your OS and continues — the command is still on stdout for you to copy by hand.

## Cheat catalog

Each Markdown file under `data/` supplies a section. The catalog covers terminal tools such as tmux, sed, tar and QMD, file transfers, network diagnostics, and security tools including Snyk. The section title and each command description can have a Spanish translation in the same file. See the [catalog authoring guide](data/README.md).

## Notes

- This is **defensive / authorized-testing material**. Only run these commands against systems you own or have explicit permission to test.
- Part of the [`binarios-go`](https://github.com/4rji/binarios-go) monorepo: each tool is an independent Go module — run `go` commands inside `mann/`, not at the repo root.

## License

Same license as the parent repo.
