package level

import (
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

type blitCall struct {
	img  *ebiten.Image
	x, y int
}

type glyphProbe struct {
	uploads int
	blits   []blitCall
}

func newProbedCache() (*glyphCache, *glyphProbe) {
	probe := &glyphProbe{}
	cache := newGlyphCache(newFace())
	cache.upload = func(image.Image) *ebiten.Image {
		probe.uploads++
		return new(ebiten.Image)
	}
	cache.blit = func(_, img *ebiten.Image, x, y int) {
		probe.blits = append(probe.blits, blitCall{img, x, y})
	}
	return cache, probe
}

func hasInk(img image.Image) bool {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if _, _, _, a := img.At(x, y).RGBA(); a != 0 {
				return true
			}
		}
	}
	return false
}

func TestRasterizeDrawsVisibleGlyph(t *testing.T) {
	img, offset, ok := rasterize(newFace(), '#')
	if !ok {
		t.Fatal("rasterize reported no glyph for '#'")
	}
	if img.Bounds().Dx() <= 0 || img.Bounds().Dy() <= 0 {
		t.Fatalf("glyph image is empty: %v", img.Bounds())
	}
	if img.Bounds().Dx() > 2*fontSize || img.Bounds().Dy() > 2*fontSize {
		t.Errorf("glyph image is implausibly large: %v", img.Bounds())
	}
	if !hasInk(img) {
		t.Error("glyph image has no visible pixels")
	}
	if offset.Y >= 0 {
		t.Errorf("offset.Y = %d, want the glyph to start above the baseline", offset.Y)
	}
}

func TestRasterizeHandlesNonASCII(t *testing.T) {
	img, _, ok := rasterize(newFace(), 'é')
	if !ok || !hasInk(img) {
		t.Error("expected ink for 'é'")
	}
}

func TestRasterizeBlankGlyph(t *testing.T) {
	if _, _, ok := rasterize(newFace(), ' '); ok {
		t.Error("rasterize produced an image for a space")
	}
}

func TestGlyphCacheUploadsEachRuneOnce(t *testing.T) {
	cache, probe := newProbedCache()

	cache.DrawGlyph(nil, '#', 0, 0)
	cache.DrawGlyph(nil, '#', 1, 0)
	cache.DrawGlyph(nil, '@', 2, 0)
	cache.DrawGlyph(nil, '#', 3, 0)

	if probe.uploads != 2 {
		t.Errorf("uploads = %d, want 2", probe.uploads)
	}
	if len(probe.blits) != 4 {
		t.Fatalf("blits = %d, want 4", len(probe.blits))
	}
	if probe.blits[0].img != probe.blits[1].img {
		t.Error("the same rune was drawn from two different images")
	}
	if probe.blits[0].img == probe.blits[2].img {
		t.Error("two different runes share one image")
	}
}

func TestGlyphCacheSkipsBlankGlyphs(t *testing.T) {
	cache, probe := newProbedCache()

	cache.DrawGlyph(nil, ' ', 0, 0)
	cache.DrawGlyph(nil, ' ', 1, 0)

	if probe.uploads != 0 || len(probe.blits) != 0 {
		t.Errorf("blank glyph caused uploads=%d blits=%d", probe.uploads, len(probe.blits))
	}
	if len(cache.sprites) != 1 {
		t.Errorf("blank glyph cached %d entries, want 1", len(cache.sprites))
	}
}

func TestGlyphCachePlacesGlyphRelativeToBaseline(t *testing.T) {
	cache, probe := newProbedCache()
	_, offset, _ := rasterize(newFace(), '#')

	cache.DrawGlyph(nil, '#', 3, 2)

	originX, originY := glyphOrigin(3, 2)
	got := probe.blits[0]
	if got.x != originX+offset.X || got.y != originY+offset.Y {
		t.Errorf("drawn at (%d,%d), want (%d,%d)", got.x, got.y, originX+offset.X, originY+offset.Y)
	}
}

func TestGlyphCacheMovesOneCellPerColumnAndRow(t *testing.T) {
	cache, probe := newProbedCache()

	cache.DrawGlyph(nil, '#', 0, 0)
	cache.DrawGlyph(nil, '#', 1, 1)

	dx := probe.blits[1].x - probe.blits[0].x
	dy := probe.blits[1].y - probe.blits[0].y
	if dx != fontSize || dy != fontSize {
		t.Errorf("cell step = (%d,%d), want (%d,%d)", dx, dy, fontSize, fontSize)
	}
}

func TestNewUsesGlyphCache(t *testing.T) {
	if _, ok := New().Glyphs.(*glyphCache); !ok {
		t.Errorf("New().Glyphs is %T, want *glyphCache", New().Glyphs)
	}
}
