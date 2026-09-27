package main

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func sendKey(m model, code rune, mod tea.KeyMod) model {
	updated, _ := m.Update(tea.KeyPressMsg{Code: code, Mod: mod})
	return updated.(model)
}

func setFields(m model, values ...string) model {
	for i, value := range values {
		m.dialog.fields[i].input.SetValue(value)
	}
	return m
}

func TestExerciseDialogCreateEditCancel(t *testing.T) {
	m := pressKeys(newModel(nil), "n")
	if m.dialog == nil || m.dialog.title != "New exercise" {
		t.Fatal("new exercise dialog did not open on empty data")
	}
	m = pressKeys(m, "isquat")
	m = sendKey(m, tea.KeyEscape, 0)
	if m.dialog == nil || m.dialog.inserting || m.dialog.fields[0].input.Focused() {
		t.Fatal("Escape did not return from insert mode to navigation")
	}
	m = pressKeys(m, "j")
	if m.dialog.active != 1 || m.dialog.fields[1].input.Focused() {
		t.Fatal("j did not select the next field in normal mode")
	}
	m = pressKeys(m, "i80")
	m = sendKey(m, tea.KeyEscape, 0)
	m = pressKeys(m, "k")
	if m.dialog.active != 0 {
		t.Fatal("k did not select the previous field")
	}
	m = sendKey(m, tea.KeyUp, 0)
	if m.dialog.active != 1 {
		t.Fatal("field focus did not wrap backwards")
	}
	m = sendKey(m, tea.KeyEnter, 0)
	if m.dialog != nil || len(m.exercises) != 1 || m.exercises[0].name != "squat" || m.exercises[0].bodyWeightKG != 80 {
		t.Fatal("new exercise was not saved")
	}
	m = pressKeys(m, "e")
	if m.dialog.fields[0].input.Value() != "squat" || m.dialog.fields[1].input.Value() != "80" {
		t.Fatal("edit did not populate the selected exercise")
	}
	m = setFields(m, "front squat", "82.5")
	m = sendKey(m, tea.KeyEscape, 0)
	if m.dialog != nil || m.exercises[0].name != "squat" || m.exercises[0].bodyWeightKG != 80 {
		t.Fatal("Escape saved a cancelled edit")
	}
	m = setFields(pressKeys(m, "e"), "front squat", "82.5")
	m = sendKey(m, tea.KeyEnter, 0)
	if m.exercises[0].name != "front squat" || m.table.SelectedRow()[0] != "front squat" || m.exercises[0].bodyWeightKG != 82.5 {
		t.Fatal("edit did not update exercise data and table")
	}
}

func TestSetDialogCreateEditAndValidation(t *testing.T) {
	m := pressKeys(sampleModel(), "gsn")
	if m.dialog == nil || m.dialog.fields[1].input.Value() != "90" {
		t.Fatal("new set did not prefill the selected set")
	}
	for _, invalid := range []struct {
		values []string
		field  int
	}{
		{[]string{"0", "90", "2:30"}, 0},
		{[]string{"1.5", "90", "2:30"}, 0},
		{[]string{"256", "90", "2:30"}, 0},
		{[]string{"8", "NaN", "2:30"}, 1},
		{[]string{"8", "Inf", "2:30"}, 1},
		{[]string{"8", "-1", "2:30"}, 1},
		{[]string{"8", "90", "2:60"}, 2},
		{[]string{"8", "90", "bad"}, 2},
	} {
		m = sendKey(setFields(m, invalid.values...), tea.KeyEnter, 0)
		if m.dialog == nil || m.dialog.err == "" || m.dialog.active != invalid.field || len(m.setLines) != 9 {
			t.Fatalf("invalid values were accepted or lost: %v", invalid.values)
		}
	}
	m = sendKey(setFields(m, "8", "92.5", "2:30"), tea.KeyEnter, 0)
	if m.dialog != nil || len(m.setLines) != 10 || !strings.Contains(m.timeline.View(), "92.5 KG") {
		t.Fatal("new set was not saved and selected")
	}
	session, index, ok := m.selectedSet()
	if !ok || m.exercises[0].sessions[session].date.Format(time.DateOnly) != time.Now().Format(time.DateOnly) {
		t.Fatal("new set did not use today's session")
	}
	date := m.exercises[0].sessions[session].date
	m = pressKeys(m, "e")
	if m.dialog.fields[2].input.Value() != "2:30" {
		t.Fatal("edit did not prefill rest")
	}
	m = sendKey(setFields(m, "10", "95", "3:05"), tea.KeyEnter, 0)
	set := m.exercises[0].sessions[session].sets[index]
	if set.reps != 10 || set.weightKG != 95 || set.rest != 185*time.Second || len(m.setLines) != 10 {
		t.Fatal("edit did not replace the selected set")
	}
	if !m.exercises[0].sessions[session].date.Equal(date) {
		t.Fatal("editing changed the set date")
	}
	if m.table.SelectedRow()[1] != "118.8%" {
		t.Fatal("exercise percentage was not refreshed")
	}
	m = sendKey(pressKeys(m, "n"), tea.KeyEnter, 0)
	if len(m.exercises[0].sessions[session].sets) != 2 || m.setCursor != 1 {
		t.Fatal("new set did not append to today's existing session")
	}
}

