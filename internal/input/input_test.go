package input

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

var _ KeyReader = Keyboard{}

func TestKeyboardReportsNothingWithoutInput(t *testing.T) {
	var reader KeyReader = Keyboard{}
	if reader.IsKeyPressed(ebiten.KeySpace) {
		t.Error("IsKeyPressed reported SPACE with no input")
	}
	if reader.IsKeyJustPressed(ebiten.KeySpace) {
		t.Error("IsKeyJustPressed reported SPACE with no input")
	}
}
