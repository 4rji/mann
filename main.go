package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
)

//go:embed data/*.md
var dataFS embed.FS

type Item struct {
	Section     string
	Description string
	Command     string
}

const (
	cReset  = "\x1b[0m"
	cCyan   = "\x1b[38;5;67m"
	cGreen  = "\x1b[38;5;108m"
	cYellow = "\x1b[38;5;143m"
	cMag    = "\x1b[38;5;132m"
	cGray   = "\x1b[38;5;245m"
	cBlue   = "\x1b[38;5;74m"
)

func have(bin string) bool {
	_, err := exec.LookPath(bin)
	return err == nil
}

// parseFile extracts items from a markdown file.
//
// Format:
//
//	# Section Title         (one per file, becomes Item.Section)
//	## Description          (becomes Item.Description)
//	`command goes here`     (first non-blank, non-heading line below the ## )
//
// Anything else is documentation and is ignored.
func parseFile(content string) []Item {
	var items []Item
	var section string
	lines := strings.Split(content, "\n")
	for i := 0; i < len(lines); i++ {
		t := strings.TrimSpace(lines[i])
		switch {
		case strings.HasPrefix(t, "## "):
			desc := strings.TrimSpace(t[3:])
			cmd, advance := readCommand(lines, i+1)
			if cmd != "" && section != "" {
				items = append(items, Item{Section: section, Description: desc, Command: cmd})
			}
			i = advance
		case strings.HasPrefix(t, "# ") && !strings.HasPrefix(t, "## "):
			section = strings.TrimSpace(t[2:])
		}
	}
	return items
}

// readCommand finds the next command line below a `## ` heading.
// Returns the parsed command and the index to resume scanning from.
func readCommand(lines []string, start int) (string, int) {
	for j := start; j < len(lines); j++ {
		nt := strings.TrimSpace(lines[j])
		if nt == "" {
			continue
		}
		if strings.HasPrefix(nt, "#") {
			return "", j - 1
		}
		if strings.HasPrefix(nt, "```") {
			var buf []string
			k := j + 1
			for k < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[k]), "```") {
				buf = append(buf, lines[k])
				k++
			}
			return strings.Join(buf, "\n"), k
		}
		return strings.Trim(nt, "`"), j
	}
	return "", len(lines)
}

func loadItems() ([]Item, error) {
	entries, err := dataFS.ReadDir("data")
	if err != nil {
		return nil, err
	}
	var all []Item
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		if strings.EqualFold(e.Name(), "README.md") {
			continue
		}
		b, err := dataFS.ReadFile("data/" + e.Name())
		if err != nil {
			return nil, err
		}
		all = append(all, parseFile(string(b))...)
	}
	sort.SliceStable(all, func(i, j int) bool {
		if all[i].Section != all[j].Section {
			return all[i].Section < all[j].Section
		}
		return all[i].Description < all[j].Description
	})
	return all, nil
}

type clipBackend struct {
	Bin  string
	Args []string
}

func clipCandidates() []clipBackend {
	switch runtime.GOOS {
	case "darwin":
		return []clipBackend{{"pbcopy", nil}}
	case "windows":
		return []clipBackend{{"clip", nil}}
	}
	var out []clipBackend
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		out = append(out, clipBackend{"wl-copy", nil})
	}
	if os.Getenv("DISPLAY") != "" {
		out = append(out,
			clipBackend{"xclip", []string{"-selection", "clipboard"}},
			clipBackend{"xsel", []string{"--clipboard", "--input"}},
		)
	}
	out = append(out,
		clipBackend{"wl-copy", nil},
		clipBackend{"xclip", []string{"-selection", "clipboard"}},
		clipBackend{"xsel", []string{"--clipboard", "--input"}},
		clipBackend{"termux-clipboard-set", nil},
	)
	return out
}

func copyToClipboard(s string) (string, error) {
	seen := map[string]bool{}
	for _, c := range clipCandidates() {
		if seen[c.Bin] {
			continue
		}
		seen[c.Bin] = true
		if !have(c.Bin) {
			continue
		}
		cmd := exec.Command(c.Bin, c.Args...)
		cmd.Stdin = strings.NewReader(s)
		if err := cmd.Run(); err != nil {
			continue
		}
		return c.Bin, nil
	}
	return "", fmt.Errorf("no clipboard tool found")
}

