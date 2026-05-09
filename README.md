# mann

Interactive cheat-sheet TUI for networking, recon, and pentesting commands. Type to filter, arrows to move, Enter to pick. The selected command is printed to stdout **and** copied to your clipboard, with the clipboard backend reported on screen.

Go port of the original `mann.sh`. Uses real `fzf` under the hood (subprocess) to keep the exact look-and-feel — full ANSI colors, preview pane, key bindings — and bundles the cheat catalog as embedded markdown files inside the binary.

## Why Go instead of bash

- Single static binary — no `wl-copy` / `xclip` / `pbcopy` shell glue, the binary detects the right backend at runtime.
- Cross-platform clipboard built in: macOS (`pbcopy`), Wayland (`wl-copy`), X11 (`xclip` / `xsel`), Windows (`clip`), Termux (`termux-clipboard-set`).
- Cheat catalog lives as plain markdown files under `data/`, embedded at build time. One file per section makes it trivial to add, edit, or remove categories.
- Same UX as the bash version: `fzf` invoked exactly the way you’d expect, just driven from Go.

## Requirements

- Go 1.22+
- `fzf` on `$PATH` (the picker uses real fzf so colors and bindings match what you already use)
- A clipboard helper for your platform (see install table below). On a headless box without one, the command still prints to stdout — the binary just reports the missing backend.

If `fzf` is missing, `mann` prints the full cheat list to stdout and exits — so it’s still usable as a static reference.

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

- Type to fuzzy-filter the list.
- `↑` / `↓` to navigate, `Enter` to select, `Esc` / `Ctrl-C` to abort.
- `Tab` toggles the preview pane.
- The preview pane shows section, description, and the full command in color.
- After selection, `mann` prints the chosen command and a line such as `Copied via wl-copy (linux)` so you know which clipboard backend was used.

Pipe-friendly:

```sh
./mann | sh        # run the picked command immediately (use with care)
./mann >> notes.md # append the command to a notes file
```

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

The catalog lives under `data/` — one markdown file per section. See `data/README.md` for the file format. Quick summary:

```markdown
# Section Title
## Description of the command
`actual command`
Optional notes — ignored by the parser.
```

Current sections:

- IP routing (`ip route del`)
- Nmap (vuln/malware scripts, decoys, source-port tricks)
- Nmap vuln+safe combinations
- Masscan (top-ports, custom rate)
- NSE script discovery
- Nvidia driver kernel-header install
- Process / port investigation (`lsof`, `pwdx`)
- Packet capture (`tcpdump`, `tshark`)
- Network discovery (`arp-scan`, `netdiscover`)
- Rsync over SSH (custom port, identity file)
- SSH keep-alive
- `find -exec` file moves
- Hydra (authorized SSH brute force)
- FTP payload extraction with tshark

To **add** a category, drop a new `data/<name>.md` and re-build. To **remove**, delete its file and re-build. To **edit** an existing command, open the matching file under `data/` — the file name maps directly to the section.

## Notes

- This is **defensive / authorized-testing material**. Only run these commands against systems you own or have explicit permission to test.
- Part of the [`binarios-go`](https://github.com/4rji/binarios-go) monorepo: each tool is an independent Go module — run `go` commands inside `mann/`, not at the repo root.

## License

Same license as the parent repo.
