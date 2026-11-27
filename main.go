package main

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"
)

//go:embed data/*.md
var dataFS embed.FS

type Item struct {
	Section       string
	Description   string
	SectionES     string
	DescriptionES string
	ListEN        string
	ListES        string
	Command       string
	Example       string
}

var placeholderRe = regexp.MustCompile(`\b(IP|INTERFACE)\b`)

const (
	cReset  = "\x1b[0m"
	cCyan   = "\x1b[38;5;67m"
	cGreen  = "\x1b[38;5;108m"
	cYellow = "\x1b[38;5;143m"
	cMag    = "\x1b[38;5;132m"
	cGray   = "\x1b[38;5;245m"
	cBlue   = "\x1b[38;5;74m"
)

const (
	// catalogAPI lists the markdown files under data/ on the default branch.
	catalogAPI = "https://api.github.com/repos/4rji/mann/contents/data"
	repoName   = "github.com/4rji/mann"
)

// dataDir returns the external catalog directory where `mann -u` downloads
// manuals and where mann reads them at runtime. It is a per-user, no-sudo
// location: MANN_DATA_DIR if set, else $XDG_DATA_HOME/mann, else
// ~/.local/share/mann.
func dataDir() string {
	if d := strings.TrimSpace(os.Getenv("MANN_DATA_DIR")); d != "" {
		return d
	}
	if d := strings.TrimSpace(os.Getenv("XDG_DATA_HOME")); d != "" {
		return filepath.Join(d, "mann")
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".local", "share", "mann")
	}
	// Last resort if there is no home directory (rare).
	return "/tmp/mann"
}

type ghContent struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
}

