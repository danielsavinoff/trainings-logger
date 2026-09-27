package main

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func sampleModel() model {
	return newModel(mockExercises()[:3])
}

func pressKeys(m model, keys string) model {
	for _, k := range keys {
		updated, _ := m.Update(tea.KeyPressMsg{Code: k, Text: string(k)})
		m = updated.(model)
	}
	return m
}

func TestSplitResize(t *testing.T) {
	for _, width := range []int{80, 101, 160} {
		updated, _ := sampleModel().Update(tea.WindowSizeMsg{Width: width, Height: 24})
		m := updated.(model)
		left := lipgloss.Width(panelSpacing.Render(baseStyle.Render(m.exercisesView())))
		right := lipgloss.Width(panelSpacing.Render(baseStyle.Render(m.setsView())))
		if left != width/2 || right != width-width/2 {
			t.Fatalf("width %d: panes are %d and %d", width, left, right)
		}
		if got := lipgloss.Height(m.View().Content); got != 24 {
			t.Fatalf("height = %d, want 24", got)
		}
		rightHeight := lipgloss.Height(baseStyle.Render(m.setsView()))
		if rightHeight != lipgloss.Height(baseStyle.Render(m.exercisesView())) {
			t.Fatal("Sets and Exercises panel heights differ")
		}
		columns := m.table.Columns()
		total := 0
		for _, column := range columns {
			total += column.Width + 2
			if delta := columns[0].Width - column.Width; delta < 0 || delta > 1 {
				t.Fatal("columns differ by more than one cell")
			}
		}
		if total != m.table.Width() {
			t.Fatal("columns do not fill table")
		}
	}
}

func TestTimelineFollowsSelectionAndFocus(t *testing.T) {
	m := sampleModel()
	if row := m.table.SelectedRow(); row[0] != "back squat" {
		t.Fatal(row)
	}
	if view := m.timeline.View(); !strings.Contains(view, "90 KG") || !strings.Contains(view, "23 September, 2026") {
		t.Fatal(view)
	}
	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = updated.(model)
	if !strings.Contains(m.timeline.View(), "65 KG") {
		t.Fatal(m.timeline.View())
	}
	m = pressKeys(m, "gs")
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = updated.(model)
	if m.activePanel != setsPanel || m.table.Cursor() != 1 || m.setCursor != 1 {
		t.Fatal("timeline scrolling changed main selection")
	}
}

func TestPanelNavigationAndCounters(t *testing.T) {
	m := sampleModel()
	for _, step := range []struct {
		keys string
		want panel
	}{{"gs", setsPanel}, {"ge", exercisesPanel}, {"gs", setsPanel}} {
		m = pressKeys(m, step.keys)
		if m.activePanel != step.want {
			t.Fatalf("active panel = %d, want %d", m.activePanel, step.want)
		}
	}
	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = updated.(model)
	if m.setCursor != 1 || m.table.Cursor() != 0 {
		t.Fatal("Sets navigation affected Exercises")
	}
	view := m.View().Content
	for _, counter := range []string{"1 of 3", "2 of 9"} {
		if !strings.Contains(view, counter) {
			t.Fatalf("missing counter %q", counter)
		}
	}
	pane := renderPane(baseStyle, "Sets", m.setsView(), "2 of 9")
	lines := strings.Split(pane, "\n")
	if !strings.Contains(lines[len(lines)-1], " 2 of 9 ┘") {
		t.Fatal("counter is not in bottom-right border")
	}
}

func TestRemovedPanelShortcuts(t *testing.T) {
	for _, active := range []panel{exercisesPanel, setsPanel} {
		m := sampleModel()
		m.focusPanel(active)
		for _, msg := range []tea.KeyPressMsg{
			{Code: tea.KeyEscape},
			{Code: tea.KeyTab},
			{Code: tea.KeyTab, Mod: tea.ModShift},
		} {
			updated, _ := m.Update(msg)
			m = updated.(model)
			if m.activePanel != active || m.table.Cursor() != 0 || m.setCursor != 0 {
				t.Fatalf("removed shortcut %s changed navigation", msg.String())
			}
		}
		for _, label := range []string{"esc", "tab"} {
			if strings.Contains(m.helpView(), label) {
				t.Fatalf("removed shortcut %s remains in help", label)
			}
		}
		for _, title := range []string{"Exercises <ge>", "Sets <gs>"} {
			if !strings.Contains(m.View().Content, title) {
				t.Fatalf("missing panel title %s", title)
			}
		}
	}
}

func TestPanelShortcutPrefix(t *testing.T) {
	m := sampleModel()
	updated, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = pressKeys(updated.(model), "g")
	if m.table.Cursor() != 1 || m.activePanel != exercisesPanel {
		t.Fatal("prefix moved the selection")
	}
	m = pressKeys(m, "s")
	if m.table.Cursor() != 1 || m.activePanel != setsPanel {
		t.Fatal("gs did not preserve the selected exercise")
	}
	m = pressKeys(m, "gx")
	m = pressKeys(m, "e")
	if m.activePanel != setsPanel {
		t.Fatal("invalid sequence did not clear the prefix")
	}
}

func TestSetsSelectionScrollsIntoView(t *testing.T) {
	updated, _ := sampleModel().Update(tea.WindowSizeMsg{Width: 100, Height: 12})
	m := updated.(model)
	m.focusPanel(setsPanel)
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnd})
	m = updated.(model)
	if m.setCursor != 8 || m.setLines[8] >= m.timeline.YOffset()+m.timeline.Height() {
		t.Fatal("last set is not selected and visible")
	}
	if !strings.Contains(m.View().Content, "9 of 9") {
		t.Fatal("Sets counter did not update")
	}
}