func TestExerciseDialogValidationAndInputIsolation(t *testing.T) {
	m := pressKeys(sampleModel(), "n")
	m = sendKey(m, tea.KeyEnter, 0)
	if m.dialog == nil || m.dialog.err == "" || len(m.exercises) != 3 {
		t.Fatal("empty name was saved")
	}
	m = sendKey(setFields(m, "test", "NaN"), tea.KeyEnter, 0)
	if m.dialog == nil || m.dialog.active != 1 {
		t.Fatal("invalid body weight was saved")
	}
	m = setFields(m, "", "")
	m.dialog.focus(0)
	m = pressKeys(m, "iqndegsjklh")
	if m.dialog.fields[0].input.Value() != "qndegsjklh" || len(m.exercises) != 3 || m.activePanel != exercisesPanel || m.table.Cursor() != 0 {
		t.Fatal("typing in a field triggered a global action")
	}
	_, cmd := m.Update(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	if cmd == nil {
		t.Fatal("Ctrl+C no longer quits from a dialog")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("Ctrl+C did not issue Quit")
	}
	m = sendKey(m, tea.KeyEnter, 0)
	if m.dialog != nil || len(m.exercises) != 4 || m.exercises[3].bodyWeightKG != 0 {
		t.Fatal("optional body weight could not be left blank")
	}
}

func TestDeleteRecordsAtPageBoundaries(t *testing.T) {
	updated, _ := sampleModel().Update(tea.WindowSizeMsg{Width: 160, Height: 12})
	m := pressKeys(updated.(model), "gsG")
	for remaining := 8; remaining >= 0; remaining-- {
		m = pressKeys(m, "d")
		if len(m.setLines) != remaining || m.setCursor != max(0, remaining-1) || m.setPager.Page >= m.setPager.TotalPages {
			t.Fatalf("invalid selection after deleting down to %d sets", remaining)
		}
		_ = m.View()
	}
	if len(m.exercises[0].sessions) != 0 {
		t.Fatal("empty date groups were not removed")
	}
	m = pressKeys(m, "ed")
	if m.dialog != nil {
		t.Fatal("editing an empty set list opened a dialog")
	}
	m = pressKeys(m, "n")
	m = sendKey(m, tea.KeyEnter, 0)
	if len(m.setLines) != 1 {
		t.Fatal("could not add a set after deleting all sets")
	}
	m = pressKeys(m, "geG")
	for remaining := 2; remaining >= 0; remaining-- {
		m = pressKeys(m, "d")
		if len(m.exercises) != remaining || m.exercisePager.Page >= m.exercisePager.TotalPages {
			t.Fatal("deleting exercise left invalid pagination")
		}
		_ = m.View()
	}
	m = pressKeys(m, "degsen")
	if m.dialog != nil || len(m.exercises) != 0 {
		t.Fatal("empty-list actions mutated data or opened an invalid dialog")
	}
}

func TestCompositedDialogFitsTerminal(t *testing.T) {
	m := pressKeys(sampleModel(), "gsn")
	for _, size := range []tea.WindowSizeMsg{{Width: 100, Height: 24}, {Width: 60, Height: 16}, {Width: 40, Height: 12}} {
		updated, _ := m.Update(size)
		m = updated.(model)
		m.dialog.fields[0].input.SetValue(strings.Repeat("1", 80))
		view := m.View()
		if lipgloss.Width(view.Content) > size.Width || lipgloss.Height(view.Content) > size.Height {
			t.Fatalf("dialog overflows terminal at %dx%d", size.Width, size.Height)
		}
		if !strings.Contains(view.Content, "New set") || !strings.Contains(view.Content, "enter") || !strings.Contains(view.Content, "esc") {
			t.Fatalf("dialog or its controls are missing at %dx%d", size.Width, size.Height)
		}
		if !view.AltScreen || view.MouseMode != tea.MouseModeCellMotion {
			t.Fatal("dialog disabled alternate screen or mouse capture")
		}
	}
}

func TestDialogVimModes(t *testing.T) {
	m := pressKeys(sampleModel(), "e")
	if m.dialog.inserting || m.dialog.fields[0].input.Focused() || m.View().Cursor == nil || m.View().Cursor.Shape != tea.CursorBlock {
		t.Fatal("dialog must open in normal mode with a block cursor")
	}
	name := m.dialog.fields[0].input.Value()
	m = pressKeys(m, "qndes")
	updated, _ := m.Update(tea.PasteMsg{Content: "ignored"})
	m = updated.(model)
	if m.dialog.fields[0].input.Value() != name || len(m.exercises) != 3 {
		t.Fatal("normal mode allowed typing, pasting, or background actions")
	}
	m = sendKey(m, tea.KeyDown, 0)
	if m.dialog.active != 1 {
		t.Fatal("Down did not select the next field")
	}
	m = sendKey(m, tea.KeyUp, 0)
	m = pressKeys(m, "i")
	if !m.dialog.inserting || !m.dialog.fields[0].input.Focused() || m.View().Cursor.Shape != tea.CursorBar {
		t.Fatal("i did not enter insert mode")
	}
	m = pressKeys(m, "jki")
	m = sendKey(m, tea.KeyDown, 0)
	m = sendKey(m, tea.KeyUp, 0)
	if m.dialog.active != 0 || m.dialog.fields[0].input.Value() != name+"jki" {
		t.Fatal("insert mode treated text or arrows as field navigation")
	}
	m = sendKey(m, tea.KeyEscape, 0)
	if m.dialog == nil || m.dialog.inserting || m.dialog.fields[0].input.Focused() || m.dialog.fields[0].input.Value() != name+"jki" {
		t.Fatal("first Escape did not retain the draft in normal mode")
	}
	m = sendKey(m, tea.KeyEscape, 0)
	if m.dialog != nil || m.exercises[0].name != name {
		t.Fatal("second Escape did not cancel the dialog")
	}
	for _, mode := range []string{"", "i"} {
		m = pressKeys(m, "e"+mode)
		m = setFields(m, "renamed", "80")
		m = sendKey(m, tea.KeyEnter, 0)
		if m.dialog != nil || m.exercises[0].name != "renamed" {
			t.Fatalf("Enter did not save in mode %q", mode)
		}
	}
}

func TestDialogCursorAndHelp(t *testing.T) {
	m := pressKeys(sampleModel(), "n")
	for _, size := range []tea.WindowSizeMsg{{Width: 100, Height: 24}, {Width: 60, Height: 24}, {Width: 40, Height: 12}} {
		updated, _ := m.Update(size)
		m = updated.(model)
		m.dialog.fields[0].input.SetValue("START" + strings.Repeat("界", 35) + "END")
		m.dialog.fields[0].input.CursorEnd()
		layers, cursor := m.dialogLayers()
		if cursor == nil || cursor.Shape != tea.CursorBlock || cursor.Blink {
			t.Fatal("normal mode needs a steady block cursor")
		}
		form := layers[0]
		if cursor.Position.X <= form.GetX() || cursor.Position.X >= form.GetX()+form.Width()-1 || cursor.Position.Y <= form.GetY() || cursor.Position.Y >= form.GetY()+form.Height()-1 {
			t.Fatal("cursor is outside the form content")
		}
		title := ansi.Strip(strings.Split(m.dialogView(), "\n")[0])
		if strings.Contains(title, "NORMAL") || strings.Contains(title, "INSERT") {
			t.Fatal("title still contains the mode")
		}
		if len(layers) != 1 || strings.Contains(m.View().Content, "Modal navigation") {
			t.Fatal("navigation guide was not removed")
		}
		lines := strings.Split(ansi.Strip(m.dialogView()), "\n")
		footer := lines[len(lines)-1]
		parts := strings.Split(footer, " Save <enter> ")
		if len(parts) != 2 || strings.Trim(parts[0], "╰─") != "" || strings.Trim(parts[1], "─╯") != "" {
			t.Fatal("bottom border should contain only Save <enter>")
		}
		if delta := lipgloss.Width(parts[1]) - lipgloss.Width(parts[0]); delta < 0 || delta > 1 {
			t.Fatal("Save <enter> is not centered")
		}
		for _, hint := range []string{"move up", "move down", "normal mode", "insert mode"} {
			if !strings.Contains(m.dialogView(), hint) {
				t.Fatalf("modal help is missing %q", hint)
			}
		}
		m = pressKeys(m, "i")
		if m.View().Cursor.Shape != tea.CursorBar || !m.View().Cursor.Blink {
			t.Fatal("insert mode needs a blinking caret")
		}
		lines = strings.Split(ansi.Strip(m.dialogView()), "\n")
		if lines[len(lines)-1] != footer {
			t.Fatal("footer changed when entering insert mode")
		}
		m = sendKey(m, tea.KeyHome, 0)
		homeX := m.View().Cursor.Position.X
		m = sendKey(m, tea.KeyRight, 0)
		if m.View().Cursor.Position.X != homeX+1 {
			t.Fatal("caret did not follow character movement")
		}
		m = sendKey(m, tea.KeyEscape, 0)
		if m.View().Cursor.Shape != tea.CursorBlock {
			t.Fatal("Escape did not restore the block cursor")
		}
		m = pressKeys(m, "?")
		view := m.View()
		if view.Cursor == nil || strings.Contains(view.Content, "Modal navigation") {
			t.Fatal("removed help shortcut still opens a guide")
		}
		if lipgloss.Width(view.Content) > size.Width || lipgloss.Height(view.Content) > size.Height {
			t.Fatal("modal exceeded terminal bounds")
		}
	}
}

func TestDialogInputDoesNotWrap(t *testing.T) {
	for _, value := range []string{
		"START" + strings.Repeat("x", 75) + "END",
		"START" + strings.Repeat("界", 35) + "END",
	} {
		m := pressKeys(sampleModel(), "ni")
		for _, size := range []tea.WindowSizeMsg{{Width: 100, Height: 24}, {Width: 48, Height: 20}, {Width: 40, Height: 12}, {Width: 100, Height: 24}} {
			updated, _ := m.Update(size)
			m = updated.(model)
			m.dialog.fields[0].input.SetValue("short")
			height := lipgloss.Height(m.dialogView())
			m.dialog.fields[0].input.SetValue(value)
			m = sendKey(m, tea.KeyEnd, 0)
			if !strings.Contains(ansi.Strip(m.dialogView()), "END") {
				t.Fatal("horizontal viewport did not reveal the end of the input")
			}
			for range 20 {
				m = sendKey(m, tea.KeyLeft, 0)
				if lipgloss.Height(m.dialogView()) != height {
					t.Fatalf("long field wrapped while moving the cursor at width %d: height %d, want %d\n%s", size.Width, lipgloss.Height(m.dialogView()), height, ansi.Strip(m.dialogView()))
				}
			}
			m = sendKey(m, tea.KeyHome, 0)
			if !strings.Contains(ansi.Strip(m.dialogView()), "START") {
				t.Fatal("horizontal viewport did not reveal the start of the input")
			}
			if m.dialog.fields[0].input.Value() != value || lipgloss.Width(m.dialogView()) > size.Width {
				t.Fatal("horizontal scrolling changed the value or overflowed the terminal")
			}
		}
	}
}
