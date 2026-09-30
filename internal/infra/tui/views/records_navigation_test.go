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
