package views

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestRecordsPanelNavigation(t *testing.T) {
	m := press(sampleModel(), "g", "s", "j")
	if m.activePanel != setsPanel || m.setCursor != 1 || m.table.Cursor() != 0 {
		t.Fatal("set navigation changed the exercise selection")
	}

	m = press(m, "l")
	if m.setPager.Page != 1 || m.setCursor != m.setPageStarts[1] {
		t.Fatal("next page did not select its first set")
	}

	m = press(m, "g", "e", "j")
	if m.table.Cursor() != 1 || len(m.setLines) != 0 {
		t.Fatal("exercise navigation did not load the selected exercise's sets")
	}
}

func TestRecordsVimRowJumpShortcuts(t *testing.T) {
	exercises := sampleModel()
	exercises.table.GotoBottom()
	exercises.syncExercisePager()
	exercises.refreshSelectedSets()

	exercises = press(
		exercises,
		"g",
		"g",
	)
	if exercises.table.Cursor() != 0 || exercises.exercisePager.Page != 0 {
		t.Fatal("gg did not select the first exercise")
	}

	exercises = press(
		exercises,
		"G",
	)
	if exercises.table.Cursor() != len(exercises.exercises)-1 {
		t.Fatal("G did not select the last exercise")
	}

	exercises = press(
		exercises,
		"home",
	)
	if exercises.table.Cursor() != 0 {
		t.Fatal("Home no longer selects the first exercise")
	}

	exercises = press(
		exercises,
		"end",
	)
	if exercises.table.Cursor() != len(exercises.exercises)-1 {
		t.Fatal("End no longer selects the last exercise")
	}

	sets := sampleModel()
	sets.focusPanel(setsPanel)
	sets = press(
		sets,
		"G",
	)
	if sets.setCursor != len(sets.setLines)-1 {
		t.Fatal("G did not select the last set")
	}

	sets = press(
		sets,
		"g",
		"g",
	)
	if sets.setCursor != 0 || sets.setPager.Page != 0 {
		t.Fatal("gg did not select the first set")
	}

	sets = press(
		sets,
		"end",
	)
	if sets.setCursor != len(sets.setLines)-1 {
		t.Fatal("End no longer selects the last set")
	}

	sets = press(
		sets,
		"home",
	)
	if sets.setCursor != 0 {
		t.Fatal("Home no longer selects the first set")
	}

	sets = press(
		sets,
		"g",
		"e",
	)
	if sets.activePanel != exercisesPanel {
		t.Fatal("ge no longer focuses exercises")
	}

	sets = press(
		sets,
		"g",
		"s",
	)
	if sets.activePanel != setsPanel {
		t.Fatal("gs no longer focuses sets")
	}
}

func TestRecordsMouseTargetsPanelAndRespectsDialog(t *testing.T) {
	m := sampleModel()
	bounds := m.itemBounds(setsPanel)
	wheel := tea.MouseWheelMsg{X: bounds.Min.X, Y: bounds.Min.Y, Button: tea.MouseWheelDown}
	updated, _ := m.Update(wheel)
	m = updated.(RecordsModel)
	if m.activePanel != setsPanel || m.setCursor != 1 || m.table.Cursor() != 0 {
		t.Fatal("wheel did not navigate only the panel under the pointer")
	}

	m = press(m, "e")
	updated, _ = m.Update(wheel)
	m = updated.(RecordsModel)
	if m.setCursor != 1 || m.dialog == nil {
		t.Fatal("wheel changed the background while the dialog was open")
	}
}
