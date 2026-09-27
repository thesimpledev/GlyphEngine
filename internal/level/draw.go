package level

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type GlyphDrawer interface {
	DrawGlyph(screen *ebiten.Image, ch rune, col, row int)
}

type sprite struct {
	image  *ebiten.Image
	offset image.Point
}

type glyphCache struct {
	face    font.Face
	sprites map[rune]sprite
	upload  func(image.Image) *ebiten.Image
	blit    func(screen, img *ebiten.Image, x, y int)
}

func newFace() font.Face {
	parsedFont, err := opentype.Parse(gomono.TTF)
	if err != nil {
		panic(fmt.Sprintf("level: parse embedded font: %v", err))
	}

	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingNone,
	})
	if err != nil {
		panic(fmt.Sprintf("level: create font face: %v", err))
	}
	return face
}

func newGlyphCache(face font.Face) *glyphCache {
	return &glyphCache{
		face:    face,
		sprites: make(map[rune]sprite),
		upload:  ebiten.NewImageFromImage,
		blit:    blitImage,
	}
}

func (g *glyphCache) DrawGlyph(screen *ebiten.Image, ch rune, col, row int) {
	s := g.sprite(ch)
	if s.image == nil {
		return
	}
	x, y := glyphOrigin(col, row)
	g.blit(screen, s.image, x+s.offset.X, y+s.offset.Y)
}

func (g *glyphCache) sprite(ch rune) sprite {
	if cached, ok := g.sprites[ch]; ok {
		return cached
	}

	var s sprite
	if img, offset, ok := rasterize(g.face, ch); ok {
		s = sprite{image: g.upload(img), offset: offset}
	}
	g.sprites[ch] = s
	return s
}

func rasterize(face font.Face, ch rune) (*image.RGBA, image.Point, bool) {
	bounds, _, ok := face.GlyphBounds(ch)
	if !ok {
		return nil, image.Point{}, false
	}

	offset := image.Pt(bounds.Min.X.Floor(), bounds.Min.Y.Floor())
	width := bounds.Max.X.Ceil() - offset.X
	height := bounds.Max.Y.Ceil() - offset.Y
	if width <= 0 || height <= 0 {
		return nil, image.Point{}, false
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	drawer := font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(-offset.X, -offset.Y),
	}
	drawer.DrawString(string(ch))
	return img, offset, true
}

func blitImage(screen, img *ebiten.Image, x, y int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(img, op)
}

func glyphOrigin(col, row int) (int, int) {
	return col * fontSize, (row + 1) * fontSize
}

func (l *Level) Draw(screen *ebiten.Image) {
	if l.Glyphs == nil {
		return
	}
	for y, row := range l.ViewGrid {
		for x, ch := range row {
			if ch == levelEmpty {
				continue
			}
			l.Glyphs.DrawGlyph(screen, ch, x, y)
		}
	}
}