func httpGet(url, accept string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// Generic client name; never leak user/personal details in the UA.
	req.Header.Set("User-Agent", "mann-cli")
	if accept != "" {
		req.Header.Set("Accept", accept)
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	return body, nil
}

func writeHint(dir string, err error) error {
	if errors.Is(err, os.ErrPermission) {
		return fmt.Errorf("%w\n  Sin permiso de escritura en %s.\n  Elige un directorio con permiso de escritura:  MANN_DATA_DIR=/mi/directorio mann -u", err, dir)
	}
	return err
}

// updateCatalog downloads every data/*.md from the repo into dataDir().
func updateCatalog() error {
	dir := dataDir()

	fmt.Printf("%s==>%s Actualizando el catálogo de mann desde %s%s%s\n", cBlue, cReset, cCyan, repoName, cReset)

	body, err := httpGet(catalogAPI, "application/vnd.github+json")
	if err != nil {
		return fmt.Errorf("no se pudo consultar el catálogo: %w", err)
	}
	var contents []ghContent
	if err := json.Unmarshal(body, &contents); err != nil {
		return fmt.Errorf("no se pudo leer el índice del catálogo: %w", err)
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return writeHint(dir, err)
	}

	count := 0
	for _, c := range contents {
		if c.Type != "file" || c.DownloadURL == "" || !strings.HasSuffix(c.Name, ".md") {
			continue
		}
		if strings.EqualFold(c.Name, "README.md") {
			continue
		}
		data, err := httpGet(c.DownloadURL, "")
		if err != nil {
			return fmt.Errorf("descarga de %s: %w", c.Name, err)
		}
		dest := filepath.Join(dir, c.Name)
		if err := os.WriteFile(dest, data, 0o644); err != nil {
			return writeHint(dir, err)
		}
		fmt.Printf("  %s✓%s %s\n", cGreen, cReset, c.Name)
		count++
	}

	fmt.Printf("%sListo%s · Manuales descargados: %d → %s%s%s\n", cGreen, cReset, count, cCyan, dir, cReset)
	return nil
}

func printUsage() {
	fmt.Printf(`mann — catálogo interactivo de comandos

Uso:
  mann            abre el selector interactivo (requiere fzf)
  mann -u         descarga los manuales de %s
                  en %s (sin sudo)
  mann -h         muestra esta ayuda

Atajos:
  Shift+S  Español (predeterminado)     Shift+E  English
  F1       Cian (predeterminado)        F2       Magenta     F3  Ámbar
  Enter    Copiar el comando            Tab      Mostrar u ocultar la vista previa
  Esc      Salir

Para cambiar la ubicación del catálogo:
  MANN_DATA_DIR=/mi/directorio mann -u

Los manuales descargados tienen prioridad sobre los incluidos con el mismo nombre.
`, repoName, dataDir())
}

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
// An optional Spanish translation comment immediately precedes its heading.
// Entry headings may also have bilingual list captions before the translation.
func parseFile(content string) []Item {
	var items []Item
	var examples []string
	var section, sectionES string
	lines := strings.Split(content, "\n")
	for i := 0; i < len(lines); i++ {
		t := strings.TrimSpace(lines[i])
		switch {
		case strings.HasPrefix(t, "## "):
			desc := strings.TrimSpace(t[3:])
			if strings.EqualFold(desc, "Examples") {
				ex, advance := readFence(lines, i+1)
				examples = append(examples, ex...)
				i = advance
				continue
			}
			cmd, advance := readCommand(lines, i+1)
			if cmd != "" && section != "" {
				descES := headingComment(lines, i, "<!-- es:")
				captionAt := i
				if descES != "" {
					captionAt--
				}
				listEN, listES := listCaptions(headingComment(lines, captionAt, "<!-- list:"))
				items = append(items, Item{
					Section: section, SectionES: sectionES, Description: desc,
					DescriptionES: descES, ListEN: listEN, ListES: listES, Command: cmd,
				})
			}
			i = advance
		case strings.HasPrefix(t, "# ") && !strings.HasPrefix(t, "## "):
			section = strings.TrimSpace(t[2:])
			sectionES = headingComment(lines, i, "<!-- es:")
		}
	}
	attachExamples(items, examples)
	return items
}

func headingComment(lines []string, heading int, prefix string) string {
	if heading == 0 {
		return ""
	}
	value, ok := strings.CutPrefix(strings.TrimSpace(lines[heading-1]), prefix)
	if !ok {
		return ""
	}
	value, ok = strings.CutSuffix(value, "-->")
	if !ok {
		return ""
	}
	return strings.TrimSpace(value)
}

func listCaptions(value string) (string, string) {
	en, es, ok := strings.Cut(value, "|")
	en, es = strings.TrimSpace(en), strings.TrimSpace(es)
	if !ok || en == "" || es == "" || strings.Contains(es, "|") {
		return "", ""
	}
	return en, es
}

// readFence reads a triple-backtick fenced block starting after a heading.
// Returns the inner lines (trimmed) and the index of the closing fence.
func readFence(lines []string, start int) ([]string, int) {
	for j := start; j < len(lines); j++ {
		nt := strings.TrimSpace(lines[j])
		if nt == "" {
			continue
		}
		if strings.HasPrefix(nt, "```") {
			var buf []string
			k := j + 1
			for k < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[k]), "```") {
				if line := strings.TrimSpace(lines[k]); line != "" {
					buf = append(buf, line)
				}
				k++
			}
			return buf, k
		}
		return nil, j - 1
	}
	return nil, len(lines)
}

// attachExamples matches each example line to the item whose command it
// realizes. Placeholders `IP` and `INTERFACE` in the command become `\S+`.
func attachExamples(items []Item, examples []string) {
	if len(examples) == 0 {
		return
	}
	for idx := range items {
		re := commandPattern(items[idx].Command)
		if re == nil {
			continue
		}
		for _, ex := range examples {
			if re.MatchString(ex) {
				items[idx].Example = ex
				break
			}
		}
	}
}

