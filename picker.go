package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type pickerText struct {
	catalog, prompt, controls, languageBar, command, example, copied, clipboardUnavailable string
}

const (
	ice        = "\x1b[38;2;100;223;211m"
	orchid     = "\x1b[38;2;239;141;215m"
	slate      = "\x1b[38;2;157;170;196m"
	brand      = "\x1b[1;38;2;16;19;30;48;2;100;223;211m MANN \x1b[0m"
	englishKey = "Shift+E English"
	spanishKey = "Shift+S Español"
)

func keyHint(key, action string) string {
	return ice + key + cReset + " " + action
}

var english = pickerText{
	catalog: "Command catalog", prompt: "Search › ",
	controls:    keyHint("Enter", "Copy") + "   " + keyHint("Tab", "Preview") + "   " + keyHint("Esc", "Quit"),
	languageBar: ice + englishKey + cReset + "  /  " + spanishKey,
	command:     "Command", example: "Example", copied: "Copied via",
	clipboardUnavailable: "Could not copy the command to the clipboard.",
}

var spanish = pickerText{
	catalog: "Catálogo de comandos", prompt: "Buscar › ",
	controls:    keyHint("Enter", "Copiar") + "   " + keyHint("Tab", "Vista previa") + "   " + keyHint("Esc", "Salir"),
	languageBar: englishKey + "  /  " + ice + spanishKey + cReset,
	command:     "Comando", example: "Ejemplo", copied: "Copiado con",
	clipboardUnavailable: "No se pudo copiar el comando al portapapeles.",
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

func runPicker(items []Item, ui pickerText, query string) ([]string, error) {
	var input strings.Builder
	for i, it := range items {
		display := it.Description + slate + "  ·  " + it.Section + cReset
		example := it.Example
		if example == it.Command {
			example = ""
		}
		fields := []string{strconv.Itoa(i), it.Section, it.Description, it.Command, example, display}
		input.WriteString(strings.Join(fields, fieldSep))
		input.WriteByte(0)
	}

	// fzf shell-quotes each placeholder. Only fixed UI labels enter the format string.
	preview := fmt.Sprintf(`printf '%s%%s%s\n\033[1m%%s%s\n\n%s▍ %s%s\n%%s\n' {2} {3} {4}; test -z {5} || printf '\n%s▍ %s%s\n%%s\n' {5}`,
		slate, cReset, cReset, ice, strings.ToUpper(ui.command), cReset, orchid, strings.ToUpper(ui.example), cReset)
	cmd := exec.Command("fzf",
		"--ansi", "--exact", "--read0", "--print0", "--delimiter="+fieldSep, "--with-nth=6",
		"--print-query", "--expect=E,S", "--query="+query,
		"--prompt="+ui.prompt, "--header-first", "--header="+brand+"  "+ui.catalog+"\n"+ui.controls,
		"--border-label= "+ui.languageBar+" ", "--border-label-pos=3:bottom",
		"--height=80%", "--layout=reverse", "--border=sharp", "--info=inline", "--pointer=▌",
		"--color=fg:#ccd5e6,bg:#10131e,fg+:#ffffff,bg+:#353049,hl:#66eada,hl+:#8affef,"+
			"pointer:#ef8dd7,prompt:#64dfd3,header:#c1cbe0,border:#7784a9,info:#9eaac5,preview-bg:#171c2b,preview-fg:#e5ecfa",
		"--preview", preview, "--preview-window=right,44%,wrap,border-left,<42(down,45%,border-top)",
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
		// A language key can be pressed when the current search has no matches.
	}
	return strings.Split(strings.TrimSuffix(string(out), "\x00"), "\x00"), nil
}

func pickCommand(items []Item) error {
	localized := spanishItems(items)
	entries, ui := items, english
	query := ""
	for {
		result, err := runPicker(entries, ui, query)
		if err != nil || len(result) < 2 {
			return err
		}
		query = result[0]
		switch result[1] {
		case "S":
			entries, ui = localized, spanish
		case "E":
			entries, ui = items, english
		default:
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
