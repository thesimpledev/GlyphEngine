package player

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

type fakeKeys struct {
	pressed map[ebiten.Key]bool
}

func (f *fakeKeys) IsKeyPressed(key ebiten.Key) bool     { return f.pressed[key] }
func (f *fakeKeys) IsKeyJustPressed(key ebiten.Key) bool { return f.pressed[key] }

type stubLevel struct {
	walkable bool
	boards   [][4]int
	cameras  [][2]int
}

func (s *stubLevel) IsWalkable(_, _ int) bool { return s.walkable }
func (s *stubLevel) UpdateBoard(fromX, fromY, toX, toY int) {
	s.boards = append(s.boards, [4]int{fromX, fromY, toX, toY})
}
func (s *stubLevel) UpdateCamera(x, y int) { s.cameras = append(s.cameras, [2]int{x, y}) }

func newTestPlayer(walkable bool, keys ...ebiten.Key) (*Player, *stubLevel, *fakeKeys) {
	pressed := make(map[ebiten.Key]bool, len(keys))
	for _, k := range keys {
		pressed[k] = true
	}
	fk := &fakeKeys{pressed: pressed}
	lvl := &stubLevel{walkable: walkable}
	p := New()
	p.Level = lvl
	p.Keys = fk
	p.SetPosition(5, 5)
	return p, lvl, fk
}

func TestNew(t *testing.T) {
	got := New()

	if got.x != 0 || got.y != 0 || got.movementCooldown != 0 {
		t.Errorf("Expected (x: 0, y: 0, movementCooldown: 0), got (x: %d, y: %d, movementCooldown: %d)",
			got.x, got.y, got.movementCooldown)
	}

	if len(got.walk) != 0 {
		t.Errorf("Expected walk to be nil or empty, got %+v", got.walk)
	}

	if got.Level != nil {
		t.Errorf("Expected Level to be nil, got %+v", got.Level)
	}

	if got.Keys == nil {
		t.Error("Expected Keys to default to the real keyboard")
	}
}

func TestUpdateMovesInEachDirection(t *testing.T) {
	cases := []struct {
		key  ebiten.Key
		x, y int
	}{
		{ebiten.KeyW, 5, 4},
		{ebiten.KeyA, 4, 5},
		{ebiten.KeyS, 5, 6},
		{ebiten.KeyD, 6, 5},
	}
	for _, tc := range cases {
		p, lvl, _ := newTestPlayer(true, tc.key)
		p.Update()
		if p.x != tc.x || p.y != tc.y {
			t.Errorf("key %v: got (%d,%d), want (%d,%d)", tc.key, p.x, p.y, tc.x, tc.y)
		}
		if len(lvl.boards) != 1 || lvl.boards[0] != [4]int{5, 5, tc.x, tc.y} {
			t.Errorf("key %v: UpdateBoard calls = %v", tc.key, lvl.boards)
		}
		if len(lvl.cameras) != 1 || lvl.cameras[0] != [2]int{tc.x, tc.y} {
			t.Errorf("key %v: UpdateCamera calls = %v", tc.key, lvl.cameras)
		}
	}
}

func TestUpdateLastBindingWins(t *testing.T) {
	p, _, _ := newTestPlayer(true, ebiten.KeyW, ebiten.KeyD)
	p.Update()
	if p.x != 6 || p.y != 5 {
		t.Errorf("got (%d,%d), want (6,5)", p.x, p.y)
	}
}

func TestUpdateCooldownBlocksRepeatMoves(t *testing.T) {
	p, lvl, _ := newTestPlayer(true, ebiten.KeyD)
	p.Update()
	for range moveCooldown {
		p.Update()
	}
	if len(lvl.boards) != 1 {
		t.Fatalf("moved %d times during cooldown, want 1", len(lvl.boards))
	}
	p.Update()
	if len(lvl.boards) != 2 {
		t.Fatalf("expected a second move after the cooldown, got %d", len(lvl.boards))
	}
}

func TestUpdateWithoutInputDoesNothing(t *testing.T) {
	p, lvl, _ := newTestPlayer(true)
	p.Update()
	if p.x != 5 || p.y != 5 || len(lvl.boards) != 0 || p.movementCooldown != 0 {
		t.Errorf("player changed without input: pos (%d,%d), boards %v, cooldown %d", p.x, p.y, lvl.boards, p.movementCooldown)
	}
}

func TestMoveBlockedByLevel(t *testing.T) {
	p, lvl, _ := newTestPlayer(false, ebiten.KeyS)
	p.Update()
	if p.x != 5 || p.y != 5 {
		t.Errorf("player moved into a blocked tile: (%d,%d)", p.x, p.y)
	}
	if len(lvl.boards) != 0 || len(lvl.cameras) != 0 {
		t.Errorf("level updated for a blocked move: boards %v cameras %v", lvl.boards, lvl.cameras)
	}
}

func TestMoveWithNilLevelDoesNotPanic(t *testing.T) {
	p := New()
	p.Keys = &fakeKeys{pressed: map[ebiten.Key]bool{ebiten.KeyW: true}}
	p.Update()
	if p.x != 0 || p.y != 0 {
		t.Errorf("player moved with no level: (%d,%d)", p.x, p.y)
	}
}

func TestPlayFootstepWithoutSoundsIsNoop(t *testing.T) {
	p := New()
	p.PlayFootstep()
}

func TestLoadFootstepsRejectsNilContext(t *testing.T) {
	p := New()
	if err := p.LoadFootsteps(nil, []byte{}); err == nil {
		t.Error("expected an error for a nil audio context")
	}
}

func TestNilKeysFallsBackToKeyboard(t *testing.T) {
	p := &Player{}
	if p.keys() == nil {
		t.Error("expected a fallback KeyReader")
	}
}
