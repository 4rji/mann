package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/rivo/uniseg"
)

type pickerText struct {
	prompt, controls, themeLabel, action, command, example, copied, clipboardUnavailable string
	spanishCaptions                                                                      bool
}

const (
	slate        = "\x1b[38;2;157;170;196m"
	englishKey   = "Shift+E English"
	spanishKey   = "Shift+S Español"
	nextThemeKey = "ctrl-t"
)

type pickerTheme struct {
	key, nameEN, nameES                      string
	bg, previewBG, accent, secondary, border uint32
}

// The picker starts with the first theme and cycles through the rest in order.
var pickerThemes = [...]pickerTheme{
	{key: "f2", nameEN: "Magenta", nameES: "Magenta", bg: 0x1a1320, previewBG: 0x261b2f, accent: 0xf3a1dc, secondary: 0xc9aeff, border: 0x75597e},
	{key: "f3", nameEN: "Amber", nameES: "Ámbar", bg: 0x191611, previewBG: 0x272016, accent: 0xf3c475, secondary: 0x9ad4bb, border: 0x786948},
	{key: "f1", nameEN: "Cyan", nameES: "Cian", bg: 0x10141c, previewBG: 0x18212d, accent: 0x7bddcf, secondary: 0xf497bf, border: 0x536579},
}

func foreground(rgb uint32) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", rgb>>16, (rgb>>8)&0xff, rgb&0xff)
}

func (theme pickerTheme) fzfColors() string {
	return fmt.Sprintf(
		"fg:#d2d7e1,fg+:#ffffff,bg:#%06x,bg+:#%06x,hl:#%06x,hl+:#%06x,"+
			"prompt:#%06x,pointer:#%06x,header:#d2d7e1,border:#%06x,info:#9caebe,preview-bg:#%06x,preview-fg:#e3ebf3",
		theme.bg, theme.bg, theme.accent, theme.accent, theme.accent, theme.secondary, theme.border, theme.previewBG,
	)
}

func keyHint(key, action string) string {
	return key + " " + slate + action + cReset
}

var english = pickerText{
	prompt:     "Search  › ",
	controls:   keyHint("Enter", "Copy") + "   " + keyHint("Tab", "Preview") + "   " + keyHint("Esc", "Quit"),
	themeLabel: "Theme",
	action:     "Action", command: "Command", example: "Example", copied: "Copied via",
	clipboardUnavailable: "Could not copy the command to the clipboard.",
}

var spanish = pickerText{
	spanishCaptions: true,
	prompt:          "Buscar  › ",
	controls:        keyHint("Enter", "Copiar") + "   " + keyHint("Tab", "Vista previa") + "   " + keyHint("Esc", "Salir"),
	themeLabel:      "Tema",
	action:          "Acción", command: "Comando", example: "Ejemplo", copied: "Copiado con",
	clipboardUnavailable: "No se pudo copiar el comando al portapapeles.",
}

func (ui pickerText) footer(active pickerTheme) string {
	accent := foreground(active.accent)
	languages := accent + englishKey + cReset + " / " + spanishKey
	name := active.nameEN
	if ui.spanishCaptions {
		languages = englishKey + " / " + accent + spanishKey + cReset
		name = active.nameES
	}
	return " " + languages + "  │  Ctrl+T " + ui.themeLabel + " · " + accent + "\x1b[1m" + name + cReset + " "
}

func spanishItems(items []Item) []Item {
	translated := append([]Item(nil), items...)
	for i := range translated {
		it := &translated[i]
		if it.SectionES != "" {
			it.Section = it.SectionES
		}
		if it.DescriptionES != "" {
			it.Description = it.DescriptionES
		}
	}
	sortItems(translated)
	return translated
}

// fzf reads NUL-terminated records so a fenced, multiline command stays one item.
// Separate fields with a unit separator so tabs inside commands survive the preview.
const fieldSep = "\x1f"

func (ui pickerText) caption(it Item) string {
	label := it.ListEN
	if ui.spanishCaptions {
		label = it.ListES
	}
	if label == "" {
		return it.Description
	}
	return label
}

// Size actions in terminal cells so an accent or a wide glyph never shifts the command.
func actionColumn(label string) string {
	const width = 34
	label = strings.ReplaceAll(label, "\t", " ")
	size := uniseg.StringWidth(label)
	if size <= width {
		return label + strings.Repeat(" ", width-size)
	}
	var out strings.Builder
	size = 0
	clusters := uniseg.NewGraphemes(label)
	for clusters.Next() {
		if size+clusters.Width() > width-1 {
			break
		}
		out.WriteString(clusters.Str())
		size += clusters.Width()
	}
	return out.String() + "…" + strings.Repeat(" ", width-1-size)
}

