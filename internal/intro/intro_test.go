package intro

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

const ticksPerScene = 120

type fakeLevel struct {
	loads int
}

func (f *fakeLevel) LoadLevel() { f.loads++ }

type fakeGame struct {
	removed []any
}

func (f *fakeGame) RemoveComponent(c any) { f.removed = append(f.removed, c) }

func newTestIntro(scenes int) (*Intro, *fakeLevel, *fakeGame) {
	lvl := &fakeLevel{}
	g := &fakeGame{}
	return &Intro{images: make([]*ebiten.Image, scenes), Level: lvl, Game: g}, lvl, g
}

func TestUpdateLoadsLevelAfterLastScene(t *testing.T) {
	i, lvl, g := newTestIntro(1)

	for range ticksPerScene - 1 {
		i.Update()
	}
	if lvl.loads != 0 || len(g.removed) != 0 {
		t.Fatalf("intro ended early: loads=%d removed=%d", lvl.loads, len(g.removed))
	}

	i.Update()
	if lvl.loads != 1 {
		t.Errorf("loads = %d, want 1", lvl.loads)
	}
	if len(g.removed) != 1 || g.removed[0] != i {
		t.Errorf("removed = %v, want the intro itself", g.removed)
	}
}

func TestUpdateAdvancesThroughScenes(t *testing.T) {
	i, lvl, _ := newTestIntro(2)

	for range ticksPerScene {
		i.Update()
	}
	if i.imageIndex != 1 || lvl.loads != 0 {
		t.Fatalf("after scene one: index=%d loads=%d", i.imageIndex, lvl.loads)
	}

	for range ticksPerScene {
		i.Update()
	}
	if lvl.loads != 1 {
		t.Errorf("loads = %d, want 1", lvl.loads)
	}
}

func TestDrawAfterLastSceneDoesNothing(t *testing.T) {
	i, _, _ := newTestIntro(1)
	i.imageIndex = 1
	i.Draw(nil)
}
