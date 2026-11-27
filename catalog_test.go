package main

import (
	"io/fs"
	"reflect"
	"testing"
)

func TestParseFileReadsTranslatedHeadings(t *testing.T) {
	const content = "<!-- es: Archivos -->\n# Files\n\n" +
		"## Empty entry\n\n" +
		"<!-- es: Mostrar dos archivos -->\n## Show two files\n" +
		"```sh\ncat a.txt\n  cat b.txt\n```\n\n" +
		"## Show a third file\n`cat c.txt`\n"

	want := []Item{
		{
			Section:       "Files",
			SectionES:     "Archivos",
			Description:   "Show two files",
			DescriptionES: "Mostrar dos archivos",
			Command:       "cat a.txt\n  cat b.txt",
		},
		{
			Section:     "Files",
			SectionES:   "Archivos",
			Description: "Show a third file",
			Command:     "cat c.txt",
		},
	}
	if got := parseFile(content); !reflect.DeepEqual(got, want) {
		t.Errorf("parseFile() = %#v, want %#v", got, want)
	}
}

func TestParseFileAcceptsManualWithoutTranslations(t *testing.T) {
	const content = "# Local commands\n\n## Show the date\n`date`\n"
	want := []Item{{Section: "Local commands", Description: "Show the date", Command: "date"}}
	if got := parseFile(content); !reflect.DeepEqual(got, want) {
		t.Errorf("parseFile() = %#v, want %#v", got, want)
	}
}

func TestEmbeddedManualsHaveBilingualHeadings(t *testing.T) {
	names, err := fs.Glob(dataFS, "data/*.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range names {
		if name == "data/README.md" {
			continue
		}
		t.Run(name, func(t *testing.T) {
			content, err := dataFS.ReadFile(name)
			if err != nil {
				t.Fatal(err)
			}
			items := parseFile(string(content))
			if len(items) == 0 {
				t.Fatal("manual contains no entries")
			}
			for _, it := range items {
				if it.Section == "" || it.SectionES == "" || it.Description == "" || it.DescriptionES == "" {
					t.Errorf("missing English or Spanish title for %q / %q", it.Section, it.Description)
				}
			}
		})
	}
}
