package main

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestWheelScrollsOnlyPanelItems(t *testing.T) {
	m := sampleModel()
	view := m.View()
	if !view.AltScreen || view.MouseMode != tea.MouseModeCellMotion {
		t.Fatal("alternate screen and mouse capture must be enabled")
	}
	for _, p := range []panel{setsPanel, exercisesPanel} {
		bounds := m.itemBounds(p)
		updated, _ := m.Update(tea.MouseWheelMsg{X: bounds.Min.X, Y: bounds.Min.Y, Button: tea.MouseWheelDown})
		m = updated.(model)
		if m.activePanel != p {
			t.Fatal("wheel did not target the hovered panel")
		}
		if p == setsPanel && (m.setCursor != 1 || m.table.Cursor() != 0) {
			t.Fatal("set wheel scroll changed the exercise")
		}
		if p == exercisesPanel && (m.table.Cursor() != 1 || m.setCursor != 0) {
			t.Fatal("exercise wheel scroll did not reset the set selection")
		}
	}
	for _, point := range [][2]int{{0, 0}, {m.width / 2, 5}, {5, m.height - 1}, {2, 2}} {
		updated, _ := m.Update(tea.MouseWheelMsg{X: point[0], Y: point[1], Button: tea.MouseWheelDown})
		next := updated.(model)
		if next.table.Cursor() != m.table.Cursor() || next.setCursor != m.setCursor || next.activePanel != m.activePanel {
			t.Fatalf("scroll outside items changed navigation at %v", point)
		}
	}
	m = pressKeys(m, "n")
	bounds := m.itemBounds(setsPanel)
	updated, _ := m.Update(tea.MouseWheelMsg{X: bounds.Min.X, Y: bounds.Min.Y, Button: tea.MouseWheelDown})
	m = updated.(model)
	if m.activePanel != exercisesPanel || m.table.Cursor() != 1 || m.setCursor != 0 || m.dialog.active != 0 {
		t.Fatal("wheel scrolled the background while a dialog was open")
	}
}
