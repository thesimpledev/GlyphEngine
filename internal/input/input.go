package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type KeyReader interface {
	IsKeyPressed(key ebiten.Key) bool
	IsKeyJustPressed(key ebiten.Key) bool
}

type Keyboard struct{}

func (Keyboard) IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}

func (Keyboard) IsKeyJustPressed(key ebiten.Key) bool {
	return inpututil.IsKeyJustPressed(key)
}
