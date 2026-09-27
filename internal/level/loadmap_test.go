package level

import (
	"reflect"
	"testing"
)

func TestParseMap(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want [][]rune
	}{
		{
			"ragged lines are padded",
			"##\n#\n",
			[][]rune{{'#', '#'}, {'#', ' '}},
		},
		{
			"windows line endings are stripped",
			"##\r\n##\r\n",
			[][]rune{{'#', '#'}, {'#', '#'}},
		},
		{
			"width counts runes not bytes",
			"#é#\n###",
			[][]rune{{'#', 'é', '#'}, {'#', '#', '#'}},
		},
		{
			"surrounding blank lines are trimmed",
			"\n\n#\n\n",
			[][]rune{{'#'}},
		},
	}
	for _, tc := range cases {
		if got := parseMap(tc.in); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%s: got %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestShippedMapsAreRectangular(t *testing.T) {
	for i, raw := range levelMaps {
		grid := parseMap(raw)
		for y, row := range grid {
			for x, ch := range row {
				if ch == 0 {
					t.Errorf("level %d: zero rune at (%d,%d)", i, x, y)
				}
			}
		}
	}
}
