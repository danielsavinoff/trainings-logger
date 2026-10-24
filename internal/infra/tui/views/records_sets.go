package views

import (
	"context"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"

	"github.com/danielsavinoff/trainings-logger/internal/app/query"
	"github.com/danielsavinoff/trainings-logger/internal/infra/tui/components"
)

func (m RecordsModel) selectedSet() (session, set int, ok bool) {
	exercise := m.selectedExercise()
	if exercise < 0 {
		return 0, 0, false
	}

	index := m.setCursor
	for s, session := range m.exercises[exercise].sessions {
		if index >= 0 && index < len(session.sets) {
			return s, index, true
		}

		index -= len(session.sets)
	}

	return 0, 0, false
}

func (m *RecordsModel) appendSet(exercise int, set workoutSet, now time.Time) {
	e := &m.exercises[exercise]
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	index, cursor := 0, 0
	for index < len(e.sessions) && e.sessions[index].date.Format(time.DateOnly) > date.Format(time.DateOnly) {
		cursor += len(e.sessions[index].sets)
		index++
	}

	if index < len(e.sessions) && e.sessions[index].date.Format(time.DateOnly) == date.Format(time.DateOnly) {
		cursor += len(e.sessions[index].sets)
		e.sessions[index].sets = append(e.sessions[index].sets, set)
	} else {
		e.sessions = append(e.sessions, workoutSession{})
		copy(e.sessions[index+1:], e.sessions[index:])
		e.sessions[index] = workoutSession{date: date, sets: []workoutSet{set}}
	}

	m.setCursor = cursor
}

func (m *RecordsModel) loadSelectedSets(ctx context.Context) error {
	m.setsExerciseID = ""
	m.setCursor = 0

	index := m.selectedExercise()
	if index < 0 || m.store == nil {
		return nil
	}

	exercise := &m.exercises[index]
	exercise.sessions = nil

	result, err := m.store.ListSets.Handle(ctx, query.ListSets{
		ExerciseID: exercise.id,
	})
	if err != nil {
		return err
	}

	records := make([]workoutSet, 0, len(result.Result))
	for _, s := range result.Result {
		record, err := setRecord(s, m.weightUnit)
		if err != nil {
			return err
		}

		records = append(records, record)
	}

	for i := len(records) - 1; i >= 0; i-- {
		s := records[i]
		m.appendSet(index, s, s.createdAt.Local())
	}

	m.setsExerciseID = exercise.id
	m.setCursor = 0

	return nil
}

func dateLine(date string, width int) string {
	date = " " + date + " "
	if width <= len(date) {
		return date[:max(0, width)]
	}

	remaining := width - len(date)

	return strings.Repeat("─", remaining/2) + date + strings.Repeat("─", remaining-remaining/2)
}

func (m RecordsModel) setColumns() []table.Column {
	return components.EvenColumns(
		[]table.Column{
			{Title: "Reps"},
			{Title: "Exercise Weight"},
			{Title: "Rest (seconds)"},
		},
		m.timeline.Width(),
	)
}

func (m RecordsModel) renderSetRow(values []string, style lipgloss.Style) string {
	columns := m.setColumns()
	cells := make([]string, 0, len(columns))
	for i, column := range columns {
		if column.Width <= 0 {
			continue
		}

		cell := lipgloss.NewStyle().Width(column.Width).MaxWidth(column.Width).Inline(true).Render(values[i])
		cells = append(cells, lipgloss.NewStyle().Padding(0, 1).Render(cell))
	}

	return style.Render(lipgloss.JoinHorizontal(lipgloss.Top, cells...))
}

func (m RecordsModel) setsView() string {
	columns := m.setColumns()
	titles := make([]string, len(columns))
	for i, column := range columns {
		titles[i] = column.Title
	}

	header := m.renderSetRow(titles, lipgloss.NewStyle().Foreground(components.ActivePanelColor))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		m.timeline.View(),
	)
}

func (m *RecordsModel) refreshTimeline() {
	m.setLines = nil
	m.setPageStarts = nil
	m.setPager.Page = 0
	m.setPager.TotalPages = 1
	m.setPager.PerPage = max(1, m.timeline.Height())

	if len(m.exercises) == 0 {
		m.setCursor = 0
		m.timeline.SetContent("")

		return
	}

	exercise := m.exercises[m.table.Cursor()]
	count := 0
	for _, session := range exercise.sessions {
		count += len(session.sets)
	}

	m.setCursor = max(0, min(m.setCursor, count-1))

	lines := []string{}
	index := 0
	for _, session := range exercise.sessions {
		separator := dateLine(session.date.Format("2 January, 2006"), m.timeline.Width())
		for sessionIndex, set := range session.sets {
			pageHeight := m.setPager.PerPage
			needsDate := sessionIndex == 0 || len(lines)%pageHeight == 0
			needed := 1
			if needsDate && pageHeight > 1 {
				needed++
			}

			if remaining := pageHeight - len(lines)%pageHeight; remaining < needed {
				lines = append(lines, make([]string, remaining)...)
				needsDate = true
			}

			if len(lines)%pageHeight == 0 {
				m.setPageStarts = append(m.setPageStarts, index)
			}

			if needsDate && pageHeight > 1 {
				lines = append(lines, lipgloss.NewStyle().Foreground(components.BaseStyle.GetBorderTopForeground()).Render(separator))
			}

			m.setLines = append(m.setLines, len(lines))

			style := lipgloss.NewStyle()
			if index == m.setCursor {
				style = components.SelectedStyle()
			}

			restSeconds := int(set.rest.Seconds())
			values := []string{
				fmt.Sprint(set.reps),
				fmt.Sprintf("%g %s", set.weight, strings.ToUpper(string(m.weightUnit))),
				fmt.Sprint(restSeconds),
			}

			lines = append(lines, m.renderSetRow(values, style))
			index++
		}
	}

	if len(m.setLines) > 0 {
		m.setCursor = min(max(0, m.setCursor), len(m.setLines)-1)
		m.setPager.SetTotalPages(len(lines))
		m.setPager.Page = m.setLines[m.setCursor] / m.setPager.PerPage

		// Pad the last page so the viewport can align it at the top.
		lines = append(lines, make([]string, m.setPager.TotalPages*m.setPager.PerPage-len(lines))...)
	}

	m.timeline.SetContent(strings.Join(lines, "\n"))
	m.timeline.SetYOffset(m.setPager.Page * m.setPager.PerPage)
}
