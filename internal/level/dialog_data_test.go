package level

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDialogDataIsConsistent(t *testing.T) {
	terminals := map[string]bool{nextLevel: true, endGame: true, endInteraction: true, endDialog: true}

	for index, raw := range levelDialogs {
		var dialog DialogState
		if err := json.Unmarshal([]byte(raw), &dialog); err != nil {
			t.Fatalf("level %d: parse dialog: %v", index, err)
		}
		grid := parseMap(levelMaps[index])
		hasPlayer := false

		for _, entity := range dialog.Entities {
			if entity.ID == "Player" {
				hasPlayer = true
			}
			if entity.Y < 0 || entity.Y >= len(grid) || entity.X < 0 || entity.X >= len(grid[entity.Y]) {
				t.Errorf("level %d: entity %q at (%d,%d) is outside the map", index, entity.ID, entity.X, entity.Y)
			}
			checkEntityStates(t, index, entity, terminals)
		}

		if !hasPlayer {
			t.Errorf("level %d: no Player entity", index)
		}
	}
}

func checkEntityStates(t *testing.T, index int, entity Entity, terminals map[string]bool) {
	t.Helper()
	ids := make(map[string]bool, len(entity.DialogStates))
	for _, state := range entity.DialogStates {
		ids[state.ID] = true
	}
	if len(entity.DialogStates) > 0 && !ids[startState] {
		t.Errorf("level %d: entity %q has dialog states but no %q state", index, entity.ID, startState)
	}

	resolves := func(next string) bool { return terminals[next] || ids[next] }

	for _, state := range entity.DialogStates {
		for _, problem := range stateProblems(state, resolves) {
			t.Errorf("level %d: entity %q state %q %s", index, entity.ID, state.ID, problem)
		}
	}
}

func stateProblems(state EntityDialogState, resolves func(string) bool) []string {
	if len(state.Options) == 0 {
		if resolves(state.NextState) {
			return nil
		}
		return []string{fmt.Sprintf("has no options and nextState %q does not resolve", state.NextState)}
	}

	var problems []string
	if state.Text == "" {
		problems = append(problems, "has options but no text")
	}
	for _, opt := range state.Options {
		if !resolves(opt.NextState) {
			problems = append(problems, fmt.Sprintf("option %q targets %q which does not exist", opt.Text, opt.NextState))
		}
	}
	return problems
}

func TestStateProblems(t *testing.T) {
	resolves := func(next string) bool { return next == "known" }

	cases := []struct {
		name  string
		state EntityDialogState
		want  int
	}{
		{"optionless resolving", EntityDialogState{NextState: "known"}, 0},
		{"optionless dangling", EntityDialogState{NextState: "nope"}, 1},
		{"optionless empty", EntityDialogState{}, 1},
		{"options all resolve", EntityDialogState{Text: "t", Options: []EntityDialogOption{{Text: "a", NextState: "known"}}}, 0},
		{"options without text", EntityDialogState{Options: []EntityDialogOption{{Text: "a", NextState: "known"}}}, 1},
		{"option dangling", EntityDialogState{Text: "t", Options: []EntityDialogOption{{Text: "a", NextState: "nope"}, {Text: "b"}}}, 2},
	}
	for _, tc := range cases {
		if got := stateProblems(tc.state, resolves); len(got) != tc.want {
			t.Errorf("%s: got %d problems %q, want %d", tc.name, len(got), got, tc.want)
		}
	}
}
