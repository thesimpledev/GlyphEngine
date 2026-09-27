package level

var wallText = []string{
	"Its a wall, you probably",
	"shouldn't lick it",
	"It doesn't look very clean",
}

func (l *Level) IsWalkable(x, y int) bool {
	if l.disableInput || l.showingIntro || !l.inBounds(x, y) {
		return false
	}

	l.dismissItem()

	if _, blocked := l.obstacles[l.MapGrid[y][x]]; !blocked {
		return true
	}

	l.showObstacle(x, y)
	return false
}

func (l *Level) dismissItem() {
	if !l.showingItem {
		return
	}
	l.showingItem = false
	if l.currentEntity != nil {
		l.currentEntity.active = false
	}
	l.currentEntity = nil
	l.clearDialogArea()
}

func (l *Level) showObstacle(x, y int) {
	l.showingItem = true
	l.clearDialogArea()

	entity, ok := l.entities[gridPos{x, y}]
	if !ok {
		l.currentEntity = nil
		l.setDialog(wallText, nil)
		return
	}

	l.currentEntity = entity
	var prompt []string
	if len(entity.DialogStates) > 0 {
		prompt = []string{"Press SPACE to interact"}
	}
	l.setDialog([]string{entity.Text}, prompt)
}

func (l *Level) UpdateBoard(fromX, fromY, toX, toY int) {
	if !l.inBounds(fromX, fromY) || !l.inBounds(toX, toY) {
		return
	}
	l.MapGrid[fromY][fromX] = levelEmpty
	l.MapGrid[toY][toX] = l.playerSymbol
}

func (l *Level) inBounds(x, y int) bool {
	return y >= 0 && y < len(l.MapGrid) && x >= 0 && x < len(l.MapGrid[y])
}
