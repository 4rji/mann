package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadItemsInheritsOnlyMatchingTranslations(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("MANN_DATA_DIR", dir)
	manuals := map[string]string{
		"qmd.md": "# QMD Commands\n" +
			"## Ver estado del índice\n`qmd status`\n" +
			"## My collection list\n`qmd collection list`\n" +
			"<!-- es: Mi búsqueda -->\n## Search by keyword\n`qmd search \"QUERY\"`\n" +
			"## Add a collection\n`qmd local add`\n",
		"custom.md": "# QMD Commands\n## List indexed files\n`qmd ls NAME`\n",
	}
	for name, content := range manuals {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	items, err := loadItems()
	if err != nil {
		t.Fatal(err)
	}
	byCommand := make(map[string]Item)
	for _, it := range items {
		byCommand[it.Command] = it
	}
	for _, want := range []Item{
		{Section: "QMD Commands", SectionES: "Comandos de QMD", Description: "Check index status", DescriptionES: "Ver estado del índice", Command: "qmd status"},
		{Section: "QMD Commands", SectionES: "Comandos de QMD", Description: "My collection list", Command: "qmd collection list"},
		{Section: "QMD Commands", SectionES: "Comandos de QMD", Description: "Search by keyword", DescriptionES: "Mi búsqueda", Command: `qmd search "QUERY"`},
		{Section: "QMD Commands", Description: "Add a collection", Command: "qmd local add"},
		{Section: "QMD Commands", Description: "List indexed files", Command: "qmd ls NAME"},
	} {
		if got := byCommand[want.Command]; got != want {
			t.Errorf("%q: got %+v, want %+v", want.Command, got, want)
		}
	}
	if _, exists := byCommand["qmd embed"]; exists {
		t.Error("embedded entry survived replacement of qmd.md")
	}
}
