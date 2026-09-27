package level

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/thesimpledev/GlyphEngine/internal/input"
)

//go:embed assets/txt/level1_dialog.json
var level1Dialog string

//go:embed assets/txt/level2_dialog.json
var level2Dialog string

//go:embed assets/txt/level3_dialog.json
var level3Dialog string

//go:embed assets/txt/level4_dialog.json
var level4Dialog string

//go:embed assets/txt/level5_dialog.json
var level5Dialog string

const (
	startState     = "start"
	nextLevel      = "next_level"
	endGame        = "end_game"
	endInteraction = "end_interaction"
	endDialog      = "end"
)

const dialogWidth = cols - colDivider - dialogIndent

type DialogState struct {
	Intro    []string `json:"intro"`
	Entities []Entity `json:"entities"`
}

type EntityDialogOption struct {
	Text      string `json:"text"`
	NextState string `json:"nextState,omitempty"`
}

type EntityDialogState struct {
	ID        string               `json:"id"`
	Text      string               `json:"text"`
	NextState string               `json:"nextState,omitempty"`
	Options   []EntityDialogOption `json:"options,omitempty"`
}

type Entity struct {
	ID           string              `json:"id"`
	Symbol       rune                `json:"symbol"`
	X            int                 `json:"x"`
	Y            int                 `json:"y"`
	Text         string              `json:"text,omitempty"`
	DialogStates []EntityDialogState `json:"dialogStates,omitempty"`
	active       bool
}

func (e *Entity) UnmarshalJSON(data []byte) error {
	type alias Entity
	var tmp struct {
		Symbol string `json:"symbol"`
		*alias
	}
	tmp.alias = (*alias)(e)

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	runes := []rune(tmp.Symbol)
	if len(runes) == 0 {
		e.Symbol = 0
	} else {
		e.Symbol = runes[0]
	}

	return nil
}

func isTerminal(next string) bool {
	return next == endInteraction || next == endDialog
}

func (l *Level) keys() input.KeyReader {
	if l.Keys == nil {
		return input.Keyboard{}
	}
	return l.Keys
}

func (l *Level) confirmPressed() bool {
	return l.keys().IsKeyJustPressed(ebiten.KeySpace) || l.keys().IsKeyJustPressed(ebiten.KeyEnter)
}

func (l *Level) loadDialog() {
	var dialogState DialogState
	if err := json.Unmarshal([]byte(levelDialogs[l.level]), &dialogState); err != nil {
		panic(fmt.Sprintf("level %d: parse dialog: %v", l.level, err))
	}

	for i := range dialogState.Entities {
		l.placeEntity(&dialogState.Entities[i])
	}

	l.setDialog(dialogState.Intro, []string{"Press SPACE to continue"})
}

func (l *Level) placeEntity(entity *Entity) {
	if !l.inBounds(entity.X, entity.Y) {
		panic(fmt.Sprintf("level %d: entity %q at (%d,%d) is outside the map", l.level, entity.ID, entity.X, entity.Y))
	}
	l.entities[gridPos{entity.X, entity.Y}] = entity
	l.obstacles[entity.Symbol] = struct{}{}
	l.MapGrid[entity.Y][entity.X] = entity.Symbol
	if entity.ID == "Player" {
		l.playerX = entity.X
		l.playerY = entity.Y
		l.playerSymbol = entity.Symbol
		l.Player.SetPosition(entity.X, entity.Y)
	}
}

func (l *Level) updateDialog() {
	if l.showingIntro {
		if l.confirmPressed() {
			l.showingIntro = false
			l.clearDialogArea()
		}
		return
	}

	if l.currentEntity == nil {
		return
	}

	if !l.currentEntity.active {
		l.activateCurrentEntity()
		return
	}

	l.processActivatedEntity()
}

func (l *Level) activateCurrentEntity() {
	if len(l.currentEntity.DialogStates) == 0 || !l.confirmPressed() {
		return
	}
	l.currentEntity.active = true
	l.enterState(startState)
}

func (l *Level) processActivatedEntity() {
	state := l.getCurrentEntityState()
	if state == nil {
		return
	}

	if l.confirmPressed() {
		l.transition(chosenNext(state, l.selectedDialog))
		return
	}

	l.moveSelection(len(state.Options))
}

func chosenNext(state *EntityDialogState, selected int) string {
	if len(state.Options) == 0 {
		return state.NextState
	}
	if selected < 0 || selected >= len(state.Options) {
		return ""
	}
	return state.Options[selected].NextState
}

