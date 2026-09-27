package level

import (
	"strings"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestLoadLevelPlacesPlayerAndShowsIntro(t *testing.T) {
	l, _, player := newTestLevel(t, 0)

	if player.x != 5 || player.y != 1 {
		t.Errorf("player placed at (%d,%d), want (5,1)", player.x, player.y)
	}
	if l.MapGrid[1][5] != '@' {
		t.Errorf("map cell = %q, want '@'", l.MapGrid[1][5])
	}
	if !l.showingIntro || dialogRow(l, dialogIndent) != "Welcome To your apartment Anon" {
		t.Errorf("intro not shown: showingIntro=%v row=%q", l.showingIntro, dialogRow(l, dialogIndent))
	}
}

func TestMovementBlockedDuringIntro(t *testing.T) {
	l, keys, _ := newTestLevel(t, 0)

	if l.IsWalkable(5, 2) {
		t.Fatal("walkable during the intro")
	}
	tap(l, keys, ebiten.KeySpace)
	if l.showingIntro || !dialogAreaBlank(l) {
		t.Fatal("intro not dismissed by SPACE")
	}
	if !l.IsWalkable(5, 2) {
		t.Error("open floor not walkable after the intro")
	}
}

func TestOutOfBoundsIsNotWalkable(t *testing.T) {
	l, keys, _ := newTestLevel(t, 0)
	tap(l, keys, ebiten.KeySpace)
	if l.IsWalkable(-1, 0) || l.IsWalkable(0, -1) || l.IsWalkable(0, len(l.MapGrid)) || l.IsWalkable(len(l.MapGrid[0]), 0) {
		t.Error("out-of-range coordinates reported walkable")
	}
}

func TestUpdateBoardMovesPlayerSymbol(t *testing.T) {
	l, keys, _ := newTestLevel(t, 0)
	tap(l, keys, ebiten.KeySpace)
	l.UpdateBoard(5, 1, 5, 2)
	if l.MapGrid[1][5] != levelEmpty || l.MapGrid[2][5] != '@' {
		t.Errorf("board after move: from=%q to=%q", l.MapGrid[1][5], l.MapGrid[2][5])
	}
	l.UpdateBoard(5, 2, 99, 99)
	if l.MapGrid[2][5] != '@' {
		t.Error("out-of-bounds move changed the board")
	}
}

func TestNextLevelResetsEntityAndDoesNotSkip(t *testing.T) {
	l, keys, _ := newTestLevel(t, 0)
	tap(l, keys, ebiten.KeySpace)

	if l.IsWalkable(5, 0) {
		t.Fatal("exit should block movement")
	}
	if l.currentEntity == nil || l.currentEntity.ID != "Exit" {
		t.Fatalf("currentEntity = %+v, want Exit", l.currentEntity)
	}
	tap(l, keys, ebiten.KeySpace)
	if got := dialogRow(l, dialogIndent); got != "Head to Acme Corp?" {
		t.Fatalf("start state text = %q", got)
	}
	tap(l, keys, ebiten.KeySpace)

	if l.level != 1 {
		t.Fatalf("level = %d after choosing Yes, want 1", l.level)
	}
	if l.currentEntity != nil || l.currentDialogStateID != "" || l.showingItem {
		t.Errorf("stale interaction state survived the level change: entity=%v id=%q showingItem=%v", l.currentEntity, l.currentDialogStateID, l.showingItem)
	}
	if !l.showingIntro || dialogRow(l, dialogIndent) != "Welcome Acme Corp Headquarters!" {
		t.Errorf("level 2 intro not shown: %q", dialogRow(l, dialogIndent))
	}

	for range 5 {
		tap(l, keys, ebiten.KeySpace)
	}
	if l.level != 1 {
		t.Errorf("holding SPACE after the intro moved to level %d", l.level)
	}
}

func TestEndGameRestartsWithIntroAndCamera(t *testing.T) {
	l, keys, _ := newTestLevel(t, 4)
	l.cameraY = 14
	tap(l, keys, ebiten.KeySpace)

	if l.IsWalkable(5, 0) {
		t.Fatal("exit should block movement")
	}
	tap(l, keys, ebiten.KeySpace)
	tap(l, keys, ebiten.KeySpace)

	if l.level != 0 {
		t.Fatalf("level = %d, want 0", l.level)
	}
	if !l.showingIntro || dialogRow(l, dialogIndent) != "Welcome To your apartment Anon" {
		t.Errorf("intro text missing after restart: %q", dialogRow(l, dialogIndent))
	}
	if l.cameraX != 0 || l.cameraY != 0 {
		t.Errorf("camera not reset: (%d,%d)", l.cameraX, l.cameraY)
	}
	if l.currentEntity != nil {
		t.Error("currentEntity not cleared on restart")
	}
}

func TestOptionlessEndInteraction(t *testing.T) {
	l, keys, _ := newTestLevel(t, 0)
	tap(l, keys, ebiten.KeySpace)

	l.IsWalkable(5, 5)
	phone := l.currentEntity
	if phone == nil || phone.ID != "Phone" {
		t.Fatalf("currentEntity = %+v, want Phone", phone)
	}
	tap(l, keys, ebiten.KeySpace)
	tap(l, keys, ebiten.KeySpace)
	if got := dialogRow(l, dialogIndent); !strings.HasPrefix(got, "You order lunch") {
		t.Fatalf("lunch state text = %q", got)
	}
	tap(l, keys, ebiten.KeySpace)

	if phone.active || !dialogAreaBlank(l) {
		t.Errorf("optionless end_interaction left the entity active=%v blank=%v", phone.active, dialogAreaBlank(l))
	}
}

func TestOptionEndKeywordEndsInteraction(t *testing.T) {
	l, keys, _ := newTestLevel(t, 4)
	tap(l, keys, ebiten.KeySpace)
	l.IsWalkable(5, 0)
	exit := l.currentEntity
	tap(l, keys, ebiten.KeySpace)
	tap(l, keys, ebiten.KeyArrowDown)
	if got := dialogRow(l, 9); got != "> No" {
		t.Fatalf("second option row = %q, want %q", got, "> No")
	}
	tap(l, keys, ebiten.KeySpace)
	if exit.active || l.level != 4 {
		t.Errorf("'end' option: active=%v level=%d", exit.active, l.level)
	}
}

func TestArrowKeysWrapSelection(t *testing.T) {
	l, keys, _ := newTestLevel(t, 0)
	tap(l, keys, ebiten.KeySpace)
	l.IsWalkable(5, 5)
	tap(l, keys, ebiten.KeySpace)

	tap(l, keys, ebiten.KeyArrowUp)
	if l.selectedDialog != 1 {
		t.Errorf("Up from 0 selected %d, want 1", l.selectedDialog)
	}
	tap(l, keys, ebiten.KeyArrowDown)
	if l.selectedDialog != 0 {
		t.Errorf("Down from last selected %d, want 0", l.selectedDialog)
	}
	if got := dialogRow(l, 8); got != "> Order Lunch" {
		t.Errorf("selected option row = %q", got)
	}
}

func TestDescriptionOnlyEntityKeepsText(t *testing.T) {
	l, keys, _ := newTestLevel(t, 2)
	tap(l, keys, ebiten.KeySpace)

	l.IsWalkable(8, 28)
	deskText := func() string {
		return dialogRow(l, dialogIndent) + " " + dialogRow(l, dialogIndent+1)
	}
	want := "The Computer On with an email showing: Executive Code 1234"
	if got := deskText(); got != want {
		t.Fatalf("desk text = %q", got)
	}
	for y := range rowDivider {
		if strings.Contains(dialogRow(l, y), "Press SPACE") {
			t.Errorf("row %d prompts to interact with a description-only entity", y)
		}
	}
	tap(l, keys, ebiten.KeySpace)
	if got := deskText(); got != want {
		t.Errorf("desk text after SPACE = %q, want it unchanged", got)
	}
}

func TestWalkingAwayDeactivatesEntity(t *testing.T) {
	l, keys, _ := newTestLevel(t, 0)
	tap(l, keys, ebiten.KeySpace)
	l.IsWalkable(5, 5)
	phone := l.currentEntity
	tap(l, keys, ebiten.KeySpace)
	if !phone.active {
		t.Fatal("phone not activated")
	}
	l.IsWalkable(5, 2)
	if phone.active || l.currentEntity != nil || !dialogAreaBlank(l) {
		t.Error("walking away did not close the interaction")
	}
	l.IsWalkable(5, 5)
	if got := dialogRow(l, dialogIndent); got != "A Telephone Sits here, call someone?" {
		t.Errorf("returning to the phone shows %q", got)
	}
}

func TestUnknownStateEndsInteraction(t *testing.T) {
	l := &Level{}
	l.resetLevelState()
	l.MapGrid = parseMap("###\n# #\n###")
	entity := &Entity{ID: "Test", DialogStates: []EntityDialogState{{
		ID:      startState,
		Text:    "?",
		Options: []EntityDialogOption{{Text: "go", NextState: "missing"}},
	}}}
	l.currentEntity = entity
	entity.active = true
	l.enterState(startState)
	l.transition("missing")
	if entity.active || !dialogAreaBlank(l) {
		t.Error("unknown state left the entity active")
	}
}

func TestEntityKeysDoNotCollide(t *testing.T) {
	l := &Level{Player: &stubPlayer{}}
	l.resetLevelState()
	l.MapGrid = parseMap(testMap(20, 20))
	a := &Entity{ID: "A", Symbol: 'A', X: 12, Y: 1}
	b := &Entity{ID: "B", Symbol: 'B', X: 2, Y: 11}
	l.placeEntity(a)
	l.placeEntity(b)
	if l.entities[gridPos{12, 1}] != a || l.entities[gridPos{2, 11}] != b {
		t.Error("entities at (12,1) and (2,11) collided")
	}
}

func TestWallShowsDescriptionWithoutPrompt(t *testing.T) {
	l, keys, _ := newTestLevel(t, 0)
	tap(l, keys, ebiten.KeySpace)
	if l.IsWalkable(0, 1) {
		t.Fatal("wall should block")
	}
	if got := dialogRow(l, dialogIndent); got != wallText[0] {
		t.Errorf("wall text = %q", got)
	}
	if l.currentEntity != nil {
		t.Error("wall set a currentEntity")
	}
}

func TestLoadLevelPanicsWithoutPlayer(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected a panic when Player is nil")
		}
	}()
	l := &Level{}
	l.LoadLevel()
}
