package game

import "testing"

type counter struct {
	updates int
}

func (c *counter) Update() { c.updates++ }

type selfRemover struct {
	game    *Game
	updates int
}

func (s *selfRemover) Update() {
	s.updates++
	s.game.RemoveComponent(s)
}

func TestUpdateCallsUpdatableComponents(t *testing.T) {
	g := &Game{}
	c := &counter{}
	g.AddComponent(c)
	g.AddComponent("not updatable")

	if err := g.Update(); err != nil {
		t.Fatalf("Update returned %v", err)
	}
	if c.updates != 1 {
		t.Errorf("updates = %d, want 1", c.updates)
	}
	if len(g.Components) != 2 {
		t.Errorf("components = %d, want 2", len(g.Components))
	}
}

func TestRemoveComponentDropsItOnNextUpdate(t *testing.T) {
	g := &Game{}
	a := &counter{}
	b := &counter{}
	g.AddComponent(a)
	g.AddComponent(b)

	g.RemoveComponent(a)
	if err := g.Update(); err != nil {
		t.Fatalf("Update returned %v", err)
	}

	if a.updates != 0 || b.updates != 1 {
		t.Errorf("updates a=%d b=%d, want 0 and 1", a.updates, b.updates)
	}
	if len(g.Components) != 1 || g.Components[0] != b {
		t.Errorf("components = %v, want only b", g.Components)
	}
}

func TestRemoveUnknownComponentIsNoop(t *testing.T) {
	g := &Game{}
	a := &counter{}
	g.AddComponent(a)
	g.RemoveComponent(&counter{})
	if len(g.Components) != 1 || g.Components[0] != a {
		t.Errorf("components = %v, want only a", g.Components)
	}
}

func TestComponentCanRemoveItselfDuringUpdate(t *testing.T) {
	g := &Game{}
	s := &selfRemover{game: g}
	after := &counter{}
	g.AddComponent(s)
	g.AddComponent(after)

	if err := g.Update(); err != nil {
		t.Fatalf("Update returned %v", err)
	}
	if err := g.Update(); err != nil {
		t.Fatalf("Update returned %v", err)
	}

	if s.updates != 1 {
		t.Errorf("self-removing component updated %d times, want 1", s.updates)
	}
	if after.updates != 2 {
		t.Errorf("later component updated %d times, want 2", after.updates)
	}
	if len(g.Components) != 1 {
		t.Errorf("components = %d, want 1", len(g.Components))
	}
}

func TestLayoutReturnsConfiguredSize(t *testing.T) {
	g := &Game{screenWidth: 1280, screenHeight: 720}
	w, h := g.Layout(1, 1)
	if w != 1280 || h != 720 {
		t.Errorf("Layout = (%d,%d), want (1280,720)", w, h)
	}
}
