package level

import (
	"reflect"
	"strings"
	"testing"
)

func TestWrapLine(t *testing.T) {
	cases := []struct {
		name  string
		in    string
		width int
		want  []string
	}{
		{"fits", "hello", 10, []string{"hello"}},
		{"breaks at a space", "aaaa bbbb cccc", 9, []string{"aaaa bbbb", "cccc"}},
		{"space exactly at the limit", "aaaa bbbb cccc", 10, []string{"aaaa bbbb", "cccc"}},
		{"hard breaks a long word", "aaaaaaaaaaaa", 5, []string{"aaaaa", "aaaaa", "aa"}},
		{"counts runes", "héllo wörld", 5, []string{"héllo", "wörld"}},
		{"empty line keeps a row", "", 5, []string{""}},
		{"trailing space does not add a row", "aaaa bbbb ", 9, []string{"aaaa bbbb"}},
	}
	for _, tc := range cases {
		if got := wrapLine(tc.in, tc.width); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: got %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestSetDialogLayout(t *testing.T) {
	l := &Level{}
	l.resetLevelState()
	l.setDialog([]string{"Hi"}, []string{"> Yes", "  No"})

	if got := dialogRow(l, dialogIndent); got != "Hi" {
		t.Errorf("text row = %q, want %q", got, "Hi")
	}
	if got := dialogRow(l, 5); got != strings.Repeat("-", dialogWidth) {
		t.Errorf("separator row = %q", got)
	}
	if got := dialogRow(l, 8); got != "> Yes" {
		t.Errorf("first option row = %q, want %q", got, "> Yes")
	}
	if got := dialogRow(l, 9); got != "  No" {
		t.Errorf("second option row = %q, want %q", got, "  No")
	}
	if l.ViewGrid[dialogIndent][colDivider+dialogIndent-1] != levelEmpty {
		t.Error("dialog text leaked left of the indent")
	}
}

func TestSetDialogWrapsLongText(t *testing.T) {
	l := &Level{}
	l.resetLevelState()
	long := strings.Repeat("word ", 20)
	l.setDialog([]string{long}, nil)

	if dialogRow(l, dialogIndent) == "" || dialogRow(l, dialogIndent+1) == "" {
		t.Error("expected the long line to wrap onto a second row")
	}
	for y := range rowDivider {
		if len([]rune(dialogRow(l, y))) > dialogWidth {
			t.Errorf("row %d exceeds the dialog width", y)
		}
	}
}

func TestSetDialogStopsAtBottom(t *testing.T) {
	l := &Level{}
	l.resetLevelState()
	lines := make([]string, rows)
	for i := range lines {
		lines[i] = "x"
	}
	l.setDialog(lines, []string{"opt"})
	for y := rowDivider; y < rows; y++ {
		if dialogRow(l, y) != "" {
			t.Errorf("row %d below the divider was written", y)
		}
	}
}
