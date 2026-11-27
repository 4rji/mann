package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInheritCaptionsRespectsCustomHeadings(t *testing.T) {
	base := Item{Description: "Display the date", DescriptionES: "Mostrar la fecha", ListEN: "Date", ListES: "Fecha", Command: "date"}
	for _, tc := range []struct {
		item   Item
		en, es string
	}{
		{Item{Description: "Mostrar la fecha", Command: "date"}, "Date", "Fecha"},
		{Item{Description: "Display the date", DescriptionES: "Mi fecha", Command: "date"}, "Date", ""},
		{Item{Description: "My date", Command: "date"}, "", ""},
		{Item{Description: "Display the date", Command: "another command"}, "", ""},
		{Item{Description: "Display the date", ListEN: "Now", ListES: "Ahora", Command: "date"}, "Now", "Ahora"},
	} {
		items := []Item{tc.item}
		inheritTranslations(items, []Item{base})
		if got := items[0]; got.ListEN != tc.en || got.ListES != tc.es {
			t.Errorf("%+v: captions = %q and %q, want %q and %q", tc.item, got.ListEN, got.ListES, tc.en, tc.es)
		}
	}
}

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