func runPicker(items []Item, ui pickerText, theme pickerTheme, query string) ([]string, error) {
	var input strings.Builder
	accent := foreground(theme.accent)
	divider := "  " + foreground(theme.border) + "│" + cReset + "  "
	// fzf applies --with-nth=6 to the sticky column headings as well.
	headings := accent + actionColumn(strings.ToUpper(ui.action)) + divider + accent + strings.ToUpper(ui.command) + cReset
	input.WriteString(strings.Repeat(fieldSep, 5) + headings + "\x00")
	for i, it := range items {
		firstLine, _, _ := strings.Cut(it.Command, "\n")
		display := actionColumn(ui.caption(it)) + divider + slate + firstLine + cReset
		example := it.Example
		if example == it.Command {
			example = ""
		}
		fields := []string{strconv.Itoa(i), it.Section, it.Description, it.Command, example, display}
		input.WriteString(strings.Join(fields, fieldSep))
		input.WriteByte(0)
	}

	// fzf shell-quotes each placeholder. Only fixed UI labels enter the format string.
	preview := fmt.Sprintf(`printf '%s%%s%s  ·  %%s\n%s%s%s  %%s\n' {2} {3} {4}; test -z {5} || printf '%s%s%s  %%s\n' {5}`,
		slate, cReset, accent, strings.ToUpper(ui.command), cReset, foreground(theme.secondary), strings.ToUpper(ui.example), cReset)
	keys := []string{"E", "S", nextThemeKey}
	for _, option := range pickerThemes {
		keys = append(keys, option.key)
	}
	brand := "\x1b[1m" + accent + "MANN" + cReset
	cmd := exec.Command("fzf",
		"--ansi", "--exact", "--read0", "--print0", "--delimiter="+fieldSep, "--with-nth=6",
		"--print-query", "--expect="+strings.Join(keys, ","), "--query="+query,
		"--prompt="+ui.prompt, "--header-first", "--header-lines=1", "--header="+brand+"   "+ui.controls,
		"--border-label="+ui.footer(theme), "--border-label-pos=3:bottom",
		"--height=95%", "--layout=reverse", "--border=rounded", "--info=default", "--pointer=▸", "--no-hscroll",
		"--color="+theme.fzfColors(),
		"--preview", preview, "--preview-window=down,5,wrap,border-top",
		"--bind=tab:toggle-preview",
	)
	cmd.Stdin = strings.NewReader(input.String())
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) || (exit.ExitCode() != 1 && exit.ExitCode() != 130) {
			return nil, fmt.Errorf("fzf: %w", err)
		}
		if exit.ExitCode() == 130 {
			return nil, nil
		}
		// A language or theme key can be pressed even when the search has no matches.
	}
	return strings.Split(strings.TrimSuffix(string(out), "\x00"), "\x00"), nil
}

func pickCommand(items []Item) error {
	localized := spanishItems(items)
	entries, ui := localized, spanish
	themeIndex := 0
	query := ""
picker:
	for {
		result, err := runPicker(entries, ui, pickerThemes[themeIndex], query)
		if err != nil || len(result) < 2 {
			return err
		}
		query = result[0]
		switch result[1] {
		case "S":
			entries, ui = localized, spanish
		case "E":
			entries, ui = items, english
		case nextThemeKey:
			themeIndex = (themeIndex + 1) % len(pickerThemes)
		default:
			for i, option := range pickerThemes {
				if result[1] == option.key {
					themeIndex = i
					continue picker
				}
			}
			if len(result) < 3 {
				return nil
			}
			index, _, _ := strings.Cut(result[2], fieldSep)
			i, err := strconv.Atoi(index)
			if err != nil || i < 0 || i >= len(entries) {
				return fmt.Errorf("fzf: invalid selection")
			}
			printSelection(entries[i], ui)
			return nil
		}
	}
}

func printSelection(picked Item, ui pickerText) {
	fmt.Printf("\n%s%s%s\n  %s%s%s\n", cGreen, ui.command, cReset, cCyan, picked.Command, cReset)
	backend, err := copyToClipboard(picked.Command)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n%s%s%s\n", cMag, ui.clipboardUnavailable, cReset)
		if hint := clipInstallHint(); hint != "" {
			fmt.Fprintln(os.Stderr, hint)
		}
		return
	}
	fmt.Printf("%s%s%s %s%s%s  %s(%s)%s\n",
		cGreen, ui.copied, cReset, cBlue, backend, cReset, cGray, runtime.GOOS, cReset)
}