func clipInstallHint() string {
	switch runtime.GOOS {
	case "darwin", "windows":
		return ""
	}
	return strings.Join([]string{
		"  Arch / Manjaro:  sudo pacman -S wl-clipboard   # Wayland",
		"                   sudo pacman -S xclip          # X11",
		"  Debian / Kali:   sudo apt install wl-clipboard # Wayland",
		"                   sudo apt install xclip        # X11",
		"  Fedora:          sudo dnf install wl-clipboard xclip",
		"  Termux:          pkg install termux-api",
	}, "\n")
}

func runFzf(items []Item) (Item, bool) {
	var in bytes.Buffer
	for _, it := range items {
		display := fmt.Sprintf("%s%s%s %s│%s %s%s%s",
			cCyan, it.Section, cReset,
			cGray, cReset,
			cYellow, it.Description, cReset)
		fmt.Fprintf(&in, "%s\t%s\t%s\t%s\n",
			it.Section, it.Description, it.Command, display)
	}

	preview := strings.Join([]string{
		`printf '`,
		`\033[38;5;132mSection\033[0m\n  %s\n\n`,
		`\033[38;5;143mDescription\033[0m\n  %s\n\n`,
		`\033[38;5;108mCommand\033[0m\n  \033[38;5;67m%s\033[0m\n`,
		`' {1} {2} {3}`,
	}, "")

	cmd := exec.Command("fzf",
		"--ansi",
		"--exact",
		"--delimiter=\t",
		"--with-nth=4",
		"--prompt=Comando> ",
		"--height=80%",
		"--layout=reverse",
		"--border",
		"--info=inline",
		"--color=fg:#b8bec8,fg+:#d0d6de,bg:-1,bg+:#262a30,hl:#8fb0d0,hl+:#a9bed6,prompt:#8fb0d0,pointer:#d0879a,marker:#8fbf9f,header:#9aa6b2,border:#555b63,spinner:#8fb0d0,info:#9aa6b2",
		"--preview", preview,
		"--preview-window=right:55%:wrap:border-left",
		"--header=ENTER select • ESC quit • TAB toggle preview",
		"--bind=tab:toggle-preview",
	)
	cmd.Stdin = &in
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			return Item{}, false
		}
		fmt.Fprintln(os.Stderr, "fzf error:", err)
		os.Exit(1)
	}

	line := strings.TrimRight(string(out), "\n")
	if line == "" {
		return Item{}, false
	}
	parts := strings.SplitN(line, "\t", 4)
	if len(parts) < 3 {
		fmt.Fprintln(os.Stderr, "unexpected fzf output:", line)
		os.Exit(1)
	}
	return Item{Section: parts[0], Description: parts[1], Command: parts[2]}, true
}

func fallbackList(items []Item) {
	fmt.Println("Install fzf for interactive filtering. Listing all entries:")
	fmt.Println()
	for _, it := range items {
		fmt.Printf("%s%s%s | %s%s%s\n  %s%s%s\n",
			cCyan, it.Section, cReset,
			cYellow, it.Description, cReset,
			cGreen, it.Command, cReset)
	}
}

func main() {
	items, err := loadItems()
	if err != nil {
		fmt.Fprintln(os.Stderr, "load error:", err)
		os.Exit(1)
	}
	if len(items) == 0 {
		fmt.Fprintln(os.Stderr, "no items found in data/")
		os.Exit(1)
	}

	if !have("fzf") {
		fallbackList(items)
		os.Exit(1)
	}

	picked, ok := runFzf(items)
	if !ok {
		os.Exit(0)
	}

	fmt.Printf("\n%sCommand%s\n  %s%s%s\n",
		cGreen, cReset, cCyan, picked.Command, cReset)

	backend, err := copyToClipboard(picked.Command)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%sclipboard%s %s%s%s\n",
			cMag, cReset, cYellow, err, cReset)
		if hint := clipInstallHint(); hint != "" {
			fmt.Fprintln(os.Stderr, hint)
		}
		return
	}
	fmt.Printf("%sCopied%s via %s%s%s  %s(%s)%s\n",
		cGreen, cReset,
		cBlue, backend, cReset,
		cGray, runtime.GOOS, cReset,
	)
}