func (l *Level) transition(next string) {
	switch {
	case next == nextLevel:
		l.advanceLevel()
	case next == endGame:
		l.level = 0
		l.LoadLevel()
	case isTerminal(next):
		l.endInteraction()
	default:
		l.enterState(next)
	}
}

func (l *Level) advanceLevel() {
	if l.level+1 >= len(levelMaps) {
		log.Printf("level %d: next_level requested on the last level, restarting", l.level)
		l.level = 0
	} else {
		l.level++
	}
	l.LoadLevel()
}

func (l *Level) endInteraction() {
	l.currentEntity.active = false
	l.currentDialogStateID = ""
	l.clearDialogArea()
}

func (l *Level) enterState(id string) {
	l.currentDialogStateID = id
	l.selectedDialog = 0
	if l.getCurrentEntityState() == nil {
		log.Printf("level %d: entity %q has no dialog state %q", l.level, l.currentEntity.ID, id)
		l.endInteraction()
		return
	}
	l.setDialogForCurrentState()
}

func (l *Level) moveSelection(optionCount int) {
	if optionCount == 0 {
		return
	}

	switch {
	case l.keys().IsKeyJustPressed(ebiten.KeyArrowDown):
		l.selectedDialog = (l.selectedDialog + 1) % optionCount
	case l.keys().IsKeyJustPressed(ebiten.KeyArrowUp):
		l.selectedDialog = (l.selectedDialog + optionCount - 1) % optionCount
	default:
		return
	}
	l.setDialogForCurrentState()
}

func (l *Level) clearDialogArea() {
	for y := range rowDivider {
		for x := colDivider; x < cols; x++ {
			l.ViewGrid[y][x] = levelEmpty
		}
	}
}

func (l *Level) setDialog(text []string, options []string) {
	row := l.writeWrapped(dialogIndent, text)
	row = l.writeSeparator(row + dialogIndent)
	l.writeOptions(row+dialogIndent, options)
}

func (l *Level) writeWrapped(row int, text []string) int {
	for _, line := range text {
		for _, part := range wrapLine(line, dialogWidth) {
			if row >= rowDivider {
				return row
			}
			l.writeDialogLine(row, part)
			row++
		}
	}
	return row
}

func (l *Level) writeSeparator(row int) int {
	if row >= rowDivider {
		return row
	}
	l.writeDialogLine(row, strings.Repeat("-", dialogWidth))
	return row + 1
}

func (l *Level) writeOptions(row int, options []string) {
	for _, opt := range options {
		if row >= rowDivider {
			return
		}
		l.writeDialogLine(row, opt)
		row++
	}
}

func (l *Level) writeDialogLine(row int, s string) {
	col := colDivider + dialogIndent
	for _, r := range s {
		if col >= cols {
			return
		}
		l.ViewGrid[row][col] = r
		col++
	}
}

func wrapLine(line string, width int) []string {
	if width <= 0 {
		return nil
	}
	runes := []rune(line)
	var parts []string
	for len(runes) > width {
		cut := breakIndex(runes, width)
		parts = append(parts, string(runes[:cut]))
		runes = trimLeadingSpaces(runes[cut:])
	}
	if len(parts) > 0 && len(runes) == 0 {
		return parts
	}
	return append(parts, string(runes))
}

func breakIndex(runes []rune, width int) int {
	for i := width; i > 0; i-- {
		if runes[i] == ' ' {
			return i
		}
	}
	return width
}

func trimLeadingSpaces(runes []rune) []rune {
	for len(runes) > 0 && runes[0] == ' ' {
		runes = runes[1:]
	}
	return runes
}

func (l *Level) getCurrentEntityState() *EntityDialogState {
	if l.currentEntity == nil {
		return nil
	}
	for i := range l.currentEntity.DialogStates {
		if l.currentEntity.DialogStates[i].ID == l.currentDialogStateID {
			return &l.currentEntity.DialogStates[i]
		}
	}
	return nil
}

func (l *Level) setDialogForCurrentState() {
	l.clearDialogArea()
	state := l.getCurrentEntityState()
	if state == nil {
		return
	}

	optionTexts := make([]string, 0, len(state.Options))
	for i, opt := range state.Options {
		prefix := "  "
		if i == l.selectedDialog {
			prefix = "> "
		}
		optionTexts = append(optionTexts, prefix+opt.Text)
	}

	l.setDialog([]string{state.Text}, optionTexts)
}
