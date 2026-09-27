package level

func (l *Level) UpdateCamera(playerX, playerY int) {
	mapHeight := len(l.MapGrid)
	if mapHeight == 0 || len(l.ViewGrid) < rowDivider {
		return
	}
	mapWidth := len(l.MapGrid[0])

	l.cameraX = clampCamera(playerX-colDivider/2, mapWidth, colDivider)
	l.cameraY = clampCamera(playerY-rowDivider/2, mapHeight, rowDivider)
	l.updateGridFromCamera()
}

func clampCamera(desired, mapSize, viewSize int) int {
	if mapSize <= viewSize {
		return 0
	}
	return min(max(desired, 0), mapSize-viewSize)
}

func (l *Level) updateGridFromCamera() {
	for screenY := range rowDivider {
		for screenX := range colDivider {
			mapX := l.cameraX + screenX
			mapY := l.cameraY + screenY

			l.ViewGrid[screenY][screenX] = levelEmpty
			if l.inBounds(mapX, mapY) {
				l.ViewGrid[screenY][screenX] = l.MapGrid[mapY][mapX]
			}
		}
	}
}
