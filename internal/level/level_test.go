package level

import (
	"strings"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

type stubPlayer struct {
	x, y int
}

func (s *stubPlayer) SetPosition(x, y int) {
	s.x = x
	s.y = y
}

type fakeKeys struct {
	pressed map[ebiten.Key]bool
}

func (f *fakeKeys) IsKeyPressed(key ebiten.Key) bool     { return f.pressed[key] }
func (f *fakeKeys) IsKeyJustPressed(key ebiten.Key) bool { return f.pressed[key] }

func (f *fakeKeys) press(key ebiten.Key) {
	f.pressed = map[ebiten.Key]bool{key: true}
}

func (f *fakeKeys) release() {
	f.pressed = nil
}

func newTestLevel(t *testing.T, index int) (*Level, *fakeKeys, *stubPlayer) {
	t.Helper()
	keys := &fakeKeys{}
	player := &stubPlayer{}
	l := &Level{Player: player, Keys: keys, level: index}
	l.LoadLevel()
	return l, keys, player
}

func tap(l *Level, keys *fakeKeys, key ebiten.Key) {
	keys.press(key)
	l.Update()
	keys.release()
}

func dialogRow(l *Level, row int) string {
	return strings.TrimRight(string(l.ViewGrid[row][colDivider+dialogIndent:]), " ")
}

func dialogAreaBlank(l *Level) bool {
	for y := range rowDivider {
		if dialogRow(l, y) != "" {
			return false
		}
	}
	return true
}
