package level

import (
	"strings"
	"testing"
)

func TestClampCamera(t *testing.T) {
	cases := []struct {
		name                       string
		desired, mapSize, viewSize int
		want                       int
	}{
		{"map smaller than view", 7, 25, 30, 0},
		{"map equal to view", 7, 30, 30, 0},
		{"desired below zero", -5, 54, 40, 0},
		{"desired past the end", 30, 54, 40, 14},
		{"desired in range", 10, 54, 40, 10},
	}
	for _, tc := range cases {
		if got := clampCamera(tc.desired, tc.mapSize, tc.viewSize); got != tc.want {
			t.Errorf("%s: got %d, want %d", tc.name, got, tc.want)
		}
	}
}

func testMap(width, height int) string {
	row := strings.Repeat("#", width)
	rows := make([]string, height)
	for i := range rows {
		rows[i] = row
	}
	return strings.Join(rows, "\n")
}

func TestUpdateCameraFollowsPlayerOnTallMap(t *testing.T) {
	l := &Level{}
	l.resetLevelState()
	l.MapGrid = parseMap(testMap(25, 54))

	cases := []struct {
		playerY int
		wantY   int
	}{
		{1, 0},
		{30, 10},
		{53, 14},
	}
	for _, tc := range cases {
		l.UpdateCamera(5, tc.playerY)
		if l.cameraX != 0 || l.cameraY != tc.wantY {
			t.Errorf("player y=%d: camera (%d,%d), want (0,%d)", tc.playerY, l.cameraX, l.cameraY, tc.wantY)
		}
	}
}

func TestUpdateCameraCopiesMapIntoView(t *testing.T) {
	l := &Level{}
	l.resetLevelState()
	l.MapGrid = parseMap("#@#\n###")
	l.UpdateCamera(1, 0)

	if got := string(l.ViewGrid[0][:3]); got != "#@#" {
		t.Errorf("view row 0 = %q, want %q", got, "#@#")
	}
	if l.ViewGrid[0][3] != levelEmpty || l.ViewGrid[2][0] != levelEmpty {
		t.Error("cells outside the map should be empty")
	}
}

func TestUpdateCameraBeforeLoadDoesNotPanic(t *testing.T) {
	l := &Level{}
	l.UpdateCamera(0, 0)
}
