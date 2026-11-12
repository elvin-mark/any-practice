package ui

import (
	"strings"
	"testing"
)

func TestVisibleLen(t *testing.T) {
	colored := "\x1b[31mHello\x1b[0m"
	if l := VisibleLen(colored); l != 5 {
		t.Errorf("expected 5, got %d", l)
	}

	plain := "Docker"
	if l := VisibleLen(plain); l != 6 {
		t.Errorf("expected 6, got %d", l)
	}
}

func TestPadRight(t *testing.T) {
	padded := PadRight("abc", 6)
	if len(padded) != 6 || padded != "abc   " {
		t.Errorf("expected 'abc   ', got '%s'", padded)
	}
}

func TestPadLeft(t *testing.T) {
	padded := PadLeft("abc", 6)
	if len(padded) != 6 || padded != "   abc" {
		t.Errorf("expected '   abc', got '%s'", padded)
	}
}

func TestRenderTable(t *testing.T) {
	cols := []TableColumn{
		{Title: "ID", Width: 4},
		{Title: "Name", Width: 6},
	}
	rows := [][]string{
		{"1", "Alpine"},
		{"2", "Debian"},
	}

	tbl := RenderTable(cols, rows)
	if !strings.Contains(tbl, "Alpine") || !strings.Contains(tbl, "Debian") {
		t.Errorf("table output missing data: %s", tbl)
	}
	if !strings.Contains(tbl, "┌") || !strings.Contains(tbl, "└") {
		t.Errorf("table missing unicode border chars: %s", tbl)
	}
}