func TestDateLine(t *testing.T) {
	line := dateLine("23 December, 2026", 30)
	if line != "───── 23 December, 2026 ──────" {
		t.Fatal(line)
	}
	for _, width := range []int{0, 5, 10, 11, 30} {
		if lipgloss.Width(dateLine("23 December, 2026", width)) != width {
			t.Fatalf("wrong width %d", width)
		}
	}
}

func TestWorkoutDataAndFixedHeadings(t *testing.T) {
	updated, _ := sampleModel().Update(tea.WindowSizeMsg{Width: 160, Height: 12})
	m := updated.(model)
	row := m.table.SelectedRow()
	if len(row) != 2 || row[1] != "112.5%" {
		t.Fatalf("exercise summary does not match latest session: %v", row)
	}
	for i, title := range []string{"Name", "Body Weight %"} {
		if m.table.Columns()[i].Title != title {
			t.Fatalf("missing heading %s", title)
		}
	}
	m.focusPanel(setsPanel)
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnd})
	m = updated.(model)
	header := strings.Split(m.setsView(), "\n")[0]
	for _, title := range []string{"Reps", "Exercise Weight", "Rest (Time)"} {
		if !strings.Contains(header, title) || strings.Count(m.setsView(), title) != 1 {
			t.Fatalf("heading %s did not stay fixed or was repeated", title)
		}
	}
	m.focusPanel(exercisesPanel)
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = updated.(model)
	if m.setCursor != 0 || m.timeline.YOffset() != 0 || !strings.Contains(m.timeline.View(), "65 KG") {
		t.Fatal("switching exercises did not reset to its latest sets")
	}
}

func TestExercisePagination(t *testing.T) {
	updated, _ := newModel(mockExercises()).Update(tea.WindowSizeMsg{Width: 160, Height: 13})
	m := updated.(model)
	if m.exercisePager.TotalPages != 3 {
		t.Fatalf("expected three exercise pages, got %d", m.exercisePager.TotalPages)
	}
	m = pressKeys(m, "l")
	if m.exercisePager.Page != 1 || m.table.Cursor() != 3 || m.setPager.Page != 0 {
		t.Fatal("next exercise page did not select its first exercise and reset sets")
	}
	view := m.exercisesView()
	if !strings.Contains(view, "overhead press") || strings.Contains(view, "back squat") || strings.Contains(view, "leg press") {
		t.Fatal("exercise page shows the wrong rows")
	}
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyPgDown})
	m = updated.(model)
	if m.exercisePager.Page != 2 || !strings.Contains(m.exercisesView(), "leg press") {
		t.Fatal("Page Down did not reach the final partial page")
	}
	m = pressKeys(m, "l")
	if m.exercisePager.Page != 2 || m.table.Cursor() != 6 {
		t.Fatal("next page moved past the last page")
	}
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m = updated.(model)
	if m.exercisePager.Page != 1 || !strings.Contains(m.exercisesView(), "romanian deadlift") {
		t.Fatal("row navigation did not update the page")
	}
}

func TestSetPaginationKeepsDatesAndAllSets(t *testing.T) {
	updated, _ := sampleModel().Update(tea.WindowSizeMsg{Width: 160, Height: 12})
	m := pressKeys(updated.(model), "gs")
	seen := make(map[int]bool)
	for page := 0; page < m.setPager.TotalPages; page++ {
		if m.setPager.Page != page || m.table.Cursor() != 0 || m.exercisePager.Page != 0 {
			t.Fatal("set pagination changed the exercise selection or missed a page")
		}
		if firstLine := strings.Split(m.timeline.View(), "\n")[0]; !strings.Contains(firstLine, "September, 2026") {
			t.Fatalf("page %d is missing its date: %s", page, firstLine)
		}
		start := m.timeline.YOffset()
		end := start + m.timeline.Height()
		if line := m.setLines[m.setCursor]; line < start || line >= end {
			t.Fatal("selected set is not on the current page")
		}
		for i, line := range m.setLines {
			if line >= start && line < end {
				if seen[i] {
					t.Fatalf("set %d appears on multiple pages", i)
				}
				seen[i] = true
			}
		}
		updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
		m = updated.(model)
	}
	if len(seen) != 9 || !m.setPager.OnLastPage() {
		t.Fatal("pagination lost sets or moved past the final page")
	}
	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyPgUp})
	m = updated.(model)
	if m.setPager.Page != m.setPager.TotalPages-2 {
		t.Fatal("Page Up did not move to the previous set page")
	}
	selected := m.setCursor
	updated, _ = m.Update(tea.WindowSizeMsg{Width: 160, Height: 24})
	m = updated.(model)
	if m.setCursor != selected || m.setPager.Page != 0 || m.setPager.TotalPages != 1 {
		t.Fatal("resize did not preserve the selected set and recalculate pages")
	}
	m = pressKeys(m, "gej")
	if m.setPager.Page != 0 || m.setCursor != 0 {
		t.Fatal("changing exercises did not reset set pagination")
	}
}

func TestEmptyPagination(t *testing.T) {
	for _, exercises := range [][]mockExercise{nil, {{name: "empty exercise"}}} {
		m := newModel(exercises)
		m = pressKeys(m, "lhgslh")
		if m.exercisePager.Page != 0 || m.setPager.Page != 0 || m.setPager.TotalPages != 1 {
			t.Fatal("empty pagination moved beyond its only page")
		}
		if !strings.Contains(m.View().Content, "0 of 0") {
			t.Fatal("empty state lost its counter")
		}
	}
}
