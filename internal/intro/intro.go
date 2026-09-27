package intro

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"math"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
)

type Intro struct {
	images     []*ebiten.Image
	timer      int
	imageIndex int
	Level      Level
	Game       Game
}

type Level interface {
	LoadLevel()
}

type Game interface {
	RemoveComponent(c any)
}

func New() *Intro {
	i := &Intro{}
	i.loadIntro()
	return i
}

//go:embed img/4.png
var aptImage []byte

func (i *Intro) loadIntro() {
	aptImageData, _, err := image.Decode(bytes.NewReader(aptImage))
	if err != nil {
		panic(fmt.Sprintf("intro: decode embedded image: %v", err))
	}

	i.images = []*ebiten.Image{
		ebiten.NewImageFromImage(aptImageData),
	}
}

func (i *Intro) Update() {
	i.timer++
	if i.timer%120 != 0 {
		return
	}
	i.imageIndex++
	if i.imageIndex >= len(i.images) {
		i.Level.LoadLevel()
		i.Game.RemoveComponent(i)
	}
}

func (i *Intro) Draw(screen *ebiten.Image) {
	if i.imageIndex >= len(i.images) {
		return
	}
	screen.Fill(color.RGBA{0x38, 0x39, 0x3a, 0xff})

	img := i.images[i.imageIndex]
	screenW := float64(screen.Bounds().Dx())
	screenH := float64(screen.Bounds().Dy())
	imgW := float64(img.Bounds().Dx())
	imgH := float64(img.Bounds().Dy())

	scale := math.Min(screenW/imgW, screenH/imgH)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate((screenW-imgW*scale)/2, (screenH-imgH*scale)/2)
	screen.DrawImage(img, op)
}