func commandPattern(cmd string) *regexp.Regexp {
	const marker = "\x00PH\x00"
	masked := placeholderRe.ReplaceAllString(cmd, marker)
	quoted := regexp.QuoteMeta(masked)
	pattern := "^" + strings.ReplaceAll(quoted, marker, `\S+`) + "$"
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}
	return re
}

// readCommand finds the next command line below a `## ` heading.
// Returns the parsed command and the index to resume scanning from.
func readCommand(lines []string, start int) (string, int) {
	for j := start; j < len(lines); j++ {
		nt := strings.TrimSpace(lines[j])
		if nt == "" {
			continue
		}
		if strings.HasPrefix(nt, "#") || strings.HasPrefix(nt, "<!-- es:") || strings.HasPrefix(nt, "<!-- list:") {
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

func inheritTranslations(items, embedded []Item) {
	for i := range items {
		it := &items[i]
		for _, base := range embedded {
			if it.Command != base.Command {
				continue
			}
			it.Section, it.SectionES = inheritHeading(it.Section, it.SectionES, base.Section, base.SectionES)
			it.Description, it.DescriptionES = inheritHeading(it.Description, it.DescriptionES, base.Description, base.DescriptionES)
			if !strings.EqualFold(it.Description, base.Description) {
				continue
			}
			if it.ListEN == "" {
				it.ListEN = base.ListEN
			}
			if it.ListES == "" && strings.EqualFold(it.DescriptionES, base.DescriptionES) {
				it.ListES = base.ListES
			}
		}
	}
}

// Older downloads may contain either language as the only heading.
func inheritHeading(title, translation, english, spanish string) (string, string) {
	if translation == "" && spanish != "" {
		switch {
		case strings.EqualFold(title, english):
			return title, spanish
		case strings.EqualFold(title, spanish):
			return english, spanish
		}
	}
	return title, translation
}

func sortItems(items []Item) {
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Section != items[j].Section {
			return items[i].Section < items[j].Section
		}
		return items[i].Description < items[j].Description
	})
}

func loadItems() ([]Item, error) {
	files := map[string][]Item{}

	entries, err := dataFS.ReadDir("data")
	if err != nil {
		return nil, err
	}
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
		files[e.Name()] = parseFile(string(b))
	}

	if dir := dataDir(); dir != "" {
		if ents, err := os.ReadDir(dir); err == nil {
			for _, e := range ents {
				if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
					continue
				}
				if strings.EqualFold(e.Name(), "README.md") {
					continue
				}
				b, err := os.ReadFile(filepath.Join(dir, e.Name()))
				if err != nil {
					continue
				}
				items := parseFile(string(b))
				inheritTranslations(items, files[e.Name()])
				files[e.Name()] = items
			}
		}
	}

	names := make([]string, 0, len(files))
	for n := range files {
		names = append(names, n)
	}
	sort.Strings(names)

	var all []Item
	for _, n := range names {
		all = append(all, files[n]...)
	}
	sortItems(all)
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

func fallbackList(items []Item) {
	fmt.Println("Instala fzf para filtrar el catálogo. Estas son todas las entradas:")
	fmt.Println()
	for _, it := range items {
		fmt.Printf("%s%s%s | %s%s%s\n  %s%s%s\n",
			cCyan, it.Section, cReset,
			cYellow, it.Description, cReset,
			cGreen, it.Command, cReset)
	}
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-u", "--update", "update":
			if err := updateCatalog(); err != nil {
				fmt.Fprintf(os.Stderr, "%serror de actualización:%s %v\n", cMag, cReset, err)
				os.Exit(1)
			}
			return
		case "-h", "--help", "help":
			printUsage()
			return
		}
	}

	items, err := loadItems()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error al cargar el catálogo:", err)
		os.Exit(1)
	}
	if len(items) == 0 {
		fmt.Fprintln(os.Stderr, "no se encontraron entradas en data/")
		os.Exit(1)
	}

	if !have("fzf") {
		fallbackList(spanishItems(items))
		os.Exit(1)
	}

	if err := pickCommand(items); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
