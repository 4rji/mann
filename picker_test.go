package main

import (
	"strings"
	"testing"

	"github.com/rivo/uniseg"
)

func TestActionColumnKeepsCommandAligned(t *testing.T) {
	for _, caption := range []string{
		"Añadir un archivo al tar",
		strings.Repeat("a", 34),
		strings.Repeat("界", 16) + "e\u0301" + "👩‍💻",
		strings.Repeat("a", 33) + "界",
	} {
		column := actionColumn(caption)
		if got := uniseg.StringWidth(column); got != 34 {
			t.Errorf("%q occupies %d cells", column, got)
		}
		if uniseg.StringWidth(caption) > 34 && !strings.Contains(column, "…") {
			t.Errorf("missing truncation marker for %q", caption)
		}
	}
}
