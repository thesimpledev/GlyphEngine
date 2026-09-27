package level

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

type drawnGlyph struct {
	ch       rune
	col, row int
}

type fakeGlyphs struct {
	drawn []drawnGlyph
}

func (f *fakeGlyphs) DrawGlyph(_ *ebiten.Image, ch rune, col, row int) {
	f.drawn = append(f.drawn, drawnGlyph{ch, col, row})
}

func TestDrawSendsOnlyVisibleCells(t *testing.T) {
	glyphs := &fakeGlyphs{}
	l := &Level{Glyphs: glyphs}
	l.resetLevelState()
	l.ViewGrid[0][0] = '#'
	l.ViewGrid[3][7] = '@'
	l.ViewGrid[rows-1][cols-1] = 'é'

	l.Draw(nil)

	want := []drawnGlyph{{'#', 0, 0}, {'@', 7, 3}, {'é', cols - 1, rows - 1}}
	if len(glyphs.drawn) != len(want) {
		t.Fatalf("drew %d glyphs %v, want %d", len(glyphs.drawn), glyphs.drawn, len(want))
	}
	for i, w := range want {
		if glyphs.drawn[i] != w {
			t.Errorf("glyph %d = %+v, want %+v", i, glyphs.drawn[i], w)
		}
	}
}

func TestDrawBeforeLoadDrawsNothing(t *testing.T) {
	glyphs := &fakeGlyphs{}
	l := &Level{Glyphs: glyphs}
	l.Draw(nil)
	if len(glyphs.drawn) != 0 {
		t.Errorf("drew %v with no grid", glyphs.drawn)
	}
}

func TestDrawWithoutGlyphDrawerDoesNotPanic(t *testing.T) {
	l := &Level{}
	l.resetLevelState()
	l.ViewGrid[0][0] = '#'
	l.Draw(nil)
}

func TestGlyphOrigin(t *testing.T) {
	cases := []struct {
		col, row     int
		wantX, wantY int
	}{
		{0, 0, 0, fontSize},
		{1, 0, fontSize, fontSize},
		{79, 44, 79 * fontSize, 45 * fontSize},
	}
	for _, tc := range cases {
		x, y := glyphOrigin(tc.col, tc.row)
		if x != tc.wantX || y != tc.wantY {
			t.Errorf("glyphOrigin(%d,%d) = (%d,%d), want (%d,%d)", tc.col, tc.row, x, y, tc.wantX, tc.wantY)
		}
	}
}

func TestNewSetsGlyphDrawer(t *testing.T) {
	if New().Glyphs == nil {
		t.Error("New left Glyphs unset")
	}
}
