package level

import (
	_ "embed"
	"fmt"
	"log"
	"strings"

	"github.com/thesimpledev/GlyphEngine/internal/input"
)

//go:embed assets/txt/level1_map.txt
var level1Map string

//go:embed assets/txt/level2_map.txt
var level2Map string

//go:embed assets/txt/level3_map.txt
var level3Map string

//go:embed assets/txt/level4_map.txt
var level4Map string

//go:embed assets/txt/level5_map.txt
var level5Map string

type Player interface {
	SetPosition(x, y int)
}

const (
	fontSize   = 16
	dpi        = 72
	cols       = 80
	rows       = 45
	colDivider = 30
	rowDivider = 40
)

const (
	levelEmpty   = ' '
	levelWall    = '#'
	dialogIndent = 2
)

var (
	levelMaps = []string{
		level1Map,
		level2Map,
		level3Map,
		level4Map,
		level5Map,
	}
	levelDialogs = []string{
		level1Dialog,
		level2Dialog,
		level3Dialog,
		level4Dialog,
		level5Dialog,
	}
)

type gridPos struct {
	x int
	y int
}

type Level struct {
	Glyphs               GlyphDrawer
	ViewGrid             [][]rune
	MapGrid              [][]rune
	Player               Player
	Keys                 input.KeyReader
	level                int
	selectedDialog       int
	showingIntro         bool
	showingItem          bool
	disableInput         bool
	entities             map[gridPos]*Entity
	obstacles            map[rune]struct{}
	currentEntity        *Entity
	currentDialogStateID string
	cameraX              int
	cameraY              int
	playerX              int
	playerY              int
	playerSymbol         rune
}

func New() *Level {
	return &Level{Glyphs: newGlyphCache(newFace()), disableInput: true, Keys: input.Keyboard{}}
}

func (l *Level) Update() {
	l.updateDialog()
}

func (l *Level) LoadLevel() {
	if l.level < 0 || l.level >= len(levelMaps) {
		panic(fmt.Sprintf("level: index %d out of range (have %d levels)", l.level, len(levelMaps)))
	}
	if l.Player == nil {
		panic("level: Player is not set")
	}
	log.Println("Loading level", l.level)
	l.resetLevelState()
	l.MapGrid = parseMap(levelMaps[l.level])
	l.loadDialog()
	l.UpdateCamera(l.playerX, l.playerY)
	l.disableInput = false
}

func (l *Level) resetLevelState() {
	l.ViewGrid = make([][]rune, rows)
	for y := range l.ViewGrid {
		l.ViewGrid[y] = make([]rune, cols)
		for x := range l.ViewGrid[y] {
			l.ViewGrid[y][x] = levelEmpty
		}
	}
	l.obstacles = map[rune]struct{}{levelWall: {}}
	l.entities = make(map[gridPos]*Entity)
	l.currentEntity = nil
	l.currentDialogStateID = ""
	l.selectedDialog = 0
	l.showingIntro = true
	l.showingItem = false
	l.cameraX = 0
	l.cameraY = 0
}

func parseMap(data string) [][]rune {
	data = strings.ReplaceAll(data, "\r\n", "\n")
	lines := strings.Split(strings.Trim(data, "\n"), "\n")

	runeLines := make([][]rune, len(lines))
	width := 0
	for i, line := range lines {
		runeLines[i] = []rune(line)
		width = max(width, len(runeLines[i]))
	}

	grid := make([][]rune, len(runeLines))
	for y, line := range runeLines {
		grid[y] = make([]rune, width)
		copy(grid[y], line)
		for x := len(line); x < width; x++ {
			grid[y][x] = levelEmpty
		}
	}
	return grid
}
